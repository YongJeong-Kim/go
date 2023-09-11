package api

import (
	"github.com/YongJeong-Kim/go/goasynq/service"
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

func newTestServer(t *testing.T, srv service.Servicer) *Server {
	return &Server{
		Service: srv,
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
