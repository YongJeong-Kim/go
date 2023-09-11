package service

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	mysqlSource = "root:1234@tcp(localhost:13306)/test?parseTime=true"
)

func newTestService(t *testing.T) Servicer {
	conn, err := sqlx.Connect("mysql", mysqlSource)
	require.NoError(t, err)

	return &Service{
		DB: conn,
	}
}
