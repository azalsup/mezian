// Package main — point d'entrée du serveur Mezian.
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
	// Charger la configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erreur chargement config: %v", err)
	}

	// Configurer le mode Gin
	gin.SetMode(cfg.Server.Mode)

	// Connexion à la base de données + migrations
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Erreur connexion base de données: %v", err)
	}
	log.Printf("Base de données connectée: %s", cfg.Database.Path)

	// Créer le répertoire d'uploads si nécessaire
	if err := os.MkdirAll(cfg.Media.UploadDir, 0755); err != nil {
		log.Fatalf("Impossible de créer le répertoire uploads: %v", err)
	}

	// --- Repositories ---
	userRepo := repository.NewUserRepo(db)
	adRepo := repository.NewAdRepo(db)
	catRepo := repository.NewCategoryRepo(db)
	shopRepo := repository.NewShopRepo(db)
	mediaRepo := repository.NewMediaRepo(db)

	// Seeder les catégories par défaut si la table est vide
	if err := catRepo.SeedDefaults(); err != nil {
		log.Printf("Avertissement: seed catégories échoué: %v", err)
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
		log.Fatalf("Erreur démarrage serveur: %v", err)
	}
}
