package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

type SpotifyOAuthService struct {
	config *oauth2.Config
}

type SpotifyUserInfo struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Country     string `json:"country"`
	Product     string `json:"product"`
	Images      []struct {
		URL string `json:"url"`
	} `json:"images"`
}

func NewSpotifyOAuthService(clientID, clientSecret, redirectURL string) *SpotifyOAuthService {
	return &SpotifyOAuthService{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"user-read-email",
				"user-read-private",
				"playlist-read-private",
				"playlist-read-collaborative",
			},
			Endpoint: spotify.Endpoint,
		},
	}
}

func (s *SpotifyOAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *SpotifyOAuthService) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

func (s *SpotifyOAuthService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*SpotifyUserInfo, error) {
	client := s.config.Client(ctx, token)

	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("spotify API error (status %d): %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var userInfo SpotifyUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}

	return &userInfo, nil
}