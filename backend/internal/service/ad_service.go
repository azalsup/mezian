package service

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"mezian/internal/models"
	"mezian/internal/repository"
)

// Erreurs métier des annonces.
var (
	ErrAdNotFound   = errors.New("annonce introuvable")
	ErrAdForbidden  = errors.New("accès non autorisé à cette annonce")
	ErrShopAdsLimit = errors.New("limite d'annonces atteinte pour ce plan")
)

// CreateAdInput regroupe les données nécessaires à la création d'une annonce.
type CreateAdInput struct {
	UserID     uint
	CategoryID uint
	ShopID     *uint
	Title      string
	Body       string
	Price      *float64
	Currency   string
	City       string
	District   *string
	Status     string
	Attributes []models.AdAttribute
}

// UpdateAdInput regroupe les champs modifiables d'une annonce.
type UpdateAdInput struct {
	Title      *string
	Body       *string
	Price      *float64
	Currency   *string
	City       *string
	District   *string
	Status     *string
	Attributes []models.AdAttribute
}

// AdService contient la logique métier des annonces.
type AdService struct {
	adRepo   *repository.AdRepo
	shopRepo *repository.ShopRepo
}

// NewAdService crée un nouveau AdService.
func NewAdService(adRepo *repository.AdRepo, shopRepo *repository.ShopRepo) *AdService {
	return &AdService{adRepo: adRepo, shopRepo: shopRepo}
}

// CreateAd crée une nouvelle annonce après vérification des droits et des limites de plan.
func (s *AdService) CreateAd(input CreateAdInput) (*models.Ad, error) {
	// Vérifier les limites de la boutique si applicable
	if input.ShopID != nil {
		if err := s.checkShopAdLimit(*input.ShopID); err != nil {
			return nil, err
		}
	}

	slug, err := s.generateUniqueSlug(input.Title)
	if err != nil {
		return nil, fmt.Errorf("génération slug: %w", err)
	}

	currency := input.Currency
	if currency == "" {
		currency = "MAD"
	}
	status := input.Status
	if status == "" {
		status = "active"
	}

	ad := &models.Ad{
		UserID:     input.UserID,
		CategoryID: input.CategoryID,
		ShopID:     input.ShopID,
		Slug:       slug,
		Title:      input.Title,
		Body:       input.Body,
		Price:      input.Price,
		Currency:   currency,
		City:       input.City,
		District:   input.District,
		Status:     status,
		Attributes: input.Attributes,
	}

	if err := s.adRepo.Create(ad); err != nil {
		return nil, fmt.Errorf("création annonce: %w", err)
	}

	// Recharger avec toutes les relations
	return s.adRepo.FindByID(ad.Model.ID)
}

// UpdateAd modifie une annonce existante (owner ou admin uniquement).
func (s *AdService) UpdateAd(slug string, requesterID uint, requesterRole string, input UpdateAdInput) (*models.Ad, error) {
	ad, err := s.adRepo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAdNotFound
		}
		return nil, err
	}

	if ad.UserID != requesterID && requesterRole != "admin" {
		return nil, ErrAdForbidden
	}

	// Appliquer les modifications (uniquement les champs fournis)
	if input.Title != nil {
		ad.Title = *input.Title
	}
	if input.Body != nil {
		ad.Body = *input.Body
	}
	if input.Price != nil {
		ad.Price = input.Price
	}
	if input.Currency != nil {
		ad.Currency = *input.Currency
	}
	if input.City != nil {
		ad.City = *input.City
	}
	if input.District != nil {
		ad.District = input.District
	}
	if input.Status != nil {
		ad.Status = *input.Status
	}

	if err := s.adRepo.Update(ad); err != nil {
		return nil, fmt.Errorf("mise à jour annonce: %w", err)
	}

	// Mettre à jour les attributs si fournis
	if input.Attributes != nil {
		if err := s.adRepo.UpdateAttributes(ad.Model.ID, input.Attributes); err != nil {
			return nil, fmt.Errorf("mise à jour attributs: %w", err)
		}
	}

	return s.adRepo.FindBySlug(slug)
}

