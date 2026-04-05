package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"mezian/internal/config"
	"mezian/internal/models"
	"mezian/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Erreurs métier de l'authentification.
var (
	ErrUserNotFound      = errors.New("utilisateur introuvable")
	ErrInvalidOTP        = errors.New("code OTP invalide ou expiré")
	ErrOTPMaxAttempts    = errors.New("trop de tentatives OTP")
	ErrOTPRateLimit      = errors.New("trop d'OTP demandés, réessayez dans 1 heure")
	ErrInvalidPassword   = errors.New("mot de passe incorrect")
	ErrInvalidToken      = errors.New("token invalide ou expiré")
	ErrPhoneAlreadyUsed  = errors.New("ce numéro de téléphone est déjà utilisé")
	ErrEmailAlreadyUsed  = errors.New("cet email est déjà utilisé")
)

// TokenPair regroupe un access token et un refresh token.
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // secondes
}

// JWTClaims étend les claims JWT standard.
type JWTClaims struct {
	UserID uint   `json:"uid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService gère l'authentification: OTP, JWT, mots de passe.
type AuthService struct {
	userRepo *repository.UserRepo
	notif    NotificationService
	cfg      *config.Config
}

// NewAuthService crée un nouveau AuthService.
func NewAuthService(userRepo *repository.UserRepo, notif NotificationService, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		notif:    notif,
		cfg:      cfg,
	}
}

// SendOTP génère un code OTP, le hache, le sauvegarde et l'envoie par notification.
func (s *AuthService) SendOTP(phone, purpose string, channel NotificationChannel) error {
	// Vérifier le rate limit
	count, err := s.userRepo.CountOTPLastHour(phone)
	if err != nil {
		return fmt.Errorf("vérification rate limit: %w", err)
	}
	if count >= int64(s.cfg.OTP.RateLimitPerHour) {
		return ErrOTPRateLimit
	}

	// Générer le code OTP
	code, err := generateOTPCode(s.cfg.OTP.Length)
	if err != nil {
		return fmt.Errorf("génération OTP: %w", err)
	}

	// Hacher le code
	hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hachage OTP: %w", err)
	}

	otp := &models.OTPCode{
		Phone:     phone,
		CodeHash:  string(hash),
		Channel:   string(channel),
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(time.Duration(s.cfg.OTP.TTLMinutes) * time.Minute),
	}

	if err := s.userRepo.SaveOTP(otp); err != nil {
		return fmt.Errorf("sauvegarde OTP: %w", err)
	}

	// Envoyer la notification
	return s.notif.SendOTP(phone, code, channel)
}

// VerifyOTP vérifie un code OTP et le marque comme utilisé si valide.
func (s *AuthService) VerifyOTP(phone, code, purpose string) error {
	otp, err := s.userRepo.FindValidOTPSQLite(phone, purpose)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrInvalidOTP
		}
		return fmt.Errorf("recherche OTP: %w", err)
	}

	// Vérifier le nombre de tentatives
	if otp.Attempts >= s.cfg.OTP.MaxAttempts {
		return ErrOTPMaxAttempts
	}

	// Incrémenter les tentatives
	otp.Attempts++
	if saveErr := s.userRepo.UpdateOTP(otp); saveErr != nil {
		return fmt.Errorf("mise à jour tentatives: %w", saveErr)
	}

	// Vérifier le code
	if err := bcrypt.CompareHashAndPassword([]byte(otp.CodeHash), []byte(code)); err != nil {
		return ErrInvalidOTP
	}

	// Marquer comme utilisé
	now := time.Now()
	otp.UsedAt = &now
	return s.userRepo.UpdateOTP(otp)
}

// RegisterWithOTP crée un compte utilisateur après vérification OTP.
func (s *AuthService) RegisterWithOTP(phone, displayName string) (*models.User, *TokenPair, error) {
	if s.userRepo.ExistsByPhone(phone) {
		return nil, nil, ErrPhoneAlreadyUsed
	}

	user := &models.User{
		Phone:       phone,
		DisplayName: displayName,
		IsVerified:  true,
		Role:        "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, fmt.Errorf("création utilisateur: %w", err)
	}

	tokens, err := s.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// LoginWithOTP connecte un utilisateur existant via OTP (téléphone).
func (s *AuthService) LoginWithOTP(phone string) (*models.User, *TokenPair, error) {
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, fmt.Errorf("recherche utilisateur: %w", err)
	}

	tokens, err := s.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// RegisterWithPassword crée un compte avec mot de passe.
func (s *AuthService) RegisterWithPassword(phone, email, password, displayName string) (*models.User, *TokenPair, error) {
	if s.userRepo.ExistsByPhone(phone) {
		return nil, nil, ErrPhoneAlreadyUsed
	}
	if email != "" && s.userRepo.ExistsByEmail(email) {
		return nil, nil, ErrEmailAlreadyUsed
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, fmt.Errorf("hachage mot de passe: %w", err)
	}
	hashStr := string(hash)

	user := &models.User{
		Phone:        phone,
		DisplayName:  displayName,
		PasswordHash: &hashStr,
		IsVerified:   false,
		Role:         "user",
	}
	if email != "" {
		user.Email = &email
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, fmt.Errorf("création utilisateur: %w", err)
	}

	tokens, err := s.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// LoginWithPassword authentifie un utilisateur par téléphone ou email + mot de passe.
func (s *AuthService) LoginWithPassword(identifier, password string) (*models.User, *TokenPair, error) {
	user, err := s.userRepo.FindByPhoneOrEmail(identifier)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, ErrUserNotFound
		}
		return nil, nil, fmt.Errorf("recherche utilisateur: %w", err)
	}

	if user.PasswordHash == nil {
		return nil, nil, ErrInvalidPassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, ErrInvalidPassword
	}

	tokens, err := s.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, tokens, nil
}

// GenerateTokenPair génère un access token JWT et un refresh token opaque.
func (s *AuthService) GenerateTokenPair(userID uint) (*TokenPair, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("utilisateur introuvable: %w", err)
	}

	accessTTL := time.Duration(s.cfg.JWT.AccessTTLMinutes) * time.Minute
	now := time.Now()

	claims := JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTTL)),
			Issuer:    "mezian",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		return nil, fmt.Errorf("signature JWT: %w", err)
	}

	// Refresh token opaque (UUID)
	rawToken := uuid.NewString()
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	refreshTTL := time.Duration(s.cfg.JWT.RefreshTTLDays) * 24 * time.Hour
	rt := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(refreshTTL),
	}
	if err := s.userRepo.SaveRefreshToken(rt); err != nil {
		return nil, fmt.Errorf("sauvegarde refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawToken,
		ExpiresIn:    int(accessTTL.Seconds()),
	}, nil
}

// RefreshTokens valide un refresh token et génère une nouvelle paire.
func (s *AuthService) RefreshTokens(rawToken string) (*TokenPair, error) {
	hash := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(hash[:])

	rt, err := s.userRepo.FindRefreshToken(tokenHash)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !rt.IsValid() {
		return nil, ErrInvalidToken
	}

	// Rotation: révoquer l'ancien token
	if err := s.userRepo.RevokeRefreshToken(rt.ID); err != nil {
		return nil, fmt.Errorf("révocation ancien token: %w", err)
	}

	return s.GenerateTokenPair(rt.UserID)
}

// ValidateAccessToken valide un access token JWT et retourne les claims.
func (s *AuthService) ValidateAccessToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algorithme inattendu: %v", t.Header["alg"])
		}
		return []byte(s.cfg.JWT.Secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// Logout révoque tous les refresh tokens de l'utilisateur.
func (s *AuthService) Logout(userID uint) error {
	return s.userRepo.RevokeAllUserTokens(userID)
}

// generateOTPCode génère un code numérique aléatoire de longueur donnée.
func generateOTPCode(length int) (string, error) {
	if length <= 0 {
		length = 6
	}
	digits := "0123456789"
	result := make([]byte, length)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		result[i] = digits[n.Int64()]
	}
	return string(result), nil
}
