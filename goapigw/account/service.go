package main

import (
	"errors"
	"time"
)

const (
	Username = "aaa"
	Password = "1234"
)

type AccountServicer interface {
	Login(username, password string, duration time.Duration) (string, error)
}

func (s *AccountServer) Login(username, password string, duration time.Duration) (string, error) {
	if username != Username || password != Password {
		return "", errors.New("invalid username or password")
	}

	token, err := s.Maker.Create(username, duration)
	if err != nil {
		return "", err
	}

	return token, nil
}
