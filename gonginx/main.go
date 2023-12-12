package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.New()
	r.GET("/", func(c *gin.Context) {
		log.Println("in")

		c.Status(http.StatusOK)
		return
	})

	r.Run()
}
