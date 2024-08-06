package service

import (
	"github.com/jmoiron/sqlx"
	"gosharding/config"
	"gosharding/repository"
	"gosharding/shard"
	"log"
	"sync"
)

func multiSelect[T any](dbs []*sqlx.DB, fnTx func(*sqlx.DB) ([]T, error)) ([]T, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	result := make([]T, 0)
	errs := make([]error, 0)
	errCh := make(chan error)
	defer close(errCh)

	for _, db := range dbs {
		wg.Add(1)
		go func(db *sqlx.DB) {
			defer wg.Done()
			ts, errTx := fnTx(db)
			switch errTx {
			case nil:
				mu.Lock()
				result = append(result, ts...)
				mu.Unlock()
				errCh <- nil
			default:
				errCh <- errTx
			}
		}(db)

		select {
		case e := <-errCh:
			if e != nil {
				errs = append(errs, e)
			}
		}
	}
	wg.Wait()
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return result, nil
}

func (s *Service) execTx(db *sqlx.DB, fn func(*repository.Tx) error) error {
	tx := db.MustBegin()
	repo := repository.NewTx(tx)
	err := fn(repo)
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			log.Println("Rollback failed:", err)
			return errRb
		}
		return err
	}
	return tx.Commit()
}

func execTxReturn[T any](db *sqlx.DB, fn func(*repository.Tx) (T, error)) (*T, error) {
	tx := db.MustBegin()
	repo := repository.NewTx(tx)
	result, err := fn(repo)
	if err != nil {
		errRb := tx.Rollback()
		if errRb != nil {
			log.Println("Rollback failed:", err)
			return nil, errRb
		}
		return nil, err
	}

	return &result, tx.Commit()
}

type Service struct {
	Repo  *repository.Repository
	Shard *shard.Shard
}

func NewService(repo *repository.Repository, shard *shard.Shard) *Service {
	return &Service{
		Repo:  repo,
		Shard: shard,
	}
}

func newTestService() *Service {
	cfg := config.LoadDBConfig("../")
	dbs := config.ConnDBs(cfg, cfg.ShardCount)
	repo := repository.NewRepository(dbs, repository.NewUser())
	shd := shard.NewShard(cfg.ShardCount)

	return &Service{
		Repo:  repo,
		Shard: shd,
	}
}
