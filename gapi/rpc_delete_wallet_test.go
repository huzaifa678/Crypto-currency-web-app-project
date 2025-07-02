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

func TestDeleteWalletRPC(t *testing.T) {
	_, wallet, _ := createRandomWallet()

	testCases := []struct {
		name          string
		req           *pb.DeleteWalletRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.DeleteWalletResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.DeleteWalletRequest{
				WalletId: wallet.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, wallet.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetWalletByID(gomock.Any(), gomock.Eq(wallet.ID)).
					Times(1).
					Return(wallet, nil)

				store.EXPECT().
					DeleteWallet(gomock.Any(), gomock.Eq(wallet.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteWalletResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "successfully delete the wallet", res.Message)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.DeleteWalletRequest{
				WalletId: wallet.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return context.Background()
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetWalletByID(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					DeleteWallet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteWalletResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidUUID",
			req: &pb.DeleteWalletRequest{
				WalletId: "invalid-uuid",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, wallet.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetWalletByID(gomock.Any(), gomock.Any()).
					Times(0)

				store.EXPECT().
					DeleteWallet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteWalletResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "WalletNotFound",
			req: &pb.DeleteWalletRequest{
				WalletId: wallet.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, wallet.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetWalletByID(gomock.Any(), gomock.Eq(wallet.ID)).
					Times(1).
					Return(wallet, nil)

				store.EXPECT().
					DeleteWallet(gomock.Any(), gomock.Eq(wallet.ID)).
					Times(1).
					Return(db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteWalletResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.DeleteWalletRequest{
				WalletId: wallet.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, wallet.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetWalletByID(gomock.Any(), gomock.Eq(wallet.ID)).
					Times(1).
					Return(wallet, nil)

				store.EXPECT().
					DeleteWallet(gomock.Any(), gomock.Eq(wallet.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteWalletResponse, err error) {
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

			res, err := server.DeleteWallet(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
} 