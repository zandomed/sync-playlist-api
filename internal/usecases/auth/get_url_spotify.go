package auth

import (
	"context"
	"time"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/providers"
	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
)

type GetUrlSpotifyResponse struct {
	URL   string
	State string
}

type GetUrlSpotifyUseCase struct {
	spotifyService   providers.SpotifyOAuthProvider
	verificationRepo repositories.VerificationRepository
	expirationState  time.Duration
}

func NewGetUrlSpotifyUseCase(
	spotifyService providers.SpotifyOAuthProvider,
	verificationRepo repositories.VerificationRepository,
	expirationState time.Duration) *GetUrlSpotifyUseCase {
	return &GetUrlSpotifyUseCase{
		spotifyService:   spotifyService,
		verificationRepo: verificationRepo,
		expirationState:  expirationState,
	}
}

func (uc *GetUrlSpotifyUseCase) Execute(ctx context.Context) (*GetUrlSpotifyResponse, error) {
	// Create a verification token for OAuth state validation (5 minutes expiration)
	// The token itself is the state parameter
	verificationToken, err := entities.NewOAuthStateToken(uc.expirationState)
	if err != nil {
		return nil, err
	}

	// Store the verification token
	if err := uc.verificationRepo.Save(ctx, verificationToken); err != nil {
		return nil, err
	}

	// Use the token as the state parameter
	state := verificationToken.Token()
	url := uc.spotifyService.GetAuthURL(state)

	return &GetUrlSpotifyResponse{
		URL: url,
	}, nil
}
