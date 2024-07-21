package repository

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

type UserManager interface {
	Create(ctx context.Context, name string) (string, error)
	Get(ctx context.Context, id string) (*GetResult, error)
}

func (u *User) Create(ctx context.Context, name string) (string, error) {
	id, err := u.sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			CREATE (u:User {id: randomuuid(), name: $name, createdDate: datetime()})
			RETURN u.id AS id
		`, map[string]any{
			"name": name,
		})
		if err != nil {
			return "", err
		}

		user, err := result.Single(ctx)
		if err != nil {
			return "", err
		}

		return user.AsMap()["id"].(string), nil
	})
	if err != nil {
		return "", err
	}

	return id.(string), nil
}

type GetResult struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"created_date"`
}

func (u *User) Get(ctx context.Context, userID string) (*GetResult, error) {
	user, err := u.sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (u:User) WHERE u.id = $userID
			RETURN u.id AS id, u.name AS name, u.createDate AS createdDate
		`, map[string]any{
			"userID": userID,
		})
		if err != nil {
			return nil, err
		}

		user, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		return &GetResult{
			ID:          user.AsMap()["id"].(string),
			Name:        user.AsMap()["name"].(string),
			CreatedDate: user.AsMap()["createdDate"].(time.Time),
		}, nil
	})
	if err != nil {
		return nil, err
	}

	return user.(*GetResult), nil
}

type User struct {
	sess neo4j.SessionWithContext
}

func NewUser(sess neo4j.SessionWithContext) *User {
	return &User{
		sess: sess,
	}
}
