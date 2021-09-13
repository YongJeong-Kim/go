package api

import (
	"context"
	"fmt"
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

func (server *Server) Login(c *gin.Context) {
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

	tokenDetails, err := server.token.CreateToken(u.Username, tokenDuration)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = server.CreateAuth(u.ID, tokenDetails)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokenDetails.Token)
}

func (server *Server) Logout(c *gin.Context) {
	authorizationHeader := c.GetHeader("authorization")

	accessToken, err := server.token.VerifyToken(authorizationHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
	}

	ctx := context.Background()
	deleted, err := server.redisClient.Del(ctx, accessToken.ID.String()).Result()
	if err != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (server *Server) CreateAuth(userID int64, details *token.PayloadDetails) error {
	ctx := context.Background()
	at := details.Payload.AccessTokenPayload.ExpiredAt
	rt := details.Payload.RefreshTokenPayload.ExpiredAt
	now := time.Now()

	accessUUID := details.Payload.AccessTokenPayload.ID.String()
	refreshUUID := details.Payload.RefreshTokenPayload.ID.String()

	errAccess := server.redisClient.Set(ctx, accessUUID, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := server.redisClient.Set(ctx, refreshUUID, strconv.Itoa(int(userID)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (server *Server) DeleteAuth(details *token.PayloadDetails) error {
	ctx := context.Background()

	at := details.Payload.AccessTokenPayload.ID.String()
	rt := details.Payload.RefreshTokenPayload.ID.String()

	deleted, err := server.redisClient.Del(ctx, at).Result()
	if err != nil || deleted == 0 {
		return fmt.Errorf("logout failed. delete access token failed")
	}

	deleted, err = server.redisClient.Del(ctx, rt).Result()
	if err != nil || deleted == 0 {
		return fmt.Errorf("logout failed. delete refresh token failed")
	}

	return nil
}
