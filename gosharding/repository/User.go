package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserManager interface {
	Create(ctx context.Context, tx *Tx, id, name string) error
	Get(ctx context.Context, tx *Tx, id string) (*GetResult, error)
	ListRangeCreatedAt(ctx context.Context, tx *Tx, start, end time.Time) ([]*ListRangeCreatedAtResult, error)
}

func (u *User) Create(ctx context.Context, tx *Tx, id, name string) error {
	result, err := tx.ExecContext(ctx, "INSERT INTO user(id, name) VALUES (uuid_to_bin(?), ?)", id, name)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("insert user failed")
	}
	return nil
}

type GetResult struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (u *User) Get(ctx context.Context, tx *Tx, id string) (*GetResult, error) {
	q := "SELECT bin_to_uuid(id) AS id, name, created_at FROM user WHERE id=uuid_to_bin(?)"
	var result GetResult
	err := tx.GetContext(ctx, &result, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found. %s", id)
		}
		return nil, err
	}
	return &result, nil
}

type ListRangeCreatedAtResult struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

func (u *User) ListRangeCreatedAt(ctx context.Context, tx *Tx, start, end time.Time) ([]*ListRangeCreatedAtResult, error) {
	q := "SELECT bin_to_uuid(id) AS id, name, created_at FROM user WHERE created_at BETWEEN ? AND ?"
	var result []*ListRangeCreatedAtResult
	err := tx.SelectContext(ctx, &result, q, start, end)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type User struct {
}

func NewUser() *User {
	return &User{}
}
