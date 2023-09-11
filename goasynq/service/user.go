package service

import (
	"context"
	"fmt"
	"github.com/YongJeong-Kim/go/goasynq/store"
)

type CreateUserParam struct {
	Name  string
	After func(name string) error
}

func (s *Service) CreateUser(param *CreateUserParam) error {
	err := s.execTx(func(q store.Store) error {
		err := q.CreateUser(context.Background(), param.Name)
		if err != nil {
			return fmt.Errorf("service create user failed: %v", err)
		}

		return param.After(param.Name)
	})
	if err != nil {
		return err
	}

	return nil
}
