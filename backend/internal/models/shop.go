package models

import (
	"time"

	"gorm.io/gorm"
)

// Shop — paid professional shop. A user can have an active shop.
type Shop struct {
	gorm.Model
	UserID      uint       `gorm:"uniqueIndex;not null"       json:"user_id"`
	Slug        string     `gorm:"uniqueIndex;not null"       json:"slug"`
	Name        string     `gorm:"not null"                   json:"name"`
	Description *string    `gorm:"type:text"                  json:"description,omitempty"` // Markdown
	LogoURL     *string    `gorm:"column:logo_url"            json:"logo_url,omitempty"`
	CoverURL    *string    `gorm:"column:cover_url"           json:"cover_url,omitempty"`
	Phone       string     `gorm:"not null"                   json:"phone"`
	City        string     `gorm:"not null"                   json:"city"`
	Plan        string     `gorm:"default:'starter';not null" json:"plan"` // starter|pro|premium
	PlanExpires *time.Time `gorm:"index"                      json:"plan_expires,omitempty"`
	IsActive    bool       `gorm:"default:true"               json:"is_active"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Ads  []Ad `gorm:"foreignKey:ShopID" json:"ads,omitempty"`
}

// IsSubscriptionValid retourne true si le plan est actif.
func (s *Shop) IsSubscriptionValid() bool {
	if s.PlanExpires == nil {
		return false
	}
	return time.Now().Before(*s.PlanExpires)
}
