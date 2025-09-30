package container

import (
	authAdapters "github.com/zandomed/sync-playlist-api/internal/adapters/auth"
	httpHandlers "github.com/zandomed/sync-playlist-api/internal/adapters/http/handlers"
	httpMappers "github.com/zandomed/sync-playlist-api/internal/adapters/http/mappers"
	repoAdapters "github.com/zandomed/sync-playlist-api/internal/adapters/repositories"
	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
	authUC "github.com/zandomed/sync-playlist-api/internal/usecases/auth"
	healthUC "github.com/zandomed/sync-playlist-api/internal/usecases/health"
	"github.com/zandomed/sync-playlist-api/pkg/database"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

type Container struct {
	// Repositories
	UserRepo    repositories.UserRepository
	AccountRepo repositories.AccountRepository
	TokenRepo   repositories.TokenRepository

	// Use Cases
	RegisterUserUC *authUC.RegisterUserUseCase
	LoginUserUC    *authUC.LoginUserUseCase
	GetStatusUC    *healthUC.GetStatusUseCase

	// Adapters
	TokenGenerator authUC.TokenGenerator

	// Mappers
	AuthMapper *httpMappers.AuthMapper

	// Handlers
	AuthHandler   *httpHandlers.AuthHandler
	HealthHandler *httpHandlers.HealthHandler
}

func NewContainer(db *database.DB, cfg *config.Config, logger *logger.Logger) *Container {
	userRepo := repoAdapters.NewPostgresUserRepository(db)
	accountRepo := repoAdapters.NewPostgresAccountRepository(db)
	tokenRepo := repoAdapters.NewPostgresTokenRepository(db)

	tokenGenerator := authAdapters.NewJWTTokenGenerator(
		cfg.JWT.Secret,
		cfg.JWT.ExpirationTime,
		cfg.JWT.RefreshExpirationTime,
	)

	registerUserUC := authUC.NewRegisterUserUseCase(userRepo, accountRepo)
	loginUserUC := authUC.NewLoginUserUseCase(userRepo, accountRepo, tokenRepo, tokenGenerator)
	getStatusUC := healthUC.NewGetStatusUseCase(cfg)

	authMapper := httpMappers.NewAuthMapper()

	authHandler := httpHandlers.NewAuthHandler(
		registerUserUC,
		loginUserUC,
		authMapper,
		logger,
	)

	healthHandler := httpHandlers.NewHealthHandler(getStatusUC)

	return &Container{
		UserRepo:       userRepo,
		AccountRepo:    accountRepo,
		TokenRepo:      tokenRepo,
		RegisterUserUC: registerUserUC,
		LoginUserUC:    loginUserUC,
		GetStatusUC:    getStatusUC,
		TokenGenerator: tokenGenerator,
		AuthMapper:     authMapper,
		AuthHandler:    authHandler,
		HealthHandler:  healthHandler,
	}
}
