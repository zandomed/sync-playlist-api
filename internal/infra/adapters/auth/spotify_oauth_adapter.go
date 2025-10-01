package auth

import (
	"context"
	"strings"
	"time"

	"github.com/zandomed/sync-playlist-api/internal/domain/providers"
	services "github.com/zandomed/sync-playlist-api/internal/infra/services/auth"
	"golang.org/x/oauth2"
)

type SpotifyOAuthAdapter struct {
	service *services.SpotifyOAuthService
}

func NewSpotifyOAuthAdapter(service *services.SpotifyOAuthService) providers.SpotifyOAuthProvider {
	return &SpotifyOAuthAdapter{
		service: service,
	}
}

func (a *SpotifyOAuthAdapter) GetAuthURL(state string) string {
	return a.service.GetAuthURL(state)
}

func (a *SpotifyOAuthAdapter) ExchangeCode(ctx context.Context, code string) (accessToken, refreshToken string, expiresAt time.Time, err error) {
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

func (a *SpotifyOAuthAdapter) GetUserInfo(ctx context.Context, accessToken string) (email, displayName string, err error) {
	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	userInfo, err := a.service.GetUserInfo(ctx, token)
	if err != nil {
		return "", "", err
	}

	// Spotify display name might be empty or contain full name
	// Extract first and last name from display name
	name := userInfo.DisplayName
	if name == "" {
		name = strings.Split(userInfo.Email, "@")[0]
	}

	return userInfo.Email, name, nil
}
