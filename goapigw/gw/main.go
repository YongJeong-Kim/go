package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	viper.SetConfigName("route")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("gw/config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("fatal error config file: ", err)
	}
	viper.AutomaticEnv()

	verifier := NewPasetoVerifier("asdfasf")
	verifier.Verify("")
	gw := Gateway{}
	err = viper.UnmarshalKey("gateway", &gw)
	if err != nil {
		panic("unmarshal failed")
	}
	fmt.Println(viper.Get("gateway"))

	r := gin.New()
	err = r.SetTrustedProxies(nil)
	if err != nil {
		log.Fatal("set trusted proxies failed")
	}

	r.Use(authMiddleware(verifier))
	for _, route := range gw.Routes {
		//proxy, err := NewProxy(route.Target)
		//if err != nil {
		//
		//}

		log.Printf("context name: %s, target: %s", route.Context, route.Target)
		//r.Any(route.Context+"*{targetPath}", Gen(route.Scheme, route.Target))

		for _, v := range route.Version {
			r.Any("/"+route.Context+"/api/"+v+"/*path", Gen(route.Scheme, route.Target))
		}
	}

	r.Run("0.0.0.0:38080")
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

func Gen(scheme, host string) gin.HandlerFunc {
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
