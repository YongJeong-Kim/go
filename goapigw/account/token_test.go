package main

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	var maker TokenMaker = NewPasetoMaker()
	token, err := maker.Create("ddd", time.Minute)
	require.NoError(t, err)
	require.NotNil(t, token)
	require.NotEmpty(t, token)
}
