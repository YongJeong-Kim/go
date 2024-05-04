package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Server struct {
	Router *gin.Engine
}

func NewServer(r *gin.Engine) *Server {
	return &Server{
		Router: r,
	}
}

func (s *Server) SetupRouter() {
	s.Router.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
		return
	})
}

func main() {
	s := NewServer(gin.New())
	s.SetupRouter()

	if err := s.Router.Run(":8080"); err != nil {
		log.Fatal("run error", err)
	}
}
