package repositories

import (
	"context"

	"github.com/zandomed/sync-playlist-api/internal/domain/entities"
	"github.com/zandomed/sync-playlist-api/internal/domain/valueobjects"
)

type UserRepository interface {
	Save(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id valueobjects.UserID) (*entities.User, error)
	FindByEmail(ctx context.Context, email valueobjects.Email) (*entities.User, error)
	Exists(ctx context.Context, email valueobjects.Email) (bool, error)
	Delete(ctx context.Context, id valueobjects.UserID) error
}