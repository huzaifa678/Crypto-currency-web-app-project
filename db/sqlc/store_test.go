package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jackc/pgconn"

	"github.com/stretchr/testify/require"
)


func TestCreateTransactionTx(t *testing.T) {
	store := NewStore(testDB)

	userArgs := CreateUserParams{
		Email: "exam1003@example.com",
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	transactionArgs := TransactionsParams {
		UserID: user.ID,
		Type: "deposit",
		Currency: "usd",
		Amount: "100.00000000",
		Status: "pending",
		Address: sql.NullString{String: "0x0000", Valid: true},
		TxHash: sql.NullString{String: "0x0000", Valid: true},
	}



	market := createRandomMarketForFee(t)

	feeArgs := CreateFeeParams {
		MarketID: market.ID,
		MakerFee: sql.NullString{String: "0.0100", Valid: true},
		TakerFee: sql.NullString{String: "0.0200", Valid: true},
	}

	feeParams := FeeParams{
		MarketID: feeArgs.MarketID,
		Amount: feeArgs.MakerFee,
		TakerFee: feeArgs.TakerFee,
	}
	err = store.CreateTransactionTx(context.Background(), transactionArgs, feeParams)

	transaction, err := testQueries.GetTransactionsByUserID(context.Background(), transactionArgs.UserID)
    require.NoError(t, err, "Failed to get transaction")
    require.NotEmpty(t, transaction, "Transaction should not be empty")
	require.Equal(t, transactionArgs.UserID, transaction[0].UserID, "UserID should match")
	require.Equal(t, transactionArgs.Type, transaction[0].Type, "Type should match")
	require.Equal(t, transactionArgs.Currency, transaction[0].Currency, "Currency should match")
	require.Equal(t, transactionArgs.Amount, transaction[0].Amount, "Amount should match")
	require.Equal(t, TransactionStatus(transactionArgs.Status), transaction[0].Status, "Status should match")
	require.Equal(t, transactionArgs.Address, transaction[0].Address, "Address should match")
	require.Equal(t, transactionArgs.TxHash, transaction[0].TxHash, "TxHash should match")

    fee, err := testQueries.GetFeeByMarketID(context.Background(), feeArgs.MarketID)
    require.NoError(t, err, "Failed to get fee")
    require.NotEmpty(t, fee, "Fee should not be empty")
    require.Equal(t, feeArgs.MarketID, fee.MarketID, "MarketID should match")
    require.Equal(t, feeArgs.MakerFee, fee.MakerFee, "MakerFee should match")
    require.Equal(t, feeArgs.TakerFee, fee.TakerFee, "TakerFee should match")
}


func TestDeadlockDetection(t *testing.T) {
	store := NewStore(testDB)

	createUserParams := CreateUserParams{
		Email:        "exam1005@example.com",
		PasswordHash: "8rrfrf4t45",
		Role:         "user",
		IsVerified:   sql.NullBool{Bool: true, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), createUserParams)
	require.NoError(t, err, "Failed to create user")

	errs := make(chan error, 2)

	transactionParams1 := TransactionsParams{
		UserID:   user.ID,
		Type:     "deposit",
		Currency: "USD",
		Amount:   "50.00000000",
		Status:   "pending",
		Address:  sql.NullString{String: "0x1111", Valid: true},
		TxHash:   sql.NullString{String: "0xhash1", Valid: true},
	}

	transactionParams2 := TransactionsParams{
		UserID:   user.ID,
		Type:     "withdrawal",
		Currency: "usd",
		Amount:   "30.00000000",
		Status:   "pending",
		Address:  sql.NullString{String: "0x2222", Valid: true},
		TxHash:   sql.NullString{String: "0xhash2", Valid: true},
	}

	market := createRandomMarketForFee(t)

	feeArgs := CreateFeeParams {
		MarketID: market.ID,
		MakerFee: sql.NullString{String: "0.0100", Valid: true},
		TakerFee: sql.NullString{String: "0.0200", Valid: true},
	}

	feeParams := FeeParams{
		MarketID: feeArgs.MarketID,
		Amount: feeArgs.MakerFee,
		TakerFee: feeArgs.TakerFee,
	}

	go func() {
		err := store.CreateTransactionTx(context.Background(), transactionParams1, feeParams)
		errs <- err
	}()

	go func() {
		err := store.CreateTransactionTx(context.Background(), transactionParams2, feeParams)
		errs <- err
	}()

	for i := 0; i < 2; i++ {
		err := <-errs
		if err != nil {
			if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "40P01" {
				t.Errorf("Deadlock detected: %v", err)
			} else {
				require.NoError(t, err, "Unexpected error")
			}
		}
	}
}
