// Package models - Ad, AdAttribute
package models

import (
    "time"

    "gorm.io/gorm"
)

// Ad — listing. The body is stored in Markdown.
type Ad struct {
    gorm.Model
    UserID       uint       `gorm:"not null;index"                   json:"user_id"`
    CategoryID   uint       `gorm:"not null;index"                   json:"category_id"`
    ShopID       *uint      `gorm:"index"                            json:"shop_id,omitempty"`
    Slug         string     `gorm:"uniqueIndex;not null"             json:"slug"`
    Title        string     `gorm:"not null"                         json:"title"`
    Body         string     `gorm:"not null;type:text"               json:"body"`          // Markdown
    Price        *float64   `gorm:"column:price"                     json:"price,omitempty"`
    Currency     string     `gorm:"default:'MAD';not null"           json:"currency"`
    City         string     `gorm:"not null;index"                   json:"city"`
    District     *string    `gorm:"index"                            json:"district,omitempty"`
    Status       string     `gorm:"default:'active';not null;index"  json:"status"` // draft|active|sold|expired|deleted
    IsBoosted    bool       `gorm:"default:false"                    json:"is_boosted"`
    BoostedUntil *time.Time `gorm:"index"                            json:"boosted_until,omitempty"`
    ViewCount    int        `gorm:"default:0"                        json:"view_count"`

    User       User          `gorm:"foreignKey:UserID"     json:"user,omitempty"`
    Category   Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
    Shop       *Shop         `gorm:"foreignKey:ShopID"     json:"shop,omitempty"`
    Media      []Media       `gorm:"foreignKey:AdID"       json:"media,omitempty"`
    Attributes []AdAttribute `gorm:"foreignKey:AdID"       json:"attributes,omitempty"`
}

// AdAttribute — value of a specific attribute (area, mileage, etc.)
type AdAttribute struct {
    ID    uint   `gorm:"primaryKey;autoIncrement"              json:"id"`
    AdID  uint   `gorm:"not null;uniqueIndex:idx_ad_key;index" json:"ad_id"`
    Key   string `gorm:"not null;uniqueIndex:idx_ad_key"       json:"key"`
    Value string `gorm:"not null"                              json:"value"` // toujours string, converti côté client
}
