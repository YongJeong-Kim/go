package token

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewJWTMaker(t *testing.T) {
	accessSecret := uuid.NewString()
	refreshSecret := uuid.NewString()
	_, err := NewJWTMaker(accessSecret, refreshSecret)
	require.NoError(t, err)
}

func TestCreateToken(t *testing.T) {
	accessSecret := uuid.NewString()
	refreshSecret := uuid.NewString()
	maker, err := NewJWTMaker(accessSecret, refreshSecret)
	require.NoError(t, err)

	duration := JWTDuration{
		AccessTokenDuration:  time.Minute,
		RefreshTokenDuration: time.Minute,
	}
	payloadDetails, err := maker.CreateToken("random username", duration)
	require.NoError(t, err)
	require.NotEmpty(t, payloadDetails.Token["access_token"])
	require.NotEmpty(t, payloadDetails.Token["refresh_token"])
}

func TestVerifyToken(t *testing.T) {
	accessSecret := uuid.NewString()
	refreshSecret := uuid.NewString()
	maker, err := NewJWTMaker(accessSecret, refreshSecret)
	require.NoError(t, err)

	duration := JWTDuration{
		AccessTokenDuration:  time.Minute,
		RefreshTokenDuration: time.Minute,
	}
	payloadDetails, err := maker.CreateToken("random username", duration)
	require.NoError(t, err)
	require.NotEmpty(t, payloadDetails.Token["access_token"])
	require.NotEmpty(t, payloadDetails.Token["refresh_token"])

	_, err = maker.VerifyAccessToken(payloadDetails.Token["access_token"])
	require.NoError(t, err)

	_, err = maker.VerifyRefreshToken(payloadDetails.Token["refresh_token"])
	require.NoError(t, err)
}

func TestExpiredToken(t *testing.T) {
	accessSecret := uuid.NewString()
	refreshSecret := uuid.NewString()
	maker, err := NewJWTMaker(accessSecret, refreshSecret)
	require.NoError(t, err)

	duration := JWTDuration{
		AccessTokenDuration:  -time.Minute,
		RefreshTokenDuration: -time.Minute,
	}
	payloadDetails, err := maker.CreateToken("random username", duration)
	require.NoError(t, err)

	accessPayload, err := maker.VerifyAccessToken(payloadDetails.Token["access_token"])
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, accessPayload)

	refreshPayload, err := maker.VerifyRefreshToken(payloadDetails.Token["refresh_token"])
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, refreshPayload)
}
