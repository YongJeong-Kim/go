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
	maker       token.Maker
	router      *gin.Engine
}

func NewServer(app config.App) (*Server, error) {
	maker, err := token.NewJWTMaker(app.AccessSecret, app.RefreshSecret)
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
		maker:       maker,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/auth/:platform", server.oauth2Login)
	router.GET("/auth/callback/:platform", server.oauth2Callback)

	router.POST("/login", server.Login)
	router.POST("/token/refresh", server.refreshToken)
	router.POST("/auth/check/google", server.checkCodeState)
	router.GET("/test", func(c *gin.Context) {
		c.SetCookie("aaa", "vxcv", 60*60*24, "/test", "/", true, true)
		c.JSON(http.StatusOK, gin.H{
			"zz": "cvv",
		})
	})

	authRoutes := router.Group("/").Use(server.authMiddleware(server.maker))
	authRoutes.POST("/logout", server.Logout)
	authRoutes.GET("/authtest", func(c *gin.Context) {
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
