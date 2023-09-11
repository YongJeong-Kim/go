package store

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Queries struct {
	*sqlx.Tx
}

type Store interface {
	CreateUser(ctx context.Context, name string) error
}

func NewQueries(tx *sqlx.Tx) Store {
	return &Queries{
		tx,
	}
}
