package repositories

import (
	"context"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
)

type VerificationRepository interface {
	// Save stores a verification token
	Save(ctx context.Context, token *entities.VerificationToken) error

	// FindByToken retrieves a verification token by its token string
	FindByToken(ctx context.Context, token string) (*entities.VerificationToken, error)

	// Update updates a verification token (e.g., marking as used)
	Update(ctx context.Context, token *entities.VerificationToken) error

	// Delete removes a verification token
	Delete(ctx context.Context, token string) error

	// CleanupExpired removes all expired tokens
	CleanupExpired(ctx context.Context) error
}
