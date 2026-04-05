package repository

import (
	"mezian/internal/models"
	"strings"

	"gorm.io/gorm"
)

// AdFilters regroupe tous les filtres applicables à la liste des annonces.
type AdFilters struct {
	CategoryID *uint
	City       string
	MinPrice   *float64
	MaxPrice   *float64
	Status     string // vide = "active" par défaut
	Search     string // recherche fulltext sur titre
	UserID     *uint
	ShopID     *uint
	Page       int
	Limit      int
	Sort       string // price_asc | price_desc | newest | oldest | views
}

// AdRepo gère les opérations base de données sur les annonces.
type AdRepo struct{ db *gorm.DB }

// NewAdRepo crée un nouveau AdRepo.
func NewAdRepo(db *gorm.DB) *AdRepo { return &AdRepo{db} }

// Create insère une nouvelle annonce avec ses attributs.
func (r *AdRepo) Create(ad *models.Ad) error {
	return r.db.Create(ad).Error
}

// Update sauvegarde les modifications d'une annonce.
func (r *AdRepo) Update(ad *models.Ad) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(ad).Error
}

// FindByID récupère une annonce complète (avec relations) par son ID.
func (r *AdRepo) FindByID(id uint) (*models.Ad, error) {
	var ad models.Ad
	err := r.db.
		Preload("User").
		Preload("Category").
		Preload("Shop").
		Preload("Media", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, created_at ASC")
		}).
		Preload("Attributes").
		First(&ad, id).Error
	return &ad, err
}

// FindBySlug récupère une annonce complète par son slug.
func (r *AdRepo) FindBySlug(slug string) (*models.Ad, error) {
	var ad models.Ad
	err := r.db.
		Preload("User").
		Preload("Category").
		Preload("Shop").
		Preload("Media", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, created_at ASC")
		}).
		Preload("Attributes").
		Where("slug = ?", slug).
		First(&ad).Error
	return &ad, err
}

// AdListResult contient les annonces paginées et le total.
type AdListResult struct {
	Ads   []models.Ad
	Total int64
	Page  int
	Limit int
}

// List retourne les annonces filtrées et paginées.
func (r *AdRepo) List(f AdFilters) (*AdListResult, error) {
	query := r.db.Model(&models.Ad{})

	// Filtre statut (défaut: active)
	status := f.Status
	if status == "" {
		status = "active"
	}
	query = query.Where("status = ?", status)

	if f.CategoryID != nil {
		query = query.Where("category_id = ?", *f.CategoryID)
	}
	if f.City != "" {
		query = query.Where("city = ?", f.City)
	}
	if f.MinPrice != nil {
		query = query.Where("price >= ?", *f.MinPrice)
	}
	if f.MaxPrice != nil {
		query = query.Where("price <= ?", *f.MaxPrice)
	}
	if f.Search != "" {
		like := "%" + strings.ToLower(f.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(body) LIKE ?", like, like)
	}
	if f.UserID != nil {
		query = query.Where("user_id = ?", *f.UserID)
	}
	if f.ShopID != nil {
		query = query.Where("shop_id = ?", *f.ShopID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Tri
	switch f.Sort {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	case "oldest":
		query = query.Order("created_at ASC")
	case "views":
		query = query.Order("view_count DESC")
	default: // newest + boosted en tête
		query = query.Order("is_boosted DESC, created_at DESC")
	}

	// Pagination
	page := f.Page
	if page < 1 {
		page = 1
	}
	limit := f.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	var ads []models.Ad
	err := query.
		Offset(offset).
		Limit(limit).
		Preload("Category").
		Preload("Media", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_cover = ?", true).Limit(1)
		}).
		Find(&ads).Error
	if err != nil {
		return nil, err
	}

	return &AdListResult{
		Ads:   ads,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// Delete supprime une annonce (soft delete via gorm.Model).
func (r *AdRepo) Delete(id uint) error {
	return r.db.Delete(&models.Ad{}, id).Error
}

// IncrementViews incrémente le compteur de vues atomiquement.
func (r *AdRepo) IncrementViews(id uint) error {
	return r.db.Model(&models.Ad{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// FindByUser retourne les annonces d'un utilisateur (toutes statuts).
func (r *AdRepo) FindByUser(userID uint, page, limit int) (*AdListResult, error) {
	return r.List(AdFilters{
		UserID: &userID,
		Status: "all",
		Page:   page,
		Limit:  limit,
		Sort:   "newest",
	})
}

// FindByShop retourne les annonces actives d'une boutique.
func (r *AdRepo) FindByShop(shopID uint, page, limit int) (*AdListResult, error) {
	return r.List(AdFilters{
		ShopID: &shopID,
		Status: "active",
		Page:   page,
		Limit:  limit,
		Sort:   "newest",
	})
}

// CountActiveByUser retourne le nombre d'annonces actives d'un utilisateur.
func (r *AdRepo) CountActiveByUser(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Ad{}).
		Where("user_id = ? AND status = 'active'", userID).
		Count(&count).Error
	return count, err
}

// CountActiveByShop retourne le nombre d'annonces actives d'une boutique.
func (r *AdRepo) CountActiveByShop(shopID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Ad{}).
		Where("shop_id = ? AND status = 'active'", shopID).
		Count(&count).Error
	return count, err
}

// UpdateAttributes remplace les attributs d'une annonce.
func (r *AdRepo) UpdateAttributes(adID uint, attrs []models.AdAttribute) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("ad_id = ?", adID).Delete(&models.AdAttribute{}).Error; err != nil {
			return err
		}
		if len(attrs) > 0 {
			for i := range attrs {
				attrs[i].AdID = adID
			}
			return tx.Create(&attrs).Error
		}
		return nil
	})
}
