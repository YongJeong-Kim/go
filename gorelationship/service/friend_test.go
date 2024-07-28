package service

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
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

func TestRequest(t *testing.T) {
	testCases := []struct {
		name string
		run  func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string)
	}{
		{
			name: "OK",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)
			},
		},
		{
			name: "already send request",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)
				err = svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.Error(t, err)
				require.Equal(t, "already send request", err.Error())
			},
		},
		{
			name: "request user not found",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, "request user not found", acceptUserID)
				require.Error(t, err)
				require.Equal(t, "invalid request user uuid", err.Error())
			},
		},
		{
			name: "accept user not found",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, "accept user not found")
				require.Error(t, err)
				require.Equal(t, "invalid accept user uuid", err.Error())
			},
		},
		{
			name: "already friend",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)
				err = svc.Friend.Accept(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)
				err = svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.Equal(t, "already friend", err.Error())
			},
		},
		{
			name: "request yourself",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, requestUserID)
				require.Error(t, err)
				require.Equal(t, "cannot request yourself", err.Error())
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

func TestFromRequestCount(t *testing.T) {
	testCases := []struct {
		name string
		run  func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string)
	}{
		{
			name: "OK",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)

				count, err := svc.Friend.FromRequestCount(ctx, acceptUserID)
				require.NoError(t, err)
				require.Equal(t, int64(1), count)
			},
		},
		{
			name: "invalid user id",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				count, err := svc.Friend.FromRequestCount(ctx, "invalid user id")
				require.Error(t, err)
				require.Equal(t, "invalid user uuid", err.Error())
				require.Equal(t, int64(0), count)
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

func TestCount(t *testing.T) {
	testCases := []struct {
		name string
		run  func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string)
	}{
		{
			name: "OK",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)
				err = svc.Friend.Accept(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)

				cnt, err := svc.Friend.Count(ctx, acceptUserID)
				require.NoError(t, err)
				require.Equal(t, int64(1), cnt)
			},
		},
		{
			name: "invalid user id",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)
				err = svc.Friend.Accept(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)

				cnt, err := svc.Friend.Count(ctx, "invalid user uuid")
				require.Equal(t, "invalid user uuid", err.Error())
				require.Equal(t, int64(0), cnt)
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

func TestList(t *testing.T) {
	testCases := []struct {
		name string
		run  func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string)
	}{
		{
			name: "OK",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)
				err = svc.Friend.Accept(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)

				fs, err := svc.Friend.List(ctx, acceptUserID)
				require.NoError(t, err)
				require.Equal(t, 1, len(fs))

				for _, f := range fs {
					require.Equal(t, requestUserID, f.ID)
					require.NotEmpty(t, f.Name)
					require.WithinDuration(t, time.Now().UTC(), f.CreatedDate, time.Second)
				}
			},
		},
		{
			name: "invalid user id",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)
				err = svc.Friend.Accept(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)

				fs, err := svc.Friend.List(ctx, "invalid user uuid")
				require.Equal(t, "invalid user uuid", err.Error())
				require.Equal(t, 0, len(fs))
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

func TestListFromRequests(t *testing.T) {
	testCases := []struct {
		name string
		run  func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string)
	}{
		{
			name: "OK",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)

				rs, err := svc.Friend.ListFromRequests(ctx, acceptUserID)
				require.NoError(t, err)
				require.Equal(t, 1, len(rs))

				for _, f := range rs {
					require.Equal(t, requestUserID, f.ID)
					require.NotEmpty(t, f.Name)
					require.WithinDuration(t, time.Now().UTC(), f.CreatedDate, time.Second)
				}
			},
		},
		{
			name: "invalid user id",
			run: func(ctx context.Context, t *testing.T, svc *Service, requestUserID, acceptUserID string) {
				err := svc.Friend.Request(ctx, requestUserID, acceptUserID)
				require.NoError(t, err)

				rs, err := svc.Friend.ListFromRequests(ctx, "invalid user uuid")
				require.Equal(t, "invalid user uuid", err.Error())
				require.Equal(t, 0, len(rs))
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

func TestMutualCount(t *testing.T) {
	testCases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "OK",
			run: func(t *testing.T) {
				ctx, svc, userID, fs := createFriendForTest(t, 3, 6)

				ctx1, svc1, userID1, _ := createUserForTest(t, 6)

				friendRequestForTest(ctx1, t, svc1, userID1, fs)
				friendAcceptForTest(ctx, t, svc, fs[0], []string{userID1})
				friendAcceptForTest(ctx, t, svc, fs[1], []string{userID1})
				friendAcceptForTest(ctx, t, svc, fs[2], []string{userID1})

				cnt, err := svc.Friend.MutualCount(ctx, userID, userID1)
				require.NoError(t, err)
				require.Equal(t, int64(3), cnt)
			},
		},
		{
			name: "invalid user id1",
			run: func(t *testing.T) {
				ctx, svc, userID, _ := createUserForTest(t, 6)
				cnt, err := svc.Friend.MutualCount(ctx, "invalid user id", userID)
				require.Error(t, err)
				require.Equal(t, "invalid user1 uuid", err.Error())
				require.Equal(t, int64(0), cnt)
			},
		},
		{
			name: "invalid user id2",
			run: func(t *testing.T) {
				ctx, svc, userID, _ := createUserForTest(t, 6)
				cnt, err := svc.Friend.MutualCount(ctx, userID, "invalid user id")
				require.Error(t, err)
				require.Equal(t, "invalid user2 uuid", err.Error())
				require.Equal(t, int64(0), cnt)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.run(t)
		})
	}
}

func createFriendRequestForTest(t *testing.T, requestNum int, usernameLen int) (context.Context, *Service, string, []string) {
	ctx, svc, createdID1, _ := createUserForTest(t, usernameLen)

	reqIDs := make([]string, requestNum)
	for range requestNum {
		_, _, createdID, _ := createUserForTest(t, usernameLen)
		err := svc.Friend.Request(ctx, createdID1, createdID)
		require.NoError(t, err)

		reqIDs = append(reqIDs, createdID)
	}
	return ctx, svc, createdID1, reqIDs
}

func friendRequestForTest(ctx context.Context, t *testing.T, svc *Service, requestUserID string, userIDs []string) {
	for _, u := range userIDs {
		err := svc.Friend.Request(ctx, requestUserID, u)
		require.NoError(t, err)
	}
}

func friendAcceptForTest(ctx context.Context, t *testing.T, svc *Service, acceptUserID string, requestUserIDs []string) {
	for _, r := range requestUserIDs {
		err := svc.Friend.Accept(ctx, r, acceptUserID)
		require.NoError(t, err)
	}
}

func createFriendForTest(t *testing.T, friendNum, usernameLen int) (context.Context, *Service, string, []string) {
	ctx, svc, createdID1, _ := createUserForTest(t, usernameLen)

	fs := make([]string, friendNum)
	for i := 0; i < friendNum; i++ {
		c, s, createdID2, _ := createUserForTest(t, usernameLen)
		err := svc.Friend.Request(ctx, createdID1, createdID2)
		require.NoError(t, err)
		err = s.Friend.Accept(c, createdID1, createdID2)
		require.NoError(t, err)
		fs[i] = createdID2
	}

	return ctx, svc, createdID1, fs
}
