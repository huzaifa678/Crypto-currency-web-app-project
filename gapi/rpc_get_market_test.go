package gapi

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetMarketRPC(t *testing.T) {
	_, market, _ := createRandomMarket()

	testCases := []struct {
		name          string
		req           *pb.GetMarketRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.GetMarketResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.GetMarketRequest{
				MarketId: market.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, market.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetMarketByID(gomock.Any(), gomock.Eq(market.ID)).
					Times(1).
					Return(market, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetMarketResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotNil(t, res.Market)
				require.Equal(t, market.ID.String(), res.Market.MarketId)
				require.Equal(t, market.BaseCurrency, res.Market.BaseCurrency)
				require.Equal(t, market.QuoteCurrency, res.Market.QuoteCurrency)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.GetMarketRequest{
				MarketId: market.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return context.Background()
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetMarketByID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.GetMarketResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidUUID",
			req: &pb.GetMarketRequest{
				MarketId: "invalid-uuid",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, market.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetMarketByID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.GetMarketResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "MarketNotFound",
			req: &pb.GetMarketRequest{
				MarketId: market.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, market.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetMarketByID(gomock.Any(), gomock.Eq(market.ID)).
					Times(1).
					Return(db.Market{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.GetMarketResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.GetMarketRequest{
				MarketId: market.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, market.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetMarketByID(gomock.Any(), gomock.Eq(market.ID)).
					Times(1).
					Return(db.Market{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.GetMarketResponse, err error) {
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

			res, err := server.GetMarket(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
} 