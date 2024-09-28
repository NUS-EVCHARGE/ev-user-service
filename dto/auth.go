package dto

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status       string `json:"status"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	IdToken      string `json:"idToken"`
	ExpiresIn    int    `json:"expiresIn"`
}