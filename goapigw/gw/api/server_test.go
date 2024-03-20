package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	tkmock "github.com/yongjeong-kim/go/goapigw/gw/token/mock"
	"go.uber.org/mock/gomock"
	"strings"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := tkmock.NewMockTokenVerifier(ctrl)
	server := newTestServer(m)
	setupRouterForTest(t, server)
}

func setupRouterForTest(t *testing.T, server *Server) {
	server.SetupRouter()
	require.NotNil(t, server.Router)
	require.Equal(t, fmt.Sprintf("%T", server.Router), fmt.Sprintf("%T", &gin.Engine{}))
}

func TestSetupReverseProxy(t *testing.T) {
	gw := loadConfigForTest(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := tkmock.NewMockTokenVerifier(ctrl)
	server := newTestServer(m)
	setupRouterForTest(t, server)
	server.SetupReverseProxy(gw)

	versionCnt := 0
	for _, g := range gw.Routes {
		versionCnt += len(g.Version)
	}

	// Total 9
	// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE
	require.Equal(t, 9*versionCnt, len(server.Router.Routes()))
}

func TestLoadConfig(t *testing.T) {
	loadConfigForTest(t)
}

func loadConfigForTest(t *testing.T) *Gateway {
	gw := LoadConfig("../config")
	require.NotEmpty(t, gw.ListenAddr)
	require.NotNil(t, gw.ListenAddr)

	for _, r := range gw.Routes {
		require.True(t, r.Scheme == "http" || r.Scheme == "https")

		ctx := strings.Split(r.Context, "/")
		require.Equal(t, len(ctx), 2)

		target := strings.Split(r.Target, ":")
		require.Equal(t, len(target), 2)

		require.GreaterOrEqual(t, len(r.Version), 1)
	}

	return gw
}
