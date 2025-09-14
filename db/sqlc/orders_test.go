package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func TestCreateOrder(t *testing.T) {

	email := createRandomEmailForOrder()

	userArgs := CreateUserParams{
		Username:     fmt.Sprintf("testuser_%s", uuid.New()),
		Email:        email,
		PasswordHash: "9009909",
		Role:         "user",
		IsVerified:   false,
	}
	user, err := testStore.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	market := createRandomMarketForOrder(t)

	args := CreateOrderParams{
		Username:  user.Username,
		UserEmail: user.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     decimal.NewFromFloat(100.50000000),
		Amount:    decimal.NewFromFloat(10.00000000),
	}

	order, err := testStore.CreateOrder(context.Background(), args)

	require.NoError(t, err, "Failed to create order")
	require.NotEmpty(t, order, "Order should not be empty")
	require.Equal(t, args.UserEmail, order.UserEmail)
	require.Equal(t, args.MarketID, order.MarketID)
	require.Equal(t, args.Type, order.Type)
	require.Equal(t, args.Status, order.Status)
	require.True(t, args.Price.Equal(order.Price), "Price mismatch")
	require.True(t, args.Amount.Equal(order.Amount), "Amount mismatch")
}

func TestGetOrderById(t *testing.T) {

	email := createRandomEmailForOrder()

	userArgs := CreateUserParams{
		Username:     utils.RandomString(10),
		Email:        email,
		PasswordHash: "9009909dddxxwd",
		Role:         "user",
		IsVerified:   false,
	}
	user, err := testStore.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	market := createRandomMarketForOrder(t)

	args := CreateOrderParams{
		Username:  user.Username,
		UserEmail: user.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     decimal.NewFromFloat(100.50000000),
		Amount:    decimal.NewFromFloat(10.00000000),
	}

	order, _ := testStore.CreateOrder(context.Background(), args)

	fetchedOrder, err := testStore.GetOrderByID(context.Background(), order.ID)
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

	userArgs := CreateUserParams{
		Username:     utils.RandomString(12),
		Email:        email,
		PasswordHash: "9009909dddxxwd",
		Role:         "user",
		IsVerified:   false,
	}
	user, err := testStore.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	market := createRandomMarketForOrder(t)

	args := CreateOrderParams{
		Username:  user.Username,
		UserEmail: user.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     decimal.NewFromFloat(100.50000000),
		Amount:    decimal.NewFromFloat(10.00000000),
	}

	order, err := testStore.CreateOrder(context.Background(), args)

	require.NoError(t, err, "Failed to create order")

	err = testStore.DeleteOrder(context.Background(), order.ID)
	require.NoError(t, err, "Failed to delete order")

	_, err = testStore.GetOrderByID(context.Background(), order.ID)
	require.Error(t, err, "Expected error when fetching deleted order")
	require.Equal(t, ErrRecordNotFound, err, "Error should be no rows found")
}

func TestUpdateOrderStatusAndFilledAmount(t *testing.T) {

	email := createRandomEmailForOrder()

	userArgs := CreateUserParams{
		Username:     fmt.Sprintf("testuser_%s", uuid.New()),
		Email:        email,
		PasswordHash: "9009909dddxxwd",
		Role:         "user",
		IsVerified:   false,
	}
	user, err := testStore.CreateUser(context.Background(), userArgs)
	require.NoError(t, err, "Failed to create user")

	market := createRandomMarketForOrder(t)

	args := CreateOrderParams{
		Username:  user.Username,
		UserEmail: user.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     decimal.NewFromFloat(100.50000000),
		Amount:    decimal.NewFromFloat(10.00000000),
	}

	order, err := testStore.CreateOrder(context.Background(), args)

	require.NoError(t, err, "Failed to create order")

	updatedArg := UpdateOrderStatusAndFilledAmountParams{
		Status:       OrderStatus("open"),
		FilledAmount: decimal.NewFromFloat(10.00000000),
		ID:           order.ID,
	}

	err = testStore.UpdateOrderStatusAndFilledAmount(context.Background(), updatedArg)
	require.NoError(t, err, "Failed to update order")

	updatedOrder, err := testStore.GetOrderByID(context.Background(), order.ID)
	require.NoError(t, err, "Failed to fetch updated order")
	require.Equal(t, updatedArg.Status, updatedOrder.Status)
	require.True(t, updatedArg.FilledAmount.Equal(updatedOrder.FilledAmount), "Filled amount mismatch")
	require.WithinDuration(t, time.Now(), updatedOrder.UpdatedAt, time.Second, "UpdatedAt should be recent")
}

func createRandomMarketForOrder(t *testing.T) CreateMarketRow {
	ctx := context.Background()

	userArgs := CreateUserParams{
		Username:     fmt.Sprintf("testuser_%d", time.Now().UnixNano()),
		Email:        fmt.Sprintf("market-%s@example.com", uuid.New().String()),
		PasswordHash: "randompassword",
		Role:         "user",
		IsVerified:   true,
	}

	user, err := testStore.CreateUser(ctx, userArgs)
	require.NoError(t, err, "Failed to create user for market")

	currencies := []string{"USD", "EUR", "BTC", "ETH", "JPY"}
	baseCurrency := currencies[rand.Intn(len(currencies))]
	quoteCurrency := currencies[rand.Intn(len(currencies))]

	for baseCurrency == quoteCurrency {
		quoteCurrency = currencies[rand.Intn(len(currencies))]
	}

	existingMarket, err := testStore.GetMarketByCurrencies(ctx, GetMarketByCurrenciesParams{
		BaseCurrency:  baseCurrency,
		QuoteCurrency: quoteCurrency,
	})

	if err == nil {
		return CreateMarketRow{
			ID:            existingMarket.ID,
			Username:      existingMarket.Username,
			BaseCurrency:  existingMarket.BaseCurrency,
			QuoteCurrency: existingMarket.QuoteCurrency,
			CreatedAt:     existingMarket.CreatedAt,
		}
	}

	arg := CreateMarketParams{
		Username:       user.Username,
		BaseCurrency:   baseCurrency,
		QuoteCurrency:  quoteCurrency,
		MinOrderAmount: decimal.NewFromFloat(0.6),
		PricePrecision: 6,
	}

	market, err := testStore.CreateMarket(ctx, arg)
	require.NoError(t, err, "Failed to create random market")
	require.NotEmpty(t, market.ID, "Market ID should not be empty")
	require.Equal(t, baseCurrency, market.BaseCurrency, "BaseCurrency should match")
	require.Equal(t, quoteCurrency, market.QuoteCurrency, "QuoteCurrency should match")

	return market
}

func createRandomEmailForOrder() string {

	return fmt.Sprintf("order-%s@example.com", uuid.New().String())
}