// DeleteAd supprime une annonce (owner ou admin uniquement).
func (s *AdService) DeleteAd(slug string, requesterID uint, requesterRole string) error {
	ad, err := s.adRepo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAdNotFound
		}
		return err
	}

	if ad.UserID != requesterID && requesterRole != "admin" {
		return ErrAdForbidden
	}

	return s.adRepo.Delete(ad.Model.ID)
}

// GetAd récupère une annonce par slug et incrémente le compteur de vues.
func (s *AdService) GetAd(slug string) (*models.Ad, error) {
	ad, err := s.adRepo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAdNotFound
		}
		return nil, err
	}

	// Incrémenter les vues de façon non bloquante
	adID := ad.Model.ID
	go s.adRepo.IncrementViews(adID) //nolint:errcheck

	ad.ViewCount++
	return ad, nil
}

// ListAds retourne les annonces filtrées et paginées.
func (s *AdService) ListAds(f repository.AdFilters) (*repository.AdListResult, error) {
	return s.adRepo.List(f)
}

// GetUserAds retourne les annonces d'un utilisateur.
func (s *AdService) GetUserAds(userID uint, page, limit int) (*repository.AdListResult, error) {
	return s.adRepo.FindByUser(userID, page, limit)
}

// checkShopAdLimit vérifie si la boutique peut encore publier des annonces.
func (s *AdService) checkShopAdLimit(shopID uint) error {
	shop, err := s.shopRepo.FindByID(shopID)
	if err != nil {
		return fmt.Errorf("boutique introuvable: %w", err)
	}

	if !shop.IsSubscriptionValid() {
		return ErrShopAdsLimit
	}

	// Les limites fines sont gérées par le ShopService
	return nil
}

// generateUniqueSlug crée un slug unique à partir du titre.
func (s *AdService) generateUniqueSlug(title string) (string, error) {
	base := slugify(title)
	if base == "" {
		base = "annonce"
	}

	// Suffixe UUID court pour garantir l'unicité sans requête DB
	suffix := uuid.NewString()[:8]
	return fmt.Sprintf("%s-%s", base, suffix), nil
}

// slugify convertit un texte en slug URL-safe (sans dépendances externes).
func slugify(s string) string {
	// Table de translittération pour les caractères courants (fr/ar → ASCII)
	replacer := strings.NewReplacer(
		"à", "a", "â", "a", "ä", "a", "á", "a", "ã", "a",
		"è", "e", "é", "e", "ê", "e", "ë", "e",
		"î", "i", "ï", "i", "í", "i", "ì", "i",
		"ô", "o", "ö", "o", "ó", "o", "ò", "o",
		"ù", "u", "û", "u", "ü", "u", "ú", "u",
		"ç", "c", "ñ", "n",
		"À", "a", "Â", "a", "Ä", "a", "Á", "a",
		"È", "e", "É", "e", "Ê", "e", "Ë", "e",
		"Î", "i", "Ï", "i",
		"Ô", "o", "Ö", "o",
		"Ù", "u", "Û", "u", "Ü", "u",
		"Ç", "c", "Ñ", "n",
		"&", "et", "+", "plus",
	)
	result := replacer.Replace(s)

	// Minuscules
	result = strings.ToLower(result)

	// Supprimer les caractères non-ASCII restants (arabe, etc.) et remplacer par des tirets
	var sb strings.Builder
	for _, r := range result {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if r < 128 { // ASCII uniquement
				sb.WriteRune(r)
			}
		} else {
			sb.WriteRune('-')
		}
	}
	result = sb.String()

	// Collapseur de tirets multiples
	re := regexp.MustCompile(`-+`)
	result = re.ReplaceAllString(result, "-")

	// Supprimer les tirets en début et fin
	result = strings.Trim(result, "-")

	// Limiter la longueur à 80 caractères
	if len(result) > 80 {
		result = result[:80]
		result = strings.TrimRight(result, "-")
	}

	return result
}
