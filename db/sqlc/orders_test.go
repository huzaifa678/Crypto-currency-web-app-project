package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)


func TestCreateOrder(t *testing.T) {

	email := createRandomEmailForOrder()

    userArgs := CreateUserParams {
		Username: utils.RandomString(15),
		Email: email,
		PasswordHash: "9009909",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	market := createRandomMarketForOrder(t)

	args := CreateOrderParams {
		UserEmail: user.Email,
		MarketID: market.ID,
		Type: OrderType("buy"),
		Status: OrderStatus("open"),
		Price: sql.NullString{String: "100.50000000", Valid: true},
		Amount: "10.00000000",
	}

	order, err := testQueries.CreateOrder(context.Background(), args)

	require.NoError(t, err, "Failed to create order")
	require.NotEmpty(t, order, "Order should not be empty")
	require.Equal(t, args.UserEmail, order.UserEmail)
	require.Equal(t, args.MarketID, order.MarketID)
	require.Equal(t, args.Type, order.Type)
	require.Equal(t, args.Status, order.Status)
	require.Equal(t, args.Price, order.Price)
	require.Equal(t, args.Amount, order.Amount)
}

func TestGetOrderById(t *testing.T) {

	email := createRandomEmailForOrder()

	userArgs := CreateUserParams {
		Username: utils.RandomString(10),
		Email: email,
		PasswordHash: "9009909dddxxwd",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")


	market := createRandomMarketForOrder(t)

	args := CreateOrderParams {
		UserEmail: user.Email,
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
	require.Equal(t, fetchedOrder.UserEmail, fetchedOrder.UserEmail)
	require.Equal(t, fetchedOrder.MarketID, fetchedOrder.MarketID)
	require.Equal(t, fetchedOrder.Type, fetchedOrder.Type)
	require.Equal(t, fetchedOrder.Status, fetchedOrder.Status)
	require.Equal(t, fetchedOrder.Price, fetchedOrder.Price)
	require.Equal(t, fetchedOrder.Amount, fetchedOrder.Amount)
}

func TestDeleteOrder(t *testing.T) {

	email := createRandomEmailForOrder()

	userArgs := CreateUserParams {
		Username: utils.RandomString(12),
		Email: email,
		PasswordHash: "9009909dddxxwd",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")


	market := createRandomMarketForOrder(t)

	args := CreateOrderParams {
		UserEmail: user.Email,
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

	email := createRandomEmailForOrder()

	userArgs := CreateUserParams {
		Username: utils.RandomString(13),
		Email: email,
		PasswordHash: "9009909dddxxwd",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}
	user, err := testQueries.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")


	market := createRandomMarketForOrder(t)

	args := CreateOrderParams {
		UserEmail: user.Email,
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

func createRandomMarketForOrder(t *testing.T) CreateMarketRow {
	ctx := context.Background()

	currencies := []string{"USD", "EUR", "BTC", "ETH", "JPY"}
	baseCurrency := currencies[rand.Intn(len(currencies))]
	quoteCurrency := currencies[rand.Intn(len(currencies))]

	for baseCurrency == quoteCurrency {
		quoteCurrency = currencies[rand.Intn(len(currencies))]
	}

	existingMarket, err := testQueries.GetMarketByCurrencies(ctx, GetMarketByCurrenciesParams{
		BaseCurrency:  baseCurrency,
		QuoteCurrency: quoteCurrency,
	})

	if err == nil {
		return CreateMarketRow{
			ID:            existingMarket.ID,
			BaseCurrency:  existingMarket.BaseCurrency,
			QuoteCurrency: existingMarket.QuoteCurrency,
			CreatedAt:     existingMarket.CreatedAt,
		}
	}

	arg := CreateMarketParams{
		BaseCurrency:  baseCurrency,
		QuoteCurrency: quoteCurrency,
		MinOrderAmount: sql.NullString{
			String: "0.1",
			Valid:  true,
		},
		PricePrecision: sql.NullInt32{
			Int32: 6,
			Valid: true,
		},
	}

	market, err := testQueries.CreateMarket(ctx, arg)
	require.NoError(t, err, "Failed to create random market")
	require.NotEmpty(t, market.ID, "Market ID should not be empty")
	require.Equal(t, baseCurrency, market.BaseCurrency, "BaseCurrency should match")
	require.Equal(t, quoteCurrency, market.QuoteCurrency, "QuoteCurrency should match")

	return market
}


func createRandomEmailForOrder() string {
	
	return fmt.Sprintf("order-%s@example.com", uuid.New().String())
}
