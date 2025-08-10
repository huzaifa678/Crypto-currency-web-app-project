package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func TestCreateTrade(t *testing.T) {

	email := createRandomEmailForTrade()

	buyerUsersArgs := CreateUserParams{
		Username:     fmt.Sprintf("buyer_%s", uuid.New().String()),
		Email:        email,
		PasswordHash: "kddeoovpds",
		Role:         "user",
		IsVerified:   true,
	}

	buyer, err := testStore.CreateUser(context.Background(), buyerUsersArgs)
	require.NoError(t, err, "Failed to create user")

	market := createRandomMarketForTrade(t)

	buyOrderArgs := CreateOrderParams{
		Username:  buyer.Username,
		UserEmail: buyer.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     "100.50000000",
		Amount:    "1000.00000000",
	}

	buyOrder, err := testStore.CreateOrder(context.Background(), buyOrderArgs)
	require.NoError(t, err, "Failed to create order for the buyer")

	email2 := createRandomEmailForTrade()

	sellerUsersArgs := CreateUserParams{
		Username:     fmt.Sprintf("buyer_%s", uuid.New().String()),
		Email:        email2,
		PasswordHash: "fvfdvrrgtg",
		Role:         "user",
		IsVerified:   true,
	}

	seller, err := testStore.CreateUser(context.Background(), sellerUsersArgs)
	require.NoError(t, err, "Failed to create user")

	sellOrderArgs := CreateOrderParams{
		Username:  seller.Username,
		UserEmail: seller.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     "110.50000000",
		Amount:    "1000.00000000",
	}

	sellOrder, err := testStore.CreateOrder(context.Background(), sellOrderArgs)
	require.NoError(t, err, "Failed to create order for the seller")

	tradeArgs := CreateTradeParams{
		Username:    market.Username,
		BuyOrderID:  buyOrder.ID,
		SellOrderID: sellOrder.ID,
		MarketID:    market.ID,
		Price:       "105.50000000",
		Amount:      "10.00000000",
		Fee:         "0.01",
	}

	trade, err := testStore.CreateTrade(context.Background(), tradeArgs)
	require.NoError(t, err, "Failed to create trade")

	require.NotEmpty(t, trade.ID, "Trade ID should not be empty")
	require.Equal(t, buyOrder.ID, trade.BuyOrderID, "BuyOrderID should match")
	require.Equal(t, sellOrder.ID, trade.SellOrderID, "SellOrderID should match")
	require.Equal(t, market.ID, trade.MarketID, "MarketID should match")
	require.Equal(t, "105.50000000", trade.Price, "Price should match the trade price")
	require.Equal(t, "10.00000000", trade.Amount, "Amount should match the trade amount")
}

func TestDeleteTrade(t *testing.T) {
	market := createRandomMarketForTrade(t)

	email := createRandomEmailForTrade()

	buyerUsersArgs := CreateUserParams{
		Username:     fmt.Sprintf("buyer_%s", uuid.New().String()),
		Email:        email,
		PasswordHash: "kddeoovpds",
		Role:         "user",
		IsVerified:   true,
	}

	buyer, err := testStore.CreateUser(context.Background(), buyerUsersArgs)
	require.NoError(t, err, "Failed to create user")

	buyOrderArgs := CreateOrderParams{
		Username:  buyer.Username,
		UserEmail: buyer.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     "100.50000000",
		Amount:    "1000.00000000",
	}

	buyOrder, err := testStore.CreateOrder(context.Background(), buyOrderArgs)
	require.NoError(t, err, "Failed to create order for the buyer")

	email2 := createRandomEmailForTrade()

	sellerUsersArgs := CreateUserParams{
		Username:     fmt.Sprintf("buyer_%s", uuid.New().String()),
		Email:        email2,
		PasswordHash: "fvfdvrrgtg",
		Role:         "user",
		IsVerified:   true,
	}

	seller, err := testStore.CreateUser(context.Background(), sellerUsersArgs)
	require.NoError(t, err, "Failed to create user")

	sellOrderArgs := CreateOrderParams{
		Username:  seller.Username,
		UserEmail: seller.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     "110.50000000",
		Amount:    "1000.00000000",
	}

	sellOrder, err := testStore.CreateOrder(context.Background(), sellOrderArgs)
	require.NoError(t, err, "Failed to create order for the seller")

	tradeArgs := CreateTradeParams{
		Username:    market.Username,
		BuyOrderID:  buyOrder.ID,
		SellOrderID: sellOrder.ID,
		MarketID:    market.ID,
		Price:       "105.50000000",
		Amount:      "10.00000000",
		Fee:         "0.01",
	}

	trade, err := testStore.CreateTrade(context.Background(), tradeArgs)
	require.NoError(t, err, "Failed to create trade")
	err = testStore.DeleteTrade(context.Background(), trade.ID)
	require.NoError(t, err, "Failed to delete trade")
	_, err = testStore.GetTradeByID(context.Background(), trade.ID)
	require.Error(t, err, "Expected error when fetching deleted trade")
	require.Equal(t, ErrRecordNotFound, err, "Error should be no rows found")
}

