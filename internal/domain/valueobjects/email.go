package valueobjects

import (
	"regexp"
	"strings"

	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
)

type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return Email{}, errors.NewDomainError("empty_email", "Email cannot be empty")
	}

	if len(email) > 254 {
		return Email{}, errors.NewDomainError("email_too_long", "Email cannot exceed 254 characters")
	}

	if !emailRegex.MatchString(email) {
		return Email{}, errors.NewDomainError("invalid_email", "Email format is invalid")
	}

	return Email{value: email}, nil
}

func (e Email) Value() string {
	return e.value
}

func (e Email) String() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

func (e Email) Domain() string {
	parts := strings.Split(e.value, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}