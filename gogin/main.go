package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yongjeong-kim/go/api"
	db "github.com/yongjeong-kim/go/db/sqlc"
	"github.com/yongjeong-kim/go/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
