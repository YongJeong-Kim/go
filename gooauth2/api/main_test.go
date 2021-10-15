package api

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gooauth2/config"
	"os"
	"testing"
)

func newTestServer(t *testing.T) *Server {
	cfg, err := config.LoadConfig("../.")
	require.NoError(t, err)

	server, err := NewServer(cfg)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
