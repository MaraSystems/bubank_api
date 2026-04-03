package db

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/MaraSystems/bubank_api/utils"
	"github.com/stretchr/testify/require"
)

func CreateTestUser(t *testing.T) User {
	arg := CreateUserParams{
		Username:       utils.RandomUsername(),
		HashedPassword: utils.RandomString(32),
		FullName:       fmt.Sprintf("%s %s", utils.RandomUsername(), utils.RandomUsername()),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(t.Context(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.NotZero(t, user.CreatedAt)

	fmt.Println(user)

	return user
}

func TestCreateUserDB(t *testing.T) {
	CreateTestUser(t)
}

func TestGetUserDB(t *testing.T) {
	testUser := CreateTestUser(t)

	user, err := testQueries.GetUser(t.Context(), testUser.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, testUser.Username, user.Username)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateTestUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 0,
	}

	users, err := testQueries.ListUsers(t.Context(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, len(users), 5)
}

func TestUpdateUser(t *testing.T) {
	testUser := CreateTestUser(t)

	email := utils.RandomEmail()
	arg := UpdateUserParams{
		Username: testUser.Username,
		Email:    email,
	}

	user, err := testQueries.UpdateUser(t.Context(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, testUser.Username, user.Username)
	require.Equal(t, user.Email, email)
}

func TestDeleteUser(t *testing.T) {
	testUser := CreateTestUser(t)

	err := testQueries.DeleteUser(t.Context(), testUser.Username)
	require.NoError(t, err)

	user, err := testQueries.GetUser(t.Context(), testUser.Username)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, user)
}
