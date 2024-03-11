package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.New()
	r.GET("/service-b/api/v1/111", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "ok1",
		})
		return
	})
	r.GET("/service-b/api/v1/111/111", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "bbb111",
		})
		return
	})

	r.Run(":8081")
}
