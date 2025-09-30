package dtos

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=50"`
	LastName string `json:"lastName" validate:"required,min=2,max=50"`
	Password string `json:"password" validate:"required,min=8,max=128"`
}

type RegisterResponse struct {
	UserID  string `json:"userID"`
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userID"`
}

type GoogleAuthURLRequest struct {
	State string `json:"state"`
}

type GoogleAuthURLResponse struct {
	URL string `json:"url"`
}

type GoogleCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state"`
}

type GoogleCallbackResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userID"`
	IsNewUser    bool   `json:"isNewUser"`
}

type SpotifyAuthURLRequest struct {
	State string `json:"state"`
}

type SpotifyAuthURLResponse struct {
	URL string `json:"url"`
}

type SpotifyCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state"`
}

type SpotifyCallbackResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UserID       string `json:"userID"`
	IsNewUser    bool   `json:"isNewUser"`
}

type LinkSpotifyRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state"`
}

type LinkSpotifyResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}
