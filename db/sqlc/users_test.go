package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	arg := CreateUserParams {
		Email: "exam124@example.com",
		PasswordHash: "rhfcjndwd3344ndd",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	users, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, arg.Email, users.Email)
	require.Equal(t, arg.Role, users.Role)
	require.Equal(t, arg.IsVerified, users.IsVerified)
	require.WithinDuration(t, time.Now(), users.CreatedAt.Time, time.Second)
}

func TestDeleteUser(t *testing.T) {
	arg := CreateUserParams{
		Email:        "deleteuser@example.com",
		PasswordHash: "hashedpassword",
		Role:         "user",
		IsVerified:   sql.NullBool{Bool: true, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	err = testQueries.DeleteUser(context.Background(), user.ID)
	require.NoError(t, err)

	_, err = testQueries.GetUserByID(context.Background(), user.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestGetUserByEmail(t *testing.T) {

	arg := CreateUserParams{
		Email:        "exam125@example.com",
		PasswordHash: "3535554frff",
		Role:         "admin",
		IsVerified:   sql.NullBool{Bool: true, Valid: true},
	}
	createdUser, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	fetchedUser, err := testQueries.GetUserByEmail(context.Background(), arg.Email)
	require.NoError(t, err)
	require.Equal(t, createdUser.ID, fetchedUser.ID)
	require.Equal(t, arg.Email, fetchedUser.Email)
	require.Equal(t, arg.Role, fetchedUser.Role)
	require.Equal(t, arg.IsVerified.Bool, fetchedUser.IsVerified.Bool)
}

func TestUpdateUser(t *testing.T) {

	arg := CreateUserParams{
		Email:        "exam126@example.com",
		PasswordHash: "54ffv895tnng",
		Role:         "user",
		IsVerified:   sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	updateArg := UpdateUserParams{
		PasswordHash: "newpasswordhash",
		IsVerified:   sql.NullBool{Bool: true, Valid: true},
		ID:           user.ID,
	}
	err = testQueries.UpdateUser(context.Background(), updateArg)
	require.NoError(t, err)

	updatedUser, err := testQueries.GetUserByID(context.Background(), user.ID)
	require.NoError(t, err)
	require.Equal(t, updateArg.PasswordHash, updatedUser.PasswordHash)
	require.Equal(t, updateArg.IsVerified.Bool, updatedUser.IsVerified.Bool)
}


