package handler

import (
    "strconv"

    "github.com/gin-gonic/gin"

    "classifieds/internal/middleware"
    "classifieds/internal/models"
    "classifieds/internal/repository"
    "classifieds/internal/service"
)

// AdHandler handles ad routes.
type AdHandler struct {
    adSvc *service.AdService
}

// NewAdHandler creates a new AdHandler.
func NewAdHandler(adSvc *service.AdService) *AdHandler {
    return &AdHandler{adSvc: adSvc}
}

// ListAds godoc
// GET /ads
// Retourne la liste paginée des ads avec filtres optionnels.
func (h *AdHandler) ListAds(c *gin.Context) {
    f := repository.AdFilters{
        City:   c.Query("city"),
        Search: c.Query("q"),
        Sort:   c.Query("sort"),
        Status: "active", // always active for public
    }

    if v := c.Query("category"); v != "" {
        id, err := strconv.ParseUint(v, 10, 64)
        if err == nil {
            uid := uint(id)
            f.CategoryID = &uid
        }
    }
    if v := c.Query("cat"); v != "" {
        f.CategorySlug = v
    }
    if v := c.Query("sub"); v != "" {
        f.SubcategorySlug = v
    }
    if v := c.Query("min_price"); v != "" {
        p, err := strconv.ParseFloat(v, 64)
        if err == nil {
            f.MinPrice = &p
        }
    }
    if v := c.Query("max_price"); v != "" {
        p, err := strconv.ParseFloat(v, 64)
        if err == nil {
            f.MaxPrice = &p
        }
    }
    if v := c.Query("page"); v != "" {
        n, err := strconv.Atoi(v)
        if err == nil {
            f.Page = n
        }
    }
    if v := c.Query("limit"); v != "" {
        n, err := strconv.Atoi(v)
        if err == nil {
            f.Limit = n
        }
    }

    result, err := h.adSvc.ListAds(f)
    if err != nil {
        respondError(c, err)
        return
    }

    c.JSON(200, gin.H{
        "data":  result.Ads,
        "total": result.Total,
        "page":  result.Page,
        "limit": result.Limit,
    })
}

// GetAd godoc
// GET /ads/:slug
// Retourne une ad complète par son slug.
func (h *AdHandler) GetAd(c *gin.Context) {
    slug := c.Param("slug")
    ad, err := h.adSvc.GetAd(slug)
    if err != nil {
        respondError(c, err)
        return
    }
    respondOK(c, ad)
}

// createAdRequest est le body de POST /ads.
type createAdRequest struct {
    CategoryID uint                   `json:"category_id" binding:"required"`
    ShopID     *uint                  `json:"shop_id"`
    Title      string                 `json:"title"       binding:"required,min=5,max=150"`
    Body       string                 `json:"body"`
    Price      *float64               `json:"price"`
    Currency   string                 `json:"currency"`
    City       string                 `json:"city"        binding:"required"`
    District   *string                `json:"district"`
    Status     string                 `json:"status"`
    Attributes []models.AdAttribute   `json:"attributes"`
}

// CreateAd godoc
// POST /ads
// Creates a new ad for the authenticated user.
func (h *AdHandler) CreateAd(c *gin.Context) {
    var req createAdRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    userID := middleware.GetUserID(c)

    input := service.CreateAdInput{
        UserID:     userID,
        CategoryID: req.CategoryID,
        ShopID:     req.ShopID,
        Title:      req.Title,
        Body:       req.Body,
        Price:      req.Price,
        Currency:   req.Currency,
        City:       req.City,
        District:   req.District,
        Status:     req.Status,
        Attributes: req.Attributes,
    }

    ad, err := h.adSvc.CreateAd(input)
    if err != nil {
        respondError(c, err)
        return
    }

    respondCreated(c, ad)
}

// updateAdRequest is the body for PUT /ads/:slug.
type updateAdRequest struct {
    Title      *string                `json:"title"`
    Body       *string                `json:"body"`
    Price      *float64               `json:"price"`
    Currency   *string                `json:"currency"`
    City       *string                `json:"city"`
    District   *string                `json:"district"`
    Status     *string                `json:"status"`
    Attributes []models.AdAttribute   `json:"attributes"`
}

// UpdateAd godoc
// PUT /ads/:slug
// Modifie une annonce existante (owner ou admin uniquement).
func (h *AdHandler) UpdateAd(c *gin.Context) {
    slug := c.Param("slug")
    userID := middleware.GetUserID(c)
    userRole := middleware.GetUserRole(c)

    var req updateAdRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    input := service.UpdateAdInput{
        Title:      req.Title,
        Body:       req.Body,
        Price:      req.Price,
        Currency:   req.Currency,
        City:       req.City,
        District:   req.District,
        Status:     req.Status,
        Attributes: req.Attributes,
    }

    ad, err := h.adSvc.UpdateAd(slug, userID, userRole, input)
    if err != nil {
        respondError(c, err)
        return
    }

    respondOK(c, ad)
}

// DeleteAd godoc
// DELETE /ads/:slug
// Supprime une annonce (owner ou admin uniquement).
func (h *AdHandler) DeleteAd(c *gin.Context) {
    slug := c.Param("slug")
    userID := middleware.GetUserID(c)
    userRole := middleware.GetUserRole(c)

    if err := h.adSvc.DeleteAd(slug, userID, userRole); err != nil {
        respondError(c, err)
        return
    }

    c.JSON(200, gin.H{"message": "ad deleted"})
}

// GetMyAds godoc
// GET /users/me/ads
// Returns the authenticated user's ads.
func (h *AdHandler) GetMyAds(c *gin.Context) {
    userID := middleware.GetUserID(c)

    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

    result, err := h.adSvc.GetUserAds(userID, page, limit)
    if err != nil {
        respondError(c, err)
        return
    }

    c.JSON(200, gin.H{
        "data":  result.Ads,
        "total": result.Total,
        "page":  result.Page,
        "limit": result.Limit,
    })
}
