package gapi

import (
	"context"
	"database/sql"
	"log"
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

func TestCreateWalletRPC(t *testing.T) {
	user, _, _, _, _, _ := createRandomUser()
	walletID := uuid.New()

	testCases := []struct {
		name          string
		req           *pb.CreateWalletRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.CreateWalletResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateWalletRequest{
				UserEmail: user.Email,
				Currency:  "BTC",
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				arg := db.CreateWalletParams{
					Username:  user.Username,
					UserEmail: user.Email,
					Currency:  "BTC",
					Balance:   decimal.NewFromFloat(0),
				}
				store.EXPECT().
					CreateWallet(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Wallet{
						ID:        walletID,
						Username:  user.Username,
						UserEmail: user.Email,
						Currency:  "BTC",
						Balance:   decimal.NewFromFloat(0),
					}, nil)
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute, token.TokenTypeAccessToken)
            },
			checkResponse: func(t *testing.T, res *pb.CreateWalletResponse, err error) {
				log.Println("ERROR: ", err)
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, walletID.String(), res.WalletId)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.CreateWalletRequest{
				UserEmail: user.Email,
				Currency:  "BTC",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return context.Background()
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateWallet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateWalletResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidEmail",
			req: &pb.CreateWalletRequest{
				UserEmail: "invalid-email",
				Currency:  "BTC",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateWallet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateWalletResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InvalidCurrency",
			req: &pb.CreateWalletRequest{
				UserEmail: user.Email,
				Currency:  "INVALID",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateWallet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateWalletResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateWalletRequest{
				UserEmail: user.Email,
				Currency:  "BTC",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, user.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				arg := db.CreateWalletParams{
					Username:  user.Username,
					UserEmail: user.Email,
					Currency:  "BTC",
					Balance:   decimal.NewFromFloat(0),
				}
				store.EXPECT().
					CreateWallet(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Wallet{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.CreateWalletResponse, err error) {
				log.Println("ERROR: ", err)
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

			res, err := server.CreateWallet(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
} 