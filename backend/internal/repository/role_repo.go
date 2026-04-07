package repository

import (
    "errors"

    "classifieds/internal/models"

    "gorm.io/gorm"
)

// RoleRepo handles CRUD for roles and permissions.
type RoleRepo struct{ db *gorm.DB }

// NewRoleRepo creates a new RoleRepo.
func NewRoleRepo(db *gorm.DB) *RoleRepo { return &RoleRepo{db} }

// DB exposes the underlying gorm.DB for ad-hoc queries in main.
func (r *RoleRepo) DB() *gorm.DB { return r.db }

// ── Permissions ───────────────────────────────────────────────────────────────

// AllPermissions returns every seeded permission.
func (r *RoleRepo) AllPermissions() ([]models.Permission, error) {
    var perms []models.Permission
    err := r.db.Order("\"group\" ASC, key ASC").Find(&perms).Error
    return perms, err
}

// SeedPermissions upserts the predefined permission set.
// Safe to call on every startup — only inserts missing rows.
func (r *RoleRepo) SeedPermissions() error {
    perms := []models.Permission{
        {Key: "categories.read",  Group: "categories", LabelFR: "Voir les catégories",         Description: "Lire l'arbre de catégories"},
        {Key: "categories.write", Group: "categories", LabelFR: "Gérer les catégories",        Description: "Créer, modifier et supprimer les catégories"},
        {Key: "users.read",       Group: "users",       LabelFR: "Voir les utilisateurs",       Description: "Lister et consulter les profils utilisateur"},
        {Key: "users.write",      Group: "users",       LabelFR: "Modifier les utilisateurs",   Description: "Éditer les informations d'un utilisateur"},
        {Key: "users.roles",      Group: "users",       LabelFR: "Gérer les rôles utilisateurs",Description: "Assigner ou retirer des rôles à un utilisateur"},
        {Key: "roles.read",       Group: "roles",       LabelFR: "Voir les rôles",              Description: "Lister les rôles et leurs permissions"},
        {Key: "roles.write",      Group: "roles",       LabelFR: "Gérer les rôles",             Description: "Créer, modifier et supprimer les rôles"},
        {Key: "ads.read",         Group: "ads",         LabelFR: "Voir toutes les annonces",    Description: "Accès aux annonces non publiées ou signalées"},
        {Key: "ads.moderate",     Group: "ads",         LabelFR: "Modérer les annonces",        Description: "Approuver, rejeter ou supprimer des annonces"},
    }
    for _, p := range perms {
        if err := r.db.Where(models.Permission{Key: p.Key}).
            FirstOrCreate(&p).Error; err != nil {
            return err
        }
    }
    return nil
}

// ── Roles ─────────────────────────────────────────────────────────────────────

// AllRoles returns all roles with their permissions.
func (r *RoleRepo) AllRoles() ([]models.Role, error) {
    var roles []models.Role
    err := r.db.Preload("Permissions").Order("id ASC").Find(&roles).Error
    return roles, err
}

// FindRoleByID returns a role with its permissions.
func (r *RoleRepo) FindRoleByID(id uint) (*models.Role, error) {
    var role models.Role
    err := r.db.Preload("Permissions").First(&role, id).Error
    return &role, err
}

// CreateRole inserts a new role and associates permissions by ID.
func (r *RoleRepo) CreateRole(role *models.Role, permIDs []uint) error {
    if err := r.db.Create(role).Error; err != nil {
        return err
    }
    return r.setPermissions(role, permIDs)
}

// UpdateRole updates a role's metadata and replaces its permissions.
func (r *RoleRepo) UpdateRole(role *models.Role, permIDs []uint) error {
    if err := r.db.Save(role).Error; err != nil {
        return err
    }
    return r.setPermissions(role, permIDs)
}

// DeleteRole removes a non-system role.
func (r *RoleRepo) DeleteRole(id uint) error {
    var role models.Role
    if err := r.db.First(&role, id).Error; err != nil {
        return err
    }
    if role.IsSystem {
        return errors.New("cannot delete a system role")
    }
    // clear join table first
    if err := r.db.Model(&role).Association("Permissions").Clear(); err != nil {
        return err
    }
    return r.db.Delete(&role).Error
}

// SeedSystemRoles upserts the built-in roles (admin, moderator).
func (r *RoleRepo) SeedSystemRoles() error {
    var allPerms []models.Permission
    if err := r.db.Find(&allPerms).Error; err != nil {
        return err
    }

    // Build a lookup by key
    permByKey := make(map[string]models.Permission, len(allPerms))
    for _, p := range allPerms {
        permByKey[p.Key] = p
    }

    type systemRole struct {
        slug  string
        name  string
        desc  string
        perms []string
    }

    definitions := []systemRole{
        {
            slug: "admin",
            name: "Administrateur",
            desc: "Accès total au panneau d'administration.",
            perms: []string{
                "categories.read", "categories.write",
                "users.read", "users.write", "users.roles",
                "roles.read", "roles.write",
                "ads.read", "ads.moderate",
            },
        },
        {
            slug: "moderator",
            name: "Modérateur",
            desc: "Modération des annonces et consultation des catégories.",
            perms: []string{"categories.read", "ads.read", "ads.moderate"},
        },
    }

    for _, def := range definitions {
        var role models.Role
        result := r.db.Where("slug = ?", def.slug).First(&role)
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            role = models.Role{
                Name: def.name, Slug: def.slug,
                Description: def.desc, IsSystem: true,
            }
            if err := r.db.Create(&role).Error; err != nil {
                return err
            }
        }
        // Always sync permissions for system roles
        perms := make([]models.Permission, 0, len(def.perms))
        for _, key := range def.perms {
            if p, ok := permByKey[key]; ok {
                perms = append(perms, p)
            }
        }
        if err := r.db.Model(&role).Association("Permissions").Replace(perms); err != nil {
            return err
        }
    }
    return nil
}

// ── Users ─────────────────────────────────────────────────────────────────────

// ListUsers returns all users with their roles (paginated).
func (r *RoleRepo) ListUsers(page, pageSize int) ([]models.User, int64, error) {
    var users []models.User
    var total int64
    r.db.Model(&models.User{}).Count(&total)
    err := r.db.Preload("Roles.Permissions").
        Order("id ASC").
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&users).Error
    return users, total, err
}

// SetUserRoles replaces the roles of a user.
func (r *RoleRepo) SetUserRoles(userID uint, roleIDs []uint) error {
    var user models.User
    if err := r.db.First(&user, userID).Error; err != nil {
        return err
    }
    var roles []models.Role
    if len(roleIDs) > 0 {
        if err := r.db.Find(&roles, roleIDs).Error; err != nil {
            return err
        }
    }
    return r.db.Model(&user).Association("Roles").Replace(roles)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func (r *RoleRepo) setPermissions(role *models.Role, permIDs []uint) error {
    var perms []models.Permission
    if len(permIDs) > 0 {
        if err := r.db.Find(&perms, permIDs).Error; err != nil {
            return err
        }
    }
    return r.db.Model(role).Association("Permissions").Replace(perms)
}
