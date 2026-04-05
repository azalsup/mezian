// Package handler contient les handlers HTTP Gin.
package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"mezian/internal/service"
)

// respondOK envoie une réponse 200 JSON.
func respondOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// respondCreated envoie une réponse 201 JSON.
func respondCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

// respondError envoie une réponse d'erreur JSON avec le bon code HTTP.
func respondError(c *gin.Context, err error) {
	code := http.StatusInternalServerError
	msg := "erreur interne du serveur"

	switch {
	case errors.Is(err, service.ErrUserNotFound),
		errors.Is(err, service.ErrAdNotFound),
		errors.Is(err, service.ErrShopNotFound),
		errors.Is(err, service.ErrMediaNotFound):
		code = http.StatusNotFound
		msg = err.Error()

	case errors.Is(err, service.ErrAdForbidden),
		errors.Is(err, service.ErrShopForbidden),
		errors.Is(err, service.ErrMediaForbidden):
		code = http.StatusForbidden
		msg = err.Error()

	case errors.Is(err, service.ErrInvalidOTP),
		errors.Is(err, service.ErrOTPMaxAttempts),
		errors.Is(err, service.ErrInvalidPassword),
		errors.Is(err, service.ErrInvalidToken):
		code = http.StatusUnauthorized
		msg = err.Error()

	case errors.Is(err, service.ErrPhoneAlreadyUsed),
		errors.Is(err, service.ErrEmailAlreadyUsed),
		errors.Is(err, service.ErrShopAlreadyExists),
		errors.Is(err, service.ErrShopSlugTaken),
		errors.Is(err, service.ErrOTPRateLimit),
		errors.Is(err, service.ErrShopAdsLimit),
		errors.Is(err, service.ErrTooManyMedia),
		errors.Is(err, service.ErrInvalidMediaType),
		errors.Is(err, service.ErrFileTooLarge),
		errors.Is(err, service.ErrInvalidYouTube):
		code = http.StatusBadRequest
		msg = err.Error()

	default:
		// Ne pas exposer les détails des erreurs internes en production
	}

	c.JSON(code, gin.H{"error": msg})
}

// respondBadRequest envoie une réponse 400 avec un message personnalisé.
func respondBadRequest(c *gin.Context, msg string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}
