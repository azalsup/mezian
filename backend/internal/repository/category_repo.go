package repository

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"classifieds/internal/models"

	"gorm.io/gorm"
)

// ── JSON data structures ──────────────────────────────────────────────────────

// vehicleCarData mirrors the structure of data/vehicles_cars.json.
type vehicleCarData struct {
	Brands []string `json:"brands"`
}

// catJSON mirrors the structure of data/categories.json.
type catJSON struct {
	Categories []catEntry `json:"categories"`
}

type catEntry struct {
	Code          string     `json:"code"`
	Icon          string     `json:"icon"`
	SortOrder     int        `json:"sort_order"`
	NameFR        string     `json:"name_fr"`
	NameAR        string     `json:"name_ar"`
	NameEN        string     `json:"name_en"`
	Featured      bool       `json:"featured"`
	Subcategories []subEntry `json:"subcategories"`
}

type subEntry struct {
	Code      string `json:"code"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
	NameFR    string `json:"name_fr"`
	NameAR    string `json:"name_ar"`
	NameEN    string `json:"name_en"`
}

// loadCategoryTree reads data/categories.json.
func loadCategoryTree() (*catJSON, error) {
	paths := []string{"data/categories.json", "../data/categories.json"}
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var d catJSON
		if err := json.Unmarshal(b, &d); err != nil {
			return nil, fmt.Errorf("parsing %s: %w", p, err)
		}
		return &d, nil
	}
	return nil, fmt.Errorf("data/categories.json not found")
}

// loadVehicleCarBrands reads brand names from data/vehicles_cars.json.
// Falls back to a minimal hardcoded list if the file is unavailable.
func loadVehicleCarBrands() []string {
	paths := []string{"data/vehicles_cars.json", "../data/vehicles_cars.json"}
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		var d vehicleCarData
		if err := json.Unmarshal(b, &d); err != nil {
			log.Printf("Warning: could not parse %s: %v", p, err)
			continue
		}
		if len(d.Brands) > 0 {
			return d.Brands
		}
	}
	log.Printf("Warning: data/vehicles_cars.json not found — using fallback brand list")
	return []string{
		"Dacia", "Renault", "Peugeot", "Citroën", "Volkswagen",
		"Hyundai", "Kia", "Toyota", "Ford", "BMW",
		"Mercedes-Benz", "Audi", "Autre",
	}
}

// vehicleDataPath returns a diagnostic path string for error messages.
func vehicleDataPath() string {
	paths := []string{"data/vehicles_cars.json", "../data/vehicles_cars.json"}
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return fmt.Sprintf("data/vehicles_cars.json (not found in %s)", func() string { d, _ := os.Getwd(); return d }())
}

// CategoryRepo handles database operations for categories.
type CategoryRepo struct{ db *gorm.DB }

// NewCategoryRepo creates a new CategoryRepo.
func NewCategoryRepo(db *gorm.DB) *CategoryRepo { return &CategoryRepo{db} }

// FindAllAdmin returns ALL root categories (including inactive) with their children.
func (r *CategoryRepo) FindAllAdmin() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.
		Where("parent_id IS NULL").
		Order("sort_order ASC, name_fr ASC").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC, name_fr ASC")
		}).
		Find(&categories).Error
	return categories, err
}

// CreateCategory inserts a new category.
func (r *CategoryRepo) CreateCategory(cat *models.Category) error {
	return r.db.Create(cat).Error
}

// UpdateCategory patches editable fields on a category by ID.
func (r *CategoryRepo) UpdateCategory(id uint, fields map[string]any) error {
	return r.db.Model(&models.Category{}).Where("id = ?", id).Updates(fields).Error
}

// DeleteCategory soft-deletes a category and its children.
func (r *CategoryRepo) DeleteCategory(id uint) error {
	r.db.Where("parent_id = ?", id).Delete(&models.Category{})
	return r.db.Delete(&models.Category{}, id).Error
}

// FindAll returns all root categories with their children and attribute definitions.
func (r *CategoryRepo) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.
		Where("parent_id IS NULL AND is_active = ?", true).
		Order("sort_order ASC, name_fr ASC").
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC, name_fr ASC")
		}).
		Preload("Children.AttributeDefinitions", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("AttributeDefinitions", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Find(&categories).Error
	return categories, err
}

// FindBySlug retrieves a category (with its children and attributes) by slug.
func (r *CategoryRepo) FindBySlug(slug string) (*models.Category, error) {
	var cat models.Category
	err := r.db.
		Where("slug = ? AND is_active = ?", slug, true).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC")
		}).
		Preload("Children.AttributeDefinitions", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		Preload("AttributeDefinitions", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&cat).Error
	return &cat, err
}

// FindByID retrieves a category by its ID (with attributes).
func (r *CategoryRepo) FindByID(id uint) (*models.Category, error) {
	var cat models.Category
	err := r.db.
		Preload("AttributeDefinitions", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&cat, id).Error
	return &cat, err
}

// SeedDefaults inserts the full category tree if the table is empty.
// If categories exist, updates featured flag for known categories.
func (r *CategoryRepo) SeedDefaults() error {
	var count int64
	r.db.Model(&models.Category{}).Count(&count)
	if count == 0 {
		return r.seedCategories()
	}
	// Update featured for existing categories
	featuredSlugs := []string{"vehicles", "real_estate", "employment"}
	return r.db.Model(&models.Category{}).Where("slug IN ?", featuredSlugs).Update("featured", true).Error
}

// ForceReseed drops and recreates all categories and attribute definitions.
// Triggered by config flag seed.force = true.
func (r *CategoryRepo) ForceReseed() error {
	// Disable FK constraints for the duration of the wipe so that the
	// self-referential parent_id on categories doesn't block the DELETEs.
	r.db.Exec("PRAGMA foreign_keys = OFF")
	r.db.Exec("DELETE FROM attribute_definitions")
	r.db.Exec("DELETE FROM categories")
	r.db.Exec("PRAGMA foreign_keys = ON")
	return r.seedCategories()
}

// ── helpers ───────────────────────────────────────────────────────────────────

func str(s string) *string { return &s }

func enumStr(values []string) *string {
	b := `[`
	for i, v := range values {
		if i > 0 {
			b += ","
		}
		b += `"` + v + `"`
	}
	b += `]`
	return &b
}

// attr builds an AttributeDefinition concisely.
// labelEN is the English label. opts: *string for unit or enum values.
func attr(key, labelFR, labelAR, labelEN, dataType string, required, filterable bool, order int, opts ...*string) models.AttributeDefinition {
	a := models.AttributeDefinition{
		Key: key, LabelFR: labelFR, LabelAR: labelAR, LabelEN: labelEN,
		DataType: dataType, IsRequired: required, IsFilterable: filterable, SortOrder: order,
	}
	if len(opts) > 0 && opts[0] != nil {
		if dataType == "enum" {
			a.EnumValues = opts[0]
		} else {
			a.Unit = opts[0]
		}
	}
	return a
}

// attrsBySlug returns the attribute definitions for a given subcategory slug.
// The category hierarchy comes from categories.json; attributes stay in Go.
func attrsBySlug(slug string, carBrands []string) []models.AttributeDefinition {
	switch slug {
	case "voitures":
		return []models.AttributeDefinition{
			attr("brand", "Marque", "العلامة التجارية", "Marca", "enum", true, true, 1, enumStr(carBrands)),
			attr("model", "Modèle", "الموديل", "Model", "string", true, true, 2),
			attr("year", "Année", "السنة", "Year", "integer", true, true, 3),
			attr("mileage_km", "Kilométrage", "عدد الكيلومترات", "Mileage", "integer", true, true, 4, str("km")),
			attr("variant_name", "Désignation commerciale", "التسمية التجارية", "Commercial name", "string", false, false, 5),
			attr("trim", "Finition", "درجة التجهيز", "Trim", "string", false, false, 6),
			attr("segment", "Segment", "الفئة", "Segment", "enum", false, true, 7,
				enumStr([]string{"Citadine", "Berline", "Break", "SUV/Crossover", "Coupé", "Cabriolet", "Monospace", "Pick-up"})),
			attr("body_type", "Carrosserie", "نوع الهيكل", "Body type", "enum", false, true, 8,
				enumStr([]string{"Citadine", "Berline", "Break", "SUV", "Crossover", "Coupé", "Cabriolet", "Monospace", "Fourgonnette", "Pick-up", "Autre"})),
			attr("doors", "Portes", "عدد الأبواب", "Doors", "integer", false, true, 9),
			attr("seats", "Places", "عدد المقاعد", "Seats", "integer", false, true, 10),
			attr("fuel_type", "Carburant", "نوع الوقود", "Fuel", "enum", true, true, 11,
				enumStr([]string{"Essence (P)", "Diesel (D)", "Hybride (HEV)", "Hybride rechargeable (PHEV)", "Électrique (BEV)", "GPL (LPG)", "GNC (CNG)", "Éthanol (E85)"})),
			attr("transmission_type", "Transmission", "ناقل الحركة", "Transmission", "enum", false, true, 12,
				enumStr([]string{"Manuelle", "Automatique", "CVT", "DSG / PDK", "Semi-automatique"})),
			attr("color", "Couleur", "اللون", "Color", "string", false, true, 13),
			attr("condition", "État", "الحالة", "Condition", "enum", true, true, 14,
				enumStr([]string{"Neuf", "Très bon", "Bon", "Correct", "Pour pièces"})),
			attr("first_owner", "Première main", "يد أولى", "First owner", "boolean", false, true, 15),
		}
	case "motos":
		return []models.AttributeDefinition{
			attr("brand", "Marque", "العلامة التجارية", "Brand", "string", true, true, 1),
			attr("model", "Modèle", "الموديل", "Model", "string", true, true, 2),
			attr("year", "Année", "السنة", "Year", "integer", true, true, 3),
			attr("mileage_km", "Kilométrage", "عدد الكيلومترات", "Mileage", "integer", true, true, 4, str("km")),
			attr("cylinder_cc", "Cylindrée", "سعة المحرك", "Engine", "integer", false, true, 5, str("cc")),
			attr("fuel_type", "Carburant", "نوع الوقود", "Fuel", "enum", false, true, 6,
				enumStr([]string{"Essence (P)", "Électrique (BEV)"})),
			attr("condition", "État", "الحالة", "Condition", "enum", true, true, 7,
				enumStr([]string{"Neuf", "Très bon", "Bon", "Correct"})),
		}
	case "utilitaires":
		return []models.AttributeDefinition{
			attr("brand", "Marque", "العلامة التجارية", "Brand", "string", true, true, 1),
			attr("model", "Modèle", "الموديل", "Model", "string", true, true, 2),
			attr("year", "Année", "السنة", "Year", "integer", true, true, 3),
			attr("mileage_km", "Kilométrage", "عدد الكيلومترات", "Mileage", "integer", true, true, 4, str("km")),
			attr("payload_kg", "Charge utile", "الحمولة", "Payload", "integer", false, true, 5, str("kg")),
			attr("fuel_type", "Carburant", "نوع الوقود", "Fuel", "enum", false, true, 6,
				enumStr([]string{"Diesel (D)", "Essence (P)", "Électrique (BEV)", "GPL (LPG)"})),
		}
	case "pieces-auto":
		return []models.AttributeDefinition{
			attr("compatible_brand", "Marque compatible", "متوافق مع", "Compatible with", "string", false, true, 1),
			attr("part_type", "Type de pièce", "نوع القطعة", "Part type", "enum", false, true, 2,
				enumStr([]string{"Moteur", "Carrosserie", "Électronique", "Suspension", "Freinage", "Autre"})),
			attr("condition", "État", "الحالة", "Condition", "enum", true, true, 3,
				enumStr([]string{"Neuf", "Occasion"})),
		}
	case "appartements-vente", "appartements-location":
		a := []models.AttributeDefinition{
			attr("surface_m2", "Surface", "المساحة", "Area", "float", true, true, 1, str("m²")),
			attr("rooms", "Pièces", "عدد الغرف", "Rooms", "enum", true, true, 2,
				enumStr([]string{"Studio", "F1", "F2", "F3", "F4", "F5", "F6+"})),
			attr("floor", "Étage", "الطابق", "Floor", "integer", false, false, 3),
			attr("has_elevator", "Ascenseur", "مصعد", "Elevator", "boolean", false, true, 4),
			attr("has_parking", "Parking", "موقف سيارة", "Parking", "boolean", false, true, 5),
		}
		if slug == "appartements-vente" {
			a = append(a, attr("condition", "État", "الحالة", "Condition", "enum", false, true, 6,
				enumStr([]string{"Neuf", "Bon état", "À rénover"})))
		} else {
			a = append(a, attr("furnished", "Meublé", "مفروش", "Furnished", "boolean", false, true, 6))
		}
		return a
	case "villas-vente", "villas-location":
		a := []models.AttributeDefinition{
			attr("surface_m2", "Surface habitable", "المساحة المبنية", "Living area", "float", true, true, 1, str("m²")),
			attr("land_m2", "Terrain", "مساحة الأرض", "Land area", "float", false, true, 2, str("m²")),
			attr("rooms", "Chambres", "عدد الغرف", "Bedrooms", "integer", false, true, 3),
			attr("has_pool", "Piscine", "حمام سباحة", "Pool", "boolean", false, true, 4),
			attr("has_garage", "Garage", "مرآب", "Garage", "boolean", false, true, 5),
		}
		if slug == "villas-location" {
			a = append(a, attr("furnished", "Meublé", "مفروش", "Furnished", "boolean", false, true, 6))
		}
		return a
	case "terrains":
		return []models.AttributeDefinition{
			attr("surface_m2", "Superficie", "المساحة", "Area", "float", true, true, 1, str("m²")),
			attr("is_constructible", "Constructible", "قابل للبناء", "Buildable", "boolean", true, true, 2),
			attr("is_connected", "Viabilisé", "مجهز بالشبكات", "Serviced", "boolean", false, true, 3),
		}
	case "bureaux-commerces":
		return []models.AttributeDefinition{
			attr("surface_m2", "Surface", "المساحة", "Area", "float", true, true, 1, str("m²")),
			attr("type", "Type", "النوع", "Type", "enum", true, true, 2,
				enumStr([]string{"Bureau", "Local commercial", "Entrepôt", "Plateau"})),
			attr("has_parking", "Parking", "موقف سيارة", "Parking", "boolean", false, true, 3),
		}
	case "telephones":
		return []models.AttributeDefinition{
			attr("brand", "Marque", "العلامة التجارية", "Brand", "enum", true, true, 1,
				enumStr([]string{"Apple", "Samsung", "Huawei", "Xiaomi", "Oppo", "Realme", "Autre"})),
			attr("model", "Modèle", "الموديل", "Model", "string", true, false, 2),
			attr("storage_gb", "Stockage", "التخزين", "Storage", "enum", false, true, 3,
				enumStr([]string{"16 GB", "32 GB", "64 GB", "128 GB", "256 GB", "512 GB", "1 TB"})),
			attr("condition", "État", "الحالة", "Condition", "enum", true, true, 4,
				enumStr([]string{"Neuf (scellé)", "Comme neuf", "Très bon", "Bon", "Acceptable"})),
		}
	case "informatique":
		return []models.AttributeDefinition{
			attr("type", "Type", "النوع", "Type", "enum", true, true, 1,
				enumStr([]string{"Laptop", "PC fixe", "Écran", "Imprimante", "Composant", "Accessoire"})),
			attr("brand", "Marque", "العلامة التجارية", "Brand", "string", false, true, 2),
			attr("ram_gb", "RAM", "الذاكرة العشوائية", "RAM", "enum", false, true, 3,
				enumStr([]string{"4 GB", "8 GB", "16 GB", "32 GB", "64 GB"})),
			attr("condition", "État", "الحالة", "Condition", "enum", true, true, 4,
				enumStr([]string{"Neuf", "Très bon", "Bon", "Acceptable"})),
		}
	case "offres-emploi":
		return []models.AttributeDefinition{
			attr("contract_type", "Type de contrat", "نوع العقد", "Contract", "enum", true, true, 1,
				enumStr([]string{"CDI", "CDD", "Freelance", "Stage", "Alternance", "Temps partiel"})),
			attr("sector", "Secteur", "القطاع", "Sector", "string", false, true, 2),
			attr("experience_years", "Expérience", "الخبرة", "Experience", "enum", false, true, 3,
				enumStr([]string{"Débutant", "1-2 ans", "3-5 ans", "5-10 ans", "+10 ans"})),
			attr("remote", "Télétravail", "عمل عن بُعد", "Remote", "enum", false, true, 4,
				enumStr([]string{"Sur site", "Hybride", "Télétravail total"})),
		}
	case "demandes-emploi":
		return []models.AttributeDefinition{
			attr("contract_type", "Type souhaité", "نوع العقد المطلوب", "Contract", "enum", false, true, 1,
				enumStr([]string{"CDI", "CDD", "Freelance", "Stage", "Temps partiel"})),
			attr("sector", "Secteur", "القطاع", "Sector", "string", false, true, 2),
			attr("experience_years", "Années d'expérience", "سنوات الخبرة", "Experience", "enum", false, true, 3,
				enumStr([]string{"Débutant", "1-2 ans", "3-5 ans", "5-10 ans", "+10 ans"})),
		}
	case "vetements-femme", "vetements-homme":
		return []models.AttributeDefinition{
			attr("size", "Taille", "المقاس", "Size", "enum", false, true, 1,
				enumStr([]string{"XS", "S", "M", "L", "XL", "XXL", "XXXL"})),
			attr("condition", "État", "الحالة", "Condition", "enum", true, true, 2,
				enumStr([]string{"Neuf avec étiquette", "Neuf sans étiquette", "Très bon", "Bon", "Acceptable"})),
		}
	case "chaussures":
		return []models.AttributeDefinition{
			attr("size_eu", "Pointure (EU)", "المقاس (EU)", "EU size", "integer", false, true, 1),
			attr("gender", "Genre", "الجنس", "Gender", "enum", false, true, 2,
				enumStr([]string{"Femme", "Homme", "Enfant", "Mixte"})),
			attr("condition", "État", "الحالة", "Condition", "enum", true, true, 3,
				enumStr([]string{"Neuf", "Très bon", "Bon", "Acceptable"})),
		}
	case "animaux":
		return []models.AttributeDefinition{
			attr("species", "Espèce", "النوع", "Species", "enum", false, true, 1,
				enumStr([]string{"Chien", "Chat", "Oiseau", "Poisson", "Reptile", "Autre"})),
			attr("vaccinated", "Vacciné", "ملقح", "Vaccinated", "boolean", false, true, 2),
		}
	case "services-education":
		return []models.AttributeDefinition{
			attr("subject", "Matière", "المادة", "Subject", "string", false, true, 1),
			attr("level", "Niveau", "المستوى", "Level", "enum", false, true, 2,
				enumStr([]string{"Primaire", "Collège", "Lycée", "Université", "Professionnel", "Tous niveaux"})),
			attr("format", "Format", "الصيغة", "Format", "enum", false, true, 3,
				enumStr([]string{"Présentiel", "En ligne", "Les deux"})),
		}
	}
	return nil
}

func (r *CategoryRepo) seedCategories() error {
	carBrands := loadVehicleCarBrands()
	log.Printf("Seeding categories — %d car brands, reading hierarchy from data/categories.json", len(carBrands))

	tree, err := loadCategoryTree()
	if err != nil {
		return fmt.Errorf("seedCategories: %w", err)
	}

	var categories []models.Category
	for _, entry := range tree.Categories {
		parent := models.Category{
			Slug:      entry.Code,
			NameFR:    entry.NameFR,
			NameAR:    entry.NameAR,
			NameEN:    entry.NameEN,
			Icon:      entry.Icon,
			SortOrder: entry.SortOrder,
			IsActive:  true,
			Featured:  entry.Featured,
		}
		parent.AttributeDefinitions = attrsBySlug(entry.Code, carBrands)

		for _, sub := range entry.Subcategories {
			child := models.Category{
				Slug:      sub.Code,
				NameFR:    sub.NameFR,
				NameAR:    sub.NameAR,
				NameEN:    sub.NameEN,
				Icon:      sub.Icon,
				SortOrder: sub.SortOrder,
				IsActive:  true,
			}
			child.AttributeDefinitions = attrsBySlug(sub.Code, carBrands)
			parent.Children = append(parent.Children, child)
		}
		categories = append(categories, parent)
	}

	return r.db.Create(&categories).Error
}
