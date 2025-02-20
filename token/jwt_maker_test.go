package token

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)



func TestJWTMaker(t *testing.T) {

	secretKey := RandomString(32)
	maker, err := NewJWTMaker(secretKey)

	require.NoError(t, err)

	username := RandomString(12)
	duration := time.Minute

	token, err := maker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)

	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)

	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestTokenExpired(t *testing.T) {
	maker, err := NewJWTMaker(RandomString(32))
	require.NoError(t, err)
	username := RandomString(12)

	token, err := maker.CreateToken(username, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidSecret(t *testing.T) {

	invalidSecret := RandomString(10)
	_, err := NewJWTMaker(invalidSecret)

	require.EqualError(t, err, "Invalid key size: The secret key size is not equal to minimum of 32 size")
}

func TestInvalidToken(t *testing.T) {
	secretKey := RandomString(32)
	maker, err := NewJWTMaker(secretKey)
	require.NoError(t, err)

	invalidToken := RandomString(20)

	payload, err := maker.VerifyToken(invalidToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func RandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTVWXYZ01234567890"

	var sb strings.Builder

	for i := 0; i < length; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}