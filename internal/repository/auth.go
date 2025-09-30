package repository

import (
	"github.com/zandomed/sync-playlist-api/internal/models"
	"github.com/zandomed/sync-playlist-api/pkg/database"
)

type AuthRepository interface {
	CreateUserWithUserpass(user *models.User, password string) (*string, error)
	GetUserByEmail(email string) (*models.User, error)
	ValidateUser(email, password string) (bool, error)
}

type authRepository struct {
	db *database.DB
}

func NewAuthRepository(db *database.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) ValidateUser(email, password string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1 FROM accounts a
			JOIN users u ON a.user_id = u.id
			WHERE a.provider = 'userpass' AND u.email = $1 AND a.password = crypt($2, a.password)
		)`
	err := r.db.QueryRow(query, email, password).Scan(&exists)
	return exists, err
}

func (r *authRepository) CreateUserWithUserpass(user *models.User, password string) (*string, error) {
	queryUser := `
		INSERT INTO users (email, name, last_name)
		VALUES ($1, $2, $3)
		RETURNING id`
	queryAccount := `
		INSERT INTO accounts (user_id, provider, password) 
		VALUES ($1, $2, $3)`

	tx, err := r.db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	var userID string
	err = tx.QueryRow(queryUser, user.Email, user.Name, user.LastName).Scan(&userID)

	if err != nil {
		return nil, err
	}

	// Password will be automatically hashed by the database trigger
	_, err = tx.Exec(queryAccount, user.ID, models.Userpass, password)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &userID, nil
}

func (r *authRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, email, name, last_name
		FROM users
		WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Name, &user.LastName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
