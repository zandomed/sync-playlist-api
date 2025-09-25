package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID `json:"id" db:"id"`
	Email           string    `json:"email" db:"email"`
	Name            string    `json:"name" db:"name"`
	LastName        string    `json:"lastName" db:"last_name"`
	IsEmailVerified bool      `json:"isEmailVerified" db:"is_email_verified"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

func (u *User) Validate() error {
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}
