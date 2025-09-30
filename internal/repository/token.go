package repository

import (
	"time"

	"github.com/zandomed/sync-playlist-api/internal/models"
	"github.com/zandomed/sync-playlist-api/pkg/database"
)

type TokenRepository interface {
	StoreRefreshToken(userID string, expiresAt time.Time) (*models.RefreshToken, error)
	ValidateRefreshToken(ID string, userID string) (bool, error)
	DeleteRefreshToken(ID string, userID string) error
	DeleteAllRefreshTokens(userID string) error
	InvalidateAllRefreshTokens(userID string) error
	InvalidateRefreshToken(ID string) error
}

type TokenRepositoryImpl struct {
	db *database.DB
}

func NewTokenRepository(db *database.DB) TokenRepository {
	return &TokenRepositoryImpl{db: db}
}

func (r *TokenRepositoryImpl) StoreRefreshToken(userID string, expiresAt time.Time) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.QueryRow(`
		INSERT INTO refresh_tokens (user_id, expires_at)
		VALUES ($1, $2)
		RETURNING id, user_id, expires_at, is_active, created_at, updated_at
	`, userID, expiresAt).Scan(&token.ID, &token.UserID, &token.ExpiresAt, &token.IsActive, &token.CreatedAt, &token.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepositoryImpl) ValidateRefreshToken(ID string, userID string) (bool, error) {
	var isActive bool
	err := r.db.QueryRow(`
		SELECT is_active FROM refresh_tokens
		WHERE id = $1 AND user_id = $2
	`, ID, userID).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func (r *TokenRepositoryImpl) DeleteRefreshToken(ID string, userID string) error {
	_, err := r.db.Exec(`
		DELETE FROM refresh_tokens
		WHERE id = $1 AND user_id = $2
	`, ID, userID)
	return err
}

func (r *TokenRepositoryImpl) DeleteAllRefreshTokens(userID string) error {
	_, err := r.db.Exec(`
		DELETE FROM refresh_tokens
		WHERE user_id = $1
	`, userID)
	return err
}

func (r *TokenRepositoryImpl) InvalidateAllRefreshTokens(userID string) error {
	_, err := r.db.Exec(`
		UPDATE refresh_tokens
		SET is_active = false
		WHERE user_id = $1
	`, userID)
	return err
}

func (r *TokenRepositoryImpl) InvalidateRefreshToken(ID string) error {
	_, err := r.db.Exec(`
		UPDATE refresh_tokens
		SET is_active = false
		WHERE id = $1
	`, ID)
	return err
}
