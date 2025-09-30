package repositories

import (
	"context"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type AccountRepository interface {
	Save(ctx context.Context, account *entities.Account) error
	FindByUserIDAndProvider(ctx context.Context, userID valueobjects.UserID, provider entities.AccountProvider) (*entities.Account, error)
	FindUserpassAccountByEmail(ctx context.Context, email valueobjects.Email) (*entities.Account, error)
	Delete(ctx context.Context, id valueobjects.AccountID) error
}