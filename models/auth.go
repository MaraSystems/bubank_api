package models

// @Description	User registration schema
type RegisterRequest struct {
	// The unique identity of the user
	Username string `json:"username" binding:"required,alphanum"`

	// The password of the user
	Password string `json:"password" binding:"required,min=6"`

	// The full name of the user
	FullName string `json:"full_name" binding:"required"`

	// The email of the user
	Email string `json:"email" binding:"required,email"`
}

// @Description User login schema
type LoginRequest struct {
	// The email of the user
	Username string `json:"username" binding:"required"`

	// The password of the user
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
