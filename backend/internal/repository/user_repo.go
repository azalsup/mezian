// Package repository encapsule les accès base de données.
package repository

import (
	"errors"
	"time"

	"mezian/internal/models"

	"gorm.io/gorm"
)

// UserRepo gère les opérations base de données sur les utilisateurs.
type UserRepo struct{ db *gorm.DB }

// NewUserRepo crée un nouveau UserRepo.
func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db} }

// FindByID récupère un utilisateur par son ID.
func (r *UserRepo) FindByID(id uint) (*models.User, error) {
	var u models.User
	err := r.db.First(&u, id).Error
	return &u, err
}

// FindByPhone récupère un utilisateur par son numéro de téléphone.
func (r *UserRepo) FindByPhone(phone string) (*models.User, error) {
	var u models.User
	err := r.db.Where("phone = ?", phone).First(&u).Error
	return &u, err
}

// FindByEmail récupère un utilisateur par son email.
func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Where("email = ?", email).First(&u).Error
	return &u, err
}

// FindByPhoneOrEmail récupère un utilisateur par téléphone ou email.
func (r *UserRepo) FindByPhoneOrEmail(identifier string) (*models.User, error) {
	var u models.User
	err := r.db.Where("phone = ? OR email = ?", identifier, identifier).First(&u).Error
	return &u, err
}

// Create insère un nouvel utilisateur.
func (r *UserRepo) Create(u *models.User) error {
	return r.db.Create(u).Error
}

// Update sauvegarde les modifications d'un utilisateur.
func (r *UserRepo) Update(u *models.User) error {
	return r.db.Save(u).Error
}

// ExistsByPhone retourne true si un utilisateur avec ce numéro existe.
func (r *UserRepo) ExistsByPhone(phone string) bool {
	var count int64
	r.db.Model(&models.User{}).Where("phone = ?", phone).Count(&count)
	return count > 0
}

// ExistsByEmail retourne true si un utilisateur avec cet email existe.
func (r *UserRepo) ExistsByEmail(email string) bool {
	if email == "" {
		return false
	}
	var count int64
	r.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// SaveOTP sauvegarde un code OTP (supprime les anciens pour ce phone+purpose).
func (r *UserRepo) SaveOTP(otp *models.OTPCode) error {
	r.db.Where("phone = ? AND purpose = ? AND used_at IS NULL", otp.Phone, otp.Purpose).
		Delete(&models.OTPCode{})
	return r.db.Create(otp).Error
}

// FindValidOTPSQLite cherche un OTP valide (non expiré, non utilisé) — compatible SQLite.
func (r *UserRepo) FindValidOTPSQLite(phone, purpose string) (*models.OTPCode, error) {
	var otp models.OTPCode
	err := r.db.Where(
		"phone = ? AND purpose = ? AND used_at IS NULL AND expires_at > datetime('now')",
		phone, purpose,
	).Order("created_at DESC").First(&otp).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &otp, err
}

// UpdateOTP sauvegarde les modifications d'un OTP (ex: marquer comme utilisé).
func (r *UserRepo) UpdateOTP(otp *models.OTPCode) error {
	return r.db.Save(otp).Error
}

// CountOTPLastHour retourne le nombre d'OTP envoyés au cours de la dernière heure.
func (r *UserRepo) CountOTPLastHour(phone string) (int64, error) {
	var count int64
	err := r.db.Model(&models.OTPCode{}).
		Where("phone = ? AND created_at > datetime('now', '-1 hour')", phone).
		Count(&count).Error
	return count, err
}

// SaveRefreshToken persiste un refresh token.
func (r *UserRepo) SaveRefreshToken(rt *models.RefreshToken) error {
	return r.db.Create(rt).Error
}

// FindRefreshToken récupère un refresh token par son hash.
func (r *UserRepo) FindRefreshToken(hash string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.Where("token_hash = ?", hash).First(&rt).Error
	return &rt, err
}

// RevokeRefreshToken révoque un refresh token par son ID.
func (r *UserRepo) RevokeRefreshToken(id uint) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).Where("id = ?", id).Update("revoked_at", &now).Error
}

// RevokeAllUserTokens révoque tous les refresh tokens actifs d'un utilisateur.
func (r *UserRepo) RevokeAllUserTokens(userID uint) error {
	now := time.Now()
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", &now).Error
}
