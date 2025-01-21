package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)


func TestCreateOrder(t *testing.T) {

    userArgs := CreateUserParams {
		Email: "exam995@example.com",
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	arg := CreateMarketParams {
		BaseCurrency: "EUR",
		QuoteCurrency: "USD",
		MinOrderAmount: sql.NullString{String: "0.01", Valid: true},
		PricePrecision: sql.NullInt32{Int32: 8, Valid: true},
	}

	market, err := testQueries.CreateMarket(context.Background(), arg)

	args := CreateOrderParams {
		UserID: user.ID,
		MarketID: market.ID,
		Type: OrderType("buy"),
		Status: OrderStatus("open"),
		Price: sql.NullString{String: "100.50000000", Valid: true},
		Amount: "10.00000000",
	}

	order, err := testQueries.CreateOrder(context.Background(), args)

	require.NoError(t, err, "Failed to create order")
	require.NotEmpty(t, order, "Order should not be empty")
	require.Equal(t, args.UserID, order.UserID)
	require.Equal(t, args.MarketID, order.MarketID)
	require.Equal(t, args.Type, order.Type)
	require.Equal(t, args.Status, order.Status)
	require.Equal(t, args.Price, order.Price)
	require.Equal(t, args.Amount, order.Amount)
}

func TestGetOrderById(t *testing.T) {
	userArgs := CreateUserParams {
		Email: "exam888@example.com",
		PasswordHash: "9009909dddxxwd",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	arg := CreateMarketParams {
		BaseCurrency: "PKR",
		QuoteCurrency: "BTC",
		MinOrderAmount: sql.NullString{String: "0.01", Valid: true},
		PricePrecision: sql.NullInt32{Int32: 8, Valid: true},
	}

	market, err := testQueries.CreateMarket(context.Background(), arg)

	args := CreateOrderParams {
		UserID: user.ID,
		MarketID: market.ID,
		Type: OrderType("buy"),
		Status: OrderStatus("open"),
		Price: sql.NullString{String: "100.50000000", Valid: true},
		Amount: "10.00000000",
	}

	order, err := testQueries.CreateOrder(context.Background(), args)

	fetchedOrder, err := testQueries.GetOrderByID(context.Background(), order.ID)
	require.NoError(t, err, "Failed to fetch order by ID")
	require.NotEmpty(t, fetchedOrder, "Fetched order should not be empty")
	require.Equal(t, fetchedOrder.ID, fetchedOrder.ID)
	require.Equal(t, fetchedOrder.UserID, fetchedOrder.UserID)
	require.Equal(t, fetchedOrder.MarketID, fetchedOrder.MarketID)
	require.Equal(t, fetchedOrder.Type, fetchedOrder.Type)
	require.Equal(t, fetchedOrder.Status, fetchedOrder.Status)
	require.Equal(t, fetchedOrder.Price, fetchedOrder.Price)
	require.Equal(t, fetchedOrder.Amount, fetchedOrder.Amount)
}

func TestDeleteOrder(t *testing.T) {
	userArgs := CreateUserParams {
		Email: "exam887@example.com",
		PasswordHash: "9009909dddxxwd",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	arg := CreateMarketParams {
		BaseCurrency: "BTC",
		QuoteCurrency: "USD",
		MinOrderAmount: sql.NullString{String: "0.01", Valid: true},
		PricePrecision: sql.NullInt32{Int32: 8, Valid: true},
	}

	market, err := testQueries.CreateMarket(context.Background(), arg)

	args := CreateOrderParams {
		UserID: user.ID,
		MarketID: market.ID,
		Type: OrderType("buy"),
		Status: OrderStatus("open"),
		Price: sql.NullString{String: "100.50000000", Valid: true},
		Amount: "10.00000000",
	}

	order, err := testQueries.CreateOrder(context.Background(), args)

	err = testQueries.DeleteOrder(context.Background(), order.ID)
	require.NoError(t, err, "Failed to delete order")

	_, err = testQueries.GetOrderByID(context.Background(), order.ID)
	require.Error(t, err, "Expected error when fetching deleted order")
	require.Equal(t, sql.ErrNoRows, err, "Error should be no rows found")
}

func TestUpdateOrderStatusAndFilledAmount(t *testing.T) {
	userArgs := CreateUserParams {
		Email: "exam879@example.com",
		PasswordHash: "9009909dddxxwd",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	arg := CreateMarketParams {
		BaseCurrency: "PKR",
		QuoteCurrency: "EUR",
		MinOrderAmount: sql.NullString{String: "0.01", Valid: true},
		PricePrecision: sql.NullInt32{Int32: 8, Valid: true},
	}

	market, err := testQueries.CreateMarket(context.Background(), arg)

	args := CreateOrderParams {
		UserID: user.ID,
		MarketID: market.ID,
		Type: OrderType("buy"),
		Status: OrderStatus("open"),
		Price: sql.NullString{String: "100.50000000", Valid: true},
		Amount: "10.00000000",
	}

	order, err := testQueries.CreateOrder(context.Background(), args)

	updatedArg := UpdateOrderStatusAndFilledAmountParams{
		Status:       OrderStatus("open"),
		FilledAmount: sql.NullString{String: "10.00000000", Valid: true},
		ID:           order.ID,
	}

	err = testQueries.UpdateOrderStatusAndFilledAmount(context.Background(), updatedArg)
	require.NoError(t, err, "Failed to update order")

	updatedOrder, err := testQueries.GetOrderByID(context.Background(), order.ID)
	require.NoError(t, err, "Failed to fetch updated order")
	require.Equal(t, updatedArg.Status, updatedOrder.Status)
	require.Equal(t, updatedArg.FilledAmount, updatedOrder.FilledAmount)
	require.WithinDuration(t, time.Now(), updatedOrder.UpdatedAt.Time, time.Second, "UpdatedAt should be recent")
}

