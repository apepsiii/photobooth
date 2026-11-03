package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"mastergin/core"
	"mastergin/models"
)

// Service encapsulates authentication business logic.
// It depends on a *gorm.DB for user persistence and the app Config for JWT signing.
type Service struct {
	db  *gorm.DB
	cfg *core.Config
}

// NewService creates a new auth Service.
func NewService(db *gorm.DB, cfg *core.Config) *Service {
	return &Service{db: db, cfg: cfg}
}

// Register creates a new user with a bcrypt-hashed password.
func (s *Service) Register(name, email, password string) (*models.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("name, email, and password are required")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	// Check if email already exists.
	var existing models.User
	if err := s.db.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, errors.New("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "user",
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Login verifies credentials and returns a signed JWT token string.
func (s *Service) Login(email, password string) (string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	return s.generateToken(&user)
}

// generateToken creates a signed JWT containing the user's ID, name, and role.
func (s *Service) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"uid":  user.ID,
		"name": user.Name,
		"role": user.Role,
		"exp":  time.Now().Add(72 * time.Hour).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

// SeedAdmin creates a default admin user if the users table is empty.
// This makes the boilerplate immediately usable on first run.
func SeedAdmin(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := models.User{
		Name:         "Admin",
		Email:        "admin@meridian.com",
		PasswordHash: string(hash),
		Role:         "admin",
	}
	return db.Create(&admin).Error
}
