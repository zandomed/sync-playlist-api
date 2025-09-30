package auth

import (
	"context"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/errors"
	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type RegisterUserRequest struct {
	Email    string
	Name     string
	LastName string
	Password string
}

type RegisterUserResponse struct {
	UserID string
}

type RegisterUserUseCase struct {
	userRepo    repositories.UserRepository
	accountRepo repositories.AccountRepository
}

func NewRegisterUserUseCase(userRepo repositories.UserRepository, accountRepo repositories.AccountRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepo:    userRepo,
		accountRepo: accountRepo,
	}
}

func (uc *RegisterUserUseCase) Execute(ctx context.Context, req RegisterUserRequest) (*RegisterUserResponse, error) {
	email, err := valueobjects.NewEmail(req.Email)
	if err != nil {
		return nil, err
	}

	exists, err := uc.userRepo.Exists(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.NewDomainError("user_already_exists", "User with this email already exists")
	}

	user, err := entities.NewUser(req.Email, req.Name, req.LastName)
	if err != nil {
		return nil, err
	}

	plainPassword, err := valueobjects.NewPlainPassword(req.Password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := plainPassword.Hash()
	if err != nil {
		return nil, err
	}

	account := entities.NewUserpassAccount(user.ID(), hashedPassword)

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	if err := uc.accountRepo.Save(ctx, account); err != nil {
		return nil, err
	}

	return &RegisterUserResponse{
		UserID: user.ID().String(),
	}, nil
}