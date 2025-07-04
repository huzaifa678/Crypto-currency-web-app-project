package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/jackc/pgconn"

	"github.com/stretchr/testify/require"
)


func TestCreateTransactionTx(t *testing.T) {
	store := NewStore(testDB)

	email := createRandomEmailForTx()

	userArgs := CreateUserParams{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email: email,
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	transactionArgs := TransactionsParams {
		Username: user.Username,
		UserEmail: user.Email,
		Type: "deposit",
		Currency: "usd",
		Amount: "100.00000000",
		Status: "pending",
		Address: sql.NullString{String: "0x0000", Valid: true},
		TxHash: sql.NullString{String: "0x0000", Valid: true},
	}



	market := createRandomMarketForFee(t)

	feeArgs := CreateFeeParams {
		Username: market.Username,
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
	require.NoError(t, err, "Failed to create transaction")
	transaction, err := testQueries.GetTransactionsByUserEmail(context.Background(), transactionArgs.UserEmail)
    require.NoError(t, err, "Failed to get transaction")
    require.NotEmpty(t, transaction, "Transaction should not be empty")
	require.Equal(t, transactionArgs.Username, transaction[0].Username, "UserID should match")
	require.Equal(t, transactionArgs.UserEmail, transaction[0].UserEmail, "UserID should match")
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


func TestDeadlockDetectionForCreateTransaction(t *testing.T) {
	store := NewStore(testDB)

	email := createRandomEmailForTx()

	createUserParams := CreateUserParams{
		Username: utils.RandomUser(),
		Email:        email,
		PasswordHash: "8rrfrf4t45",
		Role:         "user",
		IsVerified:   sql.NullBool{Bool: true, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), createUserParams)
	require.NoError(t, err, "Failed to create user")

	errs := make(chan error, 2)

	transactionParams1 := TransactionsParams{
		Username: user.Username,
		UserEmail:   user.Email,
		Type:     "deposit",
		Currency: "USD",
		Amount:   "50.00000000",
		Status:   "pending",
		Address:  sql.NullString{String: "0x1111", Valid: true},
		TxHash:   sql.NullString{String: "0xhash1", Valid: true},
	}

	transactionParams2 := TransactionsParams{
		Username: user.Username,
		UserEmail:   user.Email,
		Type:     "withdrawal",
		Currency: "usd",
		Amount:   "30.00000000",
		Status:   "pending",
		Address:  sql.NullString{String: "0x2222", Valid: true},
		TxHash:   sql.NullString{String: "0xhash2", Valid: true},
	}

	market := createRandomMarketForFee(t)

	feeArgs := CreateFeeParams {
		Username: market.Username,
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

func TestDeadLockDetectionForUpdatingAmount(t *testing.T) {
	store := NewStore(testDB)

	email := createRandomEmailForTx()

	createUser1Params := CreateUserParams{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:	email,
		PasswordHash: "cdcewcds",
		Role: "user",
		IsVerified: sql.NullBool{Bool: true, Valid: true},
	}

	user1, err := testQueries.CreateUser(context.Background(), createUser1Params)
	require.NoError(t, err, "Failed to create user")

	email2 := createRandomEmailForTx()

	createUser2Params := CreateUserParams{
		Username: fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email: email2,
		PasswordHash: "cdcewcccfvs",
		Role: "user",
		IsVerified: sql.NullBool{Bool: true, Valid: true},
	}

	user2, err := testQueries.CreateUser(context.Background(), createUser2Params)
	require.NoError(t, err, "Failed to create user")

	market := createRandomMarketForOrder(t)

	createOrder1Params := CreateOrderParams{
		Username: user1.Username,
		UserEmail: user1.Email,
		MarketID: market.ID,
		Type: "buy",
		Status: "open",
		Price: sql.NullString{String: "100.50000000", Valid: true},
		Amount: "10.00000000",
	}


	order1, err := testQueries.CreateOrder(context.Background(), createOrder1Params)
	require.NoError(t, err, "Failed to create order for user1")

	createOrder2Params := CreateOrderParams{
		Username: user2.Username,
		UserEmail:   user2.Email,
		MarketID: market.ID,
		Type:     "sell",
		Status:   "open",
		Price:    sql.NullString{String: "100.50000000", Valid: true},
		Amount:   "10.00000000",
	}

	order2, err := testQueries.CreateOrder(context.Background(), createOrder2Params)
	require.NoError(t, err, "Failed to create order for user2")

	errCh := make(chan error)

	type UpdateOrderStatusAndFilledAmountParams struct {
		Status       OrderStatus    `json:"status"`
		FilledAmount sql.NullString `json:"filled_amount"`
		ID           uuid.UUID      `json:"id"`
	}

	go func() {
		_, err := store.UpdatedOrderTx(context.Background(), UpdatedOrderParams{
			Status:	"filled",
			FilledAmount:  sql.NullString{String: "15.00000000", Valid: true},
			ID: order1.ID, 
		})
		errCh <- err
	}()

	go func() {
		_, err := store.UpdatedOrderTx(context.Background(), UpdatedOrderParams{
			Status:	"filled",
			FilledAmount:  sql.NullString{String: "5.00000000", Valid: true}, 
			ID: order2.ID,
		})
		errCh <- err
	}()

	for i := 0; i < 2; i++ {
		err := <-errCh
		if err != nil {
			require.Contains(t, err.Error(), "deadlock detected", "Expected deadlock error")
		} else {
			t.Log("Transaction succeeded")
		}
	}
}

func createRandomEmailForTx() string {
	return fmt.Sprintf("tx-%s@example.com", uuid.New().String())
}

