// Package main — Mezian server entry point.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"classifieds/internal/config"
	"classifieds/internal/database"
	"classifieds/internal/handler"
	"classifieds/internal/models"
	"classifieds/internal/repository"
	"classifieds/internal/router"
	"classifieds/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Connect to database + run migrations
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Printf("Database connected: %s", cfg.Database.Path)

	// Create uploads directory if needed
	if err := os.MkdirAll(cfg.Media.UploadDir, 0755); err != nil {
		log.Fatalf("Unable to create uploads directory: %v", err)
	}

	// --- Repositories ---
	userRepo := repository.NewUserRepo(db)
	adRepo := repository.NewAdRepo(db)
	catRepo := repository.NewCategoryRepo(db)
	shopRepo := repository.NewShopRepo(db)
	mediaRepo := repository.NewMediaRepo(db)
	roleRepo := repository.NewRoleRepo(db)

	// Seed permissions + system roles (idempotent)
	if err := roleRepo.SeedPermissions(); err != nil {
		log.Printf("Warning: permission seed failed: %v", err)
	}
	if err := roleRepo.SeedSystemRoles(); err != nil {
		log.Printf("Warning: system role seed failed: %v", err)
	}

	// Seed categories (or force-reseed if config says so)
	if cfg.Seed.Force {
		log.Println("seed.force = true → wiping and recreating all categories...")
		if err := catRepo.ForceReseed(); err != nil {
			log.Printf("Warning: force reseed failed: %v", err)
		} else {
			log.Println("Categories reseeded successfully.")
		}
	} else {
		if err := catRepo.SeedDefaults(); err != nil {
			log.Printf("Warning: category seed failed: %v", err)
		}
	}

	// Seed default admin user (idempotent — skips if already exists)
	seedAdminUser(userRepo, roleRepo, cfg)

	// --- Services ---
	notifSvc := service.NewNotificationService(cfg)
	authSvc := service.NewAuthService(userRepo, notifSvc, cfg)
	adSvc := service.NewAdService(adRepo, shopRepo)
	mediaSvc := service.NewMediaService(mediaRepo, adRepo, cfg)
	shopSvc := service.NewShopService(shopRepo, adRepo, cfg)

	// --- Handlers ---
	authHandler := handler.NewAuthHandler(authSvc, userRepo)
	adHandler := handler.NewAdHandler(adSvc)
	catHandler := handler.NewCategoryHandler(catRepo)
	mediaHandler := handler.NewMediaHandler(mediaSvc)
	shopHandler := handler.NewShopHandler(shopSvc)
	adminHandler := handler.NewAdminHandler(roleRepo)

	// --- Router ---
	deps := &router.Deps{
		AuthHandler:     authHandler,
		AdHandler:       adHandler,
		CategoryHandler: catHandler,
		MediaHandler:    mediaHandler,
		ShopHandler:     shopHandler,
		AdminHandler:    adminHandler,
		AuthService:     authSvc,
		Config:          cfg,
	}

	r := router.New(deps)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Serveur Mezian démarré sur %s (mode: %s)", addr, cfg.Server.Mode)
	log.Printf("Frontend: %s", cfg.Server.FrontendURL)

	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// seedAdminUser creates the default admin account from config if it doesn't exist yet.
func seedAdminUser(userRepo *repository.UserRepo, roleRepo *repository.RoleRepo, cfg *config.Config) {
	ac := cfg.Admin
	if ac.Email == "" && ac.Phone == "" {
		return // no admin configured
	}

	identifier := ac.Email
	if identifier == "" {
		identifier = ac.Phone
	}

	existing, err := userRepo.FindByPhoneOrEmail(identifier)
	if err == nil && existing != nil {
		return // already exists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(ac.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Warning: cannot hash admin password: %v", err)
		return
	}

	phone := ac.Phone
	if phone == "" {
		phone = "+212000000000" // placeholder if only email provided
	}
	hashStr := string(hash)
	emailPtr := (*string)(nil)
	if ac.Email != "" {
		emailPtr = &ac.Email
	}
	name := ac.DisplayName
	if name == "" {
		name = "Administrateur"
	}

	user := &models.User{
		Phone:        phone,
		Email:        emailPtr,
		PasswordHash: &hashStr,
		DisplayName:  name,
		IsVerified:   true,
		Role:         "admin",
	}
	if err := userRepo.Create(user); err != nil {
		log.Printf("Warning: admin user creation failed: %v", err)
		return
	}

	// Assign the admin system role
	var adminRole models.Role
	if err := roleRepo.DB().Where("slug = ?", "admin").First(&adminRole).Error; err == nil {
		roleRepo.SetUserRoles(user.ID, []uint{adminRole.ID}) //nolint:errcheck
	}

	log.Printf("Default admin created: %s", identifier)
}
