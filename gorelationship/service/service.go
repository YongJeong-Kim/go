package service

import (
	"context"
	"gorelationship/config"
	"gorelationship/repository"
)

type Service struct {
	Friend Friender
	User   UserManager
}

func NewService(friend Friender, user UserManager) *Service {
	return &Service{
		Friend: friend,
		User:   user,
	}
}

func newTestService() (context.Context, *Service) {
	ctx := context.Background()
	driver := config.NewDriver(ctx, config.URI, config.Username, config.Password)
	sess := config.NewSession(ctx, driver, config.DatabaseName)
	var ru repository.UserManager = repository.NewUser(sess)
	var rf repository.Friender = repository.NewFriend(sess)
	var sf Friender = NewFriend(rf)
	var su UserManager = NewUser(ru)

	return ctx, NewService(sf, su)
}
