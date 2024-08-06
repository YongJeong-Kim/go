package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestListRangeCreatedAt(t *testing.T) {
	ctx, svc, _, _ := createForTest(t, "qwer")

	end := time.Now().UTC()
	start := end.Add(-1 * time.Second)
	users, err := svc.ListRangeCreatedAt(ctx, start, end)
	require.NoError(t, err)
	require.Greater(t, len(users), 0)
}

func TestGet(t *testing.T) {
	ctx, svc, id, name := createForTest(t, "qwer")
	user, err := svc.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, id, user.ID)
	require.Equal(t, name, user.Name)
	require.WithinDuration(t, time.Now().UTC(), user.CreatedAt, time.Second)
}

func TestCreate(t *testing.T) {
	_, _, id, _ := createForTest(t, "asdf")
	require.NoError(t, uuid.Validate(id))
}

func createForTest(t *testing.T, name string) (context.Context, *Service, string, string) {
	svc := newTestService()

	id, err := NewUserID()
	require.NoError(t, err)

	ctx := context.Background()
	err = svc.Create(ctx, id, name)
	require.NoError(t, err)
	return ctx, svc, id, name
}
