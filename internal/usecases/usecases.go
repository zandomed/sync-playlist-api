package usecases

import authUC "github.com/zandomed/sync-playlist-api/internal/usecases/auth"

type AuthUseCases struct {
	RegisterUserUseCase  *authUC.RegisterUserUseCase
	LoginUserPassUseCase *authUC.LoginUserUseCase
	LoginGoogleUseCase   *authUC.LoginGoogleUseCase
	LoginSpotifyUseCase  *authUC.LoginSpotifyUseCase
	LinkSpotifyUseCase   *authUC.LinkSpotifyAccountUseCase
	GetUrlSpotifyUseCase *authUC.GetUrlSpotifyUseCase
	GetUrlGoogleUseCase  *authUC.GetUrlGoogleUseCase
}

func NewAuthUseCases(
	registerUserUC *authUC.RegisterUserUseCase,
	loginUserUC *authUC.LoginUserUseCase,
	loginGoogleUC *authUC.LoginGoogleUseCase,
	loginSpotifyUC *authUC.LoginSpotifyUseCase,
	linkSpotifyUC *authUC.LinkSpotifyAccountUseCase,
	getUrlSpotifyUC *authUC.GetUrlSpotifyUseCase,
	getUrlGoogleUC *authUC.GetUrlGoogleUseCase,
) *AuthUseCases {
	return &AuthUseCases{
		RegisterUserUseCase:  registerUserUC,
		LoginUserPassUseCase: loginUserUC,
		LoginGoogleUseCase:   loginGoogleUC,
		LoginSpotifyUseCase:  loginSpotifyUC,
		LinkSpotifyUseCase:   linkSpotifyUC,
		GetUrlSpotifyUseCase: getUrlSpotifyUC,
		GetUrlGoogleUseCase:  getUrlGoogleUC,
	}
}
