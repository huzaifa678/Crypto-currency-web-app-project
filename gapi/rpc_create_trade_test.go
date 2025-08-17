package gapi

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateTradeRPC(t *testing.T) {
	trade, createTradeParams := createRandomTrade()

	testCases := []struct {
		name          string
		req           *pb.CreateTradeRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.CreateTradeResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateTradeRequest{
				BuyOrderId:  createTradeParams.BuyOrderID.String(),
				SellOrderId: createTradeParams.SellOrderID.String(),
				MarketId:    createTradeParams.MarketID.String(),
				Price:       createTradeParams.Price,
				Amount:      createTradeParams.Amount,
				Fee:         createTradeParams.Fee,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, trade.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				arg := createTradeParams

				store.EXPECT().
					CreateTrade(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(trade, nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateTradeResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotNil(t, res.Trade)
				require.Equal(t, trade.ID.String(), res.Trade.TradeId)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.CreateTradeRequest{
				BuyOrderId:  createTradeParams.BuyOrderID.String(),
				SellOrderId: createTradeParams.SellOrderID.String(),
				MarketId:    createTradeParams.MarketID.String(),
				Price:       createTradeParams.Price,
				Amount:      createTradeParams.Amount,
				Fee:         createTradeParams.Fee,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateTrade(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateTradeResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidBuyOrderID",
			req: &pb.CreateTradeRequest{
				BuyOrderId:  "invalid-uuid",
				SellOrderId: createTradeParams.SellOrderID.String(),
				MarketId:    createTradeParams.MarketID.String(),
				Price:       createTradeParams.Price,
				Amount:      createTradeParams.Amount,
				Fee:         createTradeParams.Fee,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, trade.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateTrade(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateTradeResponse, err error) {
				log.Println("ERROR: ", err)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateTradeRequest{
				BuyOrderId:  createTradeParams.BuyOrderID.String(),
				SellOrderId: createTradeParams.SellOrderID.String(),
				MarketId:    createTradeParams.MarketID.String(),
				Price:       createTradeParams.Price,
				Amount:      createTradeParams.Amount,
				Fee:         createTradeParams.Fee,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, trade.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateTrade(gomock.Any(), gomock.Any()).
					Times(1).
					Return(trade, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.CreateTradeResponse, err error) {
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

			res, err := server.CreateTrade(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

