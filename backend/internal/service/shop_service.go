package service

import (
	"errors"
	"fmt"
	"mezian/internal/config"
	"mezian/internal/models"
	"mezian/internal/repository"

	"gorm.io/gorm"
)

// Erreurs métier des boutiques.
var (
	ErrShopNotFound      = errors.New("boutique introuvable")
	ErrShopForbidden     = errors.New("accès non autorisé à cette boutique")
	ErrShopAlreadyExists = errors.New("vous possédez déjà une boutique")
	ErrShopSlugTaken     = errors.New("ce nom de boutique est déjà pris")
)

// CreateShopInput regroupe les données nécessaires à la création d'une boutique.
type CreateShopInput struct {
	UserID      uint
	Name        string
	Description *string
	Phone       string
	City        string
	Plan        string
}

// UpdateShopInput regroupe les champs modifiables d'une boutique.
type UpdateShopInput struct {
	Name        *string
	Description *string
	Phone       *string
	City        *string
	LogoURL     *string
	CoverURL    *string
}

// ShopService gère la logique métier des boutiques.
type ShopService struct {
	shopRepo *repository.ShopRepo
	adRepo   *repository.AdRepo
	cfg      *config.Config
}

// NewShopService crée un nouveau ShopService.
func NewShopService(shopRepo *repository.ShopRepo, adRepo *repository.AdRepo, cfg *config.Config) *ShopService {
	return &ShopService{
		shopRepo: shopRepo,
		adRepo:   adRepo,
		cfg:      cfg,
	}
}

// CreateShop crée une nouvelle boutique pour un utilisateur.
func (s *ShopService) CreateShop(input CreateShopInput) (*models.Shop, error) {
	if s.shopRepo.ExistsByUserID(input.UserID) {
		return nil, ErrShopAlreadyExists
	}

	// Générer un slug unique depuis le nom
	baseSlug := slugify(input.Name)
	if baseSlug == "" {
		baseSlug = fmt.Sprintf("boutique-%d", input.UserID)
	}

	slug := baseSlug
	counter := 1
	for s.shopRepo.ExistsBySlug(slug) {
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	plan := input.Plan
	if plan == "" {
		plan = "starter"
	}

	shop := &models.Shop{
		UserID:      input.UserID,
		Slug:        slug,
		Name:        input.Name,
		Description: input.Description,
		Phone:       input.Phone,
		City:        input.City,
		Plan:        plan,
		IsActive:    true,
	}

	if err := s.shopRepo.Create(shop); err != nil {
		return nil, fmt.Errorf("création boutique: %w", err)
	}

	return shop, nil
}

// UpdateShop modifie les informations d'une boutique.
func (s *ShopService) UpdateShop(slug string, requesterID uint, input UpdateShopInput) (*models.Shop, error) {
	shop, err := s.shopRepo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrShopNotFound
		}
		return nil, err
	}

	if shop.UserID != requesterID {
		return nil, ErrShopForbidden
	}

	if input.Name != nil {
		shop.Name = *input.Name
	}
	if input.Description != nil {
		shop.Description = input.Description
	}
	if input.Phone != nil {
		shop.Phone = *input.Phone
	}
	if input.City != nil {
		shop.City = *input.City
	}
	if input.LogoURL != nil {
		shop.LogoURL = input.LogoURL
	}
	if input.CoverURL != nil {
		shop.CoverURL = input.CoverURL
	}

	if err := s.shopRepo.Update(shop); err != nil {
		return nil, fmt.Errorf("mise à jour boutique: %w", err)
	}

	return shop, nil
}

// GetShopBySlug récupère une boutique publique par son slug.
func (s *ShopService) GetShopBySlug(slug string) (*models.Shop, error) {
	shop, err := s.shopRepo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrShopNotFound
		}
		return nil, err
	}
	return shop, nil
}

// GetMyShop récupère la boutique de l'utilisateur connecté.
func (s *ShopService) GetMyShop(userID uint) (*models.Shop, error) {
	shop, err := s.shopRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrShopNotFound
		}
		return nil, err
	}
	return shop, nil
}

// CanPublishAd vérifie si la boutique peut encore publier des annonces selon son plan.
func (s *ShopService) CanPublishAd(shopID uint) (bool, error) {
	shop, err := s.shopRepo.FindByID(shopID)
	if err != nil {
		return false, ErrShopNotFound
	}

	if !shop.IsSubscriptionValid() {
		return false, nil
	}

	// Obtenir la limite selon le plan
	maxAds := s.maxAdsForPlan(shop.Plan)
	if maxAds == -1 {
		return true, nil // Illimité
	}

	count, err := s.adRepo.CountActiveByShop(shopID)
	if err != nil {
		return false, fmt.Errorf("comptage annonces: %w", err)
	}

	return count < int64(maxAds), nil
}

// maxAdsForPlan retourne la limite d'annonces selon le plan.
func (s *ShopService) maxAdsForPlan(plan string) int {
	switch plan {
	case "pro":
		return s.cfg.Plans.Pro.MaxAds
	case "premium":
		return s.cfg.Plans.Premium.MaxAds
	default: // starter
		return s.cfg.Plans.Starter.MaxAds
	}
}
