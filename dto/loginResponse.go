package dto

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"token,omitempty"`
}
