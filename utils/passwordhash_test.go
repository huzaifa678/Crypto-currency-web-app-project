package utils

import (
	"testing"
	"log"

	"github.com/stretchr/testify/require"
)



func TestHashPassword(t *testing.T) {
	password := "1234567"

	hashedPassword, err := HashPassword(password)
	log.Printf("hashedPassword: %v", hashedPassword)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
}

func TestComparePasswords(t *testing.T) {
	password := "1234567"
	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	err = ComparePasswords(hashedPassword, password)
	require.NoError(t, err)
}