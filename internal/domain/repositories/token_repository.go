package repositories

import (
	"context"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type TokenRepository interface {
	SaveRefreshToken(ctx context.Context, token *entities.RefreshToken) error
	FindRefreshToken(ctx context.Context, token string) (*entities.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteUserRefreshTokens(ctx context.Context, userID valueobjects.UserID) error
	CleanupExpiredTokens(ctx context.Context) error
}