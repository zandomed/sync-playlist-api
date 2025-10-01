package auth

import (
	"context"
	"time"

	"github.com/zandomed/sync-playlist-api/internal/domain/providers"
	"golang.org/x/oauth2"
)

type GoogleOAuthAdapter struct {
	service *GoogleOAuthService
}

func NewGoogleOAuthAdapter(service *GoogleOAuthService) providers.GoogleOAuthProvider {
	return &GoogleOAuthAdapter{
		service: service,
	}
}

func (a *GoogleOAuthAdapter) GetAuthURL(state string) string {
	return a.service.GetAuthURL(state)
}

func (a *GoogleOAuthAdapter) ExchangeCode(ctx context.Context, code string) (accessToken, refreshToken string, expiresAt time.Time, err error) {
	token, err := a.service.ExchangeCode(ctx, code)
	if err != nil {
		return "", "", time.Time{}, err
	}

	expiry := token.Expiry
	if expiry.IsZero() {
		expiry = time.Now().Add(time.Hour)
	}

	return token.AccessToken, token.RefreshToken, expiry, nil
}

func (a *GoogleOAuthAdapter) GetUserInfo(ctx context.Context, accessToken string) (email, name, familyName string, err error) {
	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	userInfo, err := a.service.GetUserInfo(ctx, token)
	if err != nil {
		return "", "", "", err
	}

	return userInfo.Email, userInfo.GivenName, userInfo.FamilyName, nil
}
