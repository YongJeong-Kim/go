package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	serviceName           = "account"
	serviceNameWithPrefix = "/account"
	accountv1             = serviceNameWithPrefix + "/api/v1"
	accountv2             = serviceNameWithPrefix + "/api/v2"
)

func main() {
	var maker TokenMaker = NewPasetoMaker()
	accountServer := NewAccountServer(maker)
	r := gin.New()
	r.POST(accountv1+"/login", accountServer.login)
	r.GET(accountv1+"/111", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "version 2 available now",
		})
		return
	})
	r.GET(accountv1+"/111", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "version 2 available now",
		})
		return
	})
	r.GET(accountv2+"/111", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "v2shortaaa111",
		})
		return
	})
	r.POST(accountv2+"/asdf", func(c *gin.Context) {
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

	r.Run(":8080")
}

type AccountServer struct {
	Maker TokenMaker
}

func NewAccountServer(maker TokenMaker) *AccountServer {
	return &AccountServer{
		Maker: maker,
	}
}

func (s *AccountServer) login(c *gin.Context) {
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

	if req.Username == "aaa" && req.Password == "1234" {

	}
	t, err := s.Maker.Create(req.Username, 10*time.Minute)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": t,
	})
}
