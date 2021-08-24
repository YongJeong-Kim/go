package token

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTMaker struct {
	accessSecret  string
	refreshSecret string
}

func (maker *JWTMaker) CreateToken(username string, atDur time.Duration, rtDur time.Duration) (interface{}, error) {
	payload, err := NewPayload(username, atDur, rtDur)
	if err != nil {
		return nil, fmt.Errorf("invalid payload : %s ", err.Error())
	}

	// td := &AccessTokenPayload{}

	// var err error

	// os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	// atClaims := jwt.MapClaims{}
	// atClaims["authorized"] = true
	// atClaims["access_uuid"] = td.AccessUUID
	// atClaims["user_id"] = userID
	// atClaims["exp"] = td.AtExpires
	// at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	// td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))

	atClaims := jwt.MapClaims{}
	atc, _ := json.Marshal(payload.accessTokenPayload)
	json.Unmarshal(atc, &atClaims)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	at.SignedString([]byte(maker.accessSecret))

	// if err != nil {
	// 	return nil, err
	// }
	// err = os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf")
	// if err != nil {
	// 	return nil, err
	// }

	rtClaims := jwt.MapClaims{}
	rtc, _ := json.Marshal(payload.refreshTokenPayload)
	json.Unmarshal(rtc, &rtClaims)
	// rtClaims["refresh_uuid"] = td.RefreshUUID
	// rtClaims["user_id"] = userID
	// rtClaims["exp"] = td.RtExpires
	// rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClamis)
	rt.SignedString([]byte(maker.refreshSecret))
	// payload.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	// if err != nil {
	// 	return nil, err
	// }
	return {
		access_token: at,
		refresh_token: rt,
	}, nil
	// return td, nil
}

func (maker *JWTMaker) VerifyToken(payload string) (*Payload, error) {
	panic("implement me")
}

const minSecretKeySize = 32

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

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
