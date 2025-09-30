package valueobjects

import (
	"github.com/google/uuid"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
)

// ID is a generic UUID-based identifier value object
type ID struct {
	value uuid.UUID
}

// NewID creates a new ID with a random UUID
func NewID() ID {
	return ID{value: uuid.New()}
}

// ReconstructID reconstructs an ID from an existing UUID
func ReconstructID(id uuid.UUID) (ID, error) {
	if id == uuid.Nil {
		return ID{}, errors.NewDomainError("invalid_id", "ID cannot be nil")
	}
	return ID{value: id}, nil
}

// ParseID parses a string representation of UUID into an ID
func ParseID(s string) (ID, error) {
	parsed, err := uuid.Parse(s)
	if err != nil {
		return ID{}, errors.NewDomainError("invalid_id_format", "Invalid ID format")
	}
	return ReconstructID(parsed)
}

// Value returns the underlying UUID value
func (id ID) Value() uuid.UUID {
	return id.value
}

// String returns the string representation of the ID
func (id ID) String() string {
	return id.value.String()
}

// Equals checks if two IDs are equal
func (id ID) Equals(other ID) bool {
	return id.value == other.value
}

// IsEmpty checks if the ID is empty (nil UUID)
func (id ID) IsEmpty() bool {
	return id.value == uuid.Nil
}

// Type-safe ID aliases for different domain entities
type (
	UserID    ID
	AccountID ID
	TokenID   ID
)

// UserID specific constructors and methods
func NewUserID() UserID {
	return UserID(NewID())
}

func ReconstructUserID(id uuid.UUID) (UserID, error) {
	baseID, err := ReconstructID(id)
	if err != nil {
		return UserID{}, err
	}
	return UserID(baseID), nil
}

func ParseUserID(s string) (UserID, error) {
	baseID, err := ParseID(s)
	if err != nil {
		return UserID{}, err
	}
	return UserID(baseID), nil
}

func (id UserID) Value() uuid.UUID {
	return ID(id).Value()
}

func (id UserID) String() string {
	return ID(id).String()
}

func (id UserID) Equals(other UserID) bool {
	return ID(id).Equals(ID(other))
}

func (id UserID) IsEmpty() bool {
	return ID(id).IsEmpty()
}

// AccountID specific constructors and methods
func NewAccountID() AccountID {
	return AccountID(NewID())
}

func ReconstructAccountID(id uuid.UUID) (AccountID, error) {
	baseID, err := ReconstructID(id)
	if err != nil {
		return AccountID{}, err
	}
	return AccountID(baseID), nil
}

func ParseAccountID(s string) (AccountID, error) {
	baseID, err := ParseID(s)
	if err != nil {
		return AccountID{}, err
	}
	return AccountID(baseID), nil
}

func (id AccountID) Value() uuid.UUID {
	return ID(id).Value()
}

func (id AccountID) String() string {
	return ID(id).String()
}

func (id AccountID) Equals(other AccountID) bool {
	return ID(id).Equals(ID(other))
}

func (id AccountID) IsEmpty() bool {
	return ID(id).IsEmpty()
}

// TokenID specific constructors and methods
func NewTokenID() TokenID {
	return TokenID(NewID())
}

func ReconstructTokenID(id uuid.UUID) (TokenID, error) {
	baseID, err := ReconstructID(id)
	if err != nil {
		return TokenID{}, err
	}
	return TokenID(baseID), nil
}

func ParseTokenID(s string) (TokenID, error) {
	baseID, err := ParseID(s)
	if err != nil {
		return TokenID{}, err
	}
	return TokenID(baseID), nil
}

func (id TokenID) Value() uuid.UUID {
	return ID(id).Value()
}

func (id TokenID) String() string {
	return ID(id).String()
}

func (id TokenID) Equals(other TokenID) bool {
	return ID(id).Equals(ID(other))
}

func (id TokenID) IsEmpty() bool {
	return ID(id).IsEmpty()
}