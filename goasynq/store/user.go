package store

import (
	"context"
	"fmt"
)

type User struct {
	Name string `db:"name" json:"name"`
}

func (q *Queries) CreateUser(ctx context.Context, name string) error {
	query := `INSERT INTO user(name) VALUES(?)`
	_, err := q.ExecContext(ctx, query, name)
	if err != nil {
		return fmt.Errorf("create user failed: %v", err)
	}

	return nil
}
