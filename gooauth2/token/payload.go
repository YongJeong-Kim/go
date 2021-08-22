package token

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid.")
	ErrExpiredToken = errors.New("token has expired.")
)

type AccessTokenPayload struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
	AccessUUID  string    `json:"access_uuid"`
	RefreshUUID string    `json:"refresh_uuid"`
	RtExpiredAt time.Time `json:"rt_expired_at"`
}

func NewAccessToken(username string, atDuration time.Duration, rtDuration time.Duration) (*AccessTokenPayload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("generate token id failed. %s ", err.Error())
	}

	return &AccessTokenPayload{
		ID:          tokenID,
		Username:    username,
		IssuedAt:    time.Now(),
		ExpiredAt:   time.Now().Add(atDuration),
		AccessUUID:  uuid.NewString(),
		RefreshUUID: uuid.NewString(),
		RtExpiredAt: time.Now().Add(rtDuration),
	}, nil
}

type RefreshTokenPayload struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
	RefreshUUID string    `json:"refresh_uuid"`
}

type Payload struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	IssuedAt    time.Time `json:"issued_at"`
	ExpiredAt   time.Time `json:"expired_at"`
	AccessUUID  string    `json:"access_uuid"`
	RefreshUUID string    `json:"refresh_uuid"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:          tokenID,
		Username:    username,
		IssuedAt:    time.Now(),
		ExpiredAt:   time.Now().Add(duration),
		AccessUUID:  uuid.NewString(),
		RefreshUUID: uuid.NewString(),
	}
	return payload, nil
}

func ExpiredValid(payload Payload) error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
