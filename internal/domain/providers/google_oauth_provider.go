package providers

import (
	"context"
	"time"
)

type GoogleOAuthProvider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (accessToken, refreshToken string, expiresAt time.Time, err error)
	GetUserInfo(ctx context.Context, accessToken string) (email, name, familyName string, err error)
}
