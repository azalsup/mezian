package repository

import (
    "mezian/internal/models"

    "gorm.io/gorm"
)

// MediaRepo gère les opérations database sur les media.
type MediaRepo struct{ db *gorm.DB }

// NewMediaRepo creates un nouveau MediaRepo.
func NewMediaRepo(db *gorm.DB) *MediaRepo { return &MediaRepo{db} }

// Create inserts a new media record.
func (r *MediaRepo) Create(m *models.Media) error {
    return r.db.Create(m).Error
}

// FindByID retrieves media by its ID.
func (r *MediaRepo) FindByID(id uint) (*models.Media, error) {
    var m models.Media
    err := r.db.First(&m, id).Error
    return &m, err
}

// FindByAdID returns all media items for an ad, sorted by sort_order.
func (r *MediaRepo) FindByAdID(adID uint) ([]models.Media, error) {
    var media []models.Media
    err := r.db.Where("ad_id = ?", adID).
        Order("sort_order ASC, created_at ASC").
        Find(&media).Error
    return media, err
}

// CountByAdID returns the number of media items for an ad.
func (r *MediaRepo) CountByAdID(adID uint) (int64, error) {
    var count int64
    err := r.db.Model(&models.Media{}).Where("ad_id = ?", adID).Count(&count).Error
    return count, err
}

// Delete removes a media record by its ID (soft delete).
func (r *MediaRepo) Delete(id uint) error {
    return r.db.Delete(&models.Media{}, id).Error
}

// SetCover définit un media comme image de couverture (et retire le flag des autres).
func (r *MediaRepo) SetCover(adID, mediaID uint) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&models.Media{}).Where("ad_id = ?", adID).
            Update("is_cover", false).Error; err != nil {
            return err
        }
        return tx.Model(&models.Media{}).Where("id = ? AND ad_id = ?", mediaID, adID).
            Update("is_cover", true).Error
    })
}

// UpdateOrder updates le sort_order d'un media.
func (r *MediaRepo) UpdateOrder(id uint, sortOrder int) error {
    return r.db.Model(&models.Media{}).Where("id = ?", id).
        Update("sort_order", sortOrder).Error
}

// Update sauvegarde les modifications d'un media.
func (r *MediaRepo) Update(m *models.Media) error {
    return r.db.Save(m).Error
}
