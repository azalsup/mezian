package repository

import (
    "classifieds/internal/models"
    "strings"

    "gorm.io/gorm"
)

// AdFilters groups all filters applicable to the ad list.
type AdFilters struct {
    CategoryID      *uint
    CategorySlug    string
    SubcategorySlug string
    City            string
    MinPrice        *float64
    MaxPrice        *float64
    Status          string // vide = "active" par défaut
    Search          string // fulltext search on title
    UserID          *uint
    ShopID          *uint
    Page            int
    Limit           int
    Sort            string // price_asc | price_desc | newest | oldest | views
}

// AdRepo handles database operations for ads.
type AdRepo struct{ db *gorm.DB }

// NewAdRepo creates a new AdRepo.
func NewAdRepo(db *gorm.DB) *AdRepo { return &AdRepo{db} }

// Create inserts a new ad with its attributes.
func (r *AdRepo) Create(ad *models.Ad) error {
    return r.db.Create(ad).Error
}

// Update saves changes to an ad.
func (r *AdRepo) Update(ad *models.Ad) error {
    return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(ad).Error
}

// FindByID retrieves a complete ad (with relations) by its ID.
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

// FindBySlug retrieves a complete ad by its slug.
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

// AdListResult contains the paginated ads and total.
type AdListResult struct {
    Ads   []models.Ad
    Total int64
    Page  int
    Limit int
}

// List returns filtered and paginated ads.
func (r *AdRepo) List(f AdFilters) (*AdListResult, error) {
    query := r.db.Model(&models.Ad{})

    // Status filter (default: active)
    status := f.Status
    if status == "" {
        status = "active"
    }
    query = query.Where("ads.status = ?", status)

    if f.CategoryID != nil {
        query = query.Where("ads.category_id = ?", *f.CategoryID)
    }
    if f.SubcategorySlug != "" {
        query = query.Joins("JOIN categories ON categories.id = ads.category_id").
            Where("categories.slug = ?", f.SubcategorySlug)
    } else if f.CategorySlug != "" {
        query = query.Joins("JOIN categories ON categories.id = ads.category_id").
            Where("categories.slug = ? OR categories.parent_id = (SELECT id FROM categories WHERE slug = ?)",
                f.CategorySlug, f.CategorySlug)
    }
    if f.City != "" {
        query = query.Where("ads.city = ?", f.City)
    }
    if f.MinPrice != nil {
        query = query.Where("ads.price >= ?", *f.MinPrice)
    }
    if f.MaxPrice != nil {
        query = query.Where("ads.price <= ?", *f.MaxPrice)
    }
    if f.Search != "" {
        like := "%" + strings.ToLower(f.Search) + "%"
        query = query.Where("LOWER(ads.title) LIKE ? OR LOWER(ads.body) LIKE ?", like, like)
    }
    if f.UserID != nil {
        query = query.Where("ads.user_id = ?", *f.UserID)
    }
    if f.ShopID != nil {
        query = query.Where("ads.shop_id = ?", *f.ShopID)
    }

    var total int64
    if err := query.Count(&total).Error; err != nil {
        return nil, err
    }

    // Tri
    switch f.Sort {
    case "price_asc":
        query = query.Order("ads.price ASC")
    case "price_desc":
        query = query.Order("ads.price DESC")
    case "oldest":
        query = query.Order("ads.created_at ASC")
    case "views":
        query = query.Order("ads.view_count DESC")
    default: // newest + boosted en tête
        query = query.Order("ads.is_boosted DESC, ads.created_at DESC")
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

// Delete deletes an ad (soft delete via gorm.Model).
func (r *AdRepo) Delete(id uint) error {
    return r.db.Delete(&models.Ad{}, id).Error
}

// IncrementViews increments the view counter atomically.
func (r *AdRepo) IncrementViews(id uint) error {
    return r.db.Model(&models.Ad{}).Where("id = ?", id).
        UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

// FindByUser returns a user's ads (all statuses).
func (r *AdRepo) FindByUser(userID uint, page, limit int) (*AdListResult, error) {
    return r.List(AdFilters{
        UserID: &userID,
        Status: "all",
        Page:   page,
        Limit:  limit,
        Sort:   "newest",
    })
}

// FindByShop returns a shop's active ads.
func (r *AdRepo) FindByShop(shopID uint, page, limit int) (*AdListResult, error) {
    return r.List(AdFilters{
        ShopID: &shopID,
        Status: "active",
        Page:   page,
        Limit:  limit,
        Sort:   "newest",
    })
}

// CountActiveByUser returns the number of active ads for a user.
func (r *AdRepo) CountActiveByUser(userID uint) (int64, error) {
    var count int64
    err := r.db.Model(&models.Ad{}).
        Where("user_id = ? AND status = 'active'", userID).
        Count(&count).Error
    return count, err
}

// CountActiveByShop returns the number of active ads for a shop.
func (r *AdRepo) CountActiveByShop(shopID uint) (int64, error) {
    var count int64
    err := r.db.Model(&models.Ad{}).
        Where("shop_id = ? AND status = 'active'", shopID).
        Count(&count).Error
    return count, err
}

// UpdateAttributes replaces an ad's attributes.
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
