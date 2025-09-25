package models

import "github.com/google/uuid"

type Verification struct {
	ID         uuid.UUID `json:"id" db:"id"`
	Identifier string    `json:"identifier" db:"identifier"`
	Value      string    `json:"value" db:"value"`
	CreatedAt  string    `json:"createdAt" db:"created_at"`
	UpdatedAt  string    `json:"updatedAt" db:"updated_at"`
	ExpiresAt  string    `json:"expiresAt" db:"expires_at"`
}
