package handler

import (
	"net/http"
	"strconv"

	"classifieds/internal/models"
	"classifieds/internal/repository"

	"github.com/gin-gonic/gin"
)

// CategoryHandler handles category routes.
type CategoryHandler struct {
	catRepo *repository.CategoryRepo
}

// NewCategoryHandler creates a new CategoryHandler.
func NewCategoryHandler(catRepo *repository.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{catRepo: catRepo}
}

// ListCategories godoc
// GET /categories
// Returns the full tree of active categories (roots + children).
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	categories, err := h.catRepo.FindAll()
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, categories)
}

// GetCategory godoc
// GET /categories/:slug
// Returns a category with its subcategories and attribute definitions.
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	slug := c.Param("slug")
	cat, err := h.catRepo.FindBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	respondOK(c, cat)
}

// ── Admin ─────────────────────────────────────────────────────────────────────

// AdminListCategories godoc
// GET /admin/categories
// Returns all categories (including inactive) for the admin panel.
func (h *CategoryHandler) AdminListCategories(c *gin.Context) {
	categories, err := h.catRepo.FindAllAdmin()
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, categories)
}

type categoryCreateReq struct {
	Slug      string `json:"slug"    binding:"required"`
	NameFR    string `json:"name_fr" binding:"required"`
	NameAR    string `json:"name_ar" binding:"required"`
	NameEN    string `json:"name_en"`
	Icon      string `json:"icon"`
	ParentID  *uint  `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
	Featured  bool   `json:"featured"`
	IsActive  bool   `json:"is_active"`
}

// AdminCreateCategory godoc
// POST /admin/categories
func (h *CategoryHandler) AdminCreateCategory(c *gin.Context) {
	var req categoryCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat := &models.Category{
		Slug:      req.Slug,
		NameFR:    req.NameFR,
		NameAR:    req.NameAR,
		NameEN:    req.NameEN,
		Icon:      req.Icon,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		Featured:  req.Featured,
		IsActive:  req.IsActive,
	}
	if err := h.catRepo.CreateCategory(cat); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, cat)
}

type categoryUpdateReq struct {
	NameFR    *string `json:"name_fr"`
	NameAR    *string `json:"name_ar"`
	NameEN    *string `json:"name_en"`
	Icon      *string `json:"icon"`
	SortOrder *int    `json:"sort_order"`
	Featured  *bool   `json:"featured"`
	IsActive  *bool   `json:"is_active"`
}

// AdminUpdateCategory godoc
// PUT /admin/categories/:id
func (h *CategoryHandler) AdminUpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req categoryUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fields := map[string]any{}
	if req.NameFR != nil    { fields["name_fr"]    = *req.NameFR }
	if req.NameAR != nil    { fields["name_ar"]    = *req.NameAR }
	if req.NameEN != nil    { fields["name_en"]    = *req.NameEN }
	if req.Icon != nil      { fields["icon"]       = *req.Icon }
	if req.SortOrder != nil { fields["sort_order"] = *req.SortOrder }
	if req.Featured != nil  { fields["featured"]   = *req.Featured }
	if req.IsActive != nil  { fields["is_active"]  = *req.IsActive }
	if len(fields) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}
	if err := h.catRepo.UpdateCategory(uint(id), fields); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}

// AdminDeleteCategory godoc
// DELETE /admin/categories/:id
func (h *CategoryHandler) AdminDeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.catRepo.DeleteCategory(uint(id)); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
