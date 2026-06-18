package services

import (
	"errors"
	"time"

	"paddle-traceability/config"
	"paddle-traceability/database"
	"paddle-traceability/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	cfg *config.JWTConfig
}

func NewAuthService(cfg *config.JWTConfig) *AuthService {
	return &AuthService{cfg: cfg}
}

func (s *AuthService) Register(username, password, role, companyName, phone string) (*models.User, error) {
	var existing models.User
	if err := database.DB.Where("username = ?", username).First(&existing).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, errors.New("password encryption failed")
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         role,
		CompanyName:  companyName,
		Phone:        phone,
		CreatedAt:    time.Now(),
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, errors.New("create user failed")
	}

	return user, nil
}

func (s *AuthService) Login(username, password string) (string, *models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, errors.New("username or password incorrect")
		}
		return "", nil, errors.New("query user failed")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("username or password incorrect")
	}

	token, err := s.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", nil, errors.New("generate token failed")
	}

	return token, &user, nil
}

func (s *AuthService) GenerateJWT(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Duration(s.cfg.Expiry) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
