package service

import (
	"fmt"
	"github.com/YongJeong-Kim/go/goasynq/store"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	DB *sqlx.DB
}

type Servicer interface {
	CreateUser(*CreateUserParam) error
}

func NewService(db *sqlx.DB) Servicer {
	return &Service{
		DB: db,
	}
}

func (s *Service) execTx(fn func(store.Store) error) error {
	tx := s.DB.MustBegin()
	q := store.NewQueries(tx)
	err := fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
