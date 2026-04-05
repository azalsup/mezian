package repository

import (
	"mezian/internal/models"

	"gorm.io/gorm"
)

// ShopRepo gère les opérations base de données sur les boutiques.
type ShopRepo struct{ db *gorm.DB }

// NewShopRepo crée un nouveau ShopRepo.
func NewShopRepo(db *gorm.DB) *ShopRepo { return &ShopRepo{db} }

// Create insère une nouvelle boutique.
func (r *ShopRepo) Create(shop *models.Shop) error {
	return r.db.Create(shop).Error
}

// Update sauvegarde les modifications d'une boutique.
func (r *ShopRepo) Update(shop *models.Shop) error {
	return r.db.Save(shop).Error
}

// FindByUserID récupère la boutique d'un utilisateur.
func (r *ShopRepo) FindByUserID(userID uint) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.Where("user_id = ?", userID).First(&shop).Error
	return &shop, err
}

// FindBySlug récupère une boutique par son slug.
func (r *ShopRepo) FindBySlug(slug string) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.
		Preload("User").
		Where("slug = ? AND is_active = ?", slug, true).
		First(&shop).Error
	return &shop, err
}

// FindByID récupère une boutique par son ID.
func (r *ShopRepo) FindByID(id uint) (*models.Shop, error) {
	var shop models.Shop
	err := r.db.Preload("User").First(&shop, id).Error
	return &shop, err
}

// ExistsByUserID retourne true si l'utilisateur possède déjà une boutique.
func (r *ShopRepo) ExistsByUserID(userID uint) bool {
	var count int64
	r.db.Model(&models.Shop{}).Where("user_id = ?", userID).Count(&count)
	return count > 0
}

// ExistsBySlug retourne true si le slug est déjà pris.
func (r *ShopRepo) ExistsBySlug(slug string) bool {
	var count int64
	r.db.Model(&models.Shop{}).Where("slug = ?", slug).Count(&count)
	return count > 0
}
