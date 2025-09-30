package valueobjects

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
)

type HashedPassword struct {
	value string
}

type PlainPassword struct {
	value string
}

func NewPlainPassword(password string) (PlainPassword, error) {
	if password == "" {
		return PlainPassword{}, errors.NewDomainError("empty_password", "Password cannot be empty")
	}

	if len(password) < 8 {
		return PlainPassword{}, errors.NewDomainError("password_too_short", "Password must be at least 8 characters")
	}

	if len(password) > 128 {
		return PlainPassword{}, errors.NewDomainError("password_too_long", "Password cannot exceed 128 characters")
	}

	return PlainPassword{value: password}, nil
}

func (p PlainPassword) Hash() (HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p.value), bcrypt.DefaultCost)
	if err != nil {
		return HashedPassword{}, errors.NewDomainError("password_hash_failed", "Failed to hash password")
	}

	return HashedPassword{value: string(hash)}, nil
}

func (p PlainPassword) Value() string {
	return p.value
}

func ReconstructHashedPassword(hash string) HashedPassword {
	return HashedPassword{value: hash}
}

func (h HashedPassword) Value() string {
	return h.value
}

func (h HashedPassword) Verify(plainPassword PlainPassword) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h.value), []byte(plainPassword.value))
	return err == nil
}

func (h HashedPassword) IsEmpty() bool {
	return h.value == ""
}