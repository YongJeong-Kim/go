package token

import (
	"encoding/json"
	"fmt"
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

func (maker *JWTMaker) VerifyToken(payload string) (*Payload, error) {
	panic("implement me")
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
