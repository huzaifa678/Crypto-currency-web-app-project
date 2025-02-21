package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {

	email := createRandomEmailForTransaction()

	userArgs := CreateUserParams {
		Username: utils.RandomString(33),
		Email: email,
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	transactionsArgs := CreateTransactionParams {
		UserEmail: user.Email,
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

	email := createRandomEmailForTransaction()

	userArgs := CreateUserParams {
		Username: utils.RandomString(30),
		Email: email,
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	transactionsArgs := CreateTransactionParams {
		UserEmail: user.Email,
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

	email := createRandomEmailForTransaction()
	userArgs := CreateUserParams {
		Username: utils.RandomString(32),
		Email: email,
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	transactionsArgs := CreateTransactionParams {
		UserEmail: user.Email,
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
	require.Equal(t, transaction.UserEmail, transactionByID.UserEmail, "UserEmail should match")
	require.Equal(t, transaction.Type, transactionByID.Type, "Type should match")
	require.Equal(t, transaction.Currency, transactionByID.Currency, "Currency should match")
	require.Equal(t, transaction.Amount, transactionByID.Amount, "Amount should match")
	require.Equal(t, transaction.Address, transactionByID.Address, "Address should match")
	require.Equal(t, transaction.TxHash, transactionByID.TxHash, "TxHash should match")
}

func TestGetTransactionsByUserID(t *testing.T) {

	email := createRandomEmailForTransaction()

	userArgs := CreateUserParams {
		Username: utils.RandomString(34),
		Email: email,
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	transactionsArgs := CreateTransactionParams {
		UserEmail: user.Email,
		Type: TransactionType("deposit"),
		Currency: "USD",
		Amount: "100.00000000",
		Address: sql.NullString{String: "0x 0000 0000 0000 0000", Valid: true},
		TxHash: sql.NullString{String: "0x 0000 0000 0000 0000", Valid: true},
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), transactionsArgs)
	require.NoError(t, err, "Failed to create transaction")

	transactionsByUserID, err := testQueries.GetTransactionsByUserEmail(context.Background(), transaction.UserEmail)
	require.NoError(t, err, "Failed to get transaction by user ID")
	require.NotEmpty(t, transactionsByUserID, "Transaction should not be empty")
	require.Equal(t, transaction.ID, transactionsByUserID[0].ID, "Transaction ID should match")
}

func createRandomEmailForTransaction() string {
	return fmt.Sprintf("example-%s@example.com", uuid.New().String())
}

