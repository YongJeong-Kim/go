package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yongjeong-kim/go/api"
	db "github.com/yongjeong-kim/go/db/sqlc"
	"log"
)

const (
	dbDriver      = "mysql"
	dbSource      = "root:1234@tcp(localhost:13306)/go?parseTime=true"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
