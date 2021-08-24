package api

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gooauth2/config"
	"gooauth2/token"
)

type Server struct {
	config      config.App
	redisClient *redis.Client
	token       token.Maker
}

func NewServer(app config.App) (*Server, error) {
	token, err := token.NewJWTMaker(app.AccessSecret, app.RefreshSecret)
	if err != nil {
		return nil, fmt.Errorf("Cannot load config. %s ", err.Error())
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
	return server, nil
}
