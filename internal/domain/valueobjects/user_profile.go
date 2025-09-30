package valueobjects

import (
	"strings"

	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
)

type UserProfile struct {
	name     string
	lastName string
}

func NewUserProfile(name, lastName string) (UserProfile, error) {
	name = strings.TrimSpace(name)
	lastName = strings.TrimSpace(lastName)

	if name == "" {
		return UserProfile{}, errors.NewDomainError("empty_name", "Name cannot be empty")
	}

	if lastName == "" {
		return UserProfile{}, errors.NewDomainError("empty_last_name", "Last name cannot be empty")
	}

	if len(name) < 2 {
		return UserProfile{}, errors.NewDomainError("name_too_short", "Name must be at least 2 characters")
	}

	if len(lastName) < 2 {
		return UserProfile{}, errors.NewDomainError("last_name_too_short", "Last name must be at least 2 characters")
	}

	if len(name) > 50 {
		return UserProfile{}, errors.NewDomainError("name_too_long", "Name cannot exceed 50 characters")
	}

	if len(lastName) > 50 {
		return UserProfile{}, errors.NewDomainError("last_name_too_long", "Last name cannot exceed 50 characters")
	}

	return UserProfile{
		name:     name,
		lastName: lastName,
	}, nil
}

func (p UserProfile) Name() string {
	return p.name
}

func (p UserProfile) LastName() string {
	return p.lastName
}

func (p UserProfile) FullName() string {
	return p.name + " " + p.lastName
}

func (p UserProfile) Equals(other UserProfile) bool {
	return p.name == other.name && p.lastName == other.lastName
}