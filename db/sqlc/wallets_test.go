package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestCreateWallet(t *testing.T) {

	userArgs := CreateUserParams {
		Email: "exam129@example.com",
		PasswordHash: "12345rtyu",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	arg := CreateWalletParams{
		UserID:   user.ID,
		Currency: "USD",
		Balance:  sql.NullString{String: "1000.00000000", Valid: true},
	}

	wallet, err := testQueries.CreateWallet(context.Background(), arg)
	require.NoError(t, err, "Failed to create wallet")

	require.NotZero(t, wallet.ID)
	require.Equal(t, arg.UserID, wallet.UserID)
	require.Equal(t, arg.Currency, wallet.Currency)
	require.Equal(t, arg.Balance.String, wallet.Balance.String)
	require.NotZero(t, wallet.CreatedAt)
	require.Equal(t, "0.00000000", wallet.LockedBalance.String) 
}

func TestDeleteWallet(t *testing.T) {
	userArgs := CreateUserParams {
		Email: "exam140@example.com",
		PasswordHash: "12345rtyuzhht",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	walletArg := CreateWalletParams{
		UserID:   user.ID,
		Currency: "USD",
		Balance:  sql.NullString{String: "1000.00000000", Valid: true},
	}

	wallet, err := testQueries.CreateWallet(context.Background(), walletArg)
	require.NoError(t, err, "Failed to create wallet")

	err = testQueries.DeleteUser(context.Background(), wallet.ID)
	require.NoError(t, err)

}

func TestGetWalletByUserIDAndCurrency(t *testing.T) {
	userArgs := CreateUserParams {
		Email: "exam141@example.com",
		PasswordHash: "vfvfe33433gtgtg",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	walletArgs := CreateWalletParams {
		UserID: user.ID,
		Currency: "USD",
		Balance: sql.NullString{String: "1000.00000000", Valid: true},
	}

	wallet, err := testQueries.CreateWallet(context.Background(), walletArgs)
	require.NoError(t, err, "Failed to create wallet")

	getWalletArgs := GetWalletByUserIDAndCurrencyParams {
		UserID: wallet.UserID,
		Currency: wallet.Currency, 
	}

	fetchedWallet, err := testQueries.GetWalletByUserIDAndCurrency(context.Background(), getWalletArgs)
	require.NoError(t, err)
	require.Equal(t, wallet.UserID, fetchedWallet.UserID)
	require.Equal(t, wallet.Currency, fetchedWallet.Currency)
}

func TestUpdateWallet(t *testing.T) {
	userArgs := CreateUserParams {
		Email: "exam143@example.com",
		PasswordHash: "vfvfe33433gtgccecdfrfr",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	walletArgs := CreateWalletParams {
		UserID: user.ID,
		Currency: "USD",
		Balance: sql.NullString{String: "1000.00000000", Valid: true},
	}

	wallet, err := testQueries.CreateWallet(context.Background(), walletArgs)
	require.NoError(t, err, "Failed to create wallet")

	updatedWalletArgs := UpdateWalletBalanceParams {
		Balance: sql.NullString{String: "2000.00000000", Valid: true},
		LockedBalance: sql.NullString{String: "0.00000000", Valid: true},
		UserID: wallet.UserID,
		Currency: "USD",
	}

	err = testQueries.UpdateWalletBalance(context.Background(), updatedWalletArgs)
    require.NoError(t, err)

	getWalletArgs := GetWalletByUserIDAndCurrencyParams {
		UserID: updatedWalletArgs.UserID,
		Currency: updatedWalletArgs.Currency, 
	}
	updatedWallet, err := testQueries.GetWalletByUserIDAndCurrency(context.Background(), getWalletArgs)

	require.NoError(t, err)
	require.Equal(t, updatedWallet.Balance, updatedWalletArgs.Balance)
	require.Equal(t, updatedWallet.LockedBalance, updatedWalletArgs.LockedBalance)
}