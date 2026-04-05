package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"mezian/internal/middleware"
	"mezian/internal/repository"
	"mezian/internal/service"
)

// AuthHandler gère les routes d'authentification.
type AuthHandler struct {
	authSvc  *service.AuthService
	userRepo *repository.UserRepo
}

// NewAuthHandler crée un nouveau AuthHandler.
func NewAuthHandler(authSvc *service.AuthService, userRepo *repository.UserRepo) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, userRepo: userRepo}
}

// sendOTPRequest est le body de POST /auth/send-otp.
type sendOTPRequest struct {
	Phone   string `json:"phone"   binding:"required"`
	Channel string `json:"channel" binding:"required,oneof=sms whatsapp"`
	Purpose string `json:"purpose" binding:"required,oneof=login signup phone_change"`
}

// SendOTP godoc
// POST /auth/send-otp
// Envoie un code OTP au numéro de téléphone indiqué.
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

	c.JSON(http.StatusOK, gin.H{"message": "code OTP envoyé"})
}

// verifyOTPRequest est le body de POST /auth/verify-otp.
type verifyOTPRequest struct {
	Phone       string `json:"phone"        binding:"required"`
	Code        string `json:"code"         binding:"required"`
	Purpose     string `json:"purpose"      binding:"required,oneof=login signup phone_change"`
	DisplayName string `json:"display_name"`
}

// VerifyOTP godoc
// POST /auth/verify-otp
// Vérifie un OTP. Si purpose=signup, crée le compte. Si purpose=login, connecte l'utilisateur.
func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req verifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	// Vérifier l'OTP
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
			"user":   user,
			"tokens": tokens,
		})

	case "login":
		user, tokens, err := h.authSvc.LoginWithOTP(req.Phone)
		if err != nil {
			respondError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"user":   user,
			"tokens": tokens,
		})

	default:
		// phone_change: juste confirmer la vérification
		c.JSON(http.StatusOK, gin.H{"message": "OTP vérifié"})
	}
}

// loginRequest est le body de POST /auth/login.
type loginRequest struct {
	Identifier string `json:"identifier" binding:"required"` // téléphone ou email
	Password   string `json:"password"   binding:"required"`
}

// Login godoc
// POST /auth/login
// Authentifie un utilisateur avec téléphone/email + mot de passe.
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
		"user":   user,
		"tokens": tokens,
	})
}

// registerRequest est le body de POST /auth/register.
type registerRequest struct {
	Phone       string `json:"phone"        binding:"required"`
	Email       string `json:"email"`
	Password    string `json:"password"     binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required"`
}

// Register godoc
// POST /auth/register
// Crée un compte utilisateur avec mot de passe.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondBadRequest(c, err.Error())
		return
	}

	user, tokens, err := h.authSvc.RegisterWithPassword(req.Phone, req.Email, req.Password, req.DisplayName)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

// refreshRequest est le body de POST /auth/refresh.
type refreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh godoc
// POST /auth/refresh
// Génère une nouvelle paire de tokens à partir d'un refresh token valide.
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
// Révoque tous les refresh tokens de l'utilisateur connecté.
func (h *AuthHandler) Logout(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := h.authSvc.Logout(userID); err != nil {
		respondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "déconnecté"})
}

// GetMe godoc
// GET /auth/me
// Retourne le profil de l'utilisateur connecté.
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
// Met à jour le profil de l'utilisateur connecté.
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
		// Vérifier que l'email n'est pas déjà pris par un autre utilisateur
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
