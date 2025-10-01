package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/domain/providers"
	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type LoginGoogleCallbackRequest struct {
	Code  string
	State string
}

type LoginGoogleCallbackResponse struct {
	AccessToken  string
	RefreshToken string
	UserID       string
	IsNewUser    bool
}

type LoginGoogleUseCase struct {
	userRepo      repositories.UserRepository
	accountRepo   repositories.AccountRepository
	tokenRepo     repositories.TokenRepository
	tokenGen      TokenGenerator
	googleService providers.GoogleOAuthProvider
}

func NewLoginGoogleUseCase(
	userRepo repositories.UserRepository,
	accountRepo repositories.AccountRepository,
	tokenRepo repositories.TokenRepository,
	tokenGen TokenGenerator,
	googleService providers.GoogleOAuthProvider,
) *LoginGoogleUseCase {
	return &LoginGoogleUseCase{
		userRepo:      userRepo,
		accountRepo:   accountRepo,
		tokenRepo:     tokenRepo,
		tokenGen:      tokenGen,
		googleService: googleService,
	}
}

func (uc *LoginGoogleUseCase) Execute(ctx context.Context, req LoginGoogleCallbackRequest) (*LoginGoogleCallbackResponse, error) {
	// Exchange code for tokens
	googleAccessToken, googleRefreshToken, expiresAt, err := uc.googleService.ExchangeCode(ctx, req.Code)
	if err != nil {
		return nil, errors.NewAuthenticationError("google_exchange_failed", fmt.Sprintf("Failed to exchange code: %v", err))
	}

	// Get user info from Google
	googleEmail, googleName, googleFamilyName, err := uc.googleService.GetUserInfo(ctx, googleAccessToken)
	if err != nil {
		return nil, errors.NewAuthenticationError("google_userinfo_failed", fmt.Sprintf("Failed to get user info: %v", err))
	}

	email, err := valueobjects.NewEmail(googleEmail)
	if err != nil {
		return nil, err
	}

	// Check if user exists
	user, err := uc.userRepo.FindByEmail(ctx, email)
	isNewUser := false

	if err != nil {
		// User doesn't exist, create new user
		user, err = entities.NewUser(googleEmail, googleName, googleFamilyName)
		if err != nil {
			return nil, err
		}

		// Verify email since Google provides verified emails
		user.VerifyEmail()

		if err := uc.userRepo.Save(ctx, user); err != nil {
			return nil, err
		}
		isNewUser = true
	}

	// Check if Google account exists for this user
	account, err := uc.accountRepo.FindByUserIDAndProvider(ctx, user.ID(), entities.GoogleProvider)
	if err != nil {
		// Create new Google account
		account, err = entities.NewOAuthAccount(user.ID(), entities.GoogleProvider)
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

	// Store Google tokens (optional - could be used for future API calls)
	_ = googleRefreshToken
	_ = expiresAt

	return &LoginGoogleCallbackResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		UserID:       user.ID().String(),
		IsNewUser:    isNewUser,
	}, nil
}

func (uc *LoginGoogleUseCase) getRefreshTokenExpirationTime() time.Time {
	return time.Now().Add(time.Duration(uc.tokenGen.GetRefreshTokenExpiration()) * time.Second)
}

func splitName(fullName string) (firstName, lastName string) {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}
