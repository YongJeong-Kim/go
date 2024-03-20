package token

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
	"github.com/yongjeong-kim/go/goapigw/gw/util"
	"time"
)

const (
	Issuer = "test.gw"
	KeyHex = "ec6208ce3e155794d647719db39c645d9fa474a75b9515ce7571c3009428be55"
)

type TokenVerifier interface {
	Verify(accessToken string) (*Payload, error)
}

type Payload struct {
	UserID    string    `json:"user_id"`
	IssuedAt  time.Time `json:"iat"`
	ExpiredAt time.Time `json:"exp"`
	Issuer    string    `json:"iss"`
	NotBefore time.Time `json:"nbf"`
}

type PasetoVerifier struct {
	secretKeyHex string
}

func (p *PasetoVerifier) Verify(accessToken string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy(Issuer))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now().UTC()))

	secretKeyHex, err := paseto.V4SymmetricKeyFromHex(p.secretKeyHex)
	if err != nil {
		return nil, fmt.Errorf("secret key from hex: %v", err)
	}

	parsed, err := parser.ParseV4Local(secretKeyHex, accessToken, nil)
	if err != nil {
		return nil, fmt.Errorf("parse failed: %v", err)
	}

	claims := parsed.Claims()
	payload, err := util.Map2Struct(&Payload{}, claims)
	if err != nil {
		return nil, fmt.Errorf("convert payload failed. %v", err)
	}

	if payload.UserID == "" {
		return nil, fmt.Errorf("cannot empty user id")
	}

	return payload, nil
}

func NewPasetoVerifier(secretKeyHex string) *PasetoVerifier {
	return &PasetoVerifier{
		secretKeyHex: secretKeyHex,
	}
}
