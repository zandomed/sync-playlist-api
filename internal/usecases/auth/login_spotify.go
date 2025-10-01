package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/domain/providers"
	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type LoginSpotifyCallbackRequest struct {
	Code  string
	State string
}

type LoginSpotifyCallbackResponse struct {
	AccessToken               string
	RefreshToken              string
	UserID                    string
	IsNewUser                 bool
	FrontendVerificationToken string
}

type LoginSpotifyUseCase struct {
	userRepo         repositories.UserRepository
	accountRepo      repositories.AccountRepository
	tokenRepo        repositories.TokenRepository
	verificationRepo repositories.VerificationRepository
	tokenGen         TokenGenerator
	spotifyService   providers.SpotifyOAuthProvider
	expirationState  time.Duration
}

func NewLoginSpotifyUseCase(
	userRepo repositories.UserRepository,
	accountRepo repositories.AccountRepository,
	tokenRepo repositories.TokenRepository,
	verificationRepo repositories.VerificationRepository,
	tokenGen TokenGenerator,
	spotifyService providers.SpotifyOAuthProvider,
	expirationState time.Duration,
) *LoginSpotifyUseCase {
	return &LoginSpotifyUseCase{
		userRepo:         userRepo,
		accountRepo:      accountRepo,
		tokenRepo:        tokenRepo,
		verificationRepo: verificationRepo,
		tokenGen:         tokenGen,
		spotifyService:   spotifyService,
		expirationState:  expirationState,
	}
}

func (uc *LoginSpotifyUseCase) Execute(ctx context.Context, req LoginSpotifyCallbackRequest) (*LoginSpotifyCallbackResponse, error) {

	// Validate the state parameter by looking it up in the verification tokens table
	// The state itself is the token for OAuth flow
	verificationToken, err := uc.verificationRepo.FindByToken(ctx, req.State)
	if err != nil {
		return nil, errors.NewAuthenticationError("invalid_state", "OAuth state parameter not found or invalid")
	}

	// Validate that the token is valid for OAuth
	if err := verificationToken.ValidateForOAuth(); err != nil {
		return nil, err
	}

	// Mark the verification token as used (one-time use)
	if err := verificationToken.MarkAsUsed(); err != nil {
		return nil, err
	}

	// Update the token in the database
	if err := uc.verificationRepo.Update(ctx, verificationToken); err != nil {
		return nil, err
	}

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

	// Create frontend verification token (10 minutes expiration)
	frontendToken, err := entities.NewFrontendVerificationToken(user.ID(), uc.expirationState)
	if err != nil {
		return nil, err
	}

	// Store the frontend verification token
	if err := uc.verificationRepo.Save(ctx, frontendToken); err != nil {
		return nil, err
	}

	return &LoginSpotifyCallbackResponse{
		AccessToken:               accessToken,
		RefreshToken:              refreshTokenStr,
		UserID:                    user.ID().String(),
		IsNewUser:                 isNewUser,
		FrontendVerificationToken: frontendToken.Token(),
	}, nil
}

func (uc *LoginSpotifyUseCase) getRefreshTokenExpirationTime() time.Time {
	return time.Now().Add(time.Duration(uc.tokenGen.GetRefreshTokenExpiration()) * time.Second)
}
