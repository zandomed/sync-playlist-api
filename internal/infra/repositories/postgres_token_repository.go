package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
	"github.com/zandomed/sync-playlist-api/pkg/database"
)

type PostgresTokenRepository struct {
	db *database.DB
}

func NewPostgresTokenRepository(db *database.DB) repositories.TokenRepository {
	return &PostgresTokenRepository{db: db}
}

func (r *PostgresTokenRepository) SaveRefreshToken(ctx context.Context, token *entities.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (token) DO UPDATE SET
			expires_at = EXCLUDED.expires_at,
			created_at = EXCLUDED.created_at`

	_, err := r.db.ExecContext(
		ctx,
		query,
		token.UserID().Value(),
		token.Token(),
		token.ExpiresAt(),
		token.CreatedAt(),
	)

	return err
}

func (r *PostgresTokenRepository) FindRefreshToken(ctx context.Context, token string) (*entities.RefreshToken, error) {
	query := `
		SELECT user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1`

	var userIDStr, tokenStr string
	var expiresAt, createdAt time.Time

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&userIDStr, &tokenStr, &expiresAt, &createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("refresh_token", "Refresh token not found")
		}
		return nil, err
	}

	userID, err := valueobjects.ReconstructUserID(uuid.MustParse(userIDStr))
	if err != nil {
		return nil, err
	}

	refreshToken, err := entities.ReconstructRefreshToken(userID, tokenStr, expiresAt, createdAt)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (r *PostgresTokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *PostgresTokenRepository) DeleteUserRefreshTokens(ctx context.Context, userID valueobjects.UserID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID.Value())
	return err
}

func (r *PostgresTokenRepository) CleanupExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM refresh_tokens WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
