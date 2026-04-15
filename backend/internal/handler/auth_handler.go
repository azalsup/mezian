package handler

import (
    "errors"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "classifieds/internal/middleware"
    "classifieds/internal/repository"
    "classifieds/internal/service"
)

// AuthHandler handles authentication routes.
type AuthHandler struct {
    authSvc  *service.AuthService
    userRepo *repository.UserRepo
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authSvc *service.AuthService, userRepo *repository.UserRepo) *AuthHandler {
    return &AuthHandler{authSvc: authSvc, userRepo: userRepo}
}

// sendOTPRequest is the body for POST /auth/send-otp.
type sendOTPRequest struct {
    Phone   string `json:"phone"   binding:"required"`
    Channel string `json:"channel" binding:"required,oneof=sms whatsapp"`
    Purpose string `json:"purpose" binding:"required,oneof=login signup phone_change"`
}

// SendOTP godoc
// POST /auth/send-otp
// Sends an OTP code to the provided phone number.
func (h *AuthHandler) SendOTP(c *gin.Context) {
    var req sendOTPRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    if err := h.authSvc.SendOTP(req.Phone, req.Purpose, service.NotificationChannel(req.Channel)); err != nil {
        respondError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "OTP code sent"})
}

// verifyOTPRequest is the body for POST /auth/verify-otp.
type verifyOTPRequest struct {
    Phone       string `json:"phone"        binding:"required"`
    Code        string `json:"code"         binding:"required"`
    Purpose     string `json:"purpose"      binding:"required,oneof=login signup phone_change"`
    DisplayName string `json:"display_name"`
}

// VerifyOTP godoc
// POST /auth/verify-otp
// Verifies an OTP. If purpose=signup, creates the account. If purpose=login, logs in the user.
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
    var req verifyOTPRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    // Verify the OTP
    if err := h.authSvc.VerifyOTP(req.Phone, req.Code, req.Purpose); err != nil {
        respondError(c, err)
        return
    }

    switch req.Purpose {
    case "signup":
        displayName := req.DisplayName
        if displayName == "" {
            displayName = req.Phone
        }
        user, tokens, err := h.authSvc.RegisterWithOTP(req.Phone, displayName)
        if err != nil {
            respondError(c, err)
            return
        }
        c.JSON(http.StatusCreated, gin.H{
            "data": gin.H{"user": user, "tokens": tokens},
        })

    case "login":
        user, tokens, err := h.authSvc.LoginWithOTP(req.Phone)
        if err != nil {
            respondError(c, err)
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "data": gin.H{"user": user, "tokens": tokens},
        })

    default:
        // phone_change: just confirm verification
        c.JSON(http.StatusOK, gin.H{"message": "OTP verified"})
    }
}

// loginRequest is the body for POST /auth/login.
type loginRequest struct {
    Identifier string `json:"identifier" binding:"required"` // phone or email
    Password   string `json:"password"   binding:"required"`
}

// Login godoc
// POST /auth/login
// Authenticates a user with phone/email + password.
func (h *AuthHandler) Login(c *gin.Context) {
    var req loginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    user, tokens, err := h.authSvc.LoginWithPassword(req.Identifier, req.Password)
    if err != nil {
        respondError(c, err)
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "data": gin.H{"user": user, "tokens": tokens},
    })
}

// registerRequest is the body for POST /auth/register.
type registerRequest struct {
    Phone       string `json:"phone"`
    Email       string `json:"email"`
    Password    string `json:"password"     binding:"required,min=8"`
    DisplayName string `json:"display_name" binding:"required"`
    Address     string `json:"address"`
    City        string `json:"city"`
    PostalCode  string `json:"postal_code"`
    Country     string `json:"country"`
}

// Register godoc
// POST /auth/register
// Creates a user account with password.
func (h *AuthHandler) Register(c *gin.Context) {
    var req registerRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    // Enforce identifier requirements from config
    authCfg := h.authSvc.AuthConfig()
    if authCfg.PhoneRequired && req.Phone == "" {
        respondBadRequest(c, "phone number is required")
        return
    }
    if authCfg.EmailRequired && req.Email == "" {
        respondBadRequest(c, "email is required")
        return
    }
    if req.Phone == "" && req.Email == "" {
        respondBadRequest(c, "phone or email is required")
        return
    }

    addr := &service.RegisterAddressFields{
        Address:    req.Address,
        City:       req.City,
        PostalCode: req.PostalCode,
        Country:    req.Country,
    }
    if addr.Country == "" {
        addr.Country = authCfg.DefaultCountry
    }

    user, tokens, err := h.authSvc.RegisterWithPassword(req.Phone, req.Email, req.Password, req.DisplayName, addr)
    if err != nil {
        respondError(c, err)
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "data": gin.H{"user": user, "tokens": tokens},
    })
}

// refreshRequest is the body for POST /auth/refresh.
type refreshRequest struct {
    RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh godoc
// POST /auth/refresh
// Generates a new token pair from a valid refresh token.
func (h *AuthHandler) Refresh(c *gin.Context) {
    var req refreshRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    tokens, err := h.authSvc.RefreshTokens(req.RefreshToken)
    if err != nil {
        respondError(c, err)
        return
    }

    respondOK(c, tokens)
}

// Logout godoc
// POST /auth/logout
// Revokes all refresh tokens of the authenticated user.
func (h *AuthHandler) Logout(c *gin.Context) {
    userID := middleware.GetUserID(c)
    if err := h.authSvc.Logout(userID); err != nil {
        respondError(c, err)
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// GetMe godoc
// GET /auth/me
// Returns the authenticated user's profile.
func (h *AuthHandler) GetMe(c *gin.Context) {
    userID := middleware.GetUserID(c)
    user, err := h.userRepo.FindByID(userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            respondError(c, service.ErrUserNotFound)
            return
        }
        respondError(c, err)
        return
    }
    respondOK(c, user)
}

// updateMeRequest est le body de PUT /auth/me.
type updateMeRequest struct {
    DisplayName string `json:"display_name"`
    Email       string `json:"email"`
    City        string `json:"city"`
    AvatarURL   string `json:"avatar_url"`
}

// UpdateMe godoc
// PUT /auth/me
// Updates the authenticated user's profile.
func (h *AuthHandler) UpdateMe(c *gin.Context) {
    userID := middleware.GetUserID(c)
    user, err := h.userRepo.FindByID(userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            respondError(c, service.ErrUserNotFound)
            return
        }
        respondError(c, err)
        return
    }

    var req updateMeRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        respondBadRequest(c, err.Error())
        return
    }

    if req.DisplayName != "" {
        user.DisplayName = req.DisplayName
    }
    if req.Email != "" {
        email := strings.ToLower(strings.TrimSpace(req.Email))
        // Check that the email is not already used by another user
        existing, err := h.userRepo.FindByEmail(email)
        if err == nil && existing.Model.ID != user.Model.ID {
            respondError(c, service.ErrEmailAlreadyUsed)
            return
        }
        user.Email = &email
    }
    if req.City != "" {
        user.City = &req.City
    }
    if req.AvatarURL != "" {
        user.AvatarURL = &req.AvatarURL
    }

    if err := h.userRepo.Update(user); err != nil {
        respondError(c, err)
        return
    }

    respondOK(c, user)
}
