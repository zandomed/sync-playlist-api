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

type PostgresAccountRepository struct {
	db *database.DB
}

func NewPostgresAccountRepository(db *database.DB) repositories.AccountRepository {
	return &PostgresAccountRepository{db: db}
}

func (r *PostgresAccountRepository) Save(ctx context.Context, account *entities.Account) error {
	query := `
		INSERT INTO accounts (id, user_id, provider, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			provider = EXCLUDED.provider,
			password = EXCLUDED.password,
			updated_at = EXCLUDED.updated_at`

	var password interface{}
	if account.IsUserpassAccount() && !account.Password().IsEmpty() {
		password = account.Password().Value()
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		account.ID().Value(),
		account.UserID().Value(),
		string(account.Provider()),
		password,
		account.CreatedAt(),
		account.UpdatedAt(),
	)

	return err
}

func (r *PostgresAccountRepository) FindByUserIDAndProvider(
	ctx context.Context,
	userID valueobjects.UserID,
	provider entities.AccountProvider,
) (*entities.Account, error) {
	query := `
		SELECT id, user_id, provider, password, created_at, updated_at
		FROM accounts
		WHERE user_id = $1 AND provider = $2`

	var accountID, userIDStr, providerStr string
	var password sql.NullString
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, userID.Value(), string(provider)).Scan(
		&accountID, &userIDStr, &providerStr, &password, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("account", "Account not found")
		}
		return nil, err
	}

	parsedAccountID, err := uuid.Parse(accountID)
	if err != nil {
		return nil, err
	}

	parsedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	passwordValue := ""
	if password.Valid {
		passwordValue = password.String
	}

	return entities.ReconstructAccount(parsedAccountID, parsedUserID, providerStr, passwordValue, createdAt, updatedAt)
}

func (r *PostgresAccountRepository) FindUserpassAccountByEmail(
	ctx context.Context,
	email valueobjects.Email,
) (*entities.Account, error) {
	query := `
		SELECT a.id, a.user_id, a.provider, a.password, a.created_at, a.updated_at
		FROM accounts a
		JOIN users u ON a.user_id = u.id
		WHERE u.email = $1 AND a.provider = 'userpass'`

	var accountID, userIDStr, providerStr string
	var password sql.NullString
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, email.Value()).Scan(
		&accountID, &userIDStr, &providerStr, &password, &createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NewNotFoundError("account", "Account not found")
		}
		return nil, err
	}

	parsedAccountID, err := uuid.Parse(accountID)
	if err != nil {
		return nil, err
	}

	parsedUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	passwordValue := ""
	if password.Valid {
		passwordValue = password.String
	}

	return entities.ReconstructAccount(parsedAccountID, parsedUserID, providerStr, passwordValue, createdAt, updatedAt)
}

func (r *PostgresAccountRepository) Delete(ctx context.Context, id valueobjects.AccountID) error {
	query := `DELETE FROM accounts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id.Value())
	return err
}