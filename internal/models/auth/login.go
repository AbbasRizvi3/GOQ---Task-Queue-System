package auth

type LoginRequest struct {
	ID       string `json:"id" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}
