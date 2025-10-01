package container

import (
	"github.com/zandomed/sync-playlist-api/internal/config"
	authAdapters "github.com/zandomed/sync-playlist-api/internal/infra/adapters/auth"
	httpHandlers "github.com/zandomed/sync-playlist-api/internal/infra/http/handlers"
	httpMappers "github.com/zandomed/sync-playlist-api/internal/infra/http/mappers"
	repoAdapters "github.com/zandomed/sync-playlist-api/internal/infra/repositories"
	services "github.com/zandomed/sync-playlist-api/internal/infra/services/auth"
	"github.com/zandomed/sync-playlist-api/internal/usecases"
	authUC "github.com/zandomed/sync-playlist-api/internal/usecases/auth"
	healthUC "github.com/zandomed/sync-playlist-api/internal/usecases/health"
	"github.com/zandomed/sync-playlist-api/pkg/database"
	"github.com/zandomed/sync-playlist-api/pkg/logger"
)

type Container struct {
	// Handlers
	AuthHandler   *httpHandlers.AuthHandler
	HealthHandler *httpHandlers.HealthHandler
}

func NewContainer(db *database.DB, cfg *config.Config, logger *logger.Logger) *Container {
	userRepo := repoAdapters.NewPostgresUserRepository(db)
	accountRepo := repoAdapters.NewPostgresAccountRepository(db)
	tokenRepo := repoAdapters.NewPostgresTokenRepository(db)
	verificationRepo := repoAdapters.NewPostgresVerificationRepository(db)

	tokenGenerator := authAdapters.NewJWTTokenGenerator(
		cfg.JWT.Secret,
		cfg.JWT.ExpirationTime,
		cfg.JWT.RefreshExpirationTime,
	)

	googleOAuthService := services.NewGoogleOAuthService(
		cfg.Google.ClientID,
		cfg.Google.ClientSecret,
		cfg.Google.RedirectURL,
	)
	googleOAuthAdapter := authAdapters.NewGoogleOAuthAdapter(googleOAuthService)

	spotifyOAuthService := services.NewSpotifyOAuthService(
		cfg.Spotify.ClientID,
		cfg.Spotify.ClientSecret,
		cfg.Spotify.RedirectURL,
		cfg.Spotify.APIUrl,
	)
	spotifyOAuthAdapter := authAdapters.NewSpotifyOAuthAdapter(spotifyOAuthService)

	// State Expiration
	expirationTimeForOAuthState := cfg.OAuth.TokenExpiration
	expirationTimeForFrontendOAuth := cfg.OAuth.FrontendTokenExpiration

	registerUserUC := authUC.NewRegisterUserUseCase(userRepo, accountRepo)
	loginUserUC := authUC.NewLoginUserUseCase(userRepo, accountRepo, tokenRepo, tokenGenerator)
	googleLoginUC := authUC.NewLoginGoogleUseCase(userRepo, accountRepo, tokenRepo, verificationRepo, tokenGenerator, googleOAuthAdapter, expirationTimeForFrontendOAuth)
	spotifyLoginUC := authUC.NewLoginSpotifyUseCase(userRepo, accountRepo, tokenRepo, verificationRepo, tokenGenerator, spotifyOAuthAdapter, expirationTimeForFrontendOAuth)
	linkSpotifyUC := authUC.NewLinkSpotifyAccountUseCase(userRepo, accountRepo, spotifyOAuthAdapter)
	getUrlSpotifyUC := authUC.NewGetUrlSpotifyUseCase(spotifyOAuthAdapter, verificationRepo, expirationTimeForOAuthState)
	getUrlGoogleUC := authUC.NewGetUrlGoogleUseCase(googleOAuthAdapter, verificationRepo, expirationTimeForOAuthState)
	verifyTokenUC := authUC.NewVerifyTokenUseCase(verificationRepo)
	getStatusUC := healthUC.NewGetStatusUseCase(cfg)

	authMapper := httpMappers.NewAuthMapper()

	authHandler := httpHandlers.NewAuthHandler(
		usecases.NewAuthUseCases(
			registerUserUC,
			loginUserUC,
			googleLoginUC,
			spotifyLoginUC,
			linkSpotifyUC,
			getUrlSpotifyUC,
			getUrlGoogleUC,
			verifyTokenUC,
		),
		authMapper,
		cfg,
		logger,
	)

	healthHandler := httpHandlers.NewHealthHandler(getStatusUC)

	return &Container{
		AuthHandler:   authHandler,
		HealthHandler: healthHandler,
	}
}
