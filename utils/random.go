//nolint:revive
package utils

import (
	"strings"

	"math/rand/v2"
)

func RandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTVWXYZ01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTVWXYZ01234567890njnkuhbnhbbiibjbhjbihbibi"

	var sb strings.Builder

	for i := 0; i < length; i++ {
		sb.WriteByte(letters[rand.IntN(len(letters))])
	}
	return sb.String()
}

func RandomUserByLength(length int) string {
	return RandomString(length)
}

func RandomUser() string {
	return RandomString(10)
}
