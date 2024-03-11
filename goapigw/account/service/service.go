package service

import (
	"errors"
	"github.com/yongjeong-kim/go/goapigw/account/token"
	"time"
)

const (
	Username = "aaa"
	Password = "1234"
)

type AccountServicer interface {
	Login(username, password string, duration time.Duration) (string, error)
}

type AccountService struct {
	Maker token.TokenMaker
}

func NewAccountService(maker token.TokenMaker) *AccountService {
	return &AccountService{
		Maker: maker,
	}
}

func (s *AccountService) Login(username, password string, duration time.Duration) (string, error) {
	if username != Username || password != Password {
		return "", errors.New("invalid username or password")
	}

	tk, err := s.Maker.Create(username, duration)
	if err != nil {
		return "", err
	}

	return tk, nil
}
