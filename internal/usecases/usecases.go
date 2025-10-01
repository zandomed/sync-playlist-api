package usecases

import authUC "github.com/zandomed/sync-playlist-api/internal/usecases/auth"

type AuthUseCases struct {
	RegisterUserUseCase *authUC.RegisterUserUseCase
	LoginUserUseCase    *authUC.LoginUserUseCase
	GoogleLoginUseCase  *authUC.GoogleLoginUseCase
	SpotifyLoginUseCase *authUC.SpotifyLoginUseCase
	LinkSpotifyUseCase  *authUC.LinkSpotifyAccountUseCase
}

func NewAuthUseCases(
	registerUserUC *authUC.RegisterUserUseCase,
	loginUserUC *authUC.LoginUserUseCase,
	googleLoginUC *authUC.GoogleLoginUseCase,
	spotifyLoginUC *authUC.SpotifyLoginUseCase,
	linkSpotifyUC *authUC.LinkSpotifyAccountUseCase,
) *AuthUseCases {
	return &AuthUseCases{
		RegisterUserUseCase: registerUserUC,
		LoginUserUseCase:    loginUserUC,
		GoogleLoginUseCase:  googleLoginUC,
		SpotifyLoginUseCase: spotifyLoginUC,
		LinkSpotifyUseCase:  linkSpotifyUC,
	}
}
