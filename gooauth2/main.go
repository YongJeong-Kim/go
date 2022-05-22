package main

import (
	"gooauth2/api"
	"gooauth2/config"
	"log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config.")
	}

	server, err := api.NewServer(cfg)
	if err != nil {
		log.Fatal("cannot create new server. ", err)
	}

	err = server.Start(cfg.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server.")
	}
}
