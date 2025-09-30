package mappers

import (
	"github.com/zandomed/sync-playlist-api/internal/adapters/http/dtos"
	"github.com/zandomed/sync-playlist-api/internal/usecases/auth"
)

type AuthMapper struct{}

func NewAuthMapper() *AuthMapper {
	return &AuthMapper{}
}

func (m *AuthMapper) ToRegisterUserRequest(dto *dtos.RegisterRequest) *auth.RegisterUserRequest {
	return &auth.RegisterUserRequest{
		Email:    dto.Email,
		Name:     dto.Name,
		LastName: dto.LastName,
		Password: dto.Password,
	}
}

func (m *AuthMapper) ToRegisterResponse(ucResponse *auth.RegisterUserResponse) *dtos.RegisterResponse {
	return &dtos.RegisterResponse{
		UserID:  ucResponse.UserID,
		Message: "User registered successfully",
	}
}

func (m *AuthMapper) ToLoginUserRequest(dto *dtos.LoginRequest) *auth.LoginUserRequest {
	return &auth.LoginUserRequest{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func (m *AuthMapper) ToLoginResponse(ucResponse *auth.LoginUserResponse) *dtos.LoginResponse {
	return &dtos.LoginResponse{
		AccessToken:  ucResponse.AccessToken,
		RefreshToken: ucResponse.RefreshToken,
		UserID:       ucResponse.UserID,
	}
}