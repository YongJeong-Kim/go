package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	DBs  []*sqlx.DB
	User UserManager
}

func NewRepository(dbs []*sqlx.DB, user UserManager) *Repository {
	return &Repository{
		DBs:  dbs,
		User: user,
	}
}

type Tx struct {
	*sqlx.Tx
}

func NewTx(tx *sqlx.Tx) *Tx {
	return &Tx{
		tx,
	}
}
