package gapi

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdateOrderStatusAndFilledAmount(t *testing.T) {
	orderID := uuid.New()
	username := "test_user"

	order := db.Order{
		ID:       orderID,
		Username: username,
		Status:   db.OrderStatus("open"),
	}

	testCases := []struct {
		name          string
		req           *pb.UpdateOrderStatusAndFilledAmountRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateOrderStatusAndFilledAmountResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateOrderStatusAndFilledAmountRequest{
				OrderId:     orderID.String(),
				OrderStatus: pb.Status_FILLED,
				FilledAmount: 10,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetOrderByID(gomock.Any(), orderID).
					Times(1).
					Return(order, nil)

				store.EXPECT().
					UpdateOrderStatusAndFilledAmount(gomock.Any(), db.UpdateOrderStatusAndFilledAmountParams{
						Status:       db.OrderStatus("filled"),
						FilledAmount: decimal.NewFromFloat(10),
						ID:           orderID,
					}).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateOrderStatusAndFilledAmountResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "successfully updated the order status and filled amount", res.Success)
			},
		},
		{
			name: "InvalidID",
			req: &pb.UpdateOrderStatusAndFilledAmountRequest{
				OrderId:     "invalid-uuid",
				OrderStatus: pb.Status_OPEN,
				FilledAmount: 5,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
			},
			checkResponse: func(t *testing.T, res *pb.UpdateOrderStatusAndFilledAmountResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "Unauthorized",
			req: &pb.UpdateOrderStatusAndFilledAmountRequest{
				OrderId:     orderID.String(),
				OrderStatus: pb.Status_CANCELLED,
				FilledAmount: 3,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
			},
			checkResponse: func(t *testing.T, res *pb.UpdateOrderStatusAndFilledAmountResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "NotAuthorizedDifferentUser",
			req: &pb.UpdateOrderStatusAndFilledAmountRequest{
				OrderId:     orderID.String(),
				OrderStatus: pb.Status_FILLED,
				FilledAmount: 7,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, "other_user", time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetOrderByID(gomock.Any(), orderID).
					Times(1).
					Return(order, nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateOrderStatusAndFilledAmountResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unknown, st.Code())
			},
		},
		{
			name: "InternalUpdateError",
			req: &pb.UpdateOrderStatusAndFilledAmountRequest{
				OrderId:     orderID.String(),
				OrderStatus: pb.Status_FILLED,
				FilledAmount: 10,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetOrderByID(gomock.Any(), orderID).
					Times(1).
					Return(order, nil)

				store.EXPECT().
					UpdateOrderStatusAndFilledAmount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("db error"))
			},
			checkResponse: func(t *testing.T, res *pb.UpdateOrderStatusAndFilledAmountResponse, err error) {
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

			res, err := server.UpdateOrderStatusAndFilledAmount(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
