package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gosharding/config"
	"gosharding/repository"
	"gosharding/service"
	"gosharding/shard"
	"log"
)

func main() {
	cfg := config.LoadDBConfig(".")
	dbs := config.ConnDBs(cfg, cfg.ShardCount)
	repo := repository.NewRepository(dbs, repository.NewUser())
	shard := shard.NewShard(cfg.ShardCount)
	svc := service.NewService(repo, shard)
	//svr := api.NewServer(dbs, svc)

	id, _ := uuid.NewV7()
	err := svc.Create(context.Background(), id.String(), "asdf")
	if err != nil {
		log.Println(err)
	}
}
