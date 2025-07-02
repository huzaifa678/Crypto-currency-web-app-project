package token

import (
	"testing"
	"time"

	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoMakerWithDifferentTokenTypes(t *testing.T) {
	secretKey := utils.RandomString(32)
	maker, err := NewPasetoMaker(secretKey)

	require.NoError(t, err)

	username := utils.RandomString(12)
	duration := time.Minute

	token, payload, err := maker.CreateToken(username, duration, TokenTypeAccessToken)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	payload, err = maker.VerifyToken(token, TokenTypeAccessToken)
	require.NoError(t, err)

	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)

	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestPasetoTokenExpired(t *testing.T) {
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)
	username := utils.RandomString(12)

	token, payload, err := maker.CreateToken(username, -time.Minute, TokenTypeAccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)

	payload, err = maker.VerifyToken(token, TokenTypeAccessToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidSymmetricKey(t *testing.T) {

	invalidSymmetricKey := utils.RandomString(10)
	_, err := NewPasetoMaker(invalidSymmetricKey)

	require.EqualError(t, err, "invalid key size: must be exactly 32 characters")
}

func TestInvalidPasetoToken(t *testing.T) {
	symmetricKey := utils.RandomString(32)
	maker, err := NewJWTMaker(symmetricKey)
	require.NoError(t, err)

	invalidToken := utils.RandomString(20)

	payload, err := maker.VerifyToken(invalidToken, TokenTypeAccessToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
