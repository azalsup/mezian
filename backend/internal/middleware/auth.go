// Package middleware contains Gin middlewares.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"mezian/internal/service"
)

const (
	// ContextUserID is the Gin context key for the user ID.
	ContextUserID = "userID"
	// ContextUserRole is the Gin context key for the user role.
	ContextUserRole = "userRole"
)

// RequireAuth is a Gin middleware that requires a valid Bearer JWT.
// It populates the context with userID and userRole.
func RequireAuth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearerToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authentication token required",
			})
			return
		}

		claims, err := authSvc.ValidateAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextUserRole, claims.Role)
		c.Next()
	}
}

// OptionalAuth is a Gin middleware that attempts to validate a Bearer JWT
// but does not block requests without a token. Useful for public routes
// that have enriched behavior for authenticated users.
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

// RequireRole verifies that the authenticated user has one of the specified roles.
// Must be used after RequireAuth.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextUserRole)
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "role required",
			})
			return
		}

		roleStr, ok := role.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "invalid role",
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

// GetUserID retrieves the user ID from the Gin context.
// Returns 0 if not authenticated.
func GetUserID(c *gin.Context) uint {
	v, exists := c.Get(ContextUserID)
	if !exists {
		return 0
	}
	id, _ := v.(uint)
	return id
}

// GetUserRole retrieves the user role from the Gin context.
func GetUserRole(c *gin.Context) string {
	v, exists := c.Get(ContextUserRole)
	if !exists {
		return ""
	}
	role, _ := v.(string)
	return role
}

// extractBearerToken extracts the token from the Authorization: Bearer <token> header.
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
