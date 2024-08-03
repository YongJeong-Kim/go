package main

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"gosharding/config"
	"gosharding/repository"
	"gosharding/service"
	"gosharding/shard"
	"log"
	"math/rand"
	"time"
)

func main() {
	cfg := config.LoadDBConfig(".")
	dbs := config.ConnDBs(cfg, cfg.ShardCount)
	repo := repository.NewRepository(dbs, repository.NewUser())
	shard := shard.NewShard(cfg.ShardCount)
	svc := service.NewService(repo, shard)
	//svr := api.NewServer(dbs, svc)
	defer func() {
		for i := range svc.Repo.DBs {
			svc.Repo.DBs[i].Close()
		}
	}()

	/*	id, _ := uuid.NewV7()
		err := svc.Create(context.Background(), id.String(), "asdf")
		if err != nil {
			log.Println(err)
		}
		log.Println(randomString(10))
		log.Println(randomString(10))
		log.Println(randomString(10))
		log.Println(randomString(10))*/

	/*var wg sync.WaitGroup
	start := time.Now()
	for range 100000 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			num := rand.Int63n(30)
			name := randomString(int(num))
			newID, _ := uuid.NewV7()

			err := svc.Create(context.Background(), newID.String(), name)
			if err != nil {
				log.Println(err)
			}
		}()
	}
	wg.Wait()
	end := time.Since(start)
	log.Println("done")
	log.Println(end)*/

	//user, err := svc.Get(context.Background(), "01911664-0412-7374-80e5-ba4b82be2e2e")
	////user, err := svc.Get(context.Background(), "01911664-0412-7374-80e5-ba4b82be2e21")
	////user, err := svc.Get222(context.Background(), "01911664-0412-7374-80e5-ba4b82be2e2e")
	//if err != nil {
	//	log.Println(err)
	//}
	//log.Println(user)

	start := time.Date(2024, time.August, 3, 3, 59, 0, 0, time.UTC)
	end := time.Date(2024, time.August, 3, 3, 59, 10, 0, time.UTC)
	users, err := svc.ListRangeCreatedAt(context.Background(), start, end)
	if err != nil {
		log.Println(err)
	}
	log.Println(len(users))
}

func randomString(num int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := ""
	for range num {
		idx := rand.Int63n(int64(len(str)))
		result = result + string(str[idx])
	}
	return result
}
