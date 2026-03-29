package models

// @Description	User registration schema
type RegisterParams struct {
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
type LoginParams struct {
	// The email of the user
	Username string `json:"username" binding:"required"`

	// The password of the user
	Password string `json:"password" binding:"required"`
}

type LoginResponseParams struct {
	Token string             `json:"token"`
	User  UserResponseParams `json:"user"`
}
