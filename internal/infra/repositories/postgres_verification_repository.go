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

type PostgresVerificationRepository struct {
	db *database.DB
}

func NewPostgresVerificationRepository(db *database.DB) repositories.VerificationRepository {
	return &PostgresVerificationRepository{db: db}
}

func (r *PostgresVerificationRepository) Save(ctx context.Context, token *entities.VerificationToken) error {
	query := `
		INSERT INTO verification_tokens (token, token_type, user_id, expires_at, created_at, used_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	var userIDValue interface{}
	if token.UserID() != nil {
		userIDValue = token.UserID().Value()
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		token.Token(),
		string(token.TokenType()),
		userIDValue,
		token.ExpiresAt(),
		token.CreatedAt(),
		token.UsedAt(),
	)

	return err
}

func (r *PostgresVerificationRepository) FindByToken(ctx context.Context, tokenStr string) (*entities.VerificationToken, error) {
	query := `
		SELECT token, token_type, user_id, expires_at, created_at, used_at
		FROM verification_tokens
		WHERE token = $1`

	var token, tokenType string
	var userIDStr sql.NullString
	var expiresAt, createdAt time.Time
	var usedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, tokenStr).Scan(
		&token, &tokenType, &userIDStr, &expiresAt, &createdAt, &usedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("verification_token", "Verification token not found")
		}
		return nil, err
	}

	var userID *valueobjects.UserID
	if userIDStr.Valid {
		parsedUUID, err := uuid.Parse(userIDStr.String)
		if err != nil {
			return nil, errors.NewDomainError("invalid_user_id", "Invalid user ID format")
		}
		reconstructedUserID, err := valueobjects.ReconstructUserID(parsedUUID)
		if err != nil {
			return nil, err
		}
		userID = &reconstructedUserID
	}

	var usedAtPtr *time.Time
	if usedAt.Valid {
		usedAtPtr = &usedAt.Time
	}

	verificationToken, err := entities.ReconstructVerificationToken(
		token,
		entities.VerificationTokenType(tokenType),
		userID,
		expiresAt,
		createdAt,
		usedAtPtr,
	)
	if err != nil {
		return nil, err
	}

	return verificationToken, nil
}

func (r *PostgresVerificationRepository) Update(ctx context.Context, token *entities.VerificationToken) error {
	query := `
		UPDATE verification_tokens
		SET used_at = $1
		WHERE token = $2`

	_, err := r.db.ExecContext(ctx, query, token.UsedAt(), token.Token())
	return err
}

func (r *PostgresVerificationRepository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM verification_tokens WHERE token = $1`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

func (r *PostgresVerificationRepository) CleanupExpired(ctx context.Context) error {
	query := `DELETE FROM verification_tokens WHERE expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
