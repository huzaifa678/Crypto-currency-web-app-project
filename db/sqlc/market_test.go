package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)


func TestCreateMarket (t *testing.T) {
	
	marketParams, _, _ := createRandomMarket()

	market, err := testQueries.CreateMarket(context.Background(), marketParams)
	require.NoError(t, err, "Failed to create market")
	require.NotEmpty(t, market.ID, "Market ID should not be empty")
	require.Equal(t, market.BaseCurrency, marketParams.BaseCurrency, "BaseCurrency should match")
	require.Equal(t, market.QuoteCurrency, marketParams.QuoteCurrency, "QuoteCurrency should match")
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

	_, _, marketRow := createRandomMarket()

	fmt.Println(marketRow.ID)

	fetchedMarket, err := testQueries.GetMarketByID(context.Background(), marketRow.ID)

	require.NoError(t, err, "Failed to fetch market by ID")
	require.Equal(t, marketRow.ID, fetchedMarket.ID, "Market ID should match")
	require.Equal(t, marketRow.BaseCurrency, fetchedMarket.BaseCurrency, "BaseCurrency should match")
	require.Equal(t, marketRow.QuoteCurrency, fetchedMarket.QuoteCurrency, "Quote currency must match")
	require.Equal(t, marketRow.CreatedAt, fetchedMarket.CreatedAt, "CreatedAt should match")
}

func TestListMarkets(t *testing.T) {
    ctx := context.Background()

    tx, err := testDB.BeginTx(ctx, nil)
    require.NoError(t, err, "Failed to begin transaction")
    defer tx.Rollback() 

    
    testQueriesWithTx := testQueries.WithTx(tx)

    
    seenPairs := make(map[string]struct{})
    for i := 0; i < 3; i++ {
        var marketParams CreateMarketParams
        var market CreateMarketRow
        for {
            marketParams, _, market = createRandomMarket()
            pairKey := marketParams.BaseCurrency + "_" + marketParams.QuoteCurrency
            if _, exists := seenPairs[pairKey]; !exists {
                seenPairs[pairKey] = struct{}{}
                break
            }
        }

        _, err := testQueriesWithTx.CreateMarket(ctx, marketParams)
        require.NoError(t, err, "Failed to create market")
        log.Println("Inserted Market:", market)
    }

    err = tx.Commit()
    require.NoError(t, err, "Failed to commit transaction")

    markets, err := testQueries.ListMarkets(ctx)
    require.NoError(t, err, "Failed to list markets")
    require.NotEmpty(t, markets, "Market list should not be empty")

    for _, m := range markets {
        log.Println("Retrieved Market:", m)
        require.NotEmpty(t, m.BaseCurrency, "BaseCurrency should not be empty")
        require.NotEmpty(t, m.QuoteCurrency, "QuoteCurrency should not be empty")
    }

}


var existingMarkets = make(map[string]struct{}) 

func createRandomMarket() (CreateMarketParams, Market, CreateMarketRow) {
	rand.Seed(uint64(time.Now().UnixNano()))
	currencies := []string{"USD", "EUR", "BTC", "ETH", "JPY"}

	var baseCurrency, quoteCurrency string
	for {
		baseCurrency = currencies[rand.Intn(len(currencies))]
		quoteCurrency = currencies[rand.Intn(len(currencies))]

		if baseCurrency != quoteCurrency {
			pairKey := baseCurrency + "_" + quoteCurrency
			if _, exists := existingMarkets[pairKey]; !exists {
				existingMarkets[pairKey] = struct{}{} 
				break
			}
		}
	}

	marketArgs := CreateMarketParams{
		BaseCurrency: baseCurrency,
		QuoteCurrency: quoteCurrency,
		MinOrderAmount: sql.NullString{
			String: "0.1",
			Valid:  true,
		},
		PricePrecision: sql.NullInt32{
			Int32: 8,
			Valid: true,
		},
	}

	market := Market{
		ID:            uuid.New(),
		BaseCurrency:  marketArgs.BaseCurrency,
		QuoteCurrency: marketArgs.QuoteCurrency,
		MinOrderAmount: marketArgs.MinOrderAmount,
		PricePrecision: marketArgs.PricePrecision,
		CreatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
	}

	marketRow := CreateMarketRow{
		ID:            market.ID,
		BaseCurrency:  market.BaseCurrency,
		QuoteCurrency: market.QuoteCurrency,
		CreatedAt:     market.CreatedAt,
	}

	log.Println("Generated Market Params:", marketArgs)
	log.Println("Generated Market Object:", market)
	log.Println("Generated Market Row:", marketRow)

	return marketArgs, market, marketRow
}


