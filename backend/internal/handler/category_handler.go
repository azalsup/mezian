package handler

import (
	"mezian/internal/repository"

	"github.com/gin-gonic/gin"
)

// CategoryHandler gère les routes des catégories.
type CategoryHandler struct {
	catRepo *repository.CategoryRepo
}

// NewCategoryHandler crée un nouveau CategoryHandler.
func NewCategoryHandler(catRepo *repository.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{catRepo: catRepo}
}

// ListCategories godoc
// GET /categories
// Retourne l'arbre complet des catégories actives (racines + enfants).
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
// Retourne une catégorie avec ses sous-catégories et ses définitions d'attributs.
func (h *CategoryHandler) GetCategory(c *gin.Context) {
	slug := c.Param("slug")
	cat, err := h.catRepo.FindBySlug(slug)
	if err != nil {
		c.JSON(404, gin.H{"error": "catégorie introuvable"})
		return
	}
	respondOK(c, cat)
}
