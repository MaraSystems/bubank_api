package users

import (
	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/models"
)

func UserResponse(user db.User) models.UserResponseParams {
	return models.UserResponseParams{
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Time,
	}
}
