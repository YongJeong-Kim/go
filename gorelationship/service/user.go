package service

import (
	"context"
	"gorelationship/repository"
)

type UserManager interface {
	Create(ctx context.Context, param map[string]any) (string, error)
}

func (u *User) Create(ctx context.Context, param map[string]any) (string, error) {
	return u.UserManager.Create(ctx, param)
}

type User struct {
	UserManager repository.UserManager
}

func NewUser(userManager repository.UserManager) *User {
	return &User{
		UserManager: userManager,
	}
}
