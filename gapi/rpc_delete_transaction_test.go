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



func TestDeleteTransactionRPC(t *testing.T) {
	transaction, _, _ := createRandomTransaction()

	testCases := []struct {
		name          string
		req           *pb.DeleteTransactionRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.DeleteTransactionResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.DeleteTransactionRequest{
				TransactionId: transaction.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, transaction.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(transaction, nil)

				store.EXPECT().
					DeleteTransaction(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTransactionResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "Successfully deleted the transaction", res.Success)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.DeleteTransactionRequest{
				TransactionId: transaction.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeleteTransaction(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTransactionResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidUUID",
			req: &pb.DeleteTransactionRequest{
				TransactionId: "invalid-uuid",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, transaction.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), gomock.Any()).
					Times(0)
				store.EXPECT().
					DeleteTransaction(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTransactionResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "TransactionNotFound",
			req: &pb.DeleteTransactionRequest{
				TransactionId: transaction.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, transaction.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(db.Transaction{}, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTransactionResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InternalErrorOnGet",
			req: &pb.DeleteTransactionRequest{
				TransactionId: transaction.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, transaction.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(db.Transaction{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTransactionResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name: "PermissionDenied",
			req: &pb.DeleteTransactionRequest{
				TransactionId: transaction.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, "differentuser", time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(transaction, nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTransactionResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unknown, st.Code())
			},
		},
		{
			name: "InternalErrorOnDelete",
			req: &pb.DeleteTransactionRequest{
				TransactionId: transaction.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, transaction.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(transaction, nil)

				store.EXPECT().
					DeleteTransaction(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteTransactionResponse, err error) {
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

			res, err := server.DeleteTransaction(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
