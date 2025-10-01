package auth

import (
	"context"
	"fmt"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/domain/providers"
	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type LinkSpotifyAccountRequest struct {
	UserID string
	Code   string
	State  string
}

type LinkSpotifyAccountResponse struct {
	Success bool
	Message string
}

type LinkSpotifyAccountUseCase struct {
	userRepo       repositories.UserRepository
	accountRepo    repositories.AccountRepository
	spotifyService providers.SpotifyOAuthProvider
}

func NewLinkSpotifyAccountUseCase(
	userRepo repositories.UserRepository,
	accountRepo repositories.AccountRepository,
	spotifyService providers.SpotifyOAuthProvider,
) *LinkSpotifyAccountUseCase {
	return &LinkSpotifyAccountUseCase{
		userRepo:       userRepo,
		accountRepo:    accountRepo,
		spotifyService: spotifyService,
	}
}

func (uc *LinkSpotifyAccountUseCase) GetAuthURL(ctx context.Context, req GetUrlSpotifyRequest) (*GetUrlSpotifyResponse, error) {
	url := uc.spotifyService.GetAuthURL(req.State)
	return &GetUrlSpotifyResponse{
		URL: url,
	}, nil
}

func (uc *LinkSpotifyAccountUseCase) Execute(ctx context.Context, req LinkSpotifyAccountRequest) (*LinkSpotifyAccountResponse, error) {
	// Validate user ID
	userID, err := valueobjects.ParseUserID(req.UserID)
	if err != nil {
		return nil, errors.NewDomainError("invalid_user_id", "Invalid user ID")
	}

	// Check if user exists
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.NewDomainError("user_not_found", "User not found")
	}

	// Exchange code for tokens
	spotifyAccessToken, spotifyRefreshToken, expiresAt, err := uc.spotifyService.ExchangeCode(ctx, req.Code)
	if err != nil {
		return nil, errors.NewAuthenticationError("spotify_exchange_failed", fmt.Sprintf("Failed to exchange code: %v", err))
	}

	// Get user info from Spotify to validate
	spotifyEmail, _, err := uc.spotifyService.GetUserInfo(ctx, spotifyAccessToken)
	if err != nil {
		return nil, errors.NewAuthenticationError("spotify_userinfo_failed", fmt.Sprintf("Failed to get user info: %v", err))
	}

	if spotifyEmail == "" {
		return nil, errors.NewAuthenticationError("spotify_no_email", "Spotify account does not have an email address")
	}

	// Check if Spotify account already exists for this user
	existingAccount, err := uc.accountRepo.FindByUserIDAndProvider(ctx, user.ID(), entities.SpotifyProvider)
	if err == nil && existingAccount != nil {
		return nil, errors.NewDomainError("account_already_linked", "Spotify account is already linked to this user")
	}

	// Create new Spotify account
	account, err := entities.NewOAuthAccount(user.ID(), entities.SpotifyProvider)
	if err != nil {
		return nil, err
	}

	if err := uc.accountRepo.Save(ctx, account); err != nil {
		return nil, err
	}

	// TODO: Store Spotify tokens (optional - for future use)
	_, _, _ = spotifyAccessToken, spotifyRefreshToken, expiresAt

	return &LinkSpotifyAccountResponse{
		Success: true,
		Message: "Spotify account linked successfully",
	}, nil
}
