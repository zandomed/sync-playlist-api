package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type AccountProvider string

const (
	UserpassProvider AccountProvider = "userpass"
	SpotifyProvider  AccountProvider = "spotify"
	AppleProvider    AccountProvider = "apple"
	GoogleProvider   AccountProvider = "google"
)

type Account struct {
	id        valueobjects.AccountID
	userID    valueobjects.UserID
	provider  AccountProvider
	password  valueobjects.HashedPassword
	createdAt time.Time
	updatedAt time.Time
}

func NewUserpassAccount(userID valueobjects.UserID, password valueobjects.HashedPassword) *Account {
	now := time.Now()
	return &Account{
		id:        valueobjects.NewAccountID(),
		userID:    userID,
		provider:  UserpassProvider,
		password:  password,
		createdAt: now,
		updatedAt: now,
	}
}

func NewOAuthAccount(userID valueobjects.UserID, provider AccountProvider) (*Account, error) {
	if provider == UserpassProvider {
		return nil, errors.NewDomainError("invalid_provider", "Cannot create OAuth account with userpass provider")
	}

	now := time.Now()
	return &Account{
		id:        valueobjects.NewAccountID(),
		userID:    userID,
		provider:  provider,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func ReconstructAccount(id, userID uuid.UUID, provider string, password string, createdAt, updatedAt time.Time) (*Account, error) {
	accountID, err := valueobjects.ReconstructAccountID(id)
	if err != nil {
		return nil, err
	}

	userIDVO, err := valueobjects.ReconstructUserID(userID)
	if err != nil {
		return nil, err
	}

	accountProvider := AccountProvider(provider)
	if !isValidProvider(accountProvider) {
		return nil, errors.NewDomainError("invalid_provider", "Invalid account provider")
	}

	var hashedPassword valueobjects.HashedPassword
	if accountProvider == UserpassProvider && password != "" {
		hashedPassword = valueobjects.ReconstructHashedPassword(password)
	}

	return &Account{
		id:        accountID,
		userID:    userIDVO,
		provider:  accountProvider,
		password:  hashedPassword,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

func (a *Account) ID() valueobjects.AccountID {
	return a.id
}

func (a *Account) UserID() valueobjects.UserID {
	return a.userID
}

func (a *Account) Provider() AccountProvider {
	return a.provider
}

func (a *Account) Password() valueobjects.HashedPassword {
	return a.password
}

func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Account) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Account) ChangePassword(newPassword valueobjects.HashedPassword) error {
	if a.provider != UserpassProvider {
		return errors.NewDomainError("invalid_operation", "Cannot change password for non-userpass account")
	}

	a.password = newPassword
	a.updatedAt = time.Now()
	return nil
}

func (a *Account) IsUserpassAccount() bool {
	return a.provider == UserpassProvider
}

func isValidProvider(provider AccountProvider) bool {
	switch provider {
	case UserpassProvider, SpotifyProvider, AppleProvider, GoogleProvider:
		return true
	default:
		return false
	}
}