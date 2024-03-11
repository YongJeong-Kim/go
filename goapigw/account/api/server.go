package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yongjeong-kim/go/goapigw/account/service"
	"net/http"
	"time"
)

type AccountServer struct {
	Service service.AccountServicer
}

func NewAccountServer(service service.AccountServicer) *AccountServer {
	return &AccountServer{
		Service: service,
	}
}

func (s *AccountServer) Login(c *gin.Context) {
	var req struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message":     "invalid username or password",
			"status_code": http.StatusBadRequest,
		})
		return
	}

	tk, err := s.Service.Login(req.Username, req.Password, time.Minute)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tk,
	})
}
