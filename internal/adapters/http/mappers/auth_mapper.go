package mappers

import (
	"github.com/zandomed/sync-playlist-api/internal/adapters/http/dtos"
	authUC "github.com/zandomed/sync-playlist-api/internal/usecases/auth"
)

type AuthMapper struct{}

func NewAuthMapper() *AuthMapper {
	return &AuthMapper{}
}

func (m *AuthMapper) ToRegisterUserRequest(dto *dtos.RegisterRequest) *authUC.RegisterUserRequest {
	return &authUC.RegisterUserRequest{
		Email:    dto.Email,
		Name:     dto.Name,
		LastName: dto.LastName,
		Password: dto.Password,
	}
}

func (m *AuthMapper) ToRegisterResponse(ucResponse *authUC.RegisterUserResponse) *dtos.RegisterResponse {
	return &dtos.RegisterResponse{
		UserID:  ucResponse.UserID,
		Message: "User registered successfully",
	}
}

func (m *AuthMapper) ToLoginUserRequest(dto *dtos.LoginRequest) *authUC.LoginUserRequest {
	return &authUC.LoginUserRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func (m *AuthMapper) ToLoginResponse(ucResponse *authUC.LoginUserResponse) *dtos.LoginResponse {
	return &dtos.LoginResponse{
		AccessToken:  ucResponse.AccessToken,
		RefreshToken: ucResponse.RefreshToken,
		UserID:       ucResponse.UserID,
	}
}

func (m *AuthMapper) ToGoogleAuthURLRequest(dto *dtos.GoogleAuthURLRequest) *authUC.GetUrlGoogleRequest {
	return &authUC.GetUrlGoogleRequest{
		State: dto.State,
	}
}

func (m *AuthMapper) ToGoogleAuthURLResponse(ucResponse *authUC.GetUrlGoogleResponse) *dtos.GoogleAuthURLResponse {
	return &dtos.GoogleAuthURLResponse{
		URL: ucResponse.URL,
	}
}

func (m *AuthMapper) ToGoogleCallbackRequest(dto *dtos.GoogleCallbackRequest) *authUC.LoginGoogleCallbackRequest {
	return &authUC.LoginGoogleCallbackRequest{
		Code:  dto.Code,
		State: dto.State,
	}
}

func (m *AuthMapper) ToGoogleCallbackResponse(ucResponse *authUC.LoginGoogleCallbackResponse) *dtos.GoogleCallbackResponse {
	return &dtos.GoogleCallbackResponse{
		AccessToken:  ucResponse.AccessToken,
		RefreshToken: ucResponse.RefreshToken,
		UserID:       ucResponse.UserID,
		IsNewUser:    ucResponse.IsNewUser,
	}
}

func (m *AuthMapper) ToSpotifyAuthURLRequest(dto *dtos.SpotifyAuthURLRequest) *authUC.GetUrlSpotifyRequest {
	return &authUC.GetUrlSpotifyRequest{
		State: dto.State,
	}
}

func (m *AuthMapper) ToSpotifyAuthURLResponse(ucResponse *authUC.GetUrlSpotifyResponse) *dtos.SpotifyAuthURLResponse {
	return &dtos.SpotifyAuthURLResponse{
		URL: ucResponse.URL,
	}
}

func (m *AuthMapper) ToSpotifyCallbackRequest(dto *dtos.SpotifyCallbackRequest) *authUC.LoginSpotifyCallbackRequest {
	return &authUC.LoginSpotifyCallbackRequest{
		Code:  dto.Code,
		State: dto.State,
	}
}

func (m *AuthMapper) ToSpotifyCallbackResponse(ucResponse *authUC.LoginSpotifyCallbackResponse) *dtos.SpotifyCallbackResponse {
	return &dtos.SpotifyCallbackResponse{
		AccessToken:  ucResponse.AccessToken,
		RefreshToken: ucResponse.RefreshToken,
		UserID:       ucResponse.UserID,
		IsNewUser:    ucResponse.IsNewUser,
	}
}

func (m *AuthMapper) ToLinkSpotifyRequest(dto *dtos.LinkSpotifyRequest, userID string) *authUC.LinkSpotifyAccountRequest {
	return &authUC.LinkSpotifyAccountRequest{
		UserID: userID,
		Code:   dto.Code,
		State:  dto.State,
	}
}

func (m *AuthMapper) ToLinkSpotifyResponse(ucResponse *authUC.LinkSpotifyAccountResponse) *dtos.LinkSpotifyResponse {
	return &dtos.LinkSpotifyResponse{
		Success: ucResponse.Success,
		Message: ucResponse.Message,
	}
}
