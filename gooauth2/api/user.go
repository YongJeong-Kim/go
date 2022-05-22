package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"gooauth2/token"
	"gooauth2/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Request

type AuthInfo struct {
	platform     string
	OAuth2Config *oauth2.Config
	url          string
	header       http.Header
	errs         authErr
}

type authErr struct {
	invalidState            string
	exchangeFailed          string
	requestFailed           string
	responseFailed          string
	readResponseFailed      string
	responseBodyCloseFailed string
	unmarshalFailed         string
}

func getOAuth2Info(platform string) (AuthInfo, error) {
	authInfo := AuthInfo{}

	errMsg := func(platform string) authErr {
		return authErr{
			invalidState:            fmt.Sprintf("%s oauth2 invalid state", platform),
			exchangeFailed:          fmt.Sprintf("%s oauth2 exchange failed.", platform),
			requestFailed:           fmt.Sprintf("%s oauth2 request failed.", platform),
			responseFailed:          fmt.Sprintf("%s oauth2 response failed.", platform),
			readResponseFailed:      fmt.Sprintf("%s oauth2 read response body failed.", platform),
			responseBodyCloseFailed: fmt.Sprintf("%s oauth2 response.body close failed.", platform),
			unmarshalFailed:         fmt.Sprintf("%s oauth2 profile unmarshal failed.", platform),
		}
	}

	switch platform {
	case "google":
		authInfo = AuthInfo{
			platform: "google",
			OAuth2Config: &oauth2.Config{
				//ClientID:     "<GOOGLE_CLIENT_ID>",
				ClientID: "118399467217-h6aj2r1i0hciv6r08iqv4k1d7jujvsps.apps.googleusercontent.com",
				//ClientSecret: "<GOOGLE_CLIENT_SECRET>",
				ClientSecret: "cPLW30VEkDCJZNJHqk6w62Za",
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://accounts.google.com/o/oauth2/auth",
					TokenURL: "https://oauth2.googleapis.com/token",
				},
				RedirectURL: "http://localhost:8080/auth/callback/google",
				Scopes: []string{
					"https://www.googleapis.com/auth/userinfo.email",
					//"https://www.googleapis.com/auth/userinfo.profile",
				},
			},
			url:    "https://www.googleapis.com/oauth2/v2/userinfo?access_token=",
			header: http.Header{},
			errs:   errMsg(platform),
		}
	case "kakao":
		authInfo = AuthInfo{
			platform: "kakao",
			OAuth2Config: &oauth2.Config{
				ClientID:     "<KAKAO_CLIENT_ID>", // kakao rest api key
				ClientSecret: "<KAKAO_CLIENT>SECRET>",
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://kauth.kakao.com/oauth/authorize",
					TokenURL: "https://kauth.kakao.com/oauth/token",
				},
				RedirectURL: "http://localhost:8080/auth/callback/kakao",
			},
			url: "https://kapi.kakao.com/v2/user/me",
			header: http.Header{
				"Host":          []string{"kapi.kakao.com"},
				"Authorization": []string{"Bearer "},
			},
			errs: errMsg(platform),
		}
	case "naver":
		authInfo = AuthInfo{
			platform: "naver",
			OAuth2Config: &oauth2.Config{
				ClientID:     "<NAVER_CLIENT_ID>",
				ClientSecret: "<NAVER_CLIENT_SECRET>",
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://nid.naver.com/oauth2.0/authorize",
					TokenURL: "https://nid.naver.com/oauth2.0/token",
				},
				RedirectURL: "http://localhost:8080/auth/callback/naver",
			},
			url: "https://openapi.naver.com/v1/nid/me",
			header: http.Header{
				"Host":          []string{"localhost:8080"},
				"Authorization": []string{"Bearer "},
			},
			errs: errMsg(platform),
		}
	default:
		authInfo = AuthInfo{}
	}

	if reflect.ValueOf(authInfo).IsZero() {
		return authInfo, errors.New("Unsupported Platform Requested\n")
	}

	return authInfo, nil
}

func (server *Server) oauth2Login(c *gin.Context) {
	platform := c.Param("platform")
	aa := c.Request.URL.String()
	fmt.Println(aa)
	fmt.Println(platform)
	authInfo, err := getOAuth2Info(platform)
	if err != nil {
		log.Panic(err.Error())
		return
	}

	state := util.RandomString(32)
	err = server.createState(state)
	if err != nil {
		log.Panic("create state failed.")
		return
	}

	oauth2URL := authInfo.OAuth2Config.AuthCodeURL(state)
	//oauth2URL := oauth2ConfGoogle.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, oauth2URL)
}

func (server *Server) oauth2Callback(c *gin.Context) {
	platform := c.Param("platform")
	authInfo, err := getOAuth2Info(platform)
	if err != nil {
		log.Panic(err.Error())
		return
	}

	state := c.Request.FormValue("state")
	err = server.getState(state)
	if err != nil {
		log.Panic("invalid state.", err.Error())
		return
	}

	code := c.Request.FormValue("code")
	oauth2Token, err := authInfo.OAuth2Config.Exchange(c, code)
	if err != nil {
		log.Panic(authInfo.errs.exchangeFailed, err.Error())
		return
	}

	var reqURL string
	switch authInfo.platform {
	case "google":
		reqURL = authInfo.url + oauth2Token.AccessToken
	case "kakao":
		reqURL = authInfo.url
		authInfo.header.Set("Authorization", "Bearer "+oauth2Token.AccessToken)
	case "naver":
		reqURL = authInfo.url
		authInfo.header.Set("Authorization", "Bearer "+oauth2Token.AccessToken)
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Panic("oauth2 request failed. ", err.Error())
		return
	}

	req.Header = authInfo.header
	resp, err := client.Do(req)
	if err != nil {
		log.Panic("oauth2 response failed. ", err.Error())
		return
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic("read response body failed. ", err.Error())
		return
	}
	defer func(closer io.ReadCloser) {
		err := closer.Close()
		if err != nil {
			log.Panic("oauth2 response.body close failed. ", err.Error())
			return
		}
	}(resp.Body)

	var profile map[string]interface{}
	err = json.Unmarshal(content, &profile)
	if err != nil {
		log.Panic("oauth2 profile unmarshal failed. ", err.Error())
		return
	}

	c.JSON(http.StatusOK, profile)
	//c.Data(http.StatusOK, "application/json", content)
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

func (server *Server) createState(key string) error {
	ctx := context.Background()
	err := server.redisClient.Set(ctx, key, "state", time.Minute).Err()
	return err
}

func (server *Server) getState(key string) error {
	ctx := context.Background()
	_, err := server.redisClient.Get(ctx, key).Result()

	return err
}

func (server *Server) checkCodeState(c *gin.Context) {
	var
	c.JSON(http.StatusOK, gin.H{
		"msg": "adad",
	})
}