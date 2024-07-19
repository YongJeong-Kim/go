package service

import (
	"context"
	"gorelationship/repository"
)

type UserManager interface {
	Create(ctx context.Context, name string) (string, error)
	Get(ctx context.Context, id string) (*repository.GetResult, error)
}

func (u *User) Create(ctx context.Context, name string) (string, error) {
	return u.User.Create(ctx, name)
}

func (u *User) Get(ctx context.Context, id string) (*repository.GetResult, error) {
	return u.User.Get(ctx, id)
}

type User struct {
	User repository.UserManager
}

func NewUser(userManager repository.UserManager) *User {
	return &User{
		User: userManager,
	}
}
