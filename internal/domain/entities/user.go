package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type User struct {
	id              valueobjects.UserID
	email           valueobjects.Email
	profile         valueobjects.UserProfile
	isEmailVerified bool
	createdAt       time.Time
	updatedAt       time.Time
}

func NewUser(email string, name string, lastName string) (*User, error) {
	userID := valueobjects.NewUserID()

	emailVO, err := valueobjects.NewEmail(email)
	if err != nil {
		return nil, err
	}

	profile, err := valueobjects.NewUserProfile(name, lastName)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		id:              userID,
		email:           emailVO,
		profile:         profile,
		isEmailVerified: false,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

func ReconstructUser(id uuid.UUID, email, name, lastName string, isEmailVerified bool, createdAt, updatedAt time.Time) (*User, error) {
	userID, err := valueobjects.ReconstructUserID(id)
	if err != nil {
		return nil, err
	}

	emailVO, err := valueobjects.NewEmail(email)
	if err != nil {
		return nil, err
	}

	profile, err := valueobjects.NewUserProfile(name, lastName)
	if err != nil {
		return nil, err
	}

	return &User{
		id:              userID,
		email:           emailVO,
		profile:         profile,
		isEmailVerified: isEmailVerified,
		createdAt:       createdAt,
		updatedAt:       updatedAt,
	}, nil
}

func (u *User) ID() valueobjects.UserID {
	return u.id
}

func (u *User) Email() valueobjects.Email {
	return u.email
}

func (u *User) Profile() valueobjects.UserProfile {
	return u.profile
}

func (u *User) IsEmailVerified() bool {
	return u.isEmailVerified
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u *User) UpdateProfile(name, lastName string) error {
	profile, err := valueobjects.NewUserProfile(name, lastName)
	if err != nil {
		return err
	}

	u.profile = profile
	u.updatedAt = time.Now()
	return nil
}

func (u *User) VerifyEmail() {
	u.isEmailVerified = true
	u.updatedAt = time.Now()
}

func (u *User) ChangeEmail(email string) error {
	emailVO, err := valueobjects.NewEmail(email)
	if err != nil {
		return err
	}

	u.email = emailVO
	u.isEmailVerified = false
	u.updatedAt = time.Now()
	return nil
}

func (u *User) CanAuthenticate() error {
	return nil // Placeholder for future authentication checks
	if !u.isEmailVerified {
		return errors.NewDomainError("email_not_verified", "Email must be verified before authentication")
	}
	return nil
}
