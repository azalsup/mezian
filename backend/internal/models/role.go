package models

import "time"

// Permission is an atomic capability that can be assigned to a Role.
// Permissions are seeded by the application and cannot be created via the API.
type Permission struct {
    ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    Key         string `gorm:"uniqueIndex;not null"     json:"key"`         // e.g. "categories.write"
    Group       string `gorm:"not null"                 json:"group"`       // e.g. "categories"
    LabelFR     string `gorm:"not null"                 json:"label_fr"`
    Description string `                                json:"description"`
}

// Role is a named set of permissions.
// System roles (IsSystem = true) cannot be deleted.
type Role struct {
    ID          uint         `gorm:"primaryKey;autoIncrement"  json:"id"`
    Name        string       `gorm:"uniqueIndex;not null"      json:"name"`
    Slug        string       `gorm:"uniqueIndex;not null"      json:"slug"`
    Description string       `                                 json:"description"`
    IsSystem    bool         `gorm:"default:false"             json:"is_system"`
    Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
    CreatedAt   time.Time    `                                 json:"created_at"`
    UpdatedAt   time.Time    `                                 json:"updated_at"`
}
