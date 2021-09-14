package api

import (
	"fmt"
	"golang.org/x/oauth2"
	"gooauth2/config"
	"gooauth2/token"
	"io"
	"io/ioutil"
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

	oauth2Conf := &oauth2.Config{
		ClientID:     "<CLIENT_ID>",
		ClientSecret: "<CLIENT_SECRET>",
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

	fmt.Println(oauth2Conf)
	router.GET("/auth/google", func(ctx *gin.Context) {
		url := oauth2Conf.AuthCodeURL("random string")
		ctx.Redirect(http.StatusTemporaryRedirect, url)
	})
	router.GET("/auth/callback", func(ctx *gin.Context) {
		if ctx.Request.FormValue("state") != "random string" {
			fmt.Println("state is not valid")
		}

		token, err := oauth2Conf.Exchange(ctx, ctx.Request.FormValue("code"))
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
			fmt.Println("response read error: %s ", err.Error())
			return
		}
		ctx.Data(http.StatusOK, "application/json", content)
		//fmt.Fprintf(ctx.Writer, "Response: %s", content)
	})

	router.POST("/login", server.Login)
	// router.POST("/token/refresh", "")
	router.GET("/test", func(c *gin.Context) {
		c.SetCookie("aaa", "vxcv", 60*60*24, "/test", "/", true, true)
		c.JSON(http.StatusOK, gin.H{
			"zz": "cvv",
		})
	})

	authRoutes := router.Group("/").Use(server.authMiddleware(server.token))
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
