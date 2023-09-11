package store

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

func newTestStore(t *testing.T) (*sqlx.Tx, Store) {
	conn, err := sqlx.Connect("mysql", "root:1234@tcp(localhost:13306)/test?parseTime=true")
	require.NoError(t, err)

	tx := conn.MustBegin()
	q := NewQueries(tx)

	return tx, q
}
