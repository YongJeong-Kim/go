package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccept(t *testing.T) {
	testCases := []struct {
		name string
		run  func(ctx context.Context, t *testing.T, svc *Service, requestUserID, approveUserID string)
	}{
		{
			name: "OK",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, approveUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, approveUserID)
				require.NoError(t, err)
				err = svc.Friend.Accept(ctx, requestUserID, approveUserID)
				require.NoError(t, err)
			},
		},
		{
			name: "accept before request",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, approveUserID string) {
				err := svc.Friend.Accept(ctx, requestUserID, approveUserID)
				require.Error(t, err)
				require.Equal(t, "You must first receive a friend request from "+requestUserID, err.Error())
			},
		},
		{
			name: "request user not found",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, approveUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, approveUserID)
				require.NoError(t, err)
				err = svc.Friend.Accept(ctx, "request user not found", approveUserID)
				require.Error(t, err)
				require.Equal(t, "invalid request user uuid", err.Error())
			},
		},
		{
			name: "approve user not found",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, approveUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, approveUserID)
				require.NoError(t, err)
				err = svc.Friend.Accept(ctx, requestUserID, "approve user not found")
				require.Error(t, err)
				require.Equal(t, "invalid approve user uuid", err.Error())
			},
		},
		{
			name: "already friend",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, approveUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, approveUserID)
				require.NoError(t, err)
				err = svc.Friend.Accept(ctx, requestUserID, approveUserID)
				require.NoError(t, err)
				err = svc.Friend.Accept(ctx, requestUserID, approveUserID)
				require.Equal(t, "already friend", err.Error())
			},
		},
		{
			name: "accept yourself",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, approveUserID string) {
				err := svc.Friend.Accept(ctx, approveUserID, approveUserID)
				require.Error(t, err)
				require.Equal(t, "cannot accept yourself", err.Error())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, svc, createdID1, _ := createUserForTest(t, 6)
			_, _, createdID2, _ := createUserForTest(t, 6)
			tc.run(ctx, t, svc, createdID1, createdID2)
		})
	}
}

/*func TestListRequests(t *testing.T) {
	testCases := []struct {
		name string
		run  func(ctx context.Context, t *testing.T, svc *Service, requestUserID, approveUserID string)
	}{
		{
			name: "OK",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, approveUserID string) {

			},
		},
	}
}*/

func createFriendForTest(t *testing.T, usernameLen int) (string, string, string, string) {
	ctx, svc, createdID1, username1 := createUserForTest(t, usernameLen)
	_, _, createdID2, username2 := createUserForTest(t, usernameLen)

	err := svc.Friend.Request(ctx, createdID1, createdID2)
	require.NoError(t, err)
	err = svc.Friend.Accept(ctx, createdID1, createdID2)
	require.NoError(t, err)

	return createdID1, username1, createdID2, username2
}
