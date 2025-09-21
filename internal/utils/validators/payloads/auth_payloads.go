package payloads

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"pass_word"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"pass_word"`
	Fullname string `json:"full_name"`
}
