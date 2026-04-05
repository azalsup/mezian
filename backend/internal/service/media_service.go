package service

import (
	"errors"
	"fmt"
	"image"
	"mezian/internal/config"
	"mezian/internal/models"
	"mezian/internal/repository"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

// Business errors for media.
var (
	ErrMediaNotFound    = errors.New("media introuvable")
	ErrMediaForbidden   = errors.New("unauthorized access to this media")
	ErrTooManyMedia     = errors.New("nombre maximum de media atteint pour cette ad")
	ErrInvalidMediaType = errors.New("unauthorized file type")
	ErrFileTooLarge     = errors.New("fichier trop volumineux")
	ErrInvalidYouTube   = errors.New("URL YouTube invalide")
)

// youtubeRegexp matche les URLs YouTube classiques et courtes.
var youtubeRegexp = regexp.MustCompile(
	`(?:youtube\.com/(?:watch\?v=|embed/|v/)|youtu\.be/)([a-zA-Z0-9_-]{11})`,
)

// MediaService gère les uploads d'images et les liens YouTube.
type MediaService struct {
	mediaRepo *repository.MediaRepo
	adRepo    *repository.AdRepo
	cfg       *config.Config
}

// NewMediaService creates un nouveau MediaService.
func NewMediaService(mediaRepo *repository.MediaRepo, adRepo *repository.AdRepo, cfg *config.Config) *MediaService {
	return &MediaService{
		mediaRepo: mediaRepo,
		adRepo:    adRepo,
		cfg:       cfg,
	}
}

// UploadImage traite l'upload d'une image, la redimensionne et génère une miniature.
func (s *MediaService) UploadImage(adID, userID uint, file multipart.File, header *multipart.FileHeader) (*models.Media, error) {
	// Verify ad ownership
	ad, err := s.adRepo.FindByID(adID)
	if err != nil {
		return nil, ErrAdNotFound
	}
	if ad.UserID != userID {
		return nil, ErrMediaForbidden
	}

	// Check the number of existing media
	count, err := s.mediaRepo.CountByAdID(adID)
	if err != nil {
		return nil, fmt.Errorf("comptage media: %w", err)
	}
	if int(count) >= s.cfg.Media.MaxPerAd {
		return nil, ErrTooManyMedia
	}

	// Check the MIME type
	contentType := header.Header.Get("Content-Type")
	if !s.isAllowedType(contentType) {
		return nil, ErrInvalidMediaType
	}

	// Check the size
	maxBytes := int64(s.cfg.Media.MaxSizeMB) * 1024 * 1024
	if header.Size > maxBytes {
		return nil, ErrFileTooLarge
	}

	// Create the upload directory if needed
	uploadDir := filepath.Join(s.cfg.Media.UploadDir, "ads", fmt.Sprintf("%d", adID))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, fmt.Errorf("creation directory: %w", err)
	}
	thumbDir := filepath.Join(s.cfg.Media.UploadDir, "thumbs", fmt.Sprintf("%d", adID))
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		return nil, fmt.Errorf("creation directory thumbs: %w", err)
	}

	// Decode the image
	img, err := decodeImage(file, contentType)
	if err != nil {
		return nil, fmt.Errorf("image decoding: %w", err)
	}

	// Redimensionner l'image principale (max 1200px de large)
	resized := imaging.Fit(img, 1200, 1200, imaging.Lanczos)

	// Generate a unique filename
	ext := contentTypeToExt(contentType)
	filename := uuid.NewString() + ext

	// Sauvegarder l'image principale
	imagePath := filepath.Join(uploadDir, filename)
	if err := imaging.Save(resized, imagePath); err != nil {
		return nil, fmt.Errorf("sauvegarde image: %w", err)
	}

	// Generate and save the thumbnail
	thumb := imaging.Fill(img,
		s.cfg.Media.ThumbnailWidth,
		s.cfg.Media.ThumbnailHeight,
		imaging.Center, imaging.Lanczos,
	)
	thumbFilename := "thumb_" + filename
	thumbPath := filepath.Join(thumbDir, thumbFilename)
	if err := imaging.Save(thumb, thumbPath); err != nil {
		// Do not block if thumbnail creation fails
		thumbPath = ""
	}

	// Construire les URLs relatives
	imageURL := fmt.Sprintf("/uploads/ads/%d/%s", adID, filename)
	var thumbURL *string
	if thumbPath != "" {
		tu := fmt.Sprintf("/uploads/thumbs/%d/%s", adID, thumbFilename)
		thumbURL = &tu
	}

	// Determine if this is the first image (default cover)
	isCover := count == 0

	media := &models.Media{
		AdID:      adID,
		Type:      "image",
		URL:       imageURL,
		ThumbURL:  thumbURL,
		SortOrder: int(count),
		IsCover:   isCover,
	}

	if err := s.mediaRepo.Create(media); err != nil {
		// Nettoyer les fichiers en cas d'erreur DB
		os.Remove(imagePath)    //nolint:errcheck
		os.Remove(thumbPath)    //nolint:errcheck
		return nil, fmt.Errorf("sauvegarde media: %w", err)
	}

	return media, nil
}

