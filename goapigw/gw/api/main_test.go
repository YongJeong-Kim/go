package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yongjeong-kim/go/goapigw/gw/token"
	"os"
	"testing"
)

func newTestServer(tokenVerifier token.TokenVerifier) *Server {
	return &Server{
		TokenVerifier: tokenVerifier,
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
