// Package models contient les structs GORM mappées sur les tables SQLite.
package models

import (
	"time"

	"gorm.io/gorm"
)

// User — le téléphone est l'identifiant principal, l'email est optionnel.
type User struct {
	gorm.Model
	Phone        string  `gorm:"uniqueIndex;not null"    json:"phone"`           // E.164: +212XXXXXXXXX
	Email        *string `gorm:"uniqueIndex"             json:"email,omitempty"` // optionnel
	PasswordHash *string `gorm:"column:password_hash"    json:"-"`               // NULL pour auth phone-only
	IsVerified   bool    `gorm:"default:false"           json:"is_verified"`
	DisplayName  string  `gorm:"not null"                json:"display_name"`
	AvatarURL    *string `gorm:"column:avatar_url"       json:"avatar_url,omitempty"`
	City         *string `gorm:"column:city"             json:"city,omitempty"`
	Role         string  `gorm:"default:'user';not null" json:"role"` // user | admin
}

// OTPCode — code à usage unique envoyé par SMS ou WhatsApp.
type OTPCode struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Phone     string     `gorm:"not null;index"           json:"phone"`
	CodeHash  string     `gorm:"not null"                 json:"-"`       // bcrypt du code
	Channel   string     `gorm:"not null"                 json:"channel"` // sms | whatsapp
	Purpose   string     `gorm:"not null"                 json:"purpose"` // login | signup | phone_change
	ExpiresAt time.Time  `gorm:"not null"                 json:"expires_at"`
	UsedAt    *time.Time `gorm:"column:used_at"           json:"used_at,omitempty"`
	Attempts  int        `gorm:"default:0"                json:"attempts"`
	CreatedAt time.Time  `gorm:"autoCreateTime"           json:"created_at"`
}

// RefreshToken — stocké hashé (SHA-256), jamais en clair.
type RefreshToken struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint       `gorm:"not null;index"           json:"-"`
	TokenHash string     `gorm:"uniqueIndex;not null"     json:"-"`
	ExpiresAt time.Time  `gorm:"not null"                 json:"expires_at"`
	RevokedAt *time.Time `gorm:"column:revoked_at"        json:"revoked_at,omitempty"`
	CreatedAt time.Time  `gorm:"autoCreateTime"           json:"created_at"`
}

func (rt *RefreshToken) IsValid() bool {
	return rt.RevokedAt == nil && time.Now().Before(rt.ExpiresAt)
}
