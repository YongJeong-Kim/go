package token

import (
	"github.com/golang-jwt/jwt/v5"
	"mingle/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomAlphabet(43))
	require.NoError(t, err)

	userID := util.CreateUUID()
	phoneNumber := util.RandomPhoneNumber()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(userID, phoneNumber, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.Equal(t, userID, payload.UserID)
	require.Equal(t, phoneNumber, payload.Issuer)

	iat := time.Unix(payload.IssuedAt, 0)
	exp := time.Unix(payload.ExpiredAt, 0)
	require.WithinDuration(t, issuedAt, iat, time.Second)
	require.WithinDuration(t, expiredAt, exp, time.Second)
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomAlphabet(43))
	require.NoError(t, err)

	token, err := maker.CreateToken(util.RandomAlphabet(6), util.RandomPhoneNumber(), -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomAlphabet(6), util.RandomPhoneNumber(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(43))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
