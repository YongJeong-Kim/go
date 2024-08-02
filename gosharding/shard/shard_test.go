package shard

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMustBe(t *testing.T) {
	s := NewShard(3)

	id, _ := uuid.NewV7()
	wantShardIdx := 1
	wantNewID, err := s.MustBe(wantShardIdx, id.String())
	require.NoError(t, err)
	idx := s.Index(wantNewID)
	require.Equal(t, wantShardIdx, int(idx))
}
