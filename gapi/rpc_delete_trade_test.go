package gapi

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteTradeRPC(t *testing.T) {

	trade, _ := createRandomTrade()

	testCases := []struct {
		name          string
		req           *pb.DeleteTradeRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.DeleteTradeResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.DeleteTradeRequest{
				TradeId: trade.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, trade.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(trade, nil)

				store.EXPECT().
					DeleteTrade(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTradeResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "Succesfully deleted the trade", res.Success)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.DeleteTradeRequest{
				TradeId: trade.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					DeleteTrade(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTradeResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidUUID",
			req: &pb.DeleteTradeRequest{
				TradeId: "invalid-uuid",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, trade.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					DeleteTrade(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTradeResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "TradeNotFound",
			req: &pb.DeleteTradeRequest{
				TradeId: trade.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, trade.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(db.Trade{}, db.ErrRecordNotFound)

				store.EXPECT().
					DeleteTrade(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTradeResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.DeleteTradeRequest{
				TradeId: trade.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, trade.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(trade, nil)

				store.EXPECT().
					DeleteTrade(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTradeResponse, err error) {
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

			res, err := server.DeleteTrade(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func createRandomTrade() (trade db.Trade, createTradeParams db.CreateTradeParams) {

	_, sellOrder, _, _ := createRandomOrder()
	_, BuyOrder, _, _ := createRandomOrder()
	_, market, _ := createRandomMarket()

	createTradeParams = db.CreateTradeParams {
		BuyOrderID: BuyOrder.ID,
    	SellOrderID: sellOrder.ID,   
    	MarketID:    market.ID,      
    	Price:       decimal.NewFromFloat(0.0),   
    	Amount:      decimal.NewFromFloat(0.0),         
    	Fee:         decimal.NewFromFloat(5),
	}

	Trade := db.Trade {
		ID: uuid.New(),
		BuyOrderID: BuyOrder.ID,
    	SellOrderID: sellOrder.ID,   
    	MarketID:    market.ID,      
    	Price:       createTradeParams.Price,   
    	Amount:      createTradeParams.Amount,         
    	Fee:         createTradeParams.Fee,
		CreatedAt:   time.Now(),
	}

	return Trade, createTradeParams
}
