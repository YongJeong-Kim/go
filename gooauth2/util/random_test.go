package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRandomString(t *testing.T) {
	size := 32
	random := RandomString(size)

	require.Equal(t, len(random), size)
}
