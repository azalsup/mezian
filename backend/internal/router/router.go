// Package router assemble toutes les routes Gin et configure les middlewares globaux.
package router

import (
	"net/http"
	"time"

	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"classifieds/internal/config"
	"classifieds/internal/handler"
	"classifieds/internal/middleware"
	"classifieds/internal/service"
)

// Deps groups all dependencies injected into the router.
type Deps struct {
	AuthHandler     *handler.AuthHandler
	AdHandler       *handler.AdHandler
	CategoryHandler *handler.CategoryHandler
	MediaHandler    *handler.MediaHandler
	ShopHandler     *handler.ShopHandler
	AdminHandler    *handler.AdminHandler
	AuthService     *service.AuthService
	Config          *config.Config
}

// New creates et configure le moteur Gin avec toutes les routes.
func New(deps *Deps) *gin.Engine {
	r := gin.New()

	// Middlewares globaux
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	origins := deps.Config.Server.CORSOrigins
	if len(origins) == 0 {
		origins = []string{"*"}
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400 * time.Second,
	}))

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

	// --- Categories (public) ---
	categories := api.Group("/categories")
	{
		categories.GET("", deps.CategoryHandler.ListCategories)
		categories.GET("/:slug", deps.CategoryHandler.GetCategory)
	}

	// --- Ads ---
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

			// Médias d'une ad (route imbriquée)
			adsProtected.POST("/:id/media", deps.MediaHandler.UploadImage)
			adsProtected.POST("/:id/media/youtube", deps.MediaHandler.AddYouTube)
		}
	}

	// --- Media (standalone) ---
	media := api.Group("/media")
	media.Use(middleware.RequireAuth(deps.AuthService))
	{
		media.DELETE("/:id", deps.MediaHandler.DeleteMedia)
		media.PUT("/:id/cover", deps.MediaHandler.SetCover)
		media.PUT("/:id/order", deps.MediaHandler.UpdateOrder)
	}

	// --- Shops ---
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

	// --- Admin (role=admin required) ---
	admin := api.Group("/admin")
	admin.Use(middleware.RequireAuth(deps.AuthService))
	admin.Use(middleware.RequireRole("admin"))
	{
		admin.GET("/permissions", deps.AdminHandler.ListPermissions)
		admin.GET("/roles", deps.AdminHandler.ListRoles)
		admin.POST("/roles", deps.AdminHandler.CreateRole)
		admin.PUT("/roles/:id", deps.AdminHandler.UpdateRole)
		admin.DELETE("/roles/:id", deps.AdminHandler.DeleteRole)
		admin.GET("/users", deps.AdminHandler.ListUsers)
		admin.PUT("/users/:id/roles", deps.AdminHandler.SetUserRoles)
	}

	return r
}
