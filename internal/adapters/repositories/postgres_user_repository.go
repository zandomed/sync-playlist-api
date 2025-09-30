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

type PostgresUserRepository struct {
	db *database.DB
}

func NewPostgresUserRepository(db *database.DB) repositories.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (id, email, name, last_name, is_email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			last_name = EXCLUDED.last_name,
			is_email_verified = EXCLUDED.is_email_verified,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID().Value(),
		user.Email().Value(),
		user.Profile().Name(),
		user.Profile().LastName(),
		user.IsEmailVerified(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)

	return err
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id valueobjects.UserID) (*entities.User, error) {
	query := `
		SELECT id, email, name, last_name, is_email_verified, created_at, updated_at
		FROM users
		WHERE id = $1`

	var userID, email, name, lastName string
	var isEmailVerified bool
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(
		&userID, &email, &name, &lastName, &isEmailVerified, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("user", "User not found")
		}
		return nil, err
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return entities.ReconstructUser(parsedID, email, name, lastName, isEmailVerified, createdAt, updatedAt)
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email valueobjects.Email) (*entities.User, error) {
	query := `
		SELECT id, email, name, last_name, is_email_verified, created_at, updated_at
		FROM users
		WHERE email = $1`

	var userID, emailStr, name, lastName string
	var isEmailVerified bool
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, email.Value()).Scan(
		&userID, &emailStr, &name, &lastName, &isEmailVerified, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("user", "User not found")
		}
		return nil, err
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return entities.ReconstructUser(parsedID, emailStr, name, lastName, isEmailVerified, createdAt, updatedAt)
}

func (r *PostgresUserRepository) Exists(ctx context.Context, email valueobjects.Email) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email.Value()).Scan(&exists)
	return exists, err
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id valueobjects.UserID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id.Value())
	return err
}