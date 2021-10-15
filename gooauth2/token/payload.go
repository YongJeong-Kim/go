package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	AccessTokenPayload  AccessTokenPayload
	RefreshTokenPayload RefreshTokenPayload
}

type PayloadDetails struct {
	Payload *Payload
	Token   map[string]string
}

type AccessTokenPayload struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (atp AccessTokenPayload) Valid() error {
	if time.Now().After(atp.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

type RefreshTokenPayload struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (rtp RefreshTokenPayload) Valid() error {
	if time.Now().After(rtp.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func NewPayload(username string, duration JWTDuration) (*Payload, error) {
	atID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("generate access token id failed. %s ", err.Error())
	}

	rtID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token id failed. %s ", err.Error())
	}

	userID := uuid.NewString()
	now := time.Now()

	return &Payload{
		AccessTokenPayload: AccessTokenPayload{
			ID:        atID,
			UserID:    userID,
			Username:  username,
			IssuedAt:  now,
			ExpiredAt: now.Add(duration.AccessTokenDuration),
		},
		RefreshTokenPayload: RefreshTokenPayload{
			ID:        rtID,
			UserID:    userID,
			Username:  username,
			IssuedAt:  now,
			ExpiredAt: now.Add(duration.RefreshTokenDuration),
		},
	}, nil
}
