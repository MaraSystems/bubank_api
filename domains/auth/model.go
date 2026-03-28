package auth

import (
	"testing"

	db "github.com/MaraSystems/graybank_api/db/sqlc"
	"github.com/MaraSystems/graybank_api/domains/users"
	"github.com/MaraSystems/graybank_api/utils"
	"github.com/stretchr/testify/require"
)

type RegisterParams struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseParams struct {
	Token string                   `json:"token"`
	User  users.UserResponseParams `json:"user"`
}

func DummyUser(t *testing.T) (user db.User, password string) {
	password = utils.RandomString(8)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       utils.RandomUsername(),
		FullName:       utils.RandomUsername(),
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
	}

	return
}
