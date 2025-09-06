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

func TestGetTransactionsByUserEmailRPC(t *testing.T) {
	transaction, _, _ := createRandomTransaction()

	testCases := []struct {
		name          string
		req           *pb.GetTransactionsByUserEmailRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.GetTransactionsByUserEmailResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.GetTransactionsByUserEmailRequest{
				UserEmail: transaction.UserEmail,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, transaction.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionsByUserEmail(gomock.Any(), gomock.Eq(transaction.UserEmail)).
					Times(1).
					Return([]db.Transaction{transaction}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetTransactionsByUserEmailResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Len(t, res.Transactions, 1)
				require.Equal(t, transaction.UserEmail, res.Transactions[0].UserEmail)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.GetTransactionsByUserEmailRequest{
				UserEmail: transaction.UserEmail,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionsByUserEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.GetTransactionsByUserEmailResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidEmail",
			req: &pb.GetTransactionsByUserEmailRequest{
				UserEmail: "not-an-email",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, "anyuser", time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionsByUserEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.GetTransactionsByUserEmailResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.GetTransactionsByUserEmailRequest{
				UserEmail: transaction.UserEmail,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, "differentuser", time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionsByUserEmail(gomock.Any(), gomock.Eq(transaction.UserEmail)).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.GetTransactionsByUserEmailResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name: "PermissionDenied",
			req: &pb.GetTransactionsByUserEmailRequest{
				UserEmail: transaction.UserEmail,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, "otheruser", time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionsByUserEmail(gomock.Any(), gomock.Eq(transaction.UserEmail)).
					Times(1).
					Return([]db.Transaction{transaction}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GetTransactionsByUserEmailResponse, err error) {
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

			res, err := server.GetTransactionsByUserEmail(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
