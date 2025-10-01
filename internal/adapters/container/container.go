package container

import (
	authAdapters "github.com/zandomed/sync-playlist-api/internal/adapters/auth"
	httpHandlers "github.com/zandomed/sync-playlist-api/internal/adapters/http/handlers"
	httpMappers "github.com/zandomed/sync-playlist-api/internal/adapters/http/mappers"
	repoAdapters "github.com/zandomed/sync-playlist-api/internal/adapters/repositories"
	"github.com/zandomed/sync-playlist-api/internal/config"
	"github.com/zandomed/sync-playlist-api/internal/domain/repositories"
	"github.com/zandomed/sync-playlist-api/internal/usecases"
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
	GoogleLoginUC  *authUC.GoogleLoginUseCase
	SpotifyLoginUC *authUC.SpotifyLoginUseCase
	LinkSpotifyUC  *authUC.LinkSpotifyAccountUseCase
	GetStatusUC    *healthUC.GetStatusUseCase

	// Adapters
	TokenGenerator authUC.TokenGenerator
	GoogleOAuth    authUC.GoogleOAuthProvider
	SpotifyOAuth   authUC.SpotifyOAuthProvider

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

	googleOAuthService := authAdapters.NewGoogleOAuthService(
		cfg.Google.ClientID,
		cfg.Google.ClientSecret,
		cfg.Google.RedirectURL,
	)
	googleOAuthAdapter := authAdapters.NewGoogleOAuthAdapter(googleOAuthService)

	spotifyOAuthService := authAdapters.NewSpotifyOAuthService(
		cfg.Spotify.ClientID,
		cfg.Spotify.ClientSecret,
		cfg.Spotify.RedirectURL,
	)
	spotifyOAuthAdapter := authAdapters.NewSpotifyOAuthAdapter(spotifyOAuthService)

	registerUserUC := authUC.NewRegisterUserUseCase(userRepo, accountRepo)
	loginUserUC := authUC.NewLoginUserUseCase(userRepo, accountRepo, tokenRepo, tokenGenerator)
	googleLoginUC := authUC.NewGoogleLoginUseCase(userRepo, accountRepo, tokenRepo, tokenGenerator, googleOAuthAdapter)
	spotifyLoginUC := authUC.NewSpotifyLoginUseCase(userRepo, accountRepo, tokenRepo, tokenGenerator, spotifyOAuthAdapter)
	linkSpotifyUC := authUC.NewLinkSpotifyAccountUseCase(userRepo, accountRepo, spotifyOAuthAdapter)
	getStatusUC := healthUC.NewGetStatusUseCase(cfg)

	authMapper := httpMappers.NewAuthMapper()

	authHandler := httpHandlers.NewAuthHandler(
		usecases.NewAuthUseCases(
			registerUserUC,
			loginUserUC,
			googleLoginUC,
			spotifyLoginUC,
			linkSpotifyUC,
		),
		authMapper,
		cfg,
		logger,
	)

	healthHandler := httpHandlers.NewHealthHandler(getStatusUC)

	return &Container{
		UserRepo:       userRepo,
		AccountRepo:    accountRepo,
		TokenRepo:      tokenRepo,
		RegisterUserUC: registerUserUC,
		LoginUserUC:    loginUserUC,
		GoogleLoginUC:  googleLoginUC,
		SpotifyLoginUC: spotifyLoginUC,
		LinkSpotifyUC:  linkSpotifyUC,
		GetStatusUC:    getStatusUC,
		TokenGenerator: tokenGenerator,
		GoogleOAuth:    googleOAuthAdapter,
		SpotifyOAuth:   spotifyOAuthAdapter,
		AuthMapper:     authMapper,
		AuthHandler:    authHandler,
		HealthHandler:  healthHandler,
	}
}
