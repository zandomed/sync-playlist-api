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

type GoogleAuthURLResponse struct {
	URL   string `json:"url"`
	State string `json:"state"` // Server-generated state for OAuth flow
}

type GoogleCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"` // OAuth state parameter
}

type GoogleCallbackResponse struct {
	AccessToken               string `json:"accessToken"`
	RefreshToken              string `json:"refreshToken"`
	UserID                    string `json:"userID"`
	IsNewUser                 bool   `json:"isNewUser"`
	FrontendVerificationToken string `json:"frontendVerificationToken"`
}

type VerifyTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

type VerifyTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"userId,omitempty"`
}

type SpotifyAuthURLResponse struct {
	URL   string `json:"url"`
	State string `json:"state"` // Server-generated state for OAuth flow
}

type SpotifyCallbackRequest struct {
	Code  string `json:"code" validate:"required"`
	State string `json:"state" validate:"required"`
}

type SpotifyCallbackResponse struct {
	AccessToken               string `json:"accessToken"`
	RefreshToken              string `json:"refreshToken"`
	UserID                    string `json:"userID"`
	IsNewUser                 bool   `json:"isNewUser"`
	FrontendVerificationToken string `json:"frontendVerificationToken"`
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
