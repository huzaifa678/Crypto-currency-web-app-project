package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {

	email := createRandomEmail()

	arg := CreateUserParams{
		Username:     utils.RandomString(10),
		Email:        email,
		PasswordHash: "rhfcjndwd3344ndd",
		Role:         "user",
		IsVerified:   false,
	}

	users, err := testStore.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, arg.Username, users.Username)
	require.Equal(t, arg.Email, users.Email)
	require.Equal(t, arg.Role, users.Role)
	require.Equal(t, arg.IsVerified, users.IsVerified)
	require.WithinDuration(t, time.Now(), users.CreatedAt, time.Second)
}

func TestDeleteUser(t *testing.T) {

	email := createRandomEmail()

	arg := CreateUserParams{
		Username:     fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:        email,
		PasswordHash: "hashedpassword",
		Role:         "user",
		IsVerified:   false,
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	err = testStore.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	_, err = testStore.GetUserByID(context.Background(), user.ID)
	require.Error(t, err)
	require.EqualError(t, err, ErrRecordNotFound.Error())
}

func TestGetUserByEmail(t *testing.T) {

	email := createRandomEmail()

	arg := CreateUserParams{
		Username:     utils.RandomString(31),
		Email:        email,
		PasswordHash: "3535554frff",
		Role:         "admin",
		IsVerified:   false,
	}
	createdUser, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	fetchedUser, err := testStore.GetUserByEmail(context.Background(), arg.Email)
	require.NoError(t, err)
	require.Equal(t, createdUser.ID, fetchedUser.ID)
	require.Equal(t, arg.Email, fetchedUser.Email)
	require.Equal(t, arg.Role, fetchedUser.Role)
	require.Equal(t, arg.IsVerified, fetchedUser.IsVerified)
}

func TestUpdateUser(t *testing.T) {

	email := createRandomEmail()

	arg := CreateUserParams{
		Username:     utils.RandomString(28),
		Email:        email,
		PasswordHash: "54ffv895tnng",
		Role:         "user",
		IsVerified:   false,
	}
	user, err := testStore.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	updateArg := UpdateUserParams{
		PasswordHash: "newpasswordhash",
		IsVerified:   true,
		ID:           user.ID,
	}
	err = testStore.UpdateUser(context.Background(), updateArg)
	require.NoError(t, err)

	updatedUser, err := testStore.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, updateArg.PasswordHash, updatedUser.PasswordHash)
	require.Equal(t, updateArg.IsVerified, updatedUser.IsVerified)
}

func createRandomEmail() string {
	return fmt.Sprintf("testing-%s@example.com", uuid.New().String())
}
