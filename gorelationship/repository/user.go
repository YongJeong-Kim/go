package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type UserManager interface {
	Create(ctx context.Context, param map[string]any) (string, error)
}

func (u *User) Create(ctx context.Context, param map[string]any) (string, error) {
	uID, err := u.sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			CREATE (u:User {id: randomuuid(), name: $name, createdDate: datetime()})
			RETURN u.id AS id
		`, param)
		if err != nil {
			return "", err
		}

		u, err := result.Single(ctx)
		if err != nil {
			return "", err
		}
		uID, _ := u.AsMap()["id"]
		return uID.(string), nil
	})
	if err != nil {
		return "", err
	}

	return uID.(string), nil
}

type User struct {
	sess neo4j.SessionWithContext
}

func NewUser(sess neo4j.SessionWithContext) *User {
	return &User{
		sess: sess,
	}
}
