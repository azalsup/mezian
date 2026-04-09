// Package models contains GORM structs mapped to SQLite tables.
package models

import (
    "time"

    "gorm.io/gorm"
)

// User — phone is the primary identifier, email is optional.
type User struct {
    gorm.Model
    Phone        string  `gorm:"uniqueIndex;not null"    json:"phone"`           // E.164: +212XXXXXXXXX
    Email        *string `gorm:"uniqueIndex"             json:"email,omitempty"` // optionnel
    PasswordHash *string `gorm:"column:password_hash"    json:"-"`               // NULL pour auth phone-only
    IsVerified   bool    `gorm:"default:false"           json:"is_verified"`
    DisplayName  string  `gorm:"not null"                json:"display_name"`
    AvatarURL    *string `gorm:"column:avatar_url"       json:"avatar_url,omitempty"`
    Address      *string `gorm:"column:address"          json:"address,omitempty"`
    City         *string `gorm:"column:city"             json:"city,omitempty"`
    PostalCode   *string `gorm:"column:postal_code"      json:"postal_code,omitempty"`
    Country      *string `gorm:"column:country"          json:"country,omitempty"` // ISO 3166-1 alpha-2, ex: MA
    Role         string  `gorm:"default:'user';not null" json:"role"`      // user | admin | moderator
    IsBanned     bool    `gorm:"default:false"           json:"is_banned"`
    Roles        []Role  `gorm:"many2many:user_roles;"   json:"roles,omitempty"`
}

// OTPCode — one-time code sent by SMS or WhatsApp.
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

// RefreshToken — stored hashed (SHA-256), never in cleartext.
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
