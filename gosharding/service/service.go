package service

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gosharding/repository"
	"gosharding/shard"
	"log"
)

func (s *Service) Create(ctx context.Context, id, name string) error {
	idx := s.Shard.Index(id)
	err := s.ExecTx(s.Repo.DBs[idx], func(tx *repository.Tx) error {
		err := s.Repo.User.Create(context.Background(), tx, id, name)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ExecTx(db *sqlx.DB, fn func(*repository.Tx) error) error {
	tx := db.MustBegin()
	repo := repository.NewTx(tx)
	err := fn(repo)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			log.Println("Rollback failed:", err)
			return err
		}
	}
	return tx.Commit()
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
