package repository

import (
	"mezian/internal/models"

	"gorm.io/gorm"
)

// ShopRepo gère les opérations database sur les shops.
type ShopRepo struct{ db *gorm.DB }

// NewShopRepo creates un nouveau ShopRepo.
func NewShopRepo(db *gorm.DB) *ShopRepo { return &ShopRepo{db} }

// Create inserts a new shop.
func (r *ShopRepo) Create(shop *models.Shop) error {
	return r.db.Create(shop).Error
}

// Update saves changes to a shop.
func (r *ShopRepo) Update(shop *models.Shop) error {
	return r.db.Save(shop).Error
}

// FindByUserID retrieves a user's shop.
func (r *ShopRepo) FindByUserID(userID uint) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.Where("user_id = ?", userID).First(&shop).Error
	return &shop, err
}

// FindBySlug retrieves a shop by its slug.
func (r *ShopRepo) FindBySlug(slug string) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.
		Preload("User").
		Where("slug = ? AND is_active = ?", slug, true).
		First(&shop).Error
	return &shop, err
}

// FindByID retrieves a shop by its ID.
func (r *ShopRepo) FindByID(id uint) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.Preload("User").First(&shop, id).Error
	return &shop, err
}

// ExistsByUserID returns true if the user already has a shop.
func (r *ShopRepo) ExistsByUserID(userID uint) bool {
	var count int64
	r.db.Model(&models.Shop{}).Where("user_id = ?", userID).Count(&count)
	return count > 0
}

// ExistsBySlug returns true if the slug is already taken.
func (r *ShopRepo) ExistsBySlug(slug string) bool {
	var count int64
	r.db.Model(&models.Shop{}).Where("slug = ?", slug).Count(&count)
	return count > 0
}
