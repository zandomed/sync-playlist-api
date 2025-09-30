package entities

import (
	"time"

	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type RefreshToken struct {
	id        valueobjects.TokenID
	userID    valueobjects.UserID
	token     string
	expiresAt time.Time
	createdAt time.Time
}

func NewRefreshToken(userID valueobjects.UserID, token string, expiresAt time.Time) (*RefreshToken, error) {
	if token == "" {
		return nil, errors.NewDomainError("empty_token", "Token cannot be empty")
	}

	if expiresAt.Before(time.Now()) {
		return nil, errors.NewDomainError("expired_token", "Token expiration time cannot be in the past")
	}

	now := time.Now()
	return &RefreshToken{
		id:        valueobjects.NewTokenID(),
		userID:    userID,
		token:     token,
		expiresAt: expiresAt,
		createdAt: now,
	}, nil
}

func ReconstructRefreshToken(userID valueobjects.UserID, token string, expiresAt, createdAt time.Time) (*RefreshToken, error) {
	if token == "" {
		return nil, errors.NewDomainError("empty_token", "Token cannot be empty")
	}

	return &RefreshToken{
		id:        valueobjects.NewTokenID(),
		userID:    userID,
		token:     token,
		expiresAt: expiresAt,
		createdAt: createdAt,
	}, nil
}

func (rt *RefreshToken) ID() valueobjects.TokenID {
	return rt.id
}

func (rt *RefreshToken) UserID() valueobjects.UserID {
	return rt.userID
}

func (rt *RefreshToken) Token() string {
	return rt.token
}

func (rt *RefreshToken) ExpiresAt() time.Time {
	return rt.expiresAt
}

func (rt *RefreshToken) CreatedAt() time.Time {
	return rt.createdAt
}

func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.expiresAt)
}

func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && rt.token != ""
}

func (rt *RefreshToken) TimeUntilExpiration() time.Duration {
	if rt.IsExpired() {
		return 0
	}
	return time.Until(rt.expiresAt)
}