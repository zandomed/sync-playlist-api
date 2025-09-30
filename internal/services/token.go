package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/repository"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // in seconds
	TokenType    string `json:"tokenType"` // e.g., "Bearer"
}

type TokenService interface {
	GenerateTokenPair(ID string, email string) (*TokenPair, error)
	// ValidateAccessToken(tokenStr string) (*jwt.RegisteredClaims, error)
	RefreshAccessToken(refreshToken string) (*TokenPair, error)
}

type TokenServiceImpl struct {
	cfg       *config.Config
	logger    *logger.Logger
	tokenRepo repository.TokenRepository
}

func NewTokenService(cfg *config.Config, logger *logger.Logger, tokenRepo repository.TokenRepository) TokenService {
	return &TokenServiceImpl{
		cfg:       cfg,
		logger:    logger,
		tokenRepo: tokenRepo,
	}
}

func (ts *TokenServiceImpl) GenerateTokenPair(ID string, email string) (*TokenPair, error) {

	accessTokenClaims := jwt.RegisteredClaims{
		Issuer:    "sync-playlist-api",
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ts.cfg.JWT.AccessTokenExpirationTime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	accessTokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	accessToken, err := accessTokenWithClaims.SignedString([]byte(ts.cfg.JWT.AccessTokenSecret))

	if err != nil {
		return nil, err
	}

	refreshTokenDB, err := ts.tokenRepo.StoreRefreshToken(ID, time.Now().Add(ts.cfg.JWT.RefreshTokenExpirationTime))
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := jwt.RegisteredClaims{
		Issuer:    "sync-playlist-api",
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ts.cfg.JWT.RefreshTokenExpirationTime)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        refreshTokenDB.ID.String(),
	}

	refreshTokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	refreshToken, err := refreshTokenWithClaims.SignedString([]byte(ts.cfg.JWT.RefreshTokenSecret))

	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(ts.cfg.JWT.AccessTokenExpirationTime.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (ts *TokenServiceImpl) RefreshAccessToken(refreshToken string) (*TokenPair, error) {

	return nil, nil

}
