package main

import (
	"aidanwoods.dev/go-paseto"
	"fmt"
	"time"
)

const (
	Issuer = "test.gw"
)

type TokenVerifier interface {
	//CreateToken(userID string, duration time.Duration) (string, error)
	Verify(token string) (*Payload, error)
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

/*func (p *PasetoVerifier) CreateToken(userID string, duration time.Duration) (string, error) {
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
}*/

func (p *PasetoVerifier) Verify(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy(Issuer))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now().UTC()))

	secretKeyHex, err := p.SecretKeyFromHex()
	if err != nil {
		return nil, fmt.Errorf("secret key from hex: %v", err)
	}

	parsed, err := parser.ParseV4Local(secretKeyHex, token, nil)
	if err != nil {
		return nil, fmt.Errorf("parse failed: %v", err)
	}

	claims := parsed.Claims()
	payload, err := Map2Struct(&Payload{}, claims)
	if err != nil {
		return nil, fmt.Errorf("convert payload failed. %v", err)
	}

	return payload, nil
}

func NewPasetoVerifier(secretKeyHex string) TokenVerifier {
	return &PasetoVerifier{
		secretKeyHex: secretKeyHex,
	}
}

func (p *PasetoVerifier) SecretKeyFromHex() (paseto.V4SymmetricKey, error) {
	return paseto.V4SymmetricKeyFromHex(p.secretKeyHex)
}
