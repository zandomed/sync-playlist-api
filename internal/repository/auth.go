package repository

import (
	"github.com/zandomed/sync-playlist-api/internal/models"
	"github.com/zandomed/sync-playlist-api/pkg/database"
)

type AuthRepository interface {
	GetAccountByProvider(userID, provider string) (string, error)
	CreateUser(user *models.User) (string, error)
	CreateUserWithOAuth(email, name, lastName, provider, accessToken, refreshToken, scope string, atExpiresIn, rtExpiresIn int64) (string, error)
	ValidateUser(email, password string) (bool, error)
	ValidateOAuthUser(provider, accessToken string) (string, error)
	LinkAccount(userID, provider, accessToken, refreshToken, scope string, atExpiresIn, rtExpiresIn int64) error
	UnlinkAccount(userID, provider string) error
	ChangePassword(email, newPassword string) error
	ResetPassword(email, token, newPassword string) error
	RefreshTokens(userID, provider, newAccessToken, newRefreshToken string, atExpiresIn, rtExpiresIn int64) error
}

type authRepository struct {
	db *database.DB
}

func NewAuthRepository(db *database.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) ValidateUser(email, password string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM accounts WHERE provider = 'userpass' AND email = $1 AND password = crypt($2, password))`
	err := r.db.QueryRow(query, email, password).Scan(&exists)
	return exists, err
}
func (r *authRepository) ValidateOAuthUser(provider, accessToken string) (string, error) {
	var email string
	query := `SELECT email FROM accounts WHERE provider = $1 AND access_token = $2`
	err := r.db.QueryRow(query, provider, accessToken).Scan(&email)
	return email, err
}
func (r *authRepository) LinkAccount(userID, provider, accessToken, refreshToken, scope string, atExpiresIn, rtExpiresIn int64) error {
	// Implementation goes here
	return nil
}
func (r *authRepository) UnlinkAccount(userID, provider string) error {
	// Implementation goes here
	return nil
}
func (r *authRepository) ChangePassword(email, newPassword string) error {
	// Implementation goes here
	return nil
}
func (r *authRepository) ResetPassword(email, token, newPassword string) error {
	// Implementation goes here
	return nil
}
func (r *authRepository) RefreshTokens(userID, provider, newAccessToken, newRefreshToken string, atExpiresIn, rtExpiresIn int64) error {
	// Implementation goes here
	return nil
}

func (r *authRepository) GetAccountByProvider(userID, provider string) (string, error) {
	var accountID string
	query := `SELECT id FROM accounts WHERE user_id = $1 AND provider = $2`
	err := r.db.QueryRow(query, userID, provider).Scan(&accountID)
	return accountID, err
}
func (r *authRepository) CreateUser(user *models.User) (string, error) {
	var userID string
	query := `INSERT INTO users (email, name, last_name) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, user.Email, user.Name, user.LastName).Scan(&userID)
	return userID, err
}
func (r *authRepository) CreateUserWithOAuth(email, name, lastName, provider, accessToken, refreshToken, scope string, atExpiresIn, rtExpiresIn int64) (string, error) {
	var userID string
	query := `INSERT INTO users (email, name, last_name) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, email, name, lastName).Scan(&userID)
	if err != nil {
		return "", err
	}

	query = `INSERT INTO accounts (user_id, provider, access_token, refresh_token, scope, at_expires_in, rt_expires_in) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = r.db.Exec(query, userID, provider, accessToken, refreshToken, scope, atExpiresIn, rtExpiresIn)
	return userID, err
}
