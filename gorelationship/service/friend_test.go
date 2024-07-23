package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccept(t *testing.T) {
	testCases := []struct {
		name string
		run  func(ctx context.Context, t *testing.T, svc *Service, fromUserID, toUserID string)
	}{
		{
			name: "OK",
			run: func(ctx context.Context, t *testing.T, svc *Service, fromUserID, toUserID string) {
				err := svc.Friend.Request(ctx, fromUserID, toUserID)
				require.NoError(t, err)
				err = svc.Friend.Accept(ctx, fromUserID, toUserID)
				require.NoError(t, err)
			},
		},
		{
			name: "not found",
			run: func(ctx context.Context, t *testing.T, svc *Service, fromUserID, toUserID string) {
				err := svc.Friend.Accept(ctx, fromUserID, toUserID)
				require.Error(t, err)
				require.Equal(t, "You must first receive a friend request from "+fromUserID, err.Error())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, svc, createdID1, _ := createUserForTest(t, 6)
			_, _, createdID2, _ := createUserForTest(t, 6)
			tc.run(ctx, t, svc, createdID1, createdID2)
		})
	}
}
func createFriendForTest(t *testing.T, usernameLen int) (string, string, string, string) {
	ctx, svc, createdID1, username1 := createUserForTest(t, usernameLen)
	_, _, createdID2, username2 := createUserForTest(t, usernameLen)

	err := svc.Friend.Request(ctx, createdID1, createdID2)
	require.NoError(t, err)
	err = svc.Friend.Accept(ctx, createdID1, createdID2)
	require.NoError(t, err)

	return createdID1, username1, createdID2, username2
}
