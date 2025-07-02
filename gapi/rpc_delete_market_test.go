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

func TestDeleteMarketRPC(t *testing.T) {
	_, market, _ := createRandomMarket()

	testCases := []struct {
		name          string
		req           *pb.DeleteMarketRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.DeleteMarketResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.DeleteMarketRequest{
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

				store.EXPECT().
					DeleteMarket(gomock.Any(), gomock.Eq(market.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteMarketResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "Market deleted successfully", res.Message)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.DeleteMarketRequest{
				MarketId: market.ID.String(),
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetMarketByID(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					DeleteMarket(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return context.Background()
            },
			checkResponse: func(t *testing.T, res *pb.DeleteMarketResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidUUID",
			req: &pb.DeleteMarketRequest{
				MarketId: "invalid-uuid",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, market.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetMarketByID(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					DeleteMarket(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteMarketResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "MarketNotFound",
			req: &pb.DeleteMarketRequest{
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

				store.EXPECT().
					DeleteMarket(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteMarketResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.DeleteMarketRequest{
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

				store.EXPECT().
					DeleteMarket(gomock.Any(), gomock.Eq(market.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteMarketResponse, err error) {
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

			server := NewTestServer(t, store)
			ctx := tc.setupAuth(t, server.tokenMaker)

			res, err := server.DeleteMarket(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
} 