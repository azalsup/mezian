package handler

import (
    "github.com/gin-gonic/gin"

    "classifieds/internal/middleware"
    "classifieds/internal/service"
)

// ShopHandler handles shop routes.
type ShopHandler struct {
    shopSvc *service.ShopService
}

// NewShopHandler creates a new ShopHandler.
func NewShopHandler(shopSvc *service.ShopService) *ShopHandler {
    return &ShopHandler{shopSvc: shopSvc}
}

// GetShop godoc
// GET /shops/:slug
// Returns a shop's public page.
func (h *ShopHandler) GetShop(c *gin.Context) {
    slug := c.Param("slug")
    shop, err := h.shopSvc.GetShopBySlug(slug)
    if err != nil {
        respondError(c, err)
        return
    }
    respondOK(c, shop)
}

// createShopRequest est le body de POST /shops.
type createShopRequest struct {
    Name        string  `json:"name"        binding:"required,min=3,max=80"`
    Description *string `json:"description"`
    Phone       string  `json:"phone"       binding:"required"`
    City        string  `json:"city"        binding:"required"`
    Plan        string  `json:"plan"`
}

// CreateShop godoc
// POST /shops
// Creates a shop for the authenticated user.
func (h *ShopHandler) CreateShop(c *gin.Context) {
    var req createShopRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    userID := middleware.GetUserID(c)

    input := service.CreateShopInput{
        UserID:      userID,
        Name:        req.Name,
        Description: req.Description,
        Phone:       req.Phone,
        City:        req.City,
        Plan:        req.Plan,
    }

    shop, err := h.shopSvc.CreateShop(input)
    if err != nil {
        respondError(c, err)
        return
    }

    respondCreated(c, shop)
}

// updateShopRequest est le body de PUT /shops/:slug.
type updateShopRequest struct {
    Name        *string `json:"name"`
    Description *string `json:"description"`
    Phone       *string `json:"phone"`
    City        *string `json:"city"`
    LogoURL     *string `json:"logo_url"`
    CoverURL    *string `json:"cover_url"`
}

// UpdateShop godoc
// PUT /shops/:slug
// Updates a shop (owner only).
func (h *ShopHandler) UpdateShop(c *gin.Context) {
    slug := c.Param("slug")
    userID := middleware.GetUserID(c)

    var req updateShopRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    input := service.UpdateShopInput{
        Name:        req.Name,
        Description: req.Description,
        Phone:       req.Phone,
        City:        req.City,
        LogoURL:     req.LogoURL,
        CoverURL:    req.CoverURL,
    }

    shop, err := h.shopSvc.UpdateShop(slug, userID, input)
    if err != nil {
        respondError(c, err)
        return
    }

    respondOK(c, shop)
}

// GetMyShop godoc
// GET /users/me/shop
// Returns the authenticated user's shop.
func (h *ShopHandler) GetMyShop(c *gin.Context) {
    userID := middleware.GetUserID(c)
    shop, err := h.shopSvc.GetMyShop(userID)
    if err != nil {
        respondError(c, err)
        return
    }
    respondOK(c, shop)
}
