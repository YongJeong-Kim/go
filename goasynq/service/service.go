package service

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"goasynq/store"
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

type CreateUserParam struct {
	Name  string
	After func(name string) error
}

func (s *Service) CreateUser(param *CreateUserParam) error {
	err := s.execTx(func(q store.Store) error {
		err := q.CreateUser(context.Background(), param.Name)
		if err != nil {
			return fmt.Errorf("service create user failed: %v", err)
		}

		return param.After(param.Name)
	})
	if err != nil {
		return err
	}

	return nil
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
