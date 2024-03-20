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
}
