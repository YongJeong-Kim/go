package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	createUserForTest(t, 6)
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name string
		run  func(ctx context.Context, t *testing.T, svc *Service, createdID, username string)
	}{
		{
			name: "OK",
			run: func(ctx context.Context, t *testing.T, svc *Service, createdID, username string) {
				user, err := svc.User.Get(ctx, createdID)
				require.NoError(t, err)
				require.Equal(t, createdID, user.ID)
				require.Equal(t, username, user.Name)
				require.WithinDuration(t, time.Now().UTC(), user.CreatedDate, time.Second)
			},
		},
		{
			name: "not found",
			run: func(ctx context.Context, t *testing.T, svc *Service, createdID, username string) {
				user, err := svc.User.Get(ctx, "not found user")
				require.Error(t, err)
				require.Equal(t, errors.New("Result contains no more records").Error(), err.Error())
				require.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, svc, createdID, username := createUserForTest(t, 6)
			tc.run(ctx, t, svc, createdID, username)
		})
	}

}

func createUserForTest(t *testing.T, usernameLen int) (context.Context, *Service, string, string) {
	ctx, svc := newTestService()
	require.NotNil(t, ctx)
	require.NotNil(t, svc)

	username := randomStringForTest(usernameLen)
	createdID, err := svc.User.Create(ctx, username)
	require.NoError(t, err)
	require.NoError(t, uuid.Validate(createdID))

	return ctx, svc, createdID, username
}

func randomStringForTest(usernameLen int) string {
	var result string
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for range usernameLen {
		r := rand.Int63n(int64(len(str)))
		result += string(str[r])
	}
	return result
}
