package store

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateUser(t *testing.T) {
	tx, q := newTestStore(t)

	err := q.CreateUser(context.Background(), "asdddddddddd")
	require.NoError(t, err)

	err = tx.Rollback()
	require.NoError(t, err)
}
