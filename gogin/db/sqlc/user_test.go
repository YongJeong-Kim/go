package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/yongjeong-kim/go/gogin/util"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	err = testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	user, err := testQueries.GetUser(context.Background(), arg.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.PasswordChangedAt)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	createdUser := createRandomUser(t)

	user, err := testQueries.GetUser(context.Background(), createdUser.Username)
	require.NoError(t, err)

	require.Equal(t, createdUser.Username, user.Username)
	require.Equal(t, createdUser.HashedPassword, user.HashedPassword)
	require.Equal(t, createdUser.FullName, user.FullName)
	require.Equal(t, createdUser.Email, user.Email)
	require.WithinDuration(t, createdUser.PasswordChangedAt, user.PasswordChangedAt, time.Second)
	require.WithinDuration(t, createdUser.CreatedAt, user.CreatedAt, time.Second)
}
