package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yongjeong-kim/go/gogin/api"
	db "github.com/yongjeong-kim/go/gogin/db/sqlc"
	"github.com/yongjeong-kim/go/gogin/util"
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
