package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"classifieds/internal/middleware"
	"classifieds/internal/models"
	"classifieds/internal/repository"
)

// AdminHandler handles /api/v1/admin/** routes.
type AdminHandler struct {
	roleRepo *repository.RoleRepo
	userRepo *repository.UserRepo
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(roleRepo *repository.RoleRepo, userRepo *repository.UserRepo) *AdminHandler {
	return &AdminHandler{roleRepo: roleRepo, userRepo: userRepo}
}

// ── Permissions ───────────────────────────────────────────────────────────────

// ListPermissions GET /admin/permissions
func (h *AdminHandler) ListPermissions(c *gin.Context) {
	perms, err := h.roleRepo.AllPermissions()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, perms)
}

// ── Roles ─────────────────────────────────────────────────────────────────────

// ListRoles GET /admin/roles
func (h *AdminHandler) ListRoles(c *gin.Context) {
	roles, err := h.roleRepo.AllRoles()
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, roles)
}

type rolePayload struct {
	Name        string `json:"name"         binding:"required"`
	Slug        string `json:"slug"         binding:"required"`
	Description string `json:"description"`
	PermIDs     []uint `json:"permission_ids"`
}

// CreateRole POST /admin/roles
func (h *AdminHandler) CreateRole(c *gin.Context) {
	var req rolePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	role := &models.Role{
		Name: req.Name, Slug: req.Slug, Description: req.Description,
	}
	if err := h.roleRepo.CreateRole(role, req.PermIDs); err != nil {
		respondError(c, err)
		return
	}
	// Reload with permissions
	full, _ := h.roleRepo.FindRoleByID(role.ID)
	c.JSON(http.StatusCreated, full)
}

// UpdateRole PUT /admin/roles/:id
func (h *AdminHandler) UpdateRole(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondBadRequest(c, "invalid id")
		return
	}
	role, err := h.roleRepo.FindRoleByID(id)
	if err != nil {
		respondNotFound(c, "role not found")
		return
	}

	var req rolePayload
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	role.Name = req.Name
	role.Description = req.Description
	// Slug is immutable for system roles
	if !role.IsSystem {
		role.Slug = req.Slug
	}
	if err := h.roleRepo.UpdateRole(role, req.PermIDs); err != nil {
		respondError(c, err)
		return
	}
	full, _ := h.roleRepo.FindRoleByID(role.ID)
	c.JSON(http.StatusOK, full)
}

// DeleteRole DELETE /admin/roles/:id
func (h *AdminHandler) DeleteRole(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondBadRequest(c, "invalid id")
		return
	}
	if err := h.roleRepo.DeleteRole(id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "role deleted"})
}

// ── Users ─────────────────────────────────────────────────────────────────────

// ListUsers GET /admin/users?page=1&page_size=20&user_type=external|internal
func (h *AdminHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	userType := c.DefaultQuery("user_type", "")
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	users, total, err := h.roleRepo.ListUsers(page, pageSize, userType)
	if err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// UpdateUser PUT /admin/users/:id  (admin only)
func (h *AdminHandler) UpdateUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondBadRequest(c, "invalid id")
		return
	}

	var req struct {
		DisplayName string  `json:"display_name"`
		Phone       string  `json:"phone"`
		Email       *string `json:"email"`
		Address     *string `json:"address"`
		City        *string `json:"city"`
		PostalCode  *string `json:"postal_code"`
		Country     *string `json:"country"`
		IsVerified  bool    `json:"is_verified"`
		Role        string  `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	fields := map[string]any{
		"display_name": req.DisplayName,
		"phone":        req.Phone,
		"email":        req.Email,
		"address":      req.Address,
		"city":         req.City,
		"postal_code":  req.PostalCode,
		"country":      req.Country,
		"is_verified":  req.IsVerified,
		"role":         req.Role,
	}

	if err := h.userRepo.UpdateProfile(id, fields); err != nil {
		respondError(c, err)
		return
	}

	user, err := h.userRepo.FindByID(id)
	if err != nil {
		respondNotFound(c, "user not found")
		return
	}
	c.JSON(http.StatusOK, user)
}

// SetUserRoles PUT /admin/users/:id/roles  (admin only)
func (h *AdminHandler) SetUserRoles(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondBadRequest(c, "invalid id")
		return
	}
	var req struct {
		RoleIDs []uint `json:"role_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	if err := h.roleRepo.SetUserRoles(id, req.RoleIDs); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "roles updated"})
}

// BanUser PUT /admin/users/:id/ban  (admin or moderator)
func (h *AdminHandler) BanUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondBadRequest(c, "invalid id")
		return
	}
	// Prevent banning another admin/moderator unless caller is admin
	target, err := h.userRepo.FindByID(id)
	if err != nil {
		respondNotFound(c, "user not found")
		return
	}
	callerRole := middleware.GetUserRole(c)
	if target.Role != "user" && callerRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "seul un administrateur peut bannir du personnel"})
		return
	}
	if err := h.userRepo.BanUser(id, true); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user banned"})
}

// UnbanUser PUT /admin/users/:id/unban  (admin or moderator)
func (h *AdminHandler) UnbanUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondBadRequest(c, "invalid id")
		return
	}
	if err := h.userRepo.BanUser(id, false); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user unbanned"})
}

// DeleteUser DELETE /admin/users/:id  (admin only)
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondBadRequest(c, "invalid id")
		return
	}
	// Prevent self-deletion
	callerID := middleware.GetUserID(c)
	if uint(id) == callerID {
		respondBadRequest(c, "cannot delete your own account")
		return
	}
	if err := h.userRepo.DeleteUser(id); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// ResetPassword PUT /admin/users/:id/reset-password  (admin only)
func (h *AdminHandler) ResetPassword(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		respondBadRequest(c, "invalid id")
		return
	}
	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		respondError(c, err)
		return
	}
	if err := h.userRepo.ResetPassword(id, string(hash)); err != nil {
		respondError(c, err)
		return
	}
	// Revoke all tokens so user must log in again
	h.userRepo.RevokeAllUserTokens(id) //nolint:errcheck
	c.JSON(http.StatusOK, gin.H{"message": "password reset"})
}

// ── helpers ───────────────────────────────────────────────────────────────────

func parseUintParam(c *gin.Context, param string) (uint, error) {
	v, err := strconv.ParseUint(c.Param(param), 10, 64)
	return uint(v), err
}
