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

type AccessTokenPayload struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
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

	now := time.Now()

	return &Payload{
		AccessTokenPayload: AccessTokenPayload{
			ID:        atID,
			UserID:    uuid.NewString(),
			Username:  username,
			IssuedAt:  now,
			ExpiredAt: now.Add(duration.AccessTokenDuration),
		},
		RefreshTokenPayload: RefreshTokenPayload{
			ID:        rtID,
			UserID:    uuid.NewString(),
			Username:  username,
			IssuedAt:  now,
			ExpiredAt: now.Add(duration.RefreshTokenDuration),
		},
	}, nil
}

type RefreshTokenPayload struct {
	ID        uuid.UUID `json:"id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type Payload struct {
	AccessTokenPayload  AccessTokenPayload
	RefreshTokenPayload RefreshTokenPayload
}

type PayloadDetails struct {
	Payload *Payload
	Token   map[string]string
}

// type Payload struct {
// 	ID          uuid.UUID `json:"id"`
// 	Username    string    `json:"username"`
// 	IssuedAt    time.Time `json:"issued_at"`
// 	ExpiredAt   time.Time `json:"expired_at"`
// 	AccessUUID  string    `json:"access_uuid"`
// 	RefreshUUID string    `json:"refresh_uuid"`
// }

// func NewPayload(username string, duration time.Duration) (*Payload, error) {
// 	tokenID, err := uuid.NewRandom()
// 	if err != nil {
// 		return nil, err
// 	}

// 	payload := &Payload{
// 		ID:          tokenID,
// 		Username:    username,
// 		IssuedAt:    time.Now(),
// 		ExpiredAt:   time.Now().Add(duration),
// 		AccessUUID:  uuid.NewString(),
// 		RefreshUUID: uuid.NewString(),
// 	}
// 	return payload, nil
// }

// func ExpiredValid(payload Payload) error {
// 	if time.Now().After(payload.ExpiredAt) {
// 		return ErrExpiredToken
// 	}
// 	return nil
// }
