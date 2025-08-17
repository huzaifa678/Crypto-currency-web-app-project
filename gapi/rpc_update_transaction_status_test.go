package gapi

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func TestUpdateTransactionStatus(t *testing.T) {

	_, _, createTxParams := createRandomTransaction()
	transactionID := uuid.New()
	transaction := db.Transaction{
    	ID:       transactionID,
    	Status:   db.TransactionStatus(pb.TransactionStatus_PENDING),
    	Username: createTxParams.Username, 
	}

	testCases := []struct {
		name      string
		req       *pb.UpdateTransactionStatusRequest
		buildStubs func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse  func(t *testing.T, res *pb.UpdateTransactionStatusResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateTransactionStatusRequest{
				TransactionId: transactionID.String(),
				Status:        pb.TransactionStatus_PENDING,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, createTxParams.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), transactionID).
					Times(1).
					Return(transaction, nil)

				store.EXPECT().
					UpdateTransactionStatus(gomock.Any(), db.UpdateTransactionStatusParams{
						Status: db.TransactionStatus(pb.TransactionStatus_PENDING),
						ID: transactionID,
					}).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateTransactionStatusResponse, err error) {
				log.Println("ERROR", err)
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "Transaction status updated successfully", res.Success)
			},
		},
		{
			name: "InvalidID",
			req: &pb.UpdateTransactionStatusRequest{
				TransactionId: "invalid-uuid",
				Status:        pb.TransactionStatus_PENDING,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, createTxParams.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().GetTransactionByID(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateTransactionStatusResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "NotFound",
			req: &pb.UpdateTransactionStatusRequest{
				TransactionId: transactionID.String(),
				Status:        pb.TransactionStatus_PENDING,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, createTxParams.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), transactionID).
					Times(1).
					Return(db.Transaction{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateTransactionStatusResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "Unauthorized",
			req: &pb.UpdateTransactionStatusRequest{
				TransactionId: transactionID.String(),
				Status:        pb.TransactionStatus_PENDING,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
			},
			checkResponse: func(t *testing.T, res *pb.UpdateTransactionStatusResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InternalUpdateError",
			req: &pb.UpdateTransactionStatusRequest{
				TransactionId: transactionID.String(),
				Status:        pb.TransactionStatus_PENDING,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, createTxParams.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), transactionID).
					Times(1).
					Return(transaction, nil)

				store.EXPECT().
					UpdateTransactionStatus(gomock.Any(), gomock.Any()).
					Times(1).
					Return(errors.New("db error"))
			},
			checkResponse: func(t *testing.T, res *pb.UpdateTransactionStatusResponse, err error) {
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

			res, err := server.UpdateTransactionStatus(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
