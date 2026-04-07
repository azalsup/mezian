// Package main — Mezian server entry point.
package main

import (
    "fmt"
    "log"
    "os"

    "github.com/gin-gonic/gin"

    "mezian/internal/config"
    "mezian/internal/database"
    "mezian/internal/handler"
    "mezian/internal/repository"
    "mezian/internal/router"
    "mezian/internal/service"
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

    // --- Router ---
    deps := &router.Deps{
        AuthHandler:     authHandler,
        AdHandler:       adHandler,
        CategoryHandler: catHandler,
        MediaHandler:    mediaHandler,
        ShopHandler:     shopHandler,
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
