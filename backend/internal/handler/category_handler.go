package handler

import (
	"mezian/internal/repository"

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
		c.JSON(404, gin.H{"error": "category not found"})
		return
	}
	respondOK(c, cat)
}
