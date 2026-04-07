// Package handler contains Gin HTTP handlers.
package handler

import (
    "errors"
    "net/http"

    "github.com/gin-gonic/gin"

    "classifieds/internal/service"
)

// respondOK sends a 200 JSON response.
func respondOK(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, gin.H{"data": data})
}

// respondCreated sends a 201 JSON response.
func respondCreated(c *gin.Context, data interface{}) {
    c.JSON(http.StatusCreated, gin.H{"data": data})
}

// respondError sends a JSON error response with the correct HTTP status code.
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
        // Do not expose internal error details in production
    }

    c.JSON(code, gin.H{"error": msg})
}

// respondBadRequest sends a 400 response with a custom message.
func respondBadRequest(c *gin.Context, msg string) {
    c.JSON(http.StatusBadRequest, gin.H{"error": msg})
}

// respondNotFound sends a 404 response with a custom message.
func respondNotFound(c *gin.Context, msg string) {
    c.JSON(http.StatusNotFound, gin.H{"error": msg})
}
