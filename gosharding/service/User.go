package service

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gosharding/repository"
	"time"
)

func (s *Service) Create(ctx context.Context, id, name string) error {
	idx := s.Shard.Index(id)
	err := s.execTx(s.Repo.DBs[idx], func(tx *repository.Tx) error {
		err := s.Repo.User.Create(ctx, tx, id, name)
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

func (s *Service) Get(ctx context.Context, id string) (*repository.GetResult, error) {
	idx := s.Shard.Index(id)
	user, err := execTxReturn(s.Repo.DBs[idx], func(tx *repository.Tx) (*repository.GetResult, error) {
		return s.Repo.User.Get(ctx, tx, id)
	})
	if err != nil {
		return nil, err
	}
	return *user, nil
}

func (s *Service) ListRangeCreatedAt(ctx context.Context, start, end time.Time) ([]*repository.ListRangeCreatedAtResult, error) {
	rr, err := multiSelect(s.Repo.DBs, func(db *sqlx.DB) ([]*repository.ListRangeCreatedAtResult, error) {
		r, err := execTxReturn(db, func(tx *repository.Tx) ([]*repository.ListRangeCreatedAtResult, error) {
			return s.Repo.User.ListRangeCreatedAt(ctx, tx, start, end)
		})
		if err != nil {
			return nil, err
		}

		return *r, nil
	})
	if err != nil {
		return nil, err
	}

	return rr, nil
}
