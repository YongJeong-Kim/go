package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yongjeong-kim/go/goapigw/account/api"
	"github.com/yongjeong-kim/go/goapigw/account/service"
	"github.com/yongjeong-kim/go/goapigw/account/token"
	"net/http"
)

const (
	serviceName           = "account"
	serviceNameWithPrefix = "/account"
	accountv1             = serviceNameWithPrefix + "/api/v1"
	accountv2             = serviceNameWithPrefix + "/api/v2"
)

func main() {
	var maker token.TokenMaker = token.NewPasetoMaker()
	var servicer service.AccountServicer = service.NewAccountService(maker)
	accountServer := api.NewAccountServer(servicer)

	r := gin.New()
	r.POST(accountv1+"/login", accountServer.Login)
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
