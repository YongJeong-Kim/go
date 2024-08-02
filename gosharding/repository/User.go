package repository

import (
	"context"
	"errors"
)

type UserManager interface {
	Create(ctx context.Context, tx *Tx, id, name string) error
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

type User struct {
}

func NewUser() *User {
	return &User{}
}
