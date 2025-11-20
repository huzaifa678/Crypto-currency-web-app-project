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

func TestGetTradeRPC(t *testing.T) {
	trade, _, getTrade := createRandomTrade()

	testCases := []struct {
		name          string
		req           *pb.GetTradeByIDRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.GetTradeByIDResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.GetTradeByIDRequest{
				TradeId: trade.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, trade.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(getTrade, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetTradeByIDResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, trade.Username, res.Trade.Username)
				require.Equal(t, trade.ID.String(), res.Trade.TradeId)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.GetTradeByIDRequest{
				TradeId: trade.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.GetTradeByIDResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidUUID",
			req: &pb.GetTradeByIDRequest{
				TradeId: "invalid-uuid",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, trade.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.GetTradeByIDResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "TradeNotFound",
			req: &pb.GetTradeByIDRequest{
				TradeId: trade.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, trade.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(db.GetTradeByIDRow{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.GetTradeByIDResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.GetTradeByIDRequest{
				TradeId: trade.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, trade.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(db.GetTradeByIDRow{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.GetTradeByIDResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name: "PermissionDenied",
			req: &pb.GetTradeByIDRequest{
				TradeId: trade.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, "testuser123", time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(getTrade, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetTradeByIDResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unknown, st.Code())
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

			res, err := server.GetTrade(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
