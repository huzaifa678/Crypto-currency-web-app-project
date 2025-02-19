package db

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)


func TestCreateMarket (t *testing.T) {
	
	marketParams, _, marketRow := createRandomMarket()

	require.NotEmpty(t, marketRow.ID, "Market ID should not be empty")
	require.Equal(t, marketRow.BaseCurrency, marketParams.BaseCurrency, "BaseCurrency should match")
	require.Equal(t, marketRow.QuoteCurrency, marketParams.QuoteCurrency, "QuoteCurrency should match")
	require.NotZero(t, marketRow.CreatedAt, "CreatedAt should not be zero")
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

	createMarketArg := CreateMarketParams {
		BaseCurrency: "BTC",
		QuoteCurrency: "INR",
		MinOrderAmount: sql.NullString{String: "0.01", Valid: true},
		PricePrecision: sql.NullInt32{Int32: 8, Valid: true},
	}

	createMarket, _ := testQueries.CreateMarket(context.Background(), createMarketArg)

	fetchedMarket, _ := testQueries.GetMarketByID(context.Background(), createMarket.ID)

	require.Equal(t, createMarket.ID, fetchedMarket.ID, "Market ID should match")
	require.Equal(t, createMarket.BaseCurrency, fetchedMarket.BaseCurrency, "BaseCurrency should match")
	require.Equal(t, createMarket.QuoteCurrency, fetchedMarket.QuoteCurrency, "Quote currency must match")
	require.Equal(t, createMarket.CreatedAt, fetchedMarket.CreatedAt, "CreatedAt should match")
}

func TestListMarkets(t *testing.T) {
    ctx := context.Background()


    
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

        
        log.Println("Inserted Market:", market)
    }


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


