package auth

import (
	"context"
	"time"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/providers"
	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
)

type GetUrlGoogleResponse struct {
	URL   string
	State string // The generated state to be used in OAuth flow
}

type GetUrlGoogleUseCase struct {
	googleService    providers.GoogleOAuthProvider
	verificationRepo repositories.VerificationRepository
	expirationState  time.Duration
}

func NewGetUrlGoogleUseCase(
	googleService providers.GoogleOAuthProvider,
	verificationRepo repositories.VerificationRepository,
	expirationState time.Duration,
) *GetUrlGoogleUseCase {
	return &GetUrlGoogleUseCase{
		googleService:    googleService,
		verificationRepo: verificationRepo,
		expirationState:  expirationState,
	}
}

func (uc *GetUrlGoogleUseCase) Execute(ctx context.Context) (*GetUrlGoogleResponse, error) {
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
	url := uc.googleService.GetAuthURL(state)

	return &GetUrlGoogleResponse{
		URL:   url,
		State: state,
	}, nil
}
