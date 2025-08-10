package gapi

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateMarketRPC(t *testing.T) {
	CreateMarketParams, market, CreateMarketRow := createRandomMarket()

	testCases := []struct {
		name          string
		req           *pb.CreateMarketRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.CreateMarketResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateMarketRequest{
				BaseCurrency:   market.BaseCurrency,
				QuoteCurrency:  market.QuoteCurrency,
				MinOrderAmount: market.MinOrderAmount,
				PricePrecision: market.PricePrecision,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, market.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				arg := CreateMarketParams

				store.EXPECT().
					CreateMarket(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(CreateMarketRow, nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateMarketResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, market.ID.String(), res.MarketId)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.CreateMarketRequest{
				BaseCurrency:   market.BaseCurrency,
				QuoteCurrency:  market.QuoteCurrency,
				MinOrderAmount: market.MinOrderAmount,
				PricePrecision: market.PricePrecision,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
				return context.Background()
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateMarket(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateMarketResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidBaseCurrency",
			req: &pb.CreateMarketRequest{
				BaseCurrency:   "", 
				QuoteCurrency:  market.QuoteCurrency,
				MinOrderAmount: market.MinOrderAmount,
				PricePrecision: market.PricePrecision,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, market.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateMarket(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateMarketResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InvalidQuoteCurrency",
			req: &pb.CreateMarketRequest{
				BaseCurrency:   market.BaseCurrency,
				QuoteCurrency:  "", 
				MinOrderAmount: market.MinOrderAmount,
				PricePrecision: market.PricePrecision,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, market.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateMarket(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateMarketResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateMarketRequest{
				BaseCurrency:   market.BaseCurrency,
				QuoteCurrency:  market.QuoteCurrency,
				MinOrderAmount: market.MinOrderAmount,
				PricePrecision: market.PricePrecision,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, market.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				arg := CreateMarketParams

				store.EXPECT().
					CreateMarket(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.CreateMarketRow{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.CreateMarketResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore_interface(ctrl)
			tc.buildStubs(store)

			server := NewTestServer(t, store, nil)
			ctx := tc.setupAuth(t, server.tokenMaker)

			res, err := server.CreateMarket(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func createRandomMarket() (db.CreateMarketParams, db.Market, db.CreateMarketRow) {
    rand.New(rand.NewSource(int64(time.Now().UnixNano())))
    currencies := []string{"USD", "EUR", "BTC", "ETH", "JPY"}
    baseCurrency := currencies[rand.Intn(len(currencies))]
    quoteCurrency := currencies[rand.Intn(len(currencies))]

    for baseCurrency == quoteCurrency {
        quoteCurrency = currencies[rand.Intn(len(currencies))]
    }

	marketArgs := db.CreateMarketParams{
		BaseCurrency:  baseCurrency,
		QuoteCurrency: quoteCurrency,
		MinOrderAmount: "0.1",
		PricePrecision: 8,
	}

    market := db.Market{
        ID:            uuid.New(),
        BaseCurrency:  marketArgs.BaseCurrency,
        QuoteCurrency: marketArgs.QuoteCurrency,
        MinOrderAmount: marketArgs.MinOrderAmount,
        PricePrecision: marketArgs.PricePrecision,
		CreatedAt:     time.Now(),
    }

	marketRow := db.CreateMarketRow {
		ID: market.ID,
		BaseCurrency: market.BaseCurrency,
		QuoteCurrency: market.QuoteCurrency,
		CreatedAt: market.CreatedAt,
	}

    return marketArgs, market, marketRow
}

