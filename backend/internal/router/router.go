// Package router assemble toutes les routes Gin et configure les middlewares globaux.
package router

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"mezian/internal/config"
	"mezian/internal/handler"
	"mezian/internal/middleware"
	"mezian/internal/service"
)

// Deps regroupe toutes les dépendances injectées dans le router.
type Deps struct {
	AuthHandler     *handler.AuthHandler
	AdHandler       *handler.AdHandler
	CategoryHandler *handler.CategoryHandler
	MediaHandler    *handler.MediaHandler
	ShopHandler     *handler.ShopHandler
	AuthService     *service.AuthService
	Config          *config.Config
}

// New crée et configure le moteur Gin avec toutes les routes.
func New(deps *Deps) *gin.Engine {
	r := gin.New()

	// Middlewares globaux
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware(deps.Config))

	// Servir les fichiers statiques (uploads)
	r.Static("/uploads", deps.Config.Media.UploadDir)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC(),
		})
	})

	api := r.Group("/api/v1")

	// --- Auth (public) ---
	auth := api.Group("/auth")
	{
		auth.POST("/send-otp", deps.AuthHandler.SendOTP)
		auth.POST("/verify-otp", deps.AuthHandler.VerifyOTP)
		auth.POST("/login", deps.AuthHandler.Login)
		auth.POST("/register", deps.AuthHandler.Register)
		auth.POST("/refresh", deps.AuthHandler.Refresh)

		// Routes protégées
		authProtected := auth.Group("")
		authProtected.Use(middleware.RequireAuth(deps.AuthService))
		{
			authProtected.POST("/logout", deps.AuthHandler.Logout)
			authProtected.GET("/me", deps.AuthHandler.GetMe)
			authProtected.PUT("/me", deps.AuthHandler.UpdateMe)
		}
	}

	// --- Catégories (publiques) ---
	categories := api.Group("/categories")
	{
		categories.GET("", deps.CategoryHandler.ListCategories)
		categories.GET("/:slug", deps.CategoryHandler.GetCategory)
	}

	// --- Annonces ---
	ads := api.Group("/ads")
	{
		// Routes publiques (avec auth optionnelle pour enrichir la réponse)
		ads.GET("", middleware.OptionalAuth(deps.AuthService), deps.AdHandler.ListAds)
		ads.GET("/:slug", middleware.OptionalAuth(deps.AuthService), deps.AdHandler.GetAd)

		// Routes protégées
		adsProtected := ads.Group("")
		adsProtected.Use(middleware.RequireAuth(deps.AuthService))
		{
			adsProtected.POST("", deps.AdHandler.CreateAd)
			adsProtected.PUT("/:slug", deps.AdHandler.UpdateAd)
			adsProtected.DELETE("/:slug", deps.AdHandler.DeleteAd)

			// Médias d'une annonce (route imbriquée)
			adsProtected.POST("/:id/media", deps.MediaHandler.UploadImage)
			adsProtected.POST("/:id/media/youtube", deps.MediaHandler.AddYouTube)
		}
	}

	// --- Médias (standalone) ---
	media := api.Group("/media")
	media.Use(middleware.RequireAuth(deps.AuthService))
	{
		media.DELETE("/:id", deps.MediaHandler.DeleteMedia)
		media.PUT("/:id/cover", deps.MediaHandler.SetCover)
		media.PUT("/:id/order", deps.MediaHandler.UpdateOrder)
	}

	// --- Boutiques ---
	shops := api.Group("/shops")
	{
		shops.GET("/:slug", deps.ShopHandler.GetShop)

		shopsProtected := shops.Group("")
		shopsProtected.Use(middleware.RequireAuth(deps.AuthService))
		{
			shopsProtected.POST("", deps.ShopHandler.CreateShop)
			shopsProtected.PUT("/:slug", deps.ShopHandler.UpdateShop)
		}
	}

	// --- Routes utilisateur ---
	users := api.Group("/users")
	users.Use(middleware.RequireAuth(deps.AuthService))
	{
		users.GET("/me/ads", deps.AdHandler.GetMyAds)
		users.GET("/me/shop", deps.ShopHandler.GetMyShop)
	}

	return r
}

// corsMiddleware configure les en-têtes CORS selon la configuration.
func corsMiddleware(cfg *config.Config) gin.HandlerFunc {
	allowedOrigins := make(map[string]bool)
	for _, origin := range cfg.Server.CORSOrigins {
		allowedOrigins[origin] = true
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Vérifier si l'origine est autorisée
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if len(cfg.Server.CORSOrigins) == 0 {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
