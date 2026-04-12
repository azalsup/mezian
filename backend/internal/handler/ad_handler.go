package handler

import (
    "log"
    "strconv"

    "github.com/gin-gonic/gin"

    "classifieds/internal/middleware"
    "classifieds/internal/models"
    "classifieds/internal/repository"
    "classifieds/internal/service"
)

// AdHandler handles ad routes.
type AdHandler struct {
    adSvc    *service.AdService
    mediaSvc *service.MediaService
}

// NewAdHandler creates a new AdHandler.
func NewAdHandler(adSvc *service.AdService, mediaSvc *service.MediaService) *AdHandler {
    return &AdHandler{adSvc: adSvc, mediaSvc: mediaSvc}
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

// CreateAd godoc
// POST /ads
// Creates a new ad for the authenticated user.
// Accepts multipart/form-data; optional "images" files are uploaded immediately.
func (h *AdHandler) CreateAd(c *gin.Context) {
    log.Printf("[CreateAd] Content-Type: %s", c.GetHeader("Content-Type"))
    log.Printf("[CreateAd] category_id=%q title=%q city=%q", c.PostForm("category_id"), c.PostForm("title"), c.PostForm("city"))

    categoryIDStr := c.PostForm("category_id")
    if categoryIDStr == "" {
        respondBadRequest(c, "category_id is required")
        return
    }
    catID, err := strconv.ParseUint(categoryIDStr, 10, 64)
    if err != nil {
        respondBadRequest(c, "invalid category_id")
        return
    }

    title := c.PostForm("title")
    if title == "" {
        respondBadRequest(c, "title is required")
        return
    }
    city := c.PostForm("city")
    if city == "" {
        respondBadRequest(c, "city is required")
        return
    }

    var price *float64
    if v := c.PostForm("price"); v != "" {
        p, err := strconv.ParseFloat(v, 64)
        if err == nil {
            price = &p
        }
    }

    currency := c.PostForm("currency")
    if currency == "" {
        currency = "MAD"
    }

    userID := middleware.GetUserID(c)

    input := service.CreateAdInput{
        UserID:     userID,
        CategoryID: uint(catID),
        Title:      title,
        Body:       c.PostForm("body"),
        Price:      price,
        Currency:   currency,
        City:       city,
    }

    ad, err := h.adSvc.CreateAd(input)
    if err != nil {
        respondError(c, err)
        return
    }

    // Upload any attached images
    form, _ := c.MultipartForm()
    if form != nil {
        for _, fh := range form.File["images"] {
            file, openErr := fh.Open()
            if openErr != nil {
                continue
            }
            h.mediaSvc.UploadImage(ad.Model.ID, userID, file, fh) //nolint:errcheck
            file.Close()
        }
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
