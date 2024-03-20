package main

import (
	"github.com/yongjeong-kim/go/goapigw/gw/api"
	"github.com/yongjeong-kim/go/goapigw/gw/token"
)

func main() {
	gw := api.LoadConfig("config")
	var tv token.TokenVerifier = token.NewPasetoVerifier(token.KeyHex)
	server := api.NewServer(tv)
	server.SetupRouter()
	server.SetupReverseProxy(gw)
	server.Router.Run(":38080")
	//viper.SetConfigName("route")
	//viper.SetConfigType("yaml")
	//viper.AddConfigPath("gw/config")
	//err := viper.ReadInConfig()
	//if err != nil {
	//	log.Fatal("fatal error config file: ", err)
	//}
	//viper.AutomaticEnv()
	//
	//verifier := token.NewPasetoVerifier("asdfasf")
	//verifier.Verify("")
	//server := api.NewServer(verifier)
	//
	//gw := Gateway{}
	//err = viper.UnmarshalKey("gateway", &gw)
	//if err != nil {
	//	panic("unmarshal failed")
	//}
	//fmt.Println(viper.Get("gateway"))
	//
	//r := gin.New()
	//err = r.SetTrustedProxies(nil)
	//if err != nil {
	//	log.Fatal("set trusted proxies failed")
	//}
	//
	//r.Use(api.authMiddleware(verifier))
	//for _, route := range gw.Routes {
	//	//proxy, err := NewProxy(route.Target)
	//	//if err != nil {
	//	//
	//	//}
	//
	//	log.Printf("context name: %s, target: %s", route.Context, route.Target)
	//	//r.Any(route.Context+"*{targetPath}", Gen(route.Scheme, route.Target))
	//
	//	for _, v := range route.Version {
	//		r.Any("/"+route.Context+"/api/"+v+"/*path", Gen(route.Scheme, route.Target))
	//	}
	//}
	//
	//r.Run("0.0.0.0:38080")
}
