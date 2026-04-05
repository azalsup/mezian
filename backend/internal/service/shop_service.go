package service

import (
	"errors"
	"fmt"
	"mezian/internal/config"
	"mezian/internal/models"
	"mezian/internal/repository"

	"gorm.io/gorm"
)

// Business errors for shops.
var (
	ErrShopNotFound      = errors.New("shop not found")
	ErrShopForbidden     = errors.New("unauthorized access to this shop")
	ErrShopAlreadyExists = errors.New("you already own a shop")
	ErrShopSlugTaken     = errors.New("this shop name is already taken")
)

// CreateShopInput groups the data required to create a shop.
type CreateShopInput struct {
	UserID      uint
	Name        string
	Description *string
	Phone       string
	City        string
	Plan        string
}

// UpdateShopInput groups the editable fields of a shop.
type UpdateShopInput struct {
	Name        *string
	Description *string
	Phone       *string
	City        *string
	LogoURL     *string
	CoverURL    *string
}

// ShopService handles the business logic for shops.
type ShopService struct {
	shopRepo *repository.ShopRepo
	adRepo   *repository.AdRepo
	cfg      *config.Config
}

// NewShopService creates a new ShopService.
func NewShopService(shopRepo *repository.ShopRepo, adRepo *repository.AdRepo, cfg *config.Config) *ShopService {
	return &ShopService{
		shopRepo: shopRepo,
		adRepo:   adRepo,
		cfg:      cfg,
	}
}

// CreateShop creates a new shop for a user.
func (s *ShopService) CreateShop(input CreateShopInput) (*models.Shop, error) {
	if s.shopRepo.ExistsByUserID(input.UserID) {
		return nil, ErrShopAlreadyExists
	}

	// Generate a unique slug from the name
	baseSlug := slugify(input.Name)
	if baseSlug == "" {
		baseSlug = fmt.Sprintf("shop-%d", input.UserID)
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
		return nil, fmt.Errorf("creation shop: %w", err)
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
		return nil, fmt.Errorf("update shop: %w", err)
	}

	return shop, nil
}

// GetShopBySlug retrieves a public shop by its slug.
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

// GetMyShop retrieves the authenticated user's shop.
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

// CanPublishAd checks whether the shop can still publish ads according to its plan.
func (s *ShopService) CanPublishAd(shopID uint) (bool, error) {
	shop, err := s.shopRepo.FindByID(shopID)
	if err != nil {
		return false, ErrShopNotFound
	}

	if !shop.IsSubscriptionValid() {
		return false, nil
	}

	// Get the limit by plan
	maxAds := s.maxAdsForPlan(shop.Plan)
	if maxAds == -1 {
		return true, nil // Unlimited
	}

	count, err := s.adRepo.CountActiveByShop(shopID)
	if err != nil {
		return false, fmt.Errorf("comptage annonces: %w", err)
	}

	return count < int64(maxAds), nil
}

// maxAdsForPlan returns the ads limit for the plan.
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
