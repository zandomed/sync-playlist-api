package auth

import (
	"context"

	"github.com/zandomed/sync-playlist-api/internal/domain/providers"
)

type GetUrlGoogleRequest struct {
	State string
}

type GetUrlGoogleResponse struct {
	URL string
}

type GetUrlGoogleUseCase struct {
	googleService providers.GoogleOAuthProvider
}

func NewGetUrlGoogleUseCase(googleService providers.GoogleOAuthProvider) *GetUrlGoogleUseCase {
	return &GetUrlGoogleUseCase{
		googleService: googleService,
	}
}

func (uc *GetUrlGoogleUseCase) Execute(ctx context.Context, req GetUrlGoogleRequest) (*GetUrlGoogleResponse, error) {
	url := uc.googleService.GetAuthURL(req.State)
	return &GetUrlGoogleResponse{
		URL: url,
	}, nil
}
