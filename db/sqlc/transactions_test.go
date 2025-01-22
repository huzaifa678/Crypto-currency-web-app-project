package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {

	userArgs := CreateUserParams {
		Email: "exam112@example.com",
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	transactionsArgs := CreateTransactionParams {
		UserID: user.ID,
		Type: TransactionType("deposit"),
		Currency: "USD",
		Amount: "100.00000000",
		Address: sql.NullString{String: "0x 0000 0000 0000 0000", Valid: true},
		TxHash: sql.NullString{String: "0x 0000 0000 0000 0000", Valid: true},
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), transactionsArgs)
	require.NoError(t, err, "Failed to create transaction")
	require.NotEmpty(t, transaction.ID, "Transaction ID should not be empty")
}

func TestDeleteTransaction(t *testing.T) {

	userArgs := CreateUserParams {
		Email: "exam113@example.com",
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	transactionsArgs := CreateTransactionParams {
		UserID: user.ID,
		Type: TransactionType("deposit"),
		Currency: "USD",
		Amount: "100.00000000",
		Address: sql.NullString{String: "0x 0000 0000 0000 0000", Valid: true},
		TxHash: sql.NullString{String: "0x 0000 0000 0000 0000", Valid: true},
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), transactionsArgs)
	require.NoError(t, err, "Failed to create transaction")

	err = testQueries.DeleteTransaction(context.Background(), transaction.ID)
	require.NoError(t, err, "Failed to delete transaction")

	_, err = testQueries.GetTransactionByID(context.Background(), transaction.ID)
	require.Error(t, err, "Transaction should be deleted")
	require.Equal(t, sql.ErrNoRows, err, "Error should be sql.ErrNoRows")
}

func TestGetTransactionById(t *testing.T) {
	userArgs := CreateUserParams {
		Email: "exam114@example.com",
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	transactionsArgs := CreateTransactionParams {
		UserID: user.ID,
		Type: TransactionType("deposit"),
		Currency: "USD",
		Amount: "100.00000000",
		Address: sql.NullString{String: "0x 0000 0000 0000 0000", Valid: true},
		TxHash: sql.NullString{String: "0x 0000 0000 0000 0000", Valid: true},
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), transactionsArgs)
	require.NoError(t, err, "Failed to create transaction")

	transactionByID, err := testQueries.GetTransactionByID(context.Background(), transaction.ID)
	require.NoError(t, err, "Failed to get transaction by ID")
	require.NotEmpty(t, transactionByID, "Transaction should not be empty")
	require.Equal(t, transaction.ID, transactionByID.ID, "Transaction ID should match")
	require.Equal(t, transaction.UserID, transactionByID.UserID, "UserID should match")
	require.Equal(t, transaction.Type, transactionByID.Type, "Type should match")
	require.Equal(t, transaction.Currency, transactionByID.Currency, "Currency should match")
	require.Equal(t, transaction.Amount, transactionByID.Amount, "Amount should match")
	require.Equal(t, transaction.Address, transactionByID.Address, "Address should match")
	require.Equal(t, transaction.TxHash, transactionByID.TxHash, "TxHash should match")
}

func TestGetTransactionsByUserID(t *testing.T) {
	userArgs := CreateUserParams {
		Email: "exam115@example.com",
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	transactionsArgs := CreateTransactionParams {
		UserID: user.ID,
		Type: TransactionType("deposit"),
		Currency: "USD",
		Amount: "100.00000000",
		Address: sql.NullString{String: "0x 0000 0000 0000 0000", Valid: true},
		TxHash: sql.NullString{String: "0x 0000 0000 0000 0000", Valid: true},
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), transactionsArgs)
	require.NoError(t, err, "Failed to create transaction")

	transactionsByUserID, err := testQueries.GetTransactionsByUserID(context.Background(), transaction.UserID)
	require.NoError(t, err, "Failed to get transaction by user ID")
	require.NotEmpty(t, transactionsByUserID, "Transaction should not be empty")
	require.Equal(t, transaction.ID, transactionsByUserID[0].ID, "Transaction ID should match")
}
