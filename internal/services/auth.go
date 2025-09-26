package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/models"
	"github.com/zandomed/sync-playlist-api/internal/repository"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

type AuthService interface {
	LoginWithPass(req LoginRequest) (*LoginResponse, error)
	LoginWithOAuth(provider, token string) (string, error)
	Register(req RegisterRequest) (*RegisterResponse, error)
}

type authService struct {
	authRepo repository.AuthRepository
	cfg      *config.Config
	logger   *logger.Logger
}

func NewAuthService(authRepo repository.AuthRepository, cfg *config.Config, logger *logger.Logger) AuthService {
	return &authService{
		authRepo: authRepo,
		cfg:      cfg,
		logger:   logger,
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (s *authService) LoginWithPass(req LoginRequest) (*LoginResponse, error) {

	s.logger.Sugar().Infof("Attempting login for user: %s", req.Email)

	valid, err := s.authRepo.ValidateUser(req.Email, req.Password)
	if err != nil {
		s.logger.Sugar().Errorf("Error validating user %s: %v", req.Email, err)
		return nil, err
	}
	if !valid {
		s.logger.Sugar().Warnf("Invalid credentials for user: %s", req.Email)
		return nil, errors.New("invalid credentials")
	}
	s.logger.Sugar().Infof("User %s successfully authenticated", req.Email)

	claims := jwt.RegisteredClaims{
		Issuer:    "sync-playlist-api",
		Subject:   req.Email,
		Audience:  []string{"sync-playlist-client"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.JWT.ExpirationTime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		s.logger.Sugar().Errorf("Error signing token for user %s: %v", req.Email, err)
		return nil, err
	}
	// Token generation logic would go here
	// For now, returning a placeholder token
	return &LoginResponse{AccessToken: signedToken}, nil
}

func (s *authService) LoginWithOAuth(provider, token string) (string, error) {
	// Implementation goes here
	return "", nil
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	LastName string `json:"lastName" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}

func (s *authService) Register(req RegisterRequest) (*RegisterResponse, error) {
	// Implementation goes here
	s.logger.Sugar().Infof("Registering user: %s", req.Email)

	creatingUser := models.User{
		Email:    req.Email,
		Name:     req.Name,
		LastName: req.LastName,
	}

	userID, err := s.authRepo.CreateUserWithUserpass(&creatingUser, req.Password)

	if err != nil {
		return nil, err
	}

	s.logger.Sugar().Infof("User registered with ID: %s", *userID)
	return &RegisterResponse{ID: *userID}, nil

}
