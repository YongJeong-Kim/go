package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yongjeong-kim/go/goapigw/account/service"
	"net/http"
	"time"
)

const (
	serviceName           = "account"
	ServiceNameWithPrefix = "/account"
	Accountv1             = ServiceNameWithPrefix + "/api/v1"
	Accountv2             = ServiceNameWithPrefix + "/api/v2"
)

type AccountServer struct {
	Service service.AccountServicer
	Router  *gin.Engine
}

func NewAccountServer(service service.AccountServicer) *AccountServer {
	return &AccountServer{
		Service: service,
	}
}

func (s *AccountServer) SetupRouter() {
	r := gin.New()
	r.POST(Accountv1+"/login", s.Login)
	r.GET(Accountv1+"/111", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "version 2 available now",
		})
		return
	})
	r.GET(Accountv2+"/111", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "version 2 available now",
		})
		return
	})
	r.POST(Accountv2+"/asdf", func(c *gin.Context) {
		var r struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}
		if err := c.ShouldBindJSON(&r); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "need name, age",
			})
			return
		}

		c.JSON(http.StatusOK, r)
	})

	s.Router = r
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
		"access_token": tk,
	})
}
