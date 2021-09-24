package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTMaker struct {
	accessSecret  string
	refreshSecret string
}

type JWTDuration struct {
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

func (maker *JWTMaker) CreateToken(username string, duration JWTDuration) (*PayloadDetails, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return nil, fmt.Errorf("invalid payload : %s ", err.Error())
	}

	atClaims := jwt.MapClaims{}
	atc, err := json.Marshal(payload.AccessTokenPayload)
	if err != nil {
		return nil, fmt.Errorf("marshal failed. access token : %s ", err.Error())
	}
	err = json.Unmarshal(atc, &atClaims)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal access token claims : %s ", err.Error())
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := at.SignedString([]byte(maker.accessSecret))
	if err != nil {
		return nil, fmt.Errorf("cannot signed string access token: %s ", err.Error())
	}

	rtClaims := jwt.MapClaims{}
	rtc, err := json.Marshal(payload.RefreshTokenPayload)
	if err != nil {
		return nil, fmt.Errorf("marshal failed. refresh token : %s ", err.Error())
	}
	err = json.Unmarshal(rtc, &rtClaims)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal refresh token claims : %s ", err.Error())
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := rt.SignedString([]byte(maker.refreshSecret))
	if err != nil {
		return nil, fmt.Errorf("cannot signed string refresh token: %s ", err.Error())
	}

	return &PayloadDetails{
		Payload: payload,
		Token: map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	}, nil
}

func (maker *JWTMaker) ExtractToken(authorizationHeader string) (string, error) {
	if len(authorizationHeader) == 0 {
		err := errors.New("authorization header is not provided")
		return "", err
	}

	fields := strings.Fields(authorizationHeader)
	if len(fields) < 2 {
		err := errors.New("invalid authorization header format")
		return "", err
	}

	authorizationType := strings.ToLower(fields[0])
	if authorizationType != "bearer" {
		err := fmt.Errorf("unsupported authorization type %s", authorizationType)
		return "", err
	}

	accessToken := fields[1]
	return accessToken, nil
}

func (maker *JWTMaker) VerifyAccessToken(accessToken string) (*AccessTokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.accessSecret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(accessToken, &AccessTokenPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*AccessTokenPayload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

func (maker *JWTMaker) VerifyRefreshToken(refreshToken string) (*RefreshTokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.refreshSecret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenPayload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*RefreshTokenPayload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}

const minSecretKeySize = 32

func NewJWTMaker(accessSecret string, refreshSecret string) (Maker, error) {
	if len(accessSecret) < minSecretKeySize && len(refreshSecret) < minSecretKeySize {
		return nil, fmt.Errorf("must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{
		accessSecret,
		refreshSecret,
	}, nil
}
