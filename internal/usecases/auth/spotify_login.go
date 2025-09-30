package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type SpotifyAuthURLRequest struct {
	State string
}

type SpotifyAuthURLResponse struct {
	URL string
}

type SpotifyCallbackRequest struct {
	Code  string
	State string
}

type SpotifyCallbackResponse struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	IsNewUser    bool
}

type SpotifyOAuthProvider interface {
	GetAuthURL(state string) string
	ExchangeCode(ctx context.Context, code string) (accessToken, refreshToken string, expiresAt time.Time, err error)
	GetUserInfo(ctx context.Context, accessToken string) (email, displayName string, err error)
}

type SpotifyLoginUseCase struct {
	userRepo       repositories.UserRepository
	accountRepo    repositories.AccountRepository
	tokenRepo      repositories.TokenRepository
	tokenGen       TokenGenerator
	spotifyService SpotifyOAuthProvider
}

func NewSpotifyLoginUseCase(
	userRepo repositories.UserRepository,
	accountRepo repositories.AccountRepository,
	tokenRepo repositories.TokenRepository,
	tokenGen TokenGenerator,
	spotifyService SpotifyOAuthProvider,
) *SpotifyLoginUseCase {
	return &SpotifyLoginUseCase{
		userRepo:       userRepo,
		accountRepo:    accountRepo,
		tokenRepo:      tokenRepo,
		tokenGen:       tokenGen,
		spotifyService: spotifyService,
	}
}

func (uc *SpotifyLoginUseCase) GetAuthURL(ctx context.Context, req SpotifyAuthURLRequest) (*SpotifyAuthURLResponse, error) {
	url := uc.spotifyService.GetAuthURL(req.State)
	return &SpotifyAuthURLResponse{
		URL: url,
	}, nil
}

func (uc *SpotifyLoginUseCase) HandleCallback(ctx context.Context, req SpotifyCallbackRequest) (*SpotifyCallbackResponse, error) {
	// Exchange code for tokens
	spotifyAccessToken, spotifyRefreshToken, expiresAt, err := uc.spotifyService.ExchangeCode(ctx, req.Code)
	if err != nil {
		return nil, errors.NewAuthenticationError("spotify_exchange_failed", fmt.Sprintf("Failed to exchange code: %v", err))
	}

	// Get user info from Spotify
	spotifyEmail, spotifyDisplayName, err := uc.spotifyService.GetUserInfo(ctx, spotifyAccessToken)
	if err != nil {
		return nil, errors.NewAuthenticationError("spotify_userinfo_failed", fmt.Sprintf("Failed to get user info: %v", err))
	}

	if spotifyEmail == "" {
		return nil, errors.NewAuthenticationError("spotify_no_email", "Spotify account does not have an email address")
	}

	email, err := valueobjects.NewEmail(spotifyEmail)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	user, err := uc.userRepo.FindByEmail(ctx, email)
	isNewUser := false

	if err != nil {
		// User doesn't exist, create new user
		firstName, lastName := splitName(spotifyDisplayName)
		user, err = entities.NewUser(spotifyEmail, firstName, lastName)
		if err != nil {
			return nil, err
		}

		// Verify email since Spotify provides verified emails
		user.VerifyEmail()

		if err := uc.userRepo.Save(ctx, user); err != nil {
			return nil, err
		}
		isNewUser = true
	}

	// Check if Spotify account exists for this user
	account, err := uc.accountRepo.FindByUserIDAndProvider(ctx, user.ID(), entities.SpotifyProvider)
	if err != nil {
		// Create new Spotify account
		account, err = entities.NewOAuthAccount(user.ID(), entities.SpotifyProvider)
		if err != nil {
			return nil, err
		}

		if err := uc.accountRepo.Save(ctx, account); err != nil {
			return nil, err
		}
	}

	// Check if user can authenticate
	if err := user.CanAuthenticate(); err != nil {
		return nil, err
	}

	// Generate JWT tokens
	accessToken, err := uc.tokenGen.GenerateAccessToken(user.ID().String(), user.Email().String())
	if err != nil {
		return nil, err
	}

	refreshTokenStr, err := uc.tokenGen.GenerateRefreshToken(user.ID().String())
	if err != nil {
		return nil, err
	}

	refreshToken, err := entities.NewRefreshToken(
		user.ID(),
		refreshTokenStr,
		uc.getRefreshTokenExpirationTime(),
	)
	if err != nil {
		return nil, err
	}

	if err := uc.tokenRepo.SaveRefreshToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	// Store Spotify tokens (optional - could be stored for future API calls)
	_ = spotifyRefreshToken
	_ = expiresAt

	return &SpotifyCallbackResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		UserID:       user.ID().String(),
		IsNewUser:    isNewUser,
	}, nil
}

func (uc *SpotifyLoginUseCase) getRefreshTokenExpirationTime() time.Time {
	return time.Now().Add(time.Duration(uc.tokenGen.GetRefreshTokenExpiration()) * time.Second)
}