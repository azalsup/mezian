package handler

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"

    "classifieds/internal/models"
    "classifieds/internal/repository"
)

// AdminHandler handles /api/v1/admin/** routes.
type AdminHandler struct {
    roleRepo *repository.RoleRepo
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(roleRepo *repository.RoleRepo) *AdminHandler {
    return &AdminHandler{roleRepo: roleRepo}
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

// ListUsers GET /admin/users?page=1&page_size=20
func (h *AdminHandler) ListUsers(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }
    users, total, err := h.roleRepo.ListUsers(page, pageSize)
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

// SetUserRoles PUT /admin/users/:id/roles
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

// ── helpers ───────────────────────────────────────────────────────────────────

func parseUintParam(c *gin.Context, param string) (uint, error) {
    v, err := strconv.ParseUint(c.Param(param), 10, 64)
    return uint(v), err
}