// AddYouTube adds a YouTube video link to an ad.
func (s *MediaService) AddYouTube(adID, userID uint, youtubeURL string) (*models.Media, error) {
	// Verify ad ownership
	ad, err := s.adRepo.FindByID(adID)
	if err != nil {
		return nil, ErrAdNotFound
	}
	if ad.UserID != userID {
		return nil, ErrMediaForbidden
	}

	// Check the number of existing media
	count, err := s.mediaRepo.CountByAdID(adID)
	if err != nil {
		return nil, fmt.Errorf("comptage media: %w", err)
	}
	if int(count) >= s.cfg.Media.MaxPerAd {
		return nil, ErrTooManyMedia
	}

	// Extract the YouTube video ID
	videoID, err := extractYouTubeID(youtubeURL)
	if err != nil {
		return nil, ErrInvalidYouTube
	}

	// URL canonique et miniature YouTube
	canonicalURL := fmt.Sprintf("https://www.youtube.com/watch?v=%s", videoID)
	thumbURL := fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", videoID)

	media := &models.Media{
		AdID:      adID,
		Type:      "youtube",
		URL:       canonicalURL,
		ThumbURL:  &thumbURL,
		SortOrder: int(count),
		IsCover:   false,
	}

	if err := s.mediaRepo.Create(media); err != nil {
		return nil, fmt.Errorf("sauvegarde media YouTube: %w", err)
	}

	return media, nil
}

// DeleteMedia removes a media item (file + DB entry).
func (s *MediaService) DeleteMedia(mediaID, userID uint, userRole string) error {
	media, err := s.mediaRepo.FindByID(mediaID)
	if err != nil {
		return ErrMediaNotFound
	}

	// Verify ownership via the ad
	ad, err := s.adRepo.FindByID(media.AdID)
	if err != nil {
		return ErrAdNotFound
	}
	if ad.UserID != userID && userRole != "admin" {
		return ErrMediaForbidden
	}

	// Delete the physical file if it is a local image
	if media.Type == "image" {
		localPath := filepath.Join(s.cfg.Media.UploadDir, strings.TrimPrefix(media.URL, "/uploads/"))
		os.Remove(localPath) //nolint:errcheck
		if media.ThumbURL != nil {
			thumbLocal := filepath.Join(s.cfg.Media.UploadDir, strings.TrimPrefix(*media.ThumbURL, "/uploads/"))
			os.Remove(thumbLocal) //nolint:errcheck
		}
	}

	return s.mediaRepo.Delete(mediaID)
}

// SetCover définit un media comme image de couverture.
func (s *MediaService) SetCover(mediaID, userID uint) error {
	media, err := s.mediaRepo.FindByID(mediaID)
	if err != nil {
		return ErrMediaNotFound
	}

	ad, err := s.adRepo.FindByID(media.AdID)
	if err != nil {
		return ErrAdNotFound
	}
	if ad.UserID != userID {
		return ErrMediaForbidden
	}

	return s.mediaRepo.SetCover(media.AdID, mediaID)
}

// UpdateOrder updates l'ordre d'affichage d'un media.
func (s *MediaService) UpdateOrder(mediaID, userID uint, sortOrder int) error {
	media, err := s.mediaRepo.FindByID(mediaID)
	if err != nil {
		return ErrMediaNotFound
	}

	ad, err := s.adRepo.FindByID(media.AdID)
	if err != nil {
		return ErrAdNotFound
	}
	if ad.UserID != userID {
		return ErrMediaForbidden
	}

	return s.mediaRepo.UpdateOrder(mediaID, sortOrder)
}

// GetAdMedia returns tous les media d'une ad.
func (s *MediaService) GetAdMedia(adID uint) ([]models.Media, error) {
	return s.mediaRepo.FindByAdID(adID)
}

// isAllowedType verifies si le type MIME est autorisé.
func (s *MediaService) isAllowedType(contentType string) bool {
	for _, t := range s.cfg.Media.AllowedTypes {
		if t == contentType {
			return true
		}
	}
	return false
}

// decodeImage décode une image depuis un reader selon son type MIME.
func decodeImage(file multipart.File, contentType string) (image.Image, error) {
	switch contentType {
	case "image/jpeg":
		return imaging.Decode(file, imaging.AutoOrientation(true))
	case "image/png", "image/webp":
		return imaging.Decode(file)
	default:
		return nil, fmt.Errorf("unsupported type: %s", contentType)
	}
}

// contentTypeToExt retourne l'extension de fichier pour un type MIME.
func contentTypeToExt(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	default:
		return ".jpg"
	}
}

// extractYouTubeID extracts a YouTube video ID from a YouTube URL.
func extractYouTubeID(rawURL string) (string, error) {
	// Essayer le regex d'abord
	matches := youtubeRegexp.FindStringSubmatch(rawURL)
	if len(matches) >= 2 {
		return matches[1], nil
	}

	// Essayer via parsing URL (cas youtu.be)
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", ErrInvalidYouTube
	}

	switch u.Host {
	case "youtu.be":
		id := strings.TrimPrefix(u.Path, "/")
		if len(id) == 11 {
			return id, nil
		}
	case "www.youtube.com", "youtube.com":
		id := u.Query().Get("v")
		if len(id) == 11 {
			return id, nil
		}
	}

	return "", ErrInvalidYouTube
}
