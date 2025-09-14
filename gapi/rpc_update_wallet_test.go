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
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdateWalletRPC(t *testing.T) {
	_, wallet, _ := createRandomWallet()

	testCases := []struct {
		name          string
		req           *pb.UpdateWalletRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateWalletResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateWalletRequest{
				WalletId:      wallet.ID.String(),
				Balance:       100,
				LockedBalance: 50,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, wallet.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				arg := db.UpdateWalletBalanceParams{
					ID:            wallet.ID,
					Balance:       decimal.NewFromFloat(100),
					LockedBalance: decimal.NewFromFloat(50),
				}
				store.EXPECT().
					UpdateWalletBalance(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(nil)

				store.EXPECT().
					GetWalletByID(gomock.Any(), gomock.Eq(wallet.ID)).
					Times(1).
					Return(wallet, nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateWalletResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "successfully updated the wallet", res.Message)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.UpdateWalletRequest{
				WalletId:      wallet.ID.String(),
				Balance:       100,
				LockedBalance: 50,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return context.Background()
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					UpdateWalletBalance(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateWalletResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidUUID",
			req: &pb.UpdateWalletRequest{
				WalletId:      "invalid-uuid",
				Balance:       100,
				LockedBalance: 50,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, wallet.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					UpdateWalletBalance(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateWalletResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.UpdateWalletRequest{
				WalletId:      wallet.ID.String(),
				Balance:       100,
				LockedBalance: 50,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, wallet.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				arg := db.UpdateWalletBalanceParams{
					ID:            wallet.ID,
					Balance:       decimal.NewFromFloat(100),
					LockedBalance: decimal.NewFromFloat(50),
				}
				store.EXPECT().
					UpdateWalletBalance(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateWalletResponse, err error) {
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

			res, err := server.UpdateWallet(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
} 