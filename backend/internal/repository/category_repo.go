package repository

import (
	"mezian/internal/models"

	"gorm.io/gorm"
)

// CategoryRepo gère les opérations base de données sur les catégories.
type CategoryRepo struct{ db *gorm.DB }

// NewCategoryRepo crée un nouveau CategoryRepo.
func NewCategoryRepo(db *gorm.DB) *CategoryRepo { return &CategoryRepo{db} }

// FindAll retourne toutes les catégories racines avec leurs enfants et définitions d'attributs.
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
		Find(&categories).Error
	return categories, err
}

// FindBySlug récupère une catégorie (avec ses attributs) par slug.
func (r *CategoryRepo) FindBySlug(slug string) (*models.Category, error) {
	var cat models.Category
	err := r.db.
		Where("slug = ? AND is_active = ?", slug, true).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_active = ?", true).Order("sort_order ASC")
		}).
		Preload("AttributeDefinitions", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&cat).Error
	return &cat, err
}

// FindByID récupère une catégorie par ID.
func (r *CategoryRepo) FindByID(id uint) (*models.Category, error) {
	var cat models.Category
	err := r.db.
		Preload("AttributeDefinitions", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order ASC")
		}).
		First(&cat, id).Error
	return &cat, err
}

// SeedDefaults insère les catégories par défaut si la table est vide.
func (r *CategoryRepo) SeedDefaults() error {
	var count int64
	r.db.Model(&models.Category{}).Count(&count)
	if count > 0 {
		return nil
	}

	categories := []models.Category{
		{Slug: "immobilier", NameFR: "Immobilier", NameAR: "عقارات", Icon: "home", SortOrder: 1, IsActive: true,
			Children: []models.Category{
				{Slug: "appartements-vente", NameFR: "Appartements à vendre", NameAR: "شقق للبيع", SortOrder: 1, IsActive: true},
				{Slug: "appartements-location", NameFR: "Appartements à louer", NameAR: "شقق للإيجار", SortOrder: 2, IsActive: true},
				{Slug: "villas-vente", NameFR: "Villas à vendre", NameAR: "فيلات للبيع", SortOrder: 3, IsActive: true},
				{Slug: "villas-location", NameFR: "Villas à louer", NameAR: "فيلات للإيجار", SortOrder: 4, IsActive: true},
				{Slug: "terrains", NameFR: "Terrains", NameAR: "أراضي", SortOrder: 5, IsActive: true},
				{Slug: "bureaux-commerces", NameFR: "Bureaux & Commerces", NameAR: "مكاتب ومحلات تجارية", SortOrder: 6, IsActive: true},
			},
		},
		{Slug: "automobiles", NameFR: "Automobiles", NameAR: "سيارات", Icon: "car", SortOrder: 2, IsActive: true,
			Children: []models.Category{
				{Slug: "voitures", NameFR: "Voitures", NameAR: "سيارات", SortOrder: 1, IsActive: true},
				{Slug: "motos", NameFR: "Motos & Scooters", NameAR: "دراجات نارية", SortOrder: 2, IsActive: true},
				{Slug: "pieces-auto", NameFR: "Pièces & Accessoires", NameAR: "قطع غيار", SortOrder: 3, IsActive: true},
				{Slug: "utilitaires", NameFR: "Utilitaires & Camions", NameAR: "شاحنات", SortOrder: 4, IsActive: true},
			},
		},
		{Slug: "electronique", NameFR: "Électronique", NameAR: "إلكترونيات", Icon: "laptop", SortOrder: 3, IsActive: true,
			Children: []models.Category{
				{Slug: "telephones", NameFR: "Téléphones & Tablettes", NameAR: "هواتف وأجهزة لوحية", SortOrder: 1, IsActive: true},
				{Slug: "informatique", NameFR: "Informatique", NameAR: "حاسوب", SortOrder: 2, IsActive: true},
				{Slug: "tv-audio", NameFR: "TV, Audio & Vidéo", NameAR: "تلفاز وصوت وفيديو", SortOrder: 3, IsActive: true},
				{Slug: "appareils-photo", NameFR: "Photo & Vidéo", NameAR: "تصوير", SortOrder: 4, IsActive: true},
			},
		},
		{Slug: "maison", NameFR: "Maison & Jardin", NameAR: "منزل وحديقة", Icon: "sofa", SortOrder: 4, IsActive: true,
			Children: []models.Category{
				{Slug: "meubles", NameFR: "Meubles", NameAR: "أثاث", SortOrder: 1, IsActive: true},
				{Slug: "electromenager", NameFR: "Électroménager", NameAR: "أجهزة منزلية", SortOrder: 2, IsActive: true},
				{Slug: "decoration", NameFR: "Décoration", NameAR: "ديكور", SortOrder: 3, IsActive: true},
				{Slug: "jardin", NameFR: "Jardin & Bricolage", NameAR: "حديقة وديكور", SortOrder: 4, IsActive: true},
			},
		},
		{Slug: "emploi", NameFR: "Emploi", NameAR: "توظيف", Icon: "briefcase", SortOrder: 5, IsActive: true,
			Children: []models.Category{
				{Slug: "offres-emploi", NameFR: "Offres d'emploi", NameAR: "عروض عمل", SortOrder: 1, IsActive: true},
				{Slug: "demandes-emploi", NameFR: "Demandes d'emploi", NameAR: "طلبات عمل", SortOrder: 2, IsActive: true},
				{Slug: "stages", NameFR: "Stages & Alternance", NameAR: "تدريب", SortOrder: 3, IsActive: true},
			},
		},
		{Slug: "services", NameFR: "Services", NameAR: "خدمات", Icon: "tools", SortOrder: 6, IsActive: true,
			Children: []models.Category{
				{Slug: "services-informatiques", NameFR: "Informatique & Web", NameAR: "خدمات إعلاميات", SortOrder: 1, IsActive: true},
				{Slug: "services-artisanat", NameFR: "Artisanat & Construction", NameAR: "بناء وصناعة تقليدية", SortOrder: 2, IsActive: true},
				{Slug: "services-education", NameFR: "Cours & Formations", NameAR: "دروس وتكوين", SortOrder: 3, IsActive: true},
				{Slug: "services-sante", NameFR: "Santé & Beauté", NameAR: "صحة وجمال", SortOrder: 4, IsActive: true},
			},
		},
		{Slug: "mode", NameFR: "Mode & Beauté", NameAR: "موضة وجمال", Icon: "shirt", SortOrder: 7, IsActive: true,
			Children: []models.Category{
				{Slug: "vetements-femme", NameFR: "Vêtements Femme", NameAR: "ملابس نسائية", SortOrder: 1, IsActive: true},
				{Slug: "vetements-homme", NameFR: "Vêtements Homme", NameAR: "ملابس رجالية", SortOrder: 2, IsActive: true},
				{Slug: "chaussures", NameFR: "Chaussures", NameAR: "أحذية", SortOrder: 3, IsActive: true},
				{Slug: "accessoires-mode", NameFR: "Accessoires", NameAR: "إكسسوارات", SortOrder: 4, IsActive: true},
			},
		},
		{Slug: "loisirs", NameFR: "Loisirs & Sport", NameAR: "ترفيه ورياضة", Icon: "football", SortOrder: 8, IsActive: true,
			Children: []models.Category{
				{Slug: "sport", NameFR: "Articles de sport", NameAR: "مستلزمات رياضية", SortOrder: 1, IsActive: true},
				{Slug: "livres-musique", NameFR: "Livres, Musique & Films", NameAR: "كتب وموسيقى وأفلام", SortOrder: 2, IsActive: true},
				{Slug: "jeux-jouets", NameFR: "Jeux & Jouets", NameAR: "ألعاب", SortOrder: 3, IsActive: true},
				{Slug: "animaux", NameFR: "Animaux", NameAR: "حيوانات", SortOrder: 4, IsActive: true},
			},
		},
	}

	// Seed des attributs pour les catégories clés (seront ajoutés via update)
	if err := r.db.Create(&categories).Error; err != nil {
		return err
	}

	// Récupérer les IDs des sous-catégories pour ajouter des attributs spécifiques
	return r.seedAttributeDefinitions()
}

