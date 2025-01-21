package db

import (
	"context"
	"database/sql"
	"testing"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func TestCreateTrade(t *testing.T) {
	
	buyerUsersArgs := CreateUserParams {
		Email: "exam369@example.com",
		PasswordHash: "kddeoovpds",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	buyer, err := testQueries.CreateUser(context.Background(), buyerUsersArgs)
	require.NoError(t, err, "Failed to create user")

	market := createRandomMarketForTrade(t)

	buyOrderArgs := CreateOrderParams {
		UserID: buyer.ID,
		MarketID: market.ID,
		Type: OrderType("buy"),
		Status: OrderStatus("open"),
		Price: sql.NullString{String: "100.50000000", Valid: true},
		Amount: "1000.00000000",
	}

	buyOrder, err := testQueries.CreateOrder(context.Background(), buyOrderArgs)
	require.NoError(t, err, "Failed to create order for the buyer")

	sellerUsersArgs := CreateUserParams {
		Email: "exam368@example.com",
		PasswordHash: "fvfdvrrgtg",
		Role: "user",
		IsVerified: sql.NullBool{Bool: false, Valid: true},
	}

	seller, err := testQueries.CreateUser(context.Background(), sellerUsersArgs)
	require.NoError(t, err, "Failed to create user")

	sellOrderArgs := CreateOrderParams {
		UserID: seller.ID,
		MarketID: market.ID,
		Type: OrderType("buy"),
		Status: OrderStatus("open"),
		Price: sql.NullString{String: "110.50000000", Valid: true},
		Amount: "1000.00000000",
	}

	sellOrder, err := testQueries.CreateOrder(context.Background(), sellOrderArgs)
	require.NoError(t, err, "Failed to create order for the seller")

	tradeArgs := CreateTradeParams {
		BuyOrderID: buyOrder.ID,
		SellOrderID: sellOrder.ID,
		MarketID: market.ID,
		Price: "105.50000000",
		Amount: "10.00000000",
	}

	trade, err := testQueries.CreateTrade(context.Background(), tradeArgs)
	require.NoError(t, err, "Failed to create trade")

	require.NotEmpty(t, trade.ID, "Trade ID should not be empty")
	require.Equal(t, buyOrder.ID, trade.BuyOrderID, "BuyOrderID should match")
	require.Equal(t, sellOrder.ID, trade.SellOrderID, "SellOrderID should match")
	require.Equal(t, market.ID, trade.MarketID, "MarketID should match")
	require.Equal(t, "105.50000000", trade.Price, "Price should match the trade price")
	require.Equal(t, "10.00000000", trade.Amount, "Amount should match the trade amount")
}

func createRandomMarketForTrade(t *testing.T) CreateMarketRow {
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

