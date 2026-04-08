package models

import "gorm.io/gorm"

// Category — hierarchical (level 1: Real estate/Automotive, level 2: subcategories).
type Category struct {
    gorm.Model
    Slug      string    `gorm:"uniqueIndex;not null"  json:"slug"`
    NameFR    string    `gorm:"not null"              json:"name_fr"`
    NameAR    string    `gorm:"not null"              json:"name_ar"`
    NameEN    string    `gorm:"not null;default:''"   json:"name_en"`
    Icon      string    `gorm:"column:icon"           json:"icon,omitempty"`
    ParentID  *uint     `gorm:"index"                 json:"parent_id,omitempty"`
    Parent    *Category `gorm:"foreignKey:ParentID"   json:"-"`
    Children  []Category `gorm:"foreignKey:ParentID"  json:"children,omitempty"`
    SortOrder int       `gorm:"default:0"             json:"sort_order"`
    IsActive  bool      `gorm:"default:true"          json:"is_active"`

    AttributeDefinitions []AttributeDefinition `gorm:"foreignKey:CategoryID" json:"attribute_definitions,omitempty"`
}

// AttributeDefinition — schema for fields specific to a category.
// Exemples: surface_m2 (immobilier), year/mileage_km (automobile).
type AttributeDefinition struct {
    gorm.Model
    CategoryID   uint    `gorm:"not null;uniqueIndex:idx_cat_key"  json:"category_id"`
    Key          string  `gorm:"not null;uniqueIndex:idx_cat_key"  json:"key"`
    LabelFR      string  `gorm:"not null"                          json:"label_fr"`
    LabelAR      string  `gorm:"not null"                          json:"label_ar"`
    LabelEN      string  `gorm:"not null;default:''"               json:"label_en"`
    DataType     string  `gorm:"not null"                          json:"data_type"` // integer|float|string|boolean|enum
    Unit         *string `gorm:"column:unit"                       json:"unit,omitempty"`
    EnumValues   *string `gorm:"column:enum_values"                json:"enum_values,omitempty"` // JSON array
    IsRequired   bool    `gorm:"default:false"                     json:"is_required"`
    IsFilterable bool    `gorm:"default:true"                      json:"is_filterable"`
    SortOrder    int     `gorm:"default:0"                         json:"sort_order"`
}
