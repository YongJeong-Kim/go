package main

import (
	"github.com/yongjeong-kim/go/goapigw/account/api"
	"github.com/yongjeong-kim/go/goapigw/account/service"
	"github.com/yongjeong-kim/go/goapigw/account/token"
)

func main() {
	var maker token.TokenMaker = token.NewPasetoMaker()
	var servicer service.AccountServicer = service.NewAccountService(maker)
	accountServer := api.NewAccountServer(servicer)
	accountServer.SetupRouter()
	accountServer.Router.Run(":8080")
}
