package model

type SignupResponse struct {
	Token string `json:"token"`
	User  struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
}

type LoginRequest struct {
	UserID string `json:"user_id"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
