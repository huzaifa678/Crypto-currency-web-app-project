package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)


func TestCreateFee(t *testing.T) {
	market := createRandomMarketForFee(t)

	feeArgs := CreateFeeParams {
        Username: market.Username,
		MarketID: market.ID,
		MakerFee: sql.NullString{String: "0.01", Valid: true},
		TakerFee: sql.NullString{String: "0.02", Valid: true},
	}

	fee, err := testQueries.CreateFee(context.Background(), feeArgs)
	require.NoError(t, err, "Failed to create fee")
	require.NotEmpty(t, fee.ID, "Fee ID should not be empty")
}

func TestDeleteFee(t *testing.T) {
	market := createRandomMarketForFee(t)

	feeArgs := CreateFeeParams {
        Username: market.Username,
        MarketID: market.ID,
		MakerFee: sql.NullString{String: "0.02", Valid: true},
		TakerFee: sql.NullString{String: "0.04", Valid: true},
	}

	fee, err := testQueries.CreateFee(context.Background(), feeArgs)

	err = testQueries.DeleteFee(context.Background(), fee.ID)
	require.NoError(t, err, "Failed to delete fee")
}

func TestGetFeeByMarketID(t *testing.T) {
	market := createRandomMarketForFee(t)

	feeArgs := CreateFeeParams {
        Username: market.Username,
		MarketID: market.ID,
		MakerFee: sql.NullString{String: "0.0100", Valid: true},
		TakerFee: sql.NullString{String: "0.0500", Valid: true},
	}

	
	fee, err := testQueries.CreateFee(context.Background(), feeArgs)
	require.NoError(t, err, "Failed to create fee")

	feeByMarketID, err := testQueries.GetFeeByMarketID(context.Background(), fee.MarketID)
	require.NoError(t, err, "Failed to get fee by market ID")
	require.NotEmpty(t, feeByMarketID, "Fee should not be empty")
	require.Equal(t, fee.MarketID, feeByMarketID.MarketID, "MarketID should match")
	require.Equal(t, fee.MakerFee, feeByMarketID.MakerFee, "The maker fee should match")
}

func createRandomMarketForFee(t *testing.T) CreateMarketRow {
	ctx := context.Background()

    userArgs := CreateUserParams{
        Username:     utils.RandomUser(),
        Email:        fmt.Sprintf("market-%s@example.com", uuid.New().String()),
        PasswordHash: "randompassword",
        Role:         "user",
        IsVerified:   sql.NullBool{Bool: true, Valid: true},
    }

    user, err := testQueries.CreateUser(ctx, userArgs)
    require.NoError(t, err, "Failed to create user for market")

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
            Username:      existingMarket.Username,
            BaseCurrency:  existingMarket.BaseCurrency,
            QuoteCurrency: existingMarket.QuoteCurrency,
            CreatedAt:     existingMarket.CreatedAt,
        }
    }

    arg := CreateMarketParams{
        Username:       user.Username, 
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

