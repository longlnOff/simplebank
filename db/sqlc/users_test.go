package db

import (
	"context"
	"database/sql"
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


func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := utils.RandomOwner()
	testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.UserName,
		FullName: sql.NullString{String: newFullName, Valid: true},
	})

	userTest, err := testQueries.GetUser(context.Background(), oldUser.UserName)

	require.NoError(t, err)
	require.NotEmpty(t, userTest)

	require.Equal(t, oldUser.UserName, userTest.UserName)
	require.Equal(t, newFullName, userTest.FullName)
	require.Equal(t, oldUser.Email, userTest.Email)
	require.Equal(t, oldUser.HashedPassword, userTest.HashedPassword)
	require.NotEqual(t, oldUser.FullName, userTest.FullName)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := utils.RandomEmail()
	testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.UserName,
		Email: sql.NullString{String: newEmail, Valid: true},
	})

	userTest, err := testQueries.GetUser(context.Background(), oldUser.UserName)

	require.NoError(t, err)
	require.NotEmpty(t, userTest)

	require.Equal(t, oldUser.UserName, userTest.UserName)
	require.Equal(t, oldUser.FullName, userTest.FullName)
	require.Equal(t, newEmail, userTest.Email)
	require.Equal(t, oldUser.HashedPassword, userTest.HashedPassword)
	require.NotEqual(t, oldUser.Email, userTest.Email)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword := utils.RandomString(6)
	hashed_password, err := utils.HashPassword(newPassword)

	require.NoError(t, err)
	testUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.UserName,
		HashedPassword: sql.NullString{String: hashed_password, Valid: true},
	})

	require.NoError(t, err)
	require.NotEmpty(t, testUser)

	require.Equal(t, oldUser.UserName, testUser.UserName)
	require.Equal(t, oldUser.FullName, testUser.FullName)
	require.Equal(t, oldUser.Email, testUser.Email)
	require.Equal(t, hashed_password, testUser.HashedPassword)
	require.NotEqual(t, oldUser.HashedPassword, testUser.HashedPassword)
}


func TestUpdateUserAllFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := utils.RandomOwner()
	newEmail := utils.RandomEmail()
	newPassword := utils.RandomString(6)
	hasedPassword, err := utils.HashPassword(newPassword)
	require.NoError(t, err)
	testUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.UserName,
		FullName: sql.NullString{String: newFullName, Valid: true},
		Email: sql.NullString{String: newEmail, Valid: true},
		HashedPassword: sql.NullString{String: hasedPassword, Valid: true},
	})

	require.NoError(t, err)
	require.NotEmpty(t, testUser)

	require.Equal(t, oldUser.UserName, testUser.UserName)
	require.Equal(t, newFullName, testUser.FullName)
	require.Equal(t, newEmail, testUser.Email)
	require.Equal(t, hasedPassword, testUser.HashedPassword)
	require.NotEqual(t, oldUser.Email, testUser.Email)
	require.NotEqual(t, oldUser.FullName, testUser.FullName)
	require.NotEqual(t, oldUser.HashedPassword, testUser.HashedPassword)
	
}