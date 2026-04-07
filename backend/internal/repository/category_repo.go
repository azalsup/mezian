package repository

import (
    "encoding/json"
    "fmt"
    "log"
    "os"

    "mezian/internal/models"

    "gorm.io/gorm"
)

// vehicleCarData mirrors the structure of data/vehicles_cars.json.
type vehicleCarData struct {
    Brands []string `json:"brands"`
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
func (r *CategoryRepo) SeedDefaults() error {
    var count int64
    r.db.Model(&models.Category{}).Count(&count)
    if count > 0 {
        return nil
    }
    return r.seedCategories()
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
// opts: pass a *string for unit (non-enum) or enum values (enum type).
func attr(key, labelFR, labelAR, dataType string, required, filterable bool, order int, opts ...*string) models.AttributeDefinition {
    a := models.AttributeDefinition{
        Key: key, LabelFR: labelFR, LabelAR: labelAR,
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

func (r *CategoryRepo) seedCategories() error {
    // Load brand list from external JSON — logged inside loadVehicleCarBrands if missing.
    carBrands := loadVehicleCarBrands()
    log.Printf("Seeding categories — loaded %d car brands from %s", len(carBrands), vehicleDataPath())

    categories := []models.Category{

        // ── Automobiles ───────────────────────────────────────────────────────
        {
            Slug: "automobiles", NameFR: "Automobiles", NameAR: "سيارات",
            Icon: "fa-car", SortOrder: 1, IsActive: true,
            Children: []models.Category{
                {
                    Slug: "voitures", NameFR: "Voitures", NameAR: "سيارات",
                    Icon: "fa-car-side", SortOrder: 1, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        // ── Identification ──────────────────────────────────────
                        attr("brand", "Marque", "العلامة التجارية", "enum", true, true, 1,
                            enumStr(carBrands)),
                        attr("model", "Modèle", "الموديل", "string", true, true, 2),
                        attr("year", "Année", "السنة", "integer", true, true, 3),
                        attr("mileage_km", "Kilométrage", "عدد الكيلومترات", "integer", true, true, 4, str("km")),
                        // ── Version / finition ──────────────────────────────────
                        attr("variant_name", "Désignation commerciale", "التسمية التجارية", "string", false, false, 5),
                        attr("trim", "Finition / Version", "درجة التجهيز", "string", false, false, 6),
                        // ── Carrosserie ─────────────────────────────────────────
                        attr("segment", "Segment", "الفئة", "enum", false, true, 7,
                            enumStr([]string{
                                "Citadine",
                                "Berline",
                                "Break",
                                "SUV/Crossover",
                                "Coupé",
                                "Cabriolet",
                                "Monospace",
                                "Pick-up",
                            })),
                        attr("body_type", "Carrosserie", "نوع الهيكل", "enum", false, true, 8,
                            enumStr([]string{
                                "Citadine",
                                "Berline",
                                "Break",
                                "SUV",
                                "Crossover",
                                "Coupé",
                                "Cabriolet",
                                "Monospace",
                                "Fourgonnette",
                                "Pick-up",
                                "Autre",
                            })),
                        attr("doors", "Nombre de portes", "عدد الأبواب", "integer", false, true, 9),
                        attr("seats", "Nombre de places", "عدد المقاعد", "integer", false, true, 10),
                        // ── Motorisation ────────────────────────────────────────
                        attr("fuel_type", "Carburant", "نوع الوقود", "enum", true, true, 11,
                            enumStr([]string{
                                "Essence (P)",
                                "Diesel (D)",
                                "Hybride (HEV)",
                                "Hybride rechargeable (PHEV)",
                                "Électrique (BEV)",
                                "GPL (LPG)",
                                "GNC (CNG)",
                                "Éthanol (E85)",
                            })),
                        attr("transmission_type", "Transmission", "ناقل الحركة", "enum", false, true, 12,
                            enumStr([]string{
                                "Manuelle",
                                "Automatique",
                                "CVT",
                                "DSG / PDK",
                                "Semi-automatique",
                            })),
                        // ── État ────────────────────────────────────────────────
                        attr("color", "Couleur", "اللون", "string", false, true, 13),
                        attr("condition", "État", "الحالة", "enum", true, true, 14,
                            enumStr([]string{"Neuf", "Très bon", "Bon", "Correct", "Pour pièces"})),
                        attr("first_owner", "Première main", "يد أولى", "boolean", false, true, 15),
                    },
                },
                {
                    Slug: "motos", NameFR: "Motos & Scooters", NameAR: "دراجات نارية",
                    Icon: "fa-motorcycle", SortOrder: 2, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("brand", "Marque", "العلامة التجارية", "string", true, true, 1),
                        attr("model", "Modèle", "الموديل", "string", true, true, 2),
                        attr("year", "Année", "السنة", "integer", true, true, 3),
                        attr("mileage_km", "Kilométrage", "عدد الكيلومترات", "integer", true, true, 4, str("km")),
                        attr("cylinder_cc", "Cylindrée", "سعة المحرك", "integer", false, true, 5, str("cc")),
                        attr("fuel_type", "Carburant", "نوع الوقود", "enum", false, true, 6,
                            enumStr([]string{"Essence (P)", "Électrique (BEV)"})),
                        attr("transmission_type", "Transmission", "ناقل الحركة", "enum", false, false, 7,
                            enumStr([]string{"Manuelle", "Automatique", "Semi-automatique"})),
                        attr("condition", "État", "الحالة", "enum", true, true, 8,
                            enumStr([]string{"Neuf", "Très bon", "Bon", "Correct"})),
                    },
                },
                {
                    Slug: "utilitaires", NameFR: "Utilitaires & Camions", NameAR: "شاحنات",
                    Icon: "fa-truck", SortOrder: 3, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("brand", "Marque", "العلامة التجارية", "string", true, true, 1),
                        attr("model", "Modèle", "الموديل", "string", true, true, 2),
                        attr("year", "Année", "السنة", "integer", true, true, 3),
                        attr("mileage_km", "Kilométrage", "عدد الكيلومترات", "integer", true, true, 4, str("km")),
                        attr("payload_kg", "Charge utile", "الحمولة", "integer", false, true, 5, str("kg")),
                        attr("fuel_type", "Carburant", "نوع الوقود", "enum", false, true, 6,
                            enumStr([]string{"Diesel (D)", "Essence (P)", "Électrique (BEV)", "GPL (LPG)"})),
                        attr("transmission_type", "Transmission", "ناقل الحركة", "enum", false, false, 7,
                            enumStr([]string{"Manuelle", "Automatique"})),
                    },
                },
                {
                    Slug: "pieces-auto", NameFR: "Pièces & Accessoires auto", NameAR: "قطع غيار",
                    Icon: "fa-gears", SortOrder: 4, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("compatible_brand", "Marque compatible", "متوافق مع", "string", false, true, 1),
                        attr("part_type", "Type de pièce", "نوع القطعة", "enum", false, true, 2,
                            enumStr([]string{"Moteur", "Carrosserie", "Électronique", "Suspension", "Freinage", "Autre"})),
                        attr("condition", "État", "الحالة", "enum", true, true, 3,
                            enumStr([]string{"Neuf", "Occasion"})),
                    },
                },
            },
        },

        // ── Immobilier ────────────────────────────────────────────────────────
        {
            Slug: "immobilier", NameFR: "Immobilier", NameAR: "عقارات",
            Icon: "fa-building", SortOrder: 2, IsActive: true,
            Children: []models.Category{
                {
                    Slug: "appartements-vente", NameFR: "Appartements à vendre", NameAR: "شقق للبيع",
                    Icon: "fa-building", SortOrder: 1, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("surface_m2", "Surface", "المساحة", "float", true, true, 1, str("m²")),
                        attr("rooms", "Pièces", "عدد الغرف", "enum", true, true, 2,
                            enumStr([]string{"Studio", "F1", "F2", "F3", "F4", "F5", "F6+"})),
                        attr("floor", "Étage", "الطابق", "integer", false, false, 3),
                        attr("total_floors", "Étages total", "عدد الطوابق", "integer", false, false, 4),
                        attr("has_elevator", "Ascenseur", "مصعد", "boolean", false, true, 5),
                        attr("has_parking", "Parking", "موقف سيارة", "boolean", false, true, 6),
                        attr("has_balcony", "Balcon / Terrasse", "شرفة", "boolean", false, true, 7),
                        attr("condition", "État", "الحالة", "enum", false, true, 8,
                            enumStr([]string{"Neuf", "Bon état", "À rénover"})),
                    },
                },
                {
                    Slug: "appartements-location", NameFR: "Appartements à louer", NameAR: "شقق للإيجار",
                    Icon: "fa-key", SortOrder: 2, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("surface_m2", "Surface", "المساحة", "float", true, true, 1, str("m²")),
                        attr("rooms", "Pièces", "عدد الغرف", "enum", true, true, 2,
                            enumStr([]string{"Studio", "F1", "F2", "F3", "F4", "F5", "F6+"})),
                        attr("floor", "Étage", "الطابق", "integer", false, false, 3),
                        attr("has_elevator", "Ascenseur", "مصعد", "boolean", false, true, 4),
                        attr("has_parking", "Parking", "موقف سيارة", "boolean", false, true, 5),
                        attr("furnished", "Meublé", "مفروش", "boolean", false, true, 6),
                    },
                },
                {
                    Slug: "villas-vente", NameFR: "Villas à vendre", NameAR: "فيلات للبيع",
                    Icon: "fa-house-chimney", SortOrder: 3, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("surface_m2", "Surface habitable", "المساحة المبنية", "float", true, true, 1, str("m²")),
                        attr("land_m2", "Superficie terrain", "مساحة الأرض", "float", false, true, 2, str("m²")),
                        attr("rooms", "Chambres", "عدد الغرف", "integer", false, true, 3),
                        attr("has_pool", "Piscine", "حمام سباحة", "boolean", false, true, 4),
                        attr("has_garage", "Garage", "مرآب", "boolean", false, true, 5),
                        attr("has_garden", "Jardin", "حديقة", "boolean", false, true, 6),
                        attr("condition", "État", "الحالة", "enum", false, true, 7,
                            enumStr([]string{"Neuf", "Bon état", "À rénover"})),
                    },
                },
                {
                    Slug: "villas-location", NameFR: "Villas à louer", NameAR: "فيلات للإيجار",
                    Icon: "fa-house-chimney-user", SortOrder: 4, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("surface_m2", "Surface habitable", "المساحة المبنية", "float", true, true, 1, str("m²")),
                        attr("rooms", "Chambres", "عدد الغرف", "integer", false, true, 2),
                        attr("has_pool", "Piscine", "حمام سباحة", "boolean", false, true, 3),
                        attr("has_garage", "Garage", "مرآب", "boolean", false, true, 4),
                        attr("furnished", "Meublé", "مفروش", "boolean", false, true, 5),
                    },
                },
                {
                    Slug: "terrains", NameFR: "Terrains", NameAR: "أراضي",
                    Icon: "fa-map", SortOrder: 5, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("surface_m2", "Superficie", "المساحة", "float", true, true, 1, str("m²")),
                        attr("is_constructible", "Constructible", "قابل للبناء", "boolean", true, true, 2),
                        attr("is_connected", "Viabilisé", "مجهز بالشبكات", "boolean", false, true, 3),
                    },
                },
                {
                    Slug: "bureaux-commerces", NameFR: "Bureaux & Commerces", NameAR: "مكاتب ومحلات",
                    Icon: "fa-store", SortOrder: 6, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("surface_m2", "Surface", "المساحة", "float", true, true, 1, str("m²")),
                        attr("type", "Type", "النوع", "enum", true, true, 2,
                            enumStr([]string{"Bureau", "Local commercial", "Entrepôt", "Plateau"})),
                        attr("has_parking", "Parking", "موقف سيارة", "boolean", false, true, 3),
                    },
                },
            },
        },

        // ── Électronique ──────────────────────────────────────────────────────
        {
            Slug: "electronique", NameFR: "Électronique", NameAR: "إلكترونيات",
            Icon: "fa-laptop", SortOrder: 3, IsActive: true,
            Children: []models.Category{
                {
                    Slug: "telephones", NameFR: "Téléphones & Tablettes", NameAR: "هواتف وأجهزة لوحية",
                    Icon: "fa-mobile-screen", SortOrder: 1, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("brand", "Marque", "العلامة التجارية", "enum", true, true, 1,
                            enumStr([]string{"Apple", "Samsung", "Huawei", "Xiaomi", "Oppo", "Realme", "Autre"})),
                        attr("model", "Modèle", "الموديل", "string", true, false, 2),
                        attr("storage_gb", "Stockage", "التخزين", "enum", false, true, 3,
                            enumStr([]string{"16 GB", "32 GB", "64 GB", "128 GB", "256 GB", "512 GB", "1 TB"})),
                        attr("color", "Couleur", "اللون", "string", false, false, 4),
                        attr("condition", "État", "الحالة", "enum", true, true, 5,
                            enumStr([]string{"Neuf (scellé)", "Comme neuf", "Très bon", "Bon", "Acceptable"})),
                    },
                },
                {
                    Slug: "informatique", NameFR: "Informatique", NameAR: "حاسوب",
                    Icon: "fa-desktop", SortOrder: 2, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("type", "Type", "النوع", "enum", true, true, 1,
                            enumStr([]string{"Laptop", "PC fixe", "Écran", "Imprimante", "Composant", "Accessoire"})),
                        attr("brand", "Marque", "العلامة التجارية", "string", false, true, 2),
                        attr("cpu", "Processeur", "المعالج", "string", false, false, 3),
                        attr("ram_gb", "RAM", "الذاكرة العشوائية", "enum", false, true, 4,
                            enumStr([]string{"4 GB", "8 GB", "16 GB", "32 GB", "64 GB"})),
                        attr("storage_gb", "Stockage", "التخزين", "string", false, false, 5),
                        attr("condition", "État", "الحالة", "enum", true, true, 6,
                            enumStr([]string{"Neuf", "Très bon", "Bon", "Acceptable"})),
                    },
                },
                {
                    Slug: "tv-audio", NameFR: "TV, Audio & Vidéo", NameAR: "تلفاز وصوت",
                    Icon: "fa-tv", SortOrder: 3, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("type", "Type", "النوع", "enum", true, true, 1,
                            enumStr([]string{"Télévision", "Home cinéma", "Enceinte", "Casque", "Autre"})),
                        attr("brand", "Marque", "العلامة التجارية", "string", false, true, 2),
                        attr("screen_inch", "Taille écran", "حجم الشاشة", "integer", false, true, 3, str("pouces")),
                        attr("condition", "État", "الحالة", "enum", true, true, 4,
                            enumStr([]string{"Neuf", "Très bon", "Bon", "Acceptable"})),
                    },
                },
                {
                    Slug: "appareils-photo", NameFR: "Photo & Vidéo", NameAR: "تصوير",
                    Icon: "fa-camera", SortOrder: 4, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("type", "Type", "النوع", "enum", true, true, 1,
                            enumStr([]string{"Appareil photo", "Caméra", "Objectif", "Accessoire"})),
                        attr("brand", "Marque", "العلامة التجارية", "string", false, true, 2),
                        attr("condition", "État", "الحالة", "enum", true, true, 3,
                            enumStr([]string{"Neuf", "Très bon", "Bon", "Acceptable"})),
                    },
                },
            },
        },

        // ── Maison & Jardin ───────────────────────────────────────────────────
        {
            Slug: "maison", NameFR: "Maison & Jardin", NameAR: "منزل وحديقة",
            Icon: "fa-couch", SortOrder: 4, IsActive: true,
            Children: []models.Category{
                {Slug: "meubles", NameFR: "Meubles", NameAR: "أثاث", Icon: "fa-couch", SortOrder: 1, IsActive: true},
                {Slug: "electromenager", NameFR: "Électroménager", NameAR: "أجهزة منزلية", Icon: "fa-blender", SortOrder: 2, IsActive: true},
                {Slug: "decoration", NameFR: "Décoration", NameAR: "ديكور", Icon: "fa-paint-roller", SortOrder: 3, IsActive: true},
                {Slug: "jardin", NameFR: "Jardin & Bricolage", NameAR: "حديقة وديكور", Icon: "fa-seedling", SortOrder: 4, IsActive: true},
            },
        },

        // ── Emploi ────────────────────────────────────────────────────────────
        {
            Slug: "emploi", NameFR: "Emploi", NameAR: "توظيف",
            Icon: "fa-briefcase", SortOrder: 5, IsActive: true,
            Children: []models.Category{
                {
                    Slug: "offres-emploi", NameFR: "Offres d'emploi", NameAR: "عروض عمل",
                    Icon: "fa-file-contract", SortOrder: 1, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("contract_type", "Type de contrat", "نوع العقد", "enum", true, true, 1,
                            enumStr([]string{"CDI", "CDD", "Freelance", "Stage", "Alternance", "Temps partiel"})),
                        attr("sector", "Secteur", "القطاع", "string", false, true, 2),
                        attr("experience_years", "Expérience requise", "الخبرة المطلوبة", "enum", false, true, 3,
                            enumStr([]string{"Débutant", "1-2 ans", "3-5 ans", "5-10 ans", "+10 ans"})),
                        attr("remote", "Télétravail", "عمل عن بُعد", "enum", false, true, 4,
                            enumStr([]string{"Sur site", "Hybride", "Télétravail total"})),
                    },
                },
                {
                    Slug: "demandes-emploi", NameFR: "Demandes d'emploi", NameAR: "طلبات عمل",
                    Icon: "fa-user-tie", SortOrder: 2, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("contract_type", "Type souhaité", "نوع العقد المطلوب", "enum", false, true, 1,
                            enumStr([]string{"CDI", "CDD", "Freelance", "Stage", "Temps partiel"})),
                        attr("sector", "Secteur", "القطاع", "string", false, true, 2),
                        attr("experience_years", "Années d'expérience", "سنوات الخبرة", "enum", false, true, 3,
                            enumStr([]string{"Débutant", "1-2 ans", "3-5 ans", "5-10 ans", "+10 ans"})),
                    },
                },
                {
                    Slug: "stages", NameFR: "Stages & Alternance", NameAR: "تدريب",
                    Icon: "fa-graduation-cap", SortOrder: 3, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("duration_months", "Durée", "المدة", "enum", false, false, 1,
                            enumStr([]string{"1 mois", "2 mois", "3 mois", "6 mois", "1 an"})),
                        attr("sector", "Secteur", "القطاع", "string", false, true, 2),
                    },
                },
            },
        },

        // ── Services ──────────────────────────────────────────────────────────
        {
            Slug: "services", NameFR: "Services", NameAR: "خدمات",
            Icon: "fa-screwdriver-wrench", SortOrder: 6, IsActive: true,
            Children: []models.Category{
                {Slug: "services-informatiques", NameFR: "Informatique & Web", NameAR: "خدمات إعلاميات", Icon: "fa-code", SortOrder: 1, IsActive: true},
                {Slug: "services-artisanat", NameFR: "Artisanat & Construction", NameAR: "بناء وصناعة تقليدية", Icon: "fa-helmet-safety", SortOrder: 2, IsActive: true},
                {
                    Slug: "services-education", NameFR: "Cours & Formations", NameAR: "دروس وتكوين",
                    Icon: "fa-chalkboard-user", SortOrder: 3, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("subject", "Matière", "المادة", "string", false, true, 1),
                        attr("level", "Niveau", "المستوى", "enum", false, true, 2,
                            enumStr([]string{"Primaire", "Collège", "Lycée", "Université", "Professionnel", "Tous niveaux"})),
                        attr("format", "Format", "الصيغة", "enum", false, true, 3,
                            enumStr([]string{"Présentiel", "En ligne", "Les deux"})),
                    },
                },
                {Slug: "services-sante", NameFR: "Santé & Beauté", NameAR: "صحة وجمال", Icon: "fa-kit-medical", SortOrder: 4, IsActive: true},
            },
        },

        // ── Mode & Beauté ─────────────────────────────────────────────────────
        {
            Slug: "mode", NameFR: "Mode & Beauté", NameAR: "موضة وجمال",
            Icon: "fa-shirt", SortOrder: 7, IsActive: true,
            Children: []models.Category{
                {
                    Slug: "vetements-femme", NameFR: "Vêtements Femme", NameAR: "ملابس نسائية",
                    Icon: "fa-person-dress", SortOrder: 1, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("size", "Taille", "المقاس", "enum", false, true, 1,
                            enumStr([]string{"XS", "S", "M", "L", "XL", "XXL", "XXXL"})),
                        attr("condition", "État", "الحالة", "enum", true, true, 2,
                            enumStr([]string{"Neuf avec étiquette", "Neuf sans étiquette", "Très bon", "Bon", "Acceptable"})),
                    },
                },
                {
                    Slug: "vetements-homme", NameFR: "Vêtements Homme", NameAR: "ملابس رجالية",
                    Icon: "fa-person", SortOrder: 2, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("size", "Taille", "المقاس", "enum", false, true, 1,
                            enumStr([]string{"XS", "S", "M", "L", "XL", "XXL", "XXXL"})),
                        attr("condition", "État", "الحالة", "enum", true, true, 2,
                            enumStr([]string{"Neuf avec étiquette", "Neuf sans étiquette", "Très bon", "Bon", "Acceptable"})),
                    },
                },
                {
                    Slug: "chaussures", NameFR: "Chaussures", NameAR: "أحذية",
                    Icon: "fa-shoe-prints", SortOrder: 3, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("size_eu", "Pointure (EU)", "المقاس (EU)", "integer", false, true, 1),
                        attr("gender", "Genre", "الجنس", "enum", false, true, 2,
                            enumStr([]string{"Femme", "Homme", "Enfant", "Mixte"})),
                        attr("condition", "État", "الحالة", "enum", true, true, 3,
                            enumStr([]string{"Neuf", "Très bon", "Bon", "Acceptable"})),
                    },
                },
                {Slug: "accessoires-mode", NameFR: "Accessoires", NameAR: "إكسسوارات", Icon: "fa-gem", SortOrder: 4, IsActive: true},
            },
        },

        // ── Loisirs & Sport ───────────────────────────────────────────────────
        {
            Slug: "loisirs", NameFR: "Loisirs & Sport", NameAR: "ترفيه ورياضة",
            Icon: "fa-futbol", SortOrder: 8, IsActive: true,
            Children: []models.Category{
                {Slug: "sport", NameFR: "Articles de sport", NameAR: "مستلزمات رياضية", Icon: "fa-dumbbell", SortOrder: 1, IsActive: true},
                {Slug: "livres-musique", NameFR: "Livres, Musique & Films", NameAR: "كتب وموسيقى وأفلام", Icon: "fa-book", SortOrder: 2, IsActive: true},
                {Slug: "jeux-jouets", NameFR: "Jeux & Jouets", NameAR: "ألعاب", Icon: "fa-gamepad", SortOrder: 3, IsActive: true},
                {
                    Slug: "animaux", NameFR: "Animaux", NameAR: "حيوانات",
                    Icon: "fa-paw", SortOrder: 4, IsActive: true,
                    AttributeDefinitions: []models.AttributeDefinition{
                        attr("species", "Espèce", "النوع", "enum", false, true, 1,
                            enumStr([]string{"Chien", "Chat", "Oiseau", "Poisson", "Reptile", "Autre"})),
                        attr("age_months", "Âge", "العمر", "integer", false, false, 2),
                        attr("vaccinated", "Vacciné", "ملقح", "boolean", false, true, 3),
                    },
                },
            },
        },

        // ── Divers ────────────────────────────────────────────────────────────
        {
            Slug: "divers", NameFR: "Divers", NameAR: "متفرقات",
            Icon: "fa-box-open", SortOrder: 9, IsActive: true,
            Children: []models.Category{
                {Slug: "divers-autres", NameFR: "Autres", NameAR: "أخرى", Icon: "fa-ellipsis", SortOrder: 1, IsActive: true},
            },
        },
    }

    return r.db.Create(&categories).Error
}
