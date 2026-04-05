// Package middleware contient les middlewares Gin.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"mezian/internal/service"
)

const (
	// ContextUserID est la clé du contexte Gin pour l'ID utilisateur.
	ContextUserID = "userID"
	// ContextUserRole est la clé du contexte Gin pour le rôle utilisateur.
	ContextUserRole = "userRole"
)

// RequireAuth est un middleware Gin qui exige un Bearer JWT valide.
// Il peuple le contexte avec userID et userRole.
func RequireAuth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token d'authentification requis",
			})
			return
		}

		claims, err := authSvc.ValidateAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token invalide ou expiré",
			})
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUserRole, claims.Role)
		c.Next()
	}
}

// OptionalAuth est un middleware Gin qui tente de valider un Bearer JWT
// mais ne bloque pas les requêtes sans token. Utile pour les routes publiques
// qui ont un comportement enrichi pour les utilisateurs connectés.
func OptionalAuth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := authSvc.ValidateAccessToken(token)
		if err == nil {
			c.Set(ContextUserID, claims.UserID)
			c.Set(ContextUserRole, claims.Role)
		}

		c.Next()
	}
}

// RequireRole vérifie que l'utilisateur connecté possède l'un des rôles spécifiés.
// Doit être utilisé après RequireAuth.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextUserRole)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "rôle requis",
			})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "rôle invalide",
			})
			return
		}

		for _, r := range roles {
			if r == roleStr {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "droits insuffisants",
		})
	}
}

// GetUserID récupère l'ID de l'utilisateur depuis le contexte Gin.
// Retourne 0 si non authentifié.
func GetUserID(c *gin.Context) uint {
	v, exists := c.Get(ContextUserID)
	if !exists {
		return 0
	}
	id, _ := v.(uint)
	return id
}

// GetUserRole récupère le rôle de l'utilisateur depuis le contexte Gin.
func GetUserRole(c *gin.Context) string {
	v, exists := c.Get(ContextUserRole)
	if !exists {
		return ""
	}
	role, _ := v.(string)
	return role
}

// extractBearerToken extrait le token depuis le header Authorization: Bearer <token>.
func extractBearerToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
