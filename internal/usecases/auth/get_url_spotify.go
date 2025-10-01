package auth

import (
	"context"

	"github.com/zandomed/sync-playlist-api/internal/domain/providers"
)

type GetUrlSpotifyRequest struct {
	State string
}

type GetUrlSpotifyResponse struct {
	URL string
}

type GetUrlSpotifyUseCase struct {
	spotifyService providers.SpotifyOAuthProvider
}

func NewGetUrlSpotifyUseCase(spotifyService providers.SpotifyOAuthProvider) *GetUrlSpotifyUseCase {
	return &GetUrlSpotifyUseCase{
		spotifyService: spotifyService,
	}
}

func (uc *GetUrlSpotifyUseCase) Execute(ctx context.Context, req GetUrlSpotifyRequest) (*GetUrlSpotifyResponse, error) {
	url := uc.spotifyService.GetAuthURL(req.State)
	return &GetUrlSpotifyResponse{
		URL: url,
	}, nil
}
