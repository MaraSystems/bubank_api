package users

import db "github.com/MaraSystems/graybank_api/db/sqlc"

func UserResponse(user db.User) UserResponseParams {
	return UserResponseParams{
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Time,
	}
}
