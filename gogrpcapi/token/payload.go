package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	UserID string `json:"user_id"`
	Issuer string `json:"iss"`
	//IssuedAt  time.Time `json:"iat"`
	IssuedAt int64 `json:"iat"`
	//ExpiredAt time.Time `json:"exp"`
	ExpiredAt int64 `json:"exp"`
	jwt.RegisteredClaims
}

func NewPayload(userID string, phoneNumber string, duration time.Duration) (*Payload, error) {
	//payload := &Payload{
	//	UserID:    userID,
	//	Issuer:    phoneNumber,
	//	IssuedAt:  time.Now().Unix(),
	//	ExpiredAt: time.Now().Add(duration).Unix(),
	//}
	now := time.Now().UTC()
	payload := &Payload{
		UserID: userID,
		Issuer: phoneNumber,
		//IssuedAt:  time.Now().UTC(),
		IssuedAt:  now.Unix(),
		ExpiredAt: now.Add(duration).Unix(),
	}
	return payload, nil
}

func (payload *Payload) Valid() error {
	//expiredAt := time.Unix(0, payload.ExpiredAt)

	now := time.Now().UTC().Unix()
	if now > payload.ExpiredAt {
		return ErrExpiredToken
	}
	//if time.Now().After(payload.ExpiredAt) {
	//	return ErrExpiredToken
	//}
	return nil
}
