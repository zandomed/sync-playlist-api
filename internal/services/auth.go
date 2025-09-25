package services

import (
	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/models"
	"github.com/zandomed/sync-playlist-api/internal/repository"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

type AuthService interface {
	LoginWithPass(email, password string) (string, error)
	LoginWithOAuth(provider, token string) (string, error)
	Register(req RegisterRequest) (string, error)
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
func (s *authService) LoginWithPass(email, password string) (string, error) {
	// Implementation goes here
	return "", nil
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

func (s *authService) Register(req RegisterRequest) (string, error) {
	// Implementation goes here
	s.logger.Sugar().Infof("Registering user: %s", req.Email)

	user := models.User{
		Email:    req.Email,
		Name:     req.Name,
		LastName: req.LastName,
	}

	userID, err := s.authRepo.CreateUser(&user)

	if err != nil {
		return "", err
	}

	account := models.Account{
		Provider: models.Userpass,
		Password: req.Password,
		UserID:   userID,
	}

	if err := s.authRepo.CreateAccount(&account); err == nil {
		return userID, nil
	} else {
		return "", err
	}
}
