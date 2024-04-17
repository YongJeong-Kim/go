package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yongjeong-kim/go/goapigw/gw/token"
	"net/http"
	"net/http/httputil"
)

type Server struct {
	TokenVerifier token.TokenVerifier
	Router        *gin.Engine
}

func NewServer(tokenVerifier token.TokenVerifier) *Server {
	return &Server{
		TokenVerifier: tokenVerifier,
	}
}

func addProxy(scheme, host string) gin.HandlerFunc {
	return func(c *gin.Context) {
		director := func(r *http.Request) {
			req := c.Request
			r.URL.Scheme = scheme
			r.URL.Host = host
			r.Header["my-header"] = []string{req.Header.Get("my-header")}
			delete(r.Header, "my-header")
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) SetupRouter() {
	r := gin.New()
	s.Router = r
}

func (s *Server) SetupReverseProxy(gateway *Gateway) {
	for _, r := range gateway.Routes {
		for _, v := range r.Version {
			path := fmt.Sprintf("/%s/%s/%s", v, r.Context, "*path")
			if r.Context == "account" {
				s.Router.Any(path, addProxy(r.Scheme, r.Target))
				continue
			}
			s.Router.Use(authMiddleware(s.TokenVerifier)).Any(path, addProxy(r.Scheme, r.Target))
		}
	}
}

type Gateway struct {
	ListenAddr string  `mapstructure:"listenAddr"`
	Routes     []Route `mapstructure:"routes"`
}

type Route struct {
	Version []string `mapstructure:"version"`
	Scheme  string   `mapstructure:"scheme"`
	Context string   `mapstructure:"context"`
	Target  string   `mapstructure:"target"`
}

func LoadConfig(path string) *Gateway {
	viper.SetConfigType("yml")
	viper.AddConfigPath(path)
	viper.SetConfigName("route")
	//viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	gateway := &Gateway{}
	err = viper.UnmarshalKey("gateway", gateway)
	if err != nil {
		panic(fmt.Errorf("Fatal error unmarshal gateway: %s \n", err))
	}

	return gateway
}
