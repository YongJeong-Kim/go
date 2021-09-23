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
			return nil, fmt.Errorf("invalid token. ")
		}
		return []byte(maker.accessSecret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(accessToken, jwt.MapClaims{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errors.New("token has expired")) {
			return nil, errors.New("token has expired")
		}
		return nil, errors.New("invalid token")
	}

	claims, err := json.Marshal(jwtToken.Claims)
	if err != nil {
		return nil, fmt.Errorf("marshal failed. claims : %s ", err.Error())
	}
	atp := AccessTokenPayload{}
	err = json.Unmarshal(claims, &atp)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal verify token claims : %s ", err.Error())
	}

	return &atp, nil
}

func (maker *JWTMaker) VerifyRefreshToken(refreshToken string) (*RefreshTokenPayload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token. ")
		}
		return []byte(maker.refreshSecret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(refreshToken, jwt.MapClaims{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, errors.New("token has expired")) {
			return nil, errors.New("token has expired")
		}
		return nil, errors.New("invalid token")
	}

	claims, err := json.Marshal(jwtToken.Claims)
	if err != nil {
		return nil, fmt.Errorf("marshal failed. claims : %s ", err.Error())
	}
	rtp := RefreshTokenPayload{}
	err = json.Unmarshal(claims, &rtp)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal verify token claims : %s ", err.Error())
	}

	return &rtp, nil
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

/*type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}*/

//func CreateToken(userID uint64) (*TokenDetails, error) {
//	td := &TokenDetails{}
//	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
//	td.AccessUUID = uuid.NewString()
//	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
//	td.RefreshUUID = uuid.NewString()
//
//	var err error
//
//	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
//	atClaims := jwt.MapClaims{}
//	atClaims["authorized"] = true
//	atClaims["access_uuid"] = td.AccessUUID
//	atClaims["user_id"] = userID
//	atClaims["exp"] = td.AtExpires
//	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
//	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
//	if err != nil {
//		return nil, err
//	}
//	err = os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf")
//	if err != nil {
//		return nil, err
//	}
//
//	rtClaims := jwt.MapClaims{}
//	rtClaims["refresh_uuid"] = td.RefreshUUID
//	rtClaims["user_id"] = userID
//	rtClaims["exp"] = td.RtExpires
//	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
//	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
//	if err != nil {
//		return nil, err
//	}
//	return td, nil
//}
