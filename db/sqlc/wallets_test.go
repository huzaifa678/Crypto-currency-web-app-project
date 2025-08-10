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

func TestCreateWallet(t *testing.T) {

	email := createRandomEmailForWallet()

	userArgs := CreateUserParams{
		Username:     fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:        email,
		PasswordHash: "12345rtyu",
		Role:         "user",
		IsVerified:   true,
	}
	user, err := testStore.CreateUser(context.Background(), userArgs)

	require.NoError(t, err, "Failed to create user")

	arg := CreateWalletParams{
		Username:  user.Username,
		UserEmail: user.Email,
		Currency:  "USD",
		Balance:   "1000.00000000",
	}

	wallet, err := testStore.CreateWallet(context.Background(), arg)
	require.NoError(t, err, "Failed to create wallet")

	require.NotZero(t, wallet.ID)
	require.Equal(t, arg.Username, wallet.Username)
	require.Equal(t, arg.UserEmail, wallet.UserEmail)
	require.Equal(t, arg.Currency, wallet.Currency)
	require.Equal(t, arg.Balance, wallet.Balance)
	require.NotZero(t, wallet.CreatedAt)
	require.Equal(t, "0", wallet.LockedBalance)
}

func TestDeleteWallet(t *testing.T) {

	email := createRandomEmailForWallet()

	userArgs := CreateUserParams{
		Username:     utils.RandomString(36),
		Email:        email,
		PasswordHash: "12345rtyuzhht",
		Role:         "user",
		IsVerified:   true,
	}
	user, err := testStore.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	walletArg := CreateWalletParams{
		Username:  user.Username,
		UserEmail: user.Email,
		Currency:  "USD",
		Balance:   "1000.00000000",
	}

	wallet, err := testStore.CreateWallet(context.Background(), walletArg)
	require.NoError(t, err, "Failed to create wallet")

	err = testStore.DeleteUser(context.Background(), wallet.ID)
	require.NoError(t, err)

}

func TestGetWalletByUserEmailAndCurrency(t *testing.T) {

	email := createRandomEmailForWallet()

	userArgs := CreateUserParams{
		Username:     utils.RandomString(37),
		Email:        email,
		PasswordHash: "vfvfe33433gtgtg",
		Role:         "user",
		IsVerified:   true,
	}
	user, err := testStore.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	walletArgs := CreateWalletParams{
		Username:  user.Username,
		UserEmail: user.Email,
		Currency:  "USD",
		Balance:   "1000.00000000",
	}

	wallet, err := testStore.CreateWallet(context.Background(), walletArgs)
	require.NoError(t, err, "Failed to create wallet")

	fetchedWallet, err := testStore.GetWalletByID(context.Background(), wallet.ID)
	require.NoError(t, err)
	require.Equal(t, wallet.UserEmail, fetchedWallet.UserEmail)
	require.Equal(t, wallet.Currency, fetchedWallet.Currency)
}

func TestUpdateWallet(t *testing.T) {

	email := createRandomEmailForWallet()

	userArgs := CreateUserParams{
		Username:     fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:        email,
		PasswordHash: "vfvfe33433gtgccecdfrfr",
		Role:         "user",
		IsVerified:   true,
	}
	user, err := testStore.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	walletArgs := CreateWalletParams{
		Username:  user.Username,
		UserEmail: user.Email,
		Currency:  "USD",
		Balance:   "1000.00000000",
	}

	wallet, err := testStore.CreateWallet(context.Background(), walletArgs)
	require.NoError(t, err, "Failed to create wallet")

	updatedWalletArgs := UpdateWalletBalanceParams{
		ID:            wallet.ID,
		Balance:       "2000.00000000",
		LockedBalance: "0",
	}

	err = testStore.UpdateWalletBalance(context.Background(), updatedWalletArgs)
	require.NoError(t, err)

	updatedWallet, err := testStore.GetWalletByID(context.Background(), wallet.ID)

	require.NoError(t, err)
	require.Equal(t, updatedWallet.Balance, updatedWalletArgs.Balance)
	require.Equal(t, updatedWallet.LockedBalance, updatedWalletArgs.LockedBalance)
}

func createRandomEmailForWallet() string {
	return fmt.Sprintf("test-%s@example.com", uuid.New().String())
}
