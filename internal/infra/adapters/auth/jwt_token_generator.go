package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	usecases_auth "github.com/zandomed/sync-playlist-api/internal/usecases/auth"
)

type JWTTokenGenerator struct {
	secretKey              string
	accessTokenExpiration  time.Duration
	refreshTokenExpiration time.Duration
}

func NewJWTTokenGenerator(secretKey string, accessTokenExp, refreshTokenExp time.Duration) usecases_auth.TokenGenerator {
	return &JWTTokenGenerator{
		secretKey:              secretKey,
		accessTokenExpiration:  accessTokenExp,
		refreshTokenExpiration: refreshTokenExp,
	}
}

type AccessTokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (g *JWTTokenGenerator) GenerateAccessToken(userID string, email string) (string, error) {

	claims := &AccessTokenClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.accessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(g.secretKey))
	if err != nil {
		return "", errors.NewDomainError("token_generation_failed", "Failed to generate access token")
	}

	return tokenString, nil
}

func (g *JWTTokenGenerator) GenerateRefreshToken(userID string) (string, error) {
	claims := &RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.refreshTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(g.secretKey))
	if err != nil {
		return "", errors.NewDomainError("token_generation_failed", "Failed to generate refresh token")
	}

	return tokenString, nil
}

func (g *JWTTokenGenerator) GetRefreshTokenExpiration() int64 {
	return int64(g.refreshTokenExpiration.Seconds())
}

func (g *JWTTokenGenerator) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewDomainError("invalid_signing_method", "Invalid signing method")
		}
		return []byte(g.secretKey), nil
	})

	if err != nil {
		return nil, errors.NewDomainError("invalid_token", "Invalid or expired token")
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.NewDomainError("invalid_token", "Invalid token claims")
}
