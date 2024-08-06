package service

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
	"gosharding/config"
	"gosharding/repository"
	"gosharding/shard"
	"log"
	"testing"
	"time"
)

func TestGet1(t *testing.T) {
	defer goleak.VerifyNone(t)

	cfg := config.LoadDBConfig("../")
	dbs := config.ConnDBs(cfg, cfg.ShardCount)
	repo := repository.NewRepository(dbs, repository.NewUser())
	shard := shard.NewShard(cfg.ShardCount)
	svc := NewService(repo, shard)
	defer func() {
		for i := range svc.Repo.DBs {
			svc.Repo.DBs[i].Close()
		}
	}()

	start := time.Date(2024, time.August, 3, 3, 59, 0, 0, time.UTC)
	end := time.Date(2024, time.August, 3, 3, 59, 10, 0, time.UTC)
	users, err := svc.ListRangeCreatedAt(context.Background(), start, end)
	require.NoError(t, err)
	log.Println(len(users))
}
