package api

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"gooauth2/token"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

var oauth2ConfNaver = &oauth2.Config{
	ClientID:     "<NAVER_CLIENT_ID>",
	ClientSecret: "<NAVER_CLIENT_SECRET>",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://nid.naver.com/oauth2.0/authorize",
		TokenURL: "https://nid.naver.com/oauth2.0/token",
	},
	RedirectURL: "http://localhost:8080/auth/callback/naver",
}

func (server *Server) oauth2LoginNaver(c *gin.Context) {
	url := oauth2ConfNaver.AuthCodeURL("random string")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (server *Server) oauth2CallbackNaver(c *gin.Context) {
	if c.Request.FormValue("state") != "random string" {
		fmt.Println("state is not valid")
	}

	token, err := oauth2ConfNaver.Exchange(c, c.Request.FormValue("code"))
	if err != nil {
		fmt.Printf("could not get token: %s\n", err.Error())
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://openapi.naver.com/v1/nid/me", nil)
	if err != nil {
		log.Fatal("get naver user info failed.")
		return
	}

	req.Header = http.Header{
		"Host":          []string{"localhost:8080"},
		"Authorization": []string{"Bearer " + token.AccessToken},
	}

	//c.Request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	//resp, err := http.Get("https://openapi.naver.com/v1/nid/me" + token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("could not create get request: %s\n", err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Close response error: %s", err.Error())
			return
		}
	}(resp.Body)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("response read error: %s ", err.Error())
		return
	}
	c.Data(http.StatusOK, "application/json", content)
}

var oauth2ConfGoogle = &oauth2.Config{
	ClientID:     "<GOOGLE_CLIENT_ID>",
	ClientSecret: "<GOOGLE_CLIENT_SECRET>",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://oauth2.googleapis.com/token",
	},
	RedirectURL: "http://localhost:8080/auth/callback",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email",
		//"https://www.googleapis.com/auth/userinfo.profile",
	},
}

func (server *Server) oauth2LoginGoogle(c *gin.Context) {
	url := oauth2ConfGoogle.AuthCodeURL("random string")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (server *Server) oauth2CallbackGoogle(c *gin.Context) {
	if c.Request.FormValue("state") != "random string" {
		fmt.Println("state is not valid")
	}

	token, err := oauth2ConfGoogle.Exchange(c, c.Request.FormValue("code"))
	if err != nil {
		fmt.Printf("could not get token: %s\n", err.Error())
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		fmt.Printf("could not create get request: %s\n", err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Close response error: %s", err.Error())
			return
		}
	}(resp.Body)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("response read error: %s ", err.Error())
		return
	}
	c.Data(http.StatusOK, "application/json", content)
	//fmt.Fprintf(ctx.Writer, "Response: %s", content)
}

var oauthConfKakao = &oauth2.Config{
	ClientID:     "<KAKAO_CLIENT_ID>", // kakao rest api key
	ClientSecret: "<KAKAO_CLIENT_SECRET>",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://kauth.kakao.com/oauth/authorize",
		TokenURL: "https://kauth.kakao.com/oauth/token",
	},
	RedirectURL: "http://localhost:8080/auth/callback/kakao",
}

func (server *Server) oauth2LoginKakao(c *gin.Context) {
	url := oauthConfKakao.AuthCodeURL("random string")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (server *Server) oauth2CallbackKakao(c *gin.Context) {
	if c.Request.FormValue("state") != "random string" {
		fmt.Println("state is not valid")
	}

	token, err := oauthConfKakao.Exchange(c, c.Request.FormValue("code"))
	if err != nil {
		fmt.Printf("could not get token: %s\n", err.Error())
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		log.Fatal("get kakao user info failed.")
		return
	}

	req.Header = http.Header{
		"Host":          []string{"kapi.kakao.com"},
		"Authorization": []string{"Bearer " + token.AccessToken},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("could not create get request: %s\n", err.Error())
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Close response error: %s", err.Error())
			return
		}
	}(resp.Body)
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("response read error: %s ", err.Error())
		return
	}
	c.Data(http.StatusOK, "application/json", content)
	//fmt.Fprintf(ctx.Writer, "Response: %s", content)
}

func (server *Server) Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid JSON provided.")
		return
	}

	sampleUser := User{
		ID:       "uuid",
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

	tokenDetails, err := server.maker.CreateToken(u.Username, tokenDuration)
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

	accessToken, err := server.maker.ExtractToken(authorizationHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	accessTokenPayload, err := server.maker.VerifyAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	ctx := context.Background()
	deleted, err := server.redisClient.Del(ctx, accessTokenPayload.ID.String()).Result()
	if err != nil || deleted == 0 {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (server *Server) refreshToken(c *gin.Context) {
	authorizationHeader := c.GetHeader("authorization")
	refreshToken, err := server.maker.ExtractToken(authorizationHeader)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	refreshTokenPayload, err := server.maker.VerifyRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()
	deleted, err := server.redisClient.Del(ctx, refreshTokenPayload.ID.String()).Result()
	if err != nil || deleted == 0 {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	tokenDuration := token.JWTDuration{
		AccessTokenDuration:  server.config.AccessTokenDuration,
		RefreshTokenDuration: server.config.RefreshTokenDuration,
	}
	payloadDetails, err := server.maker.CreateToken(refreshTokenPayload.Username, tokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	err = server.CreateAuth(payloadDetails.Payload.AccessTokenPayload.UserID, payloadDetails)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, payloadDetails.Token)
}

func (server *Server) CreateAuth(userID string, details *token.PayloadDetails) error {
	ctx := context.Background()
	at := details.Payload.AccessTokenPayload.ExpiredAt
	rt := details.Payload.RefreshTokenPayload.ExpiredAt
	now := time.Now()

	accessUUID := details.Payload.AccessTokenPayload.ID.String()
	refreshUUID := details.Payload.RefreshTokenPayload.ID.String()

	errAccess := server.redisClient.Set(ctx, accessUUID, userID, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := server.redisClient.Set(ctx, refreshUUID, userID, rt.Sub(now)).Err()
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

func (server *Server) DeleteToken(tokenID string) error {
	ctx := context.Background()
	deleted, err := server.redisClient.Del(ctx, tokenID).Result()
	if err != nil || deleted == 0 {
		return fmt.Errorf("delete refresh token failed")
	}
	return nil
}
