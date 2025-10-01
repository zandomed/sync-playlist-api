package auth

import (
	"context"
	"time"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type LoginUserRequest struct {
	Email    string
	Password string
}

type LoginUserResponse struct {
	AccessToken  string
	RefreshToken string
	UserID       string
}

type TokenGenerator interface {
	GenerateAccessToken(userID string, email string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	GetRefreshTokenExpiration() int64
}

type LoginUserUseCase struct {
	userRepo    repositories.UserRepository
	accountRepo repositories.AccountRepository
	tokenRepo   repositories.TokenRepository
	tokenGen    TokenGenerator
}

func NewLoginUserUseCase(
	userRepo repositories.UserRepository,
	accountRepo repositories.AccountRepository,
	tokenRepo repositories.TokenRepository,
	tokenGen TokenGenerator,
) *LoginUserUseCase {
	return &LoginUserUseCase{
		userRepo:    userRepo,
		accountRepo: accountRepo,
		tokenRepo:   tokenRepo,
		tokenGen:    tokenGen,
	}
}

func (uc *LoginUserUseCase) Execute(ctx context.Context, req LoginUserRequest) (*LoginUserResponse, error) {
	email, err := valueobjects.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	plainPassword, err := valueobjects.NewPlainPassword(req.Password)
	if err != nil {
		return nil, err
	}

	account, err := uc.accountRepo.FindUserpassAccountByEmail(ctx, email)
	if err != nil {
		return nil, errors.NewAuthenticationError("invalid_credentials", "Invalid email or password")
	}

	if !account.Password().Verify(plainPassword) {
		return nil, errors.NewAuthenticationError("invalid_credentials", "Invalid email or password")
	}

	user, err := uc.userRepo.FindByID(ctx, account.UserID())
	if err != nil {
		return nil, errors.NewAuthenticationError("user_not_found", "User not found")
	}

	if err := user.CanAuthenticate(); err != nil {
		return nil, err
	}

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

	return &LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		UserID:       user.ID().String(),
	}, nil
}

func (uc *LoginUserUseCase) getRefreshTokenExpirationTime() time.Time {
	return time.Now().Add(time.Duration(uc.tokenGen.GetRefreshTokenExpiration()) * time.Second)
}