func TestGetTradeById(t *testing.T) {
	market := createRandomMarketForTrade(t)

	email := createRandomEmailForTrade()

	buyerUsersArgs := CreateUserParams{
		Username:     fmt.Sprintf("buyer_%s", uuid.New().String()),
		Email:        email,
		PasswordHash: "kddeoovpds",
		Role:         "user",
		IsVerified:   true,
	}

	buyer, err := testStore.CreateUser(context.Background(), buyerUsersArgs)
	require.NoError(t, err, "Failed to create user")

	buyOrderArgs := CreateOrderParams{
		Username:  buyer.Username,
		UserEmail: buyer.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     "100.50000000",
		Amount:    "1000.00000000",
	}

	buyOrder, err := testStore.CreateOrder(context.Background(), buyOrderArgs)
	require.NoError(t, err, "Failed to create order for the buyer")

	email2 := createRandomEmailForTrade()

	sellerUsersArgs := CreateUserParams{
		Username:     fmt.Sprintf("buyer_%s", uuid.New().String()),
		Email:        email2,
		PasswordHash: "fvfdvrrgtg",
		Role:         "user",
		IsVerified:   true,
	}

	seller, err := testStore.CreateUser(context.Background(), sellerUsersArgs)
	require.NoError(t, err, "Failed to create user")

	sellOrderArgs := CreateOrderParams{
		Username:  seller.Username,
		UserEmail: seller.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     "110.50000000",
		Amount:    "1000.00000000",
	}

	sellOrder, err := testStore.CreateOrder(context.Background(), sellOrderArgs)
	require.NoError(t, err, "Failed to create order for the seller")

	tradeArgs := CreateTradeParams{
		Username:    market.Username,
		BuyOrderID:  buyOrder.ID,
		SellOrderID: sellOrder.ID,
		MarketID:    market.ID,
		Price:       "105.50000000",
		Amount:      "10.00000000",
		Fee:         "0.01",
	}

	trade, err := testStore.CreateTrade(context.Background(), tradeArgs)
	require.NoError(t, err, "Failed to create trade")
	fetchedTrade, err := testStore.GetTradeByID(context.Background(), trade.ID)
	require.NoError(t, err, "failed to get the trade by ID")
	require.Equal(t, trade.ID, fetchedTrade.ID, "The Trade ID is matched")
	require.Equal(t, trade.BuyOrderID, fetchedTrade.BuyOrderID, "The BuyOrderID is matched")
	require.Equal(t, trade.SellOrderID, fetchedTrade.SellOrderID, "The SellOrderID is matched")
	require.Equal(t, trade.MarketID, fetchedTrade.MarketID, "The MarketID is matched")
	require.Equal(t, trade.Price, fetchedTrade.Price, "The Price is matched")
	require.Equal(t, trade.Amount, fetchedTrade.Amount, "The Amount is matched")
}

func TestGetTradeByMarketID(t *testing.T) {
	market := createRandomMarketForTrade(t)

	email := createRandomEmailForTrade()

	buyerUsersArgs := CreateUserParams{
		Username:     fmt.Sprintf("buyer_%s", uuid.New().String()),
		Email:        email,
		PasswordHash: "kddeoovpds",
		Role:         "user",
		IsVerified:   true,
	}

	buyer, err := testStore.CreateUser(context.Background(), buyerUsersArgs)
	require.NoError(t, err, "Failed to create user")

	buyOrderArgs := CreateOrderParams{
		Username:  buyer.Username,
		UserEmail: buyer.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     "100.50000000",
		Amount:    "1000.00000000",
	}

	buyOrder, err := testStore.CreateOrder(context.Background(), buyOrderArgs)
	require.NoError(t, err, "Failed to create order for the buyer")

	email2 := createRandomEmailForTrade()

	sellerUsersArgs := CreateUserParams{
		Username:     utils.RandomUser(),
		Email:        email2,
		PasswordHash: "fvfdvrrgtg",
		Role:         "user",
		IsVerified:   true,
	}

	seller, err := testStore.CreateUser(context.Background(), sellerUsersArgs)
	require.NoError(t, err, "Failed to create user")

	sellOrderArgs := CreateOrderParams{
		Username:  seller.Username,
		UserEmail: seller.Email,
		MarketID:  market.ID,
		Type:      OrderType("buy"),
		Status:    OrderStatus("open"),
		Price:     "110.50000000",
		Amount:    "1000.00000000",
	}

	sellOrder, err := testStore.CreateOrder(context.Background(), sellOrderArgs)
	require.NoError(t, err, "Failed to create order for the seller")

	tradeArgs := CreateTradeParams{
		Username:    market.Username,
		BuyOrderID:  buyOrder.ID,
		SellOrderID: sellOrder.ID,
		MarketID:    market.ID,
		Price:       "105.50000000",
		Amount:      "10.00000000",
		Fee:         "0.01",
	}

	trade, err := testStore.CreateTrade(context.Background(), tradeArgs)
	require.NoError(t, err, "Failed to create trade")
	tradeByMarketID, err := testStore.GetTradesByMarketID(context.Background(), trade.MarketID)
	require.NoError(t, err, "failed to get the trade by Market ID")
	require.NotEmpty(t, trade, "Trade should not be empty")
	require.Equal(t, trade.MarketID, tradeByMarketID[0].MarketID, "The MarketID is matched")
}

func createRandomMarketForTrade(t *testing.T) CreateMarketRow {
	ctx := context.Background()

	userArgs := CreateUserParams{
		Username:     fmt.Sprintf("buyer_%s", uuid.New().String()),
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
		MinOrderAmount: "0.1",
		PricePrecision: 6,
	}

	market, err := testStore.CreateMarket(ctx, arg)
	require.NoError(t, err, "Failed to create random market")
	require.NotEmpty(t, market.ID, "Market ID should not be empty")
	require.Equal(t, baseCurrency, market.BaseCurrency, "BaseCurrency should match")
	require.Equal(t, quoteCurrency, market.QuoteCurrency, "QuoteCurrency should match")

	return market
}

func createRandomEmailForTrade() string {
	return fmt.Sprintf("trade-%s@example.com", uuid.New().String())
}
