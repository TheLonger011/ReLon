package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/TheLonger011/ReLon/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	repo      AuthInterface
	jwtsecret string
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

type AuthInterface interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetByEmailOrUsername(ctx context.Context, identifier string) (*models.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}

func NewAuthService(repo AuthInterface, jwtSecret string) *AuthService {
	return &AuthService{repo: repo, jwtsecret: jwtSecret}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}

	user := models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}

	err = s.repo.CreateUser(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) Login(ctx context.Context, identifier, password string) (*models.User, string, error) {
	user, err := s.repo.GetByEmailOrUsername(ctx, identifier)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.generateJWT(user.ID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate JWT: %w", err)
	}

	return user, token, nil
}

func (s *AuthService) generateJWT(userID uuid.UUID) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtsecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(ctx, userID)
}
