package db

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func TestCreateMarket(t *testing.T) {

	marketParams, _, marketRow := createRandomMarket()

	require.NotEmpty(t, marketRow.ID, "Market ID should not be empty")
	require.Equal(t, marketRow.BaseCurrency, marketParams.BaseCurrency, "BaseCurrency should match")
	require.Equal(t, marketRow.QuoteCurrency, marketParams.QuoteCurrency, "QuoteCurrency should match")
	require.NotZero(t, marketRow.CreatedAt, "CreatedAt should not be zero")
}

func TestDeleteMarket(t *testing.T) {


	marketParams, _, _ := createRandomMarket()

	existingMarkets, err := testStore.ListMarketsByUsername(context.Background(), marketParams.Username)
	require.NoError(t, err)
	
	for _, market := range existingMarkets {
		err := testStore.DeleteMarket(context.Background(), market.ID)
		require.NoError(t, err)
	}


	createMarketArg := CreateMarketParams{
		Username:       marketParams.Username,
		BaseCurrency:   "USD",
		QuoteCurrency:  "INR",
		MinOrderAmount: decimal.NewFromFloat(0.01),
		PricePrecision: 8,
	}

	log.Println("Creating market with args:", createMarketArg)

	createMarket, err := testStore.CreateMarket(context.Background(), createMarketArg)

	if err != nil {
		t.Fatalf("Failed to create market: %v", err)
	}

	err = testStore.DeleteMarket(context.Background(), createMarket.ID)
	require.NoError(t, err, "Failed to delete market")
	deletedMarket, err := testStore.GetMarketByID(context.Background(), createMarket.ID)
	require.Error(t, err, "Expected error for non-existent market")
	require.Equal(t, ErrRecordNotFound, err, "Error should be pgx error")
	require.Empty(t, deletedMarket, "Deleted market should be empty")
}

func TestGetMarketById(t *testing.T) {

	marketParams, _, _ := createRandomMarket()

	createMarketArg := CreateMarketParams{
		Username:       marketParams.Username,
		BaseCurrency:   "BTC",
		QuoteCurrency:  "INR",
		MinOrderAmount: decimal.NewFromFloat(0.01),
		PricePrecision: 8,
	}

	createMarket, _ := testStore.CreateMarket(context.Background(), createMarketArg)

	fetchedMarket, _ := testStore.GetMarketByID(context.Background(), createMarket.ID)

	require.Equal(t, createMarket.ID, fetchedMarket.ID, "Market ID should match")
	require.Equal(t, createMarket.BaseCurrency, fetchedMarket.BaseCurrency, "BaseCurrency should match")
	require.Equal(t, createMarket.QuoteCurrency, fetchedMarket.QuoteCurrency, "Quote currency must match")
	require.Equal(t, createMarket.CreatedAt, fetchedMarket.CreatedAt, "CreatedAt should match")
}

func TestListMarkets(t *testing.T) {
	ctx := context.Background()


	seenPairs := make(map[string]struct{})
	var marketParams CreateMarketParams
	var market CreateMarketRow
	for i := 0; i < 3; i++ {
		for {
			marketParams, _, market = createRandomMarket()
			pairKey := marketParams.BaseCurrency + "_" + marketParams.QuoteCurrency
			if _, exists := seenPairs[pairKey]; !exists {
				seenPairs[pairKey] = struct{}{}
				break
			}
		}

		_, err := testStore.CreateMarket(ctx, marketParams)
    	require.NoError(t, err, "Failed to create market")

		log.Println("Inserted Market:", market)
	}

	markets, err := testStore.ListMarketsByUsername(ctx, marketParams.Username)
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
	ctx := context.Background()

	userArgs := CreateUserParams{
		Username:     fmt.Sprintf("markets-%s-username", uuid.New().String()),
		Email:        fmt.Sprintf("market-%s@example.com", uuid.New().String()),
		PasswordHash: "randompassword",
		Role:         "user",
		IsVerified:   true,
	}

	user, _ := testStore.CreateUser(ctx, userArgs)

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
		Username:       user.Username,
		BaseCurrency:   baseCurrency + fmt.Sprintf("_%d", rand.Intn(100000)),
		QuoteCurrency:  quoteCurrency + fmt.Sprintf("_%d", rand.Intn(100000)),
		MinOrderAmount: decimal.NewFromFloat(0.1),
		PricePrecision: 8,
	}

	market := Market{
		ID:             uuid.New(),
		Username:       marketArgs.Username,
		BaseCurrency:   marketArgs.BaseCurrency,
		QuoteCurrency:  marketArgs.QuoteCurrency,
		MinOrderAmount: marketArgs.MinOrderAmount,
		PricePrecision: marketArgs.PricePrecision,
		CreatedAt:      time.Now(),
	}

	marketRow := CreateMarketRow{
		ID:            market.ID,
		Username:      market.Username,
		BaseCurrency:  market.BaseCurrency,
		QuoteCurrency: market.QuoteCurrency,
		CreatedAt:     market.CreatedAt,
	}

	log.Println("Generated Market Params:", marketArgs)
	log.Println("Generated Market Object:", market)
	log.Println("Generated Market Row:", marketRow)

	return marketArgs, market, marketRow
}
