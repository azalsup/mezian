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

// Business errors for authentication.
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrInvalidOTP        = errors.New("invalid or expired OTP code")
    ErrOTPMaxAttempts    = errors.New("too many OTP attempts")
    ErrOTPRateLimit      = errors.New("too many OTP requests, try again in 1 hour")
    ErrInvalidPassword   = errors.New("mot de passe incorrect")
    ErrInvalidToken      = errors.New("invalid or expired token")
    ErrPhoneAlreadyUsed  = errors.New("this phone number is already used")
    ErrEmailAlreadyUsed  = errors.New("this email is already used")
)

// RegisterAddressFields holds optional address data for registration.
type RegisterAddressFields struct {
    Address    string
    City       string
    PostalCode string
    Country    string
}

// TokenPair groups an access token and a refresh token.
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int    `json:"expires_in"` // secondes
}

// JWTClaims extends the standard JWT claims.
type JWTClaims struct {
    UserID uint   `json:"uid"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

// AuthService handles authentication: OTP, JWT, and passwords.
type AuthService struct {
    userRepo *repository.UserRepo
    notif    NotificationService
    cfg      *config.Config
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo *repository.UserRepo, notif NotificationService, cfg *config.Config) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        notif:    notif,
        cfg:      cfg,
    }
}

// AuthConfig exposes the auth section of the config for use by handlers.
func (s *AuthService) AuthConfig() config.AuthConfig {
    return s.cfg.Auth
}

// SendOTP generates an OTP code, hashes it, saves it, and sends it by notification.
func (s *AuthService) SendOTP(phone, purpose string, channel NotificationChannel) error {
    // Check the rate limit
    count, err := s.userRepo.CountOTPLastHour(phone)
    if err != nil {
        return fmt.Errorf("rate limit verification: %w", err)
    }
    if count >= int64(s.cfg.OTP.RateLimitPerHour) {
        return ErrOTPRateLimit
    }

    // Generate the OTP code
    code, err := generateOTPCode(s.cfg.OTP.Length)
    if err != nil {
        return fmt.Errorf("OTP generation: %w", err)
    }

    // Hacher le code
    hash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("OTP hashing: %w", err)
    }

    otp := &models.OTPCode{
        Phone:     phone,
        CodeHash:  string(hash),
        Channel:   string(channel),
        Purpose:   purpose,
        ExpiresAt: time.Now().Add(time.Duration(s.cfg.OTP.TTLMinutes) * time.Minute),
    }

    if err := s.userRepo.SaveOTP(otp); err != nil {
        return fmt.Errorf("OTP saving: %w", err)
    }

    // Envoyer la notification
    return s.notif.SendOTP(phone, code, channel)
}

// VerifyOTP verifies an OTP code and marks it as used if valid.
func (s *AuthService) VerifyOTP(phone, code, purpose string) error {
    otp, err := s.userRepo.FindValidOTPSQLite(phone, purpose)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrInvalidOTP
        }
        return fmt.Errorf("OTP search: %w", err)
    }

    // Check the number of attempts
    if otp.Attempts >= s.cfg.OTP.MaxAttempts {
        return ErrOTPMaxAttempts
    }

    // Increment attempts
    otp.Attempts++
    if saveErr := s.userRepo.UpdateOTP(otp); saveErr != nil {
        return fmt.Errorf("update attempts: %w", saveErr)
    }

    // Verify the code
    if err := bcrypt.CompareHashAndPassword([]byte(otp.CodeHash), []byte(code)); err != nil {
        return ErrInvalidOTP
    }

    // Mark as used
    now := time.Now()
    otp.UsedAt = &now
    return s.userRepo.UpdateOTP(otp)
}

// RegisterWithOTP creates a user after OTP verification.
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
        return nil, nil, fmt.Errorf("user creation: %w", err)
    }

    tokens, err := s.GenerateTokenPair(user.ID)
    if err != nil {
        return nil, nil, err
    }

    return user, tokens, nil
}

// LoginWithOTP connects an existing user via OTP (phone).
func (s *AuthService) LoginWithOTP(phone string) (*models.User, *TokenPair, error) {
    user, err := s.userRepo.FindByPhone(phone)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil, ErrUserNotFound
        }
        return nil, nil, fmt.Errorf("user lookup: %w", err)
    }

    tokens, err := s.GenerateTokenPair(user.ID)
    if err != nil {
        return nil, nil, err
    }

    return user, tokens, nil
}

// RegisterWithPassword creates an account with a password.
func (s *AuthService) RegisterWithPassword(phone, email, password, displayName string, addr *RegisterAddressFields) (*models.User, *TokenPair, error) {
    if phone != "" && s.userRepo.ExistsByPhone(phone) {
        return nil, nil, ErrPhoneAlreadyUsed
    }
    if email != "" && s.userRepo.ExistsByEmail(email) {
        return nil, nil, ErrEmailAlreadyUsed
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, nil, fmt.Errorf("password hashing: %w", err)
    }
    hashStr := string(hash)

    // Use phone as the stored identifier; if phone is empty, store email as phone substitute.
    // The model requires phone not null, so we use a placeholder when phone is absent.
    storedPhone := phone
    if storedPhone == "" {
        storedPhone = email // email-only accounts: store email in phone field temporarily
    }

    user := &models.User{
        Phone:        storedPhone,
        DisplayName:  displayName,
        PasswordHash: &hashStr,
        IsVerified:   false,
        Role:         "user",
    }
    if email != "" {
        user.Email = &email
    }
    if addr != nil {
        if addr.Address != "" {
            user.Address = &addr.Address
        }
        if addr.City != "" {
            user.City = &addr.City
        }
        if addr.PostalCode != "" {
            user.PostalCode = &addr.PostalCode
        }
        country := addr.Country
        if country == "" {
            country = "MA"
        }
        user.Country = &country
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, nil, fmt.Errorf("user creation: %w", err)
    }

    tokens, err := s.GenerateTokenPair(user.ID)
    if err != nil {
        return nil, nil, err
    }

    return user, tokens, nil
}

// LoginWithPassword authenticates a user by phone or email and password.
func (s *AuthService) LoginWithPassword(identifier, password string) (*models.User, *TokenPair, error) {
    user, err := s.userRepo.FindByPhoneOrEmail(identifier)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil, ErrUserNotFound
        }
        return nil, nil, fmt.Errorf("user lookup: %w", err)
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
        return nil, fmt.Errorf("old token revocation: %w", err)
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

// Logout révoque tous les refresh tokens de l'user.
func (s *AuthService) Logout(userID uint) error {
    return s.userRepo.RevokeAllUserTokens(userID)
}

// generateOTPCode generates a random numeric code of the given length.
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
