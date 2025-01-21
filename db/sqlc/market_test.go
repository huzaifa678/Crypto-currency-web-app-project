package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)


func TestCreateMarket (t *testing.T) {
	arg := CreateMarketParams {
		BaseCurrency: "BTC",
		QuoteCurrency: "USD",
		MinOrderAmount: sql.NullString{String: "0.01", Valid: true},
		PricePrecision: sql.NullInt32{Int32: 8, Valid: true},
	}

	market, err := testQueries.CreateMarket(context.Background(), arg)
	require.NoError(t, err, "Failed to create market")
	require.NotEmpty(t, market.ID, "Market ID should not be empty")
	require.Equal(t, arg.BaseCurrency, market.BaseCurrency, "BaseCurrency should match")
	require.Equal(t, arg.QuoteCurrency, market.QuoteCurrency, "QuoteCurrency should match")
	require.NotZero(t, market.CreatedAt, "CreatedAt should not be zero")
}

func TestDeleteMarket(t *testing.T) {
	
	createMarketArg := CreateMarketParams {
		BaseCurrency: "BTC",
		QuoteCurrency: "USD",
		MinOrderAmount: sql.NullString{String: "0.01", Valid: true},
		PricePrecision: sql.NullInt32{Int32: 8, Valid: true},
	}

	createMarket, err := testQueries.CreateMarket(context.Background(), createMarketArg)


	err = testQueries.DeleteMarket(context.Background(), createMarket.ID)
	require.NoError(t, err, "Failed to delete market")
	deletedMarket, err := testQueries.GetMarketByID(context.Background(), createMarket.ID)
	require.Error(t, err, "Expected error for non-existent market")
	require.Equal(t, sql.ErrNoRows, err, "Error should be sql.ErrNoRows")
	require.Empty(t, deletedMarket, "Deleted market should be empty")
}

func TestGetMarketById(t *testing.T) {

	arg := CreateMarketParams {
		BaseCurrency: "BTC",
		QuoteCurrency: "PKR",
		MinOrderAmount: sql.NullString{String: "0.01", Valid: true},
		PricePrecision: sql.NullInt32{Int32: 8, Valid: true},
	}

	market, err := testQueries.CreateMarket(context.Background(), arg)

	fmt.Println(market.ID)

	fetchedMarket, err := testQueries.GetMarketByID(context.Background(), market.ID)

	require.NoError(t, err, "Failed to fetch market by ID")
	require.Equal(t, market.ID, fetchedMarket.ID, "Market ID should match")
	require.Equal(t, market.BaseCurrency, fetchedMarket.BaseCurrency, "BaseCurrency should match")
	require.Equal(t, market.QuoteCurrency, fetchedMarket.QuoteCurrency, "Quote currency must match")
	require.Equal(t, market.CreatedAt, fetchedMarket.CreatedAt, "CreatedAt should match")
	
}

func TestListMarkets(t *testing.T) {
	
    for i := 0; i < 5; i++ {
        arg := CreateMarketParams{
            BaseCurrency:  fmt.Sprintf("BASE%d", i),
            QuoteCurrency: fmt.Sprintf("QUOTE%d", i),
            MinOrderAmount: sql.NullString{String: "0.01", Valid: true},
            PricePrecision: sql.NullInt32{Int32: 8, Valid: true},
        }

        _, err := testQueries.CreateMarket(context.Background(), arg)
        require.NoError(t, err, "Failed to create market")
    }

    markets, err := testQueries.ListMarkets(context.Background())
    require.NoError(t, err, "Failed to list markets")
    require.NotEmpty(t, markets, "Markets list should not be empty")

    for i, market := range markets {
        require.NotEmpty(t, market.ID, "Market ID should not be empty")
        require.NotEmpty(t, market.BaseCurrency, "BaseCurrency should not be empty")
        require.NotEmpty(t, market.QuoteCurrency, "QuoteCurrency should not be empty")
        require.NotZero(t, market.CreatedAt, "CreatedAt should not be zero")

        if i < 5 {
            require.Equal(t, fmt.Sprintf("BASE%d", 4-i), market.BaseCurrency, "BaseCurrency should match")
            require.Equal(t, fmt.Sprintf("QUOTE%d", 4-i), market.QuoteCurrency, "QuoteCurrency should match")
        }
    }
}
