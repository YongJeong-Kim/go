package token

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
	"time"
)

const (
	Issuer = "test.gw"
	keyHex = "ec6208ce3e155794d647719db39c645d9fa474a75b9515ce7571c3009428be55"
)

type TokenMaker interface {
	Create(userID string, duration time.Duration) (string, error)
}

type PasetoMaker struct {
	secretKeyHex string
}

func NewPasetoMaker() *PasetoMaker {
	return &PasetoMaker{
		secretKeyHex: keyHex,
	}
}

func (p *PasetoMaker) Create(userID string, duration time.Duration) (string, error) {
	tk := paseto.NewToken()

	now := time.Now().UTC()
	tk.SetExpiration(now.Add(duration))
	tk.SetIssuedAt(now)
	tk.SetNotBefore(now)
	tk.SetIssuer(Issuer)
	tk.SetString("user_id", userID)

	secretKey, err := p.SecretKeyFromHex()
	if err != nil {
		return "", fmt.Errorf("invalid hex: %v", err)
	}

	return tk.V4Encrypt(secretKey, nil), nil
}

func (p *PasetoMaker) SecretKeyFromHex() (paseto.V4SymmetricKey, error) {
	return paseto.V4SymmetricKeyFromHex(p.secretKeyHex)
}
