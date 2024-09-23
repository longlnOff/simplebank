package db

import (
	"context"
	"testing"
	"time"

	"github.com/longln/simplebank/utils"
	"github.com/stretchr/testify/require"
)




func createRandomUser(t *testing.T) User {
	password := utils.RandomOwner()
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	arg := CreateUserParams{
		UserName: utils.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName: utils.RandomOwner(),
		Email: utils.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.UserName, user.UserName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}


func TestGetUser(t *testing.T) {
	// create random user
	user := createRandomUser(t)

	// get user and compare
	userTest, err := testQueries.GetUser(context.Background(), user.UserName)

	require.NoError(t, err)
	require.NotEmpty(t, userTest)

	require.Equal(t, user.UserName, userTest.UserName)
	require.Equal(t, user.HashedPassword, userTest.HashedPassword)
	require.Equal(t, user.FullName, userTest.FullName)
	require.Equal(t, user.Email, userTest.Email)

	require.WithinDuration(t, user.CreatedAt, userTest.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, userTest.PasswordChangedAt, time.Second)
}