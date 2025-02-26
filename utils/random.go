package utils

import (
	"strings"

	"golang.org/x/exp/rand"
)


func RandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTVWXYZ01234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTVWXYZ01234567890njnkuhbnhbbiibjbhjbihbibi"

	var sb strings.Builder

	for i := 0; i < length; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}

func RandomUser() string {
	return RandomString(10)
}