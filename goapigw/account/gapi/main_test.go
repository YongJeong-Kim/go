package gapi

import (
	"github.com/gin-gonic/gin"
	"github.com/yongjeong-kim/go/goapigw/account/service"
	"os"
	"testing"
)

func newTestServer(svc service.AccountServicer) *AccountServer {
	return &AccountServer{
		servicer: svc,
	}
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func runTestGRPCServer() {

}

func runTestGatewayServer() {

}
