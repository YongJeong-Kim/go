package api

import (
	"fmt"
	"gooauth2/config"
	"gooauth2/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	config      config.App
	redisClient *redis.Client
	token       token.Maker
	router      *gin.Engine
}

func NewServer(app config.App) (*Server, error) {
	token, err := token.NewJWTMaker(app.AccessSecret, app.RefreshSecret)
	if err != nil {
		return nil, fmt.Errorf("cannot load config. %s ", err.Error())
	}

	redisClient, err := config.ConnectRedis()
	if err != nil {
		return nil, fmt.Errorf("connection failed. redis : %s ", err.Error())
	}

	server := &Server{
		config:      app,
		redisClient: redisClient,
		token:       token,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/login", server.Login)
	// router.POST("/token/refresh", "")
	router.GET("/test", func(c *gin.Context) {
		c.SetCookie("aaa", "vxcv", 60*60*24, "/test", "/", true, true)
		c.JSON(http.StatusOK, gin.H{
			"zz": "cvv",
		})
	})

	authRoutes := router.Group("/auth").Use(server.authMiddleware(server.token))
	authRoutes.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"aa": "vcx",
		})
	})
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
