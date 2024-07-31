package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"gorelationship/repository"
)

type UserManager interface {
	Create(ctx context.Context, name string) (string, error)
	Get(ctx context.Context, userID string) (*repository.GetResult, error)
}

func (u *User) Create(ctx context.Context, name string) (string, error) {
	createdID, err := neo4j.ExecuteWrite(ctx, u.Sess, func(tx neo4j.ManagedTransaction) (string, error) {
		r, err := u.User.Create(ctx, tx, name)
		if err != nil {
			return "", err
		}

		return r, nil
	})
	if err != nil {
		return "", err
	}

	return createdID, nil
}

func (u *User) Get(ctx context.Context, userID string) (*repository.GetResult, error) {
	if err := uuid.Validate(userID); err != nil {
		return nil, errors.New("invalid user uuid")
	}

	user, err := neo4j.ExecuteRead(ctx, u.Sess, func(tx neo4j.ManagedTransaction) (*repository.GetResult, error) {
		user, err := u.User.Get(ctx, tx, userID)
		if err != nil {
			return nil, err
		}

		return user, nil
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

type User struct {
	Sess neo4j.SessionWithContext
	User repository.UserManager
}

func NewUser(sess neo4j.SessionWithContext, userManager repository.UserManager) *User {
	return &User{
		Sess: sess,
		User: userManager,
	}
}
