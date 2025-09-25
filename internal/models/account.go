package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	UserID                uuid.UUID `json:"userId" db:"user_id"`
	Provider              Provider  `json:"provider" db:"provider"`
	AccessToken           string    `json:"accessToken" db:"access_token"`
	RefreshToken          string    `json:"refreshToken" db:"refresh_token"`
	Scope                 string    `json:"scope" db:"scope"`
	AccessTokenExpiresAt  time.Time `json:"accessTokenExpiresAt" db:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time `json:"refreshTokenExpiresAt" db:"refresh_token_expires_at"`
	Password              string    `json:"-" db:"password"`
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`
}
