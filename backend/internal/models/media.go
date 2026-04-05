package models

import "gorm.io/gorm"

// Media — photo (stored locally) or YouTube video linked to an ad.
type Media struct {
	gorm.Model
	AdID      uint    `gorm:"not null;index"   json:"ad_id"`
	Type      string  `gorm:"not null"         json:"type"`                 // image | youtube
	URL       string  `gorm:"not null"         json:"url"`                  // path relatif (image) ou URL YouTube
	ThumbURL  *string `gorm:"column:thumb_url" json:"thumb_url,omitempty"`
	SortOrder int     `gorm:"default:0"        json:"sort_order"`
	IsCover   bool    `gorm:"default:false"    json:"is_cover"`
}
