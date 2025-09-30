package services

import (
	"errors"

	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/models"
	"github.com/zandomed/sync-playlist-api/internal/repository"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

type AuthService interface {
	LoginWithPass(req LoginRequest) (*TokenPair, error)
	LoginWithOAuth(provider, token string) (string, error)
	Register(req RegisterRequest) (*RegisterResponse, error)
}

type AuthServiceImpl struct {
	authRepo repository.AuthRepository
	tokenSvc TokenService
	cfg      *config.Config
	logger   *logger.Logger
}

func NewAuthService(authRepo repository.AuthRepository, tokenSvc TokenService, cfg *config.Config, logger *logger.Logger) AuthService {
	return &AuthServiceImpl{
		authRepo: authRepo,
		tokenSvc: tokenSvc,
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

func (s *AuthServiceImpl) LoginWithPass(req LoginRequest) (*TokenPair, error) {

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

	user, err := s.authRepo.GetUserByEmail(req.Email)
	if err != nil {
		s.logger.Sugar().Errorf("Error fetching user %s: %v", req.Email, err)
		return nil, err
	}

	tokens, err := s.tokenSvc.GenerateTokenPair(user.ID.String(), req.Email)
	if err != nil {
		s.logger.Sugar().Errorf("Error generating tokens for user %s: %v", req.Email, err)
		return nil, err
	}

	// Token generation logic would go here
	// For now, returning a placeholder token
	return tokens, nil
}

func (s *AuthServiceImpl) LoginWithOAuth(provider, token string) (string, error) {
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

func (s *AuthServiceImpl) Register(req RegisterRequest) (*RegisterResponse, error) {
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
