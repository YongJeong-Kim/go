package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yongjeong-kim/go/goapigw/account/service"
	"os"
	"testing"
)

func newTestServer(service service.AccountServicer) *AccountServer {
	s := &AccountServer{
		Service: service,
	}
	s.SetupRouter()

	return s
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
