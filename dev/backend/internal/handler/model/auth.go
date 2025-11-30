package model

type SignupResponse struct {
	Token string `json:"token"`
	User  struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"user"`
}

type LoginRequest struct {
	UserUUID string `json:"user_uuid"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
