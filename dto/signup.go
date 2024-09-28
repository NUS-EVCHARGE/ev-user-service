package dto

type SignUpResendRequest struct {
	Email string `json:"email"`
}

type SignUpRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type SignUpResponse struct {
	Status string `json:"status"`
	Message string `json:"message,omitempty"`
}

type ConfirmUser struct {
	Email            string `json:"email"`
	ConfirmationCode string `json:"confirmationCode"`
}
