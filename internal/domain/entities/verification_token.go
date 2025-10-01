package entities

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type VerificationTokenType string

const (
	// OAuth state verification - short lived (5 minutes)
	OAuthStateToken VerificationTokenType = "oauth_state"
	// Frontend verification - medium lived (10 minutes)
	FrontendVerificationToken VerificationTokenType = "frontend_verification"
)

type VerificationToken struct {
	id        valueobjects.TokenID
	token     string // For OAuth: this is the state; For Frontend: this is a generated token
	tokenType VerificationTokenType
	userID    *valueobjects.UserID // Only set for frontend verification tokens
	expiresAt time.Time
	createdAt time.Time
	usedAt    *time.Time
}

// NewOAuthStateToken creates a verification token for OAuth state validation
// The state parameter itself serves as the token
func NewOAuthStateToken(expiration time.Duration) (*VerificationToken, error) {
	// Generate a secure state token
	state, err := generateSecureToken()
	if err != nil {
		return nil, errors.NewDomainError("token_generation_failed", "Failed to generate state token")
	}

	now := time.Now()
	expiresAt := now.Add(expiration)

	return &VerificationToken{
		id:        valueobjects.NewTokenID(),
		token:     state, // The state IS the token
		tokenType: OAuthStateToken,
		userID:    nil,
		expiresAt: expiresAt,
		createdAt: now,
		usedAt:    nil,
	}, nil
}

// NewFrontendVerificationToken creates a verification token for frontend validation
func NewFrontendVerificationToken(userID valueobjects.UserID, expiration time.Duration) (*VerificationToken, error) {
	token, err := generateSecureToken()
	if err != nil {
		return nil, errors.NewDomainError("token_generation_failed", "Failed to generate verification token")
	}

	now := time.Now()
	expiresAt := now.Add(expiration)

	return &VerificationToken{
		id:        valueobjects.NewTokenID(),
		token:     token,
		tokenType: FrontendVerificationToken,
		userID:    &userID,
		expiresAt: expiresAt,
		createdAt: now,
		usedAt:    nil,
	}, nil
}

// ReconstructVerificationToken reconstructs a verification token from persistence
func ReconstructVerificationToken(
	token string,
	tokenType VerificationTokenType,
	userID *valueobjects.UserID,
	expiresAt time.Time,
	createdAt time.Time,
	usedAt *time.Time,
) (*VerificationToken, error) {
	if token == "" {
		return nil, errors.NewDomainError("empty_token", "Token cannot be empty")
	}

	return &VerificationToken{
		id:        valueobjects.NewTokenID(),
		token:     token,
		tokenType: tokenType,
		userID:    userID,
		expiresAt: expiresAt,
		createdAt: createdAt,
		usedAt:    usedAt,
	}, nil
}

func (vt *VerificationToken) ID() valueobjects.TokenID {
	return vt.id
}

func (vt *VerificationToken) Token() string {
	return vt.token
}

func (vt *VerificationToken) TokenType() VerificationTokenType {
	return vt.tokenType
}

// State returns the token value for OAuth state tokens
// For OAuth tokens, the token itself is the state
func (vt *VerificationToken) State() string {
	if vt.tokenType == OAuthStateToken {
		return vt.token
	}
	return ""
}

func (vt *VerificationToken) UserID() *valueobjects.UserID {
	return vt.userID
}

func (vt *VerificationToken) ExpiresAt() time.Time {
	return vt.expiresAt
}

func (vt *VerificationToken) CreatedAt() time.Time {
	return vt.createdAt
}

func (vt *VerificationToken) UsedAt() *time.Time {
	return vt.usedAt
}

func (vt *VerificationToken) IsExpired() bool {
	return time.Now().After(vt.expiresAt)
}

func (vt *VerificationToken) IsUsed() bool {
	return vt.usedAt != nil
}

func (vt *VerificationToken) IsValid() bool {
	return !vt.IsExpired() && !vt.IsUsed() && vt.token != ""
}

// MarkAsUsed marks the token as used (for one-time use tokens)
func (vt *VerificationToken) MarkAsUsed() error {
	if vt.IsUsed() {
		return errors.NewDomainError("token_already_used", "Token has already been used")
	}

	if vt.IsExpired() {
		return errors.NewDomainError("token_expired", "Token has expired")
	}

	now := time.Now()
	vt.usedAt = &now
	return nil
}

// ValidateForOAuth validates the token for OAuth flow
// The state parameter should match the token itself
func (vt *VerificationToken) ValidateForOAuth() error {
	if vt.tokenType != OAuthStateToken {
		return errors.NewAuthenticationError("invalid_token_type", "Token is not an OAuth state token")
	}

	if !vt.IsValid() {
		if vt.IsExpired() {
			return errors.NewAuthenticationError("token_expired", "Verification token has expired")
		}
		if vt.IsUsed() {
			return errors.NewAuthenticationError("token_used", "Verification token has already been used")
		}
		return errors.NewAuthenticationError("invalid_token", "Invalid verification token")
	}

	return nil
}

// ValidateForFrontend validates the token for frontend verification
func (vt *VerificationToken) ValidateForFrontend() error {
	if vt.tokenType != FrontendVerificationToken {
		return errors.NewAuthenticationError("invalid_token_type", "Token is not a frontend verification token")
	}

	if !vt.IsValid() {
		if vt.IsExpired() {
			return errors.NewAuthenticationError("token_expired", "Verification token has expired")
		}
		if vt.IsUsed() {
			return errors.NewAuthenticationError("token_used", "Verification token has already been used")
		}
		return errors.NewAuthenticationError("invalid_token", "Invalid verification token")
	}

	if vt.userID == nil {
		return errors.NewAuthenticationError("missing_user_id", "Frontend verification token must have a user ID")
	}

	return nil
}

// generateSecureToken generates a cryptographically secure random token
// for OAuth state parameters and verification tokens.
// Uses 256 bits of entropy to prevent CSRF and replay attacks.
func generateSecureToken() (string, error) {
	b := make([]byte, 32) // 256 bits of entropy
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	// URL-safe encoding (no padding needed for OAuth state)
	return base64.URLEncoding.EncodeToString(b), nil
}