func (r *CategoryRepo) seedAttributeDefinitions() error {
	// Appartements à vendre
	aptVente := &models.Category{}
	if err := r.db.Where("slug = ?", "appartements-vente").First(aptVente).Error; err == nil {
		enumValuesString := `["Studio","F1","F2","F3","F4","F5","F6+"]`
		attrs := []models.AttributeDefinition{
			{CategoryID: aptVente.ID, Key: "surface_m2", LabelFR: "Surface (m²)", LabelAR: "المساحة (م²)", DataType: "float", IsRequired: true, IsFilterable: true, SortOrder: 1},
			{CategoryID: aptVente.ID, Key: "rooms", LabelFR: "Nombre de pièces", LabelAR: "عدد الغرف", DataType: "enum", EnumValues: &enumValuesString, IsRequired: true, IsFilterable: true, SortOrder: 2},
			{CategoryID: aptVente.ID, Key: "floor", LabelFR: "Étage", LabelAR: "الطابق", DataType: "integer", IsRequired: false, IsFilterable: false, SortOrder: 3},
			{CategoryID: aptVente.ID, Key: "has_elevator", LabelFR: "Ascenseur", LabelAR: "مصعد", DataType: "boolean", IsRequired: false, IsFilterable: true, SortOrder: 4},
			{CategoryID: aptVente.ID, Key: "has_parking", LabelFR: "Parking", LabelAR: "موقف سيارة", DataType: "boolean", IsRequired: false, IsFilterable: true, SortOrder: 5},
		}
		r.db.Create(&attrs)
	}

	// Voitures
	voitures := &models.Category{}
	if err := r.db.Where("slug = ?", "voitures").First(voitures).Error; err == nil {
		fuelValues := `["Essence","Diesel","Hybride","Électrique","GPL"]`
		gearValues := `["Manuelle","Automatique"]`
		condValues := `["Excellent","Très bon","Bon","Passable"]`
		attrs := []models.AttributeDefinition{
			{CategoryID: voitures.ID, Key: "brand", LabelFR: "Marque", LabelAR: "العلامة التجارية", DataType: "string", IsRequired: true, IsFilterable: true, SortOrder: 1},
			{CategoryID: voitures.ID, Key: "model", LabelFR: "Modèle", LabelAR: "الموديل", DataType: "string", IsRequired: true, IsFilterable: false, SortOrder: 2},
			{CategoryID: voitures.ID, Key: "year", LabelFR: "Année", LabelAR: "السنة", DataType: "integer", IsRequired: true, IsFilterable: true, SortOrder: 3},
			{CategoryID: voitures.ID, Key: "mileage_km", LabelFR: "Kilométrage", LabelAR: "المسافة المقطوعة", DataType: "integer", IsRequired: true, IsFilterable: true, SortOrder: 4},
			{CategoryID: voitures.ID, Key: "fuel_type", LabelFR: "Carburant", LabelAR: "نوع الوقود", DataType: "enum", EnumValues: &fuelValues, IsRequired: true, IsFilterable: true, SortOrder: 5},
			{CategoryID: voitures.ID, Key: "gearbox", LabelFR: "Boîte de vitesse", LabelAR: "ناقل الحركة", DataType: "enum", EnumValues: &gearValues, IsRequired: false, IsFilterable: true, SortOrder: 6},
			{CategoryID: voitures.ID, Key: "condition", LabelFR: "État", LabelAR: "الحالة", DataType: "enum", EnumValues: &condValues, IsRequired: false, IsFilterable: false, SortOrder: 7},
		}
		r.db.Create(&attrs)
	}

	return nil
}
