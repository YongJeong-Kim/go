package api

import (
	"context"
	"gooauth2/token"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (server *Server) login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid JSON provided.")
		return
	}

	sampleUser := User{
		ID:       10,
		Username: "aaa name",
	}

	if u.Username != sampleUser.Username || u.ID != sampleUser.ID {
		c.JSON(http.StatusUnauthorized, "Login failed.")
		return
	}

	tokenDuration := token.JWTDuration{
		AccessTokenDuration:  server.config.AccessTokenDuration,
		RefreshTokenDuration: server.config.RefreshTokenDuration,
	}

	token, err := server.token.CreateToken(u.Username, tokenDuration)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, token)
}

func (server *Server) CreateAuth(userID int64, payload *token.AccessTokenPayload) error {
	at := payload.ExpiredAt
	rt := payload.RtExpiredAt
	now := time.Now()

	errAccess := server.redisClient.Set(context.Background(), payload.AccessUUID, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := server.redisClient.Set(context.Background(), payload.RefreshUUID, strconv.Itoa(int(userID)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}
