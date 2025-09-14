package gapi

import (
	"context"
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



 func TestGetWalletRPC(t *testing.T) {

    _, wallet, _ := createRandomWallet()

    testCases := []struct {
        name 		  string
        req        	  *pb.GetWalletRequest
        buildStubs 	  func(store *mockdb.MockStore_interface)
        setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
        checkResponse func(t *testing.T, res *pb.GetWalletResponse, err error)
    }{
        {
            name: "OK",
            req: &pb.GetWalletRequest{
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
            },
            checkResponse: func(t *testing.T, res *pb.GetWalletResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotNil(t, res.Wallet)
				require.Equal(t, wallet.ID.String(), res.Wallet.Id)
				require.Equal(t, wallet.UserEmail, res.Wallet.UserEmail)
			},
        },
        {
            name: "NotFound",
            req: &pb.GetWalletRequest{
				WalletId: wallet.ID.String(),
			},
            setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, wallet.Username, time.Minute, token.TokenTypeAccessToken)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Eq(wallet.ID)).
                    Times(1).
                    Return(db.Wallet{}, db.ErrRecordNotFound)
            },
            checkResponse: func(t *testing.T, res *pb.GetWalletResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
        },
        {
			name: "Unauthorized",
			req: &pb.GetWalletRequest{
				WalletId: wallet.ID.String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return context.Background()
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateWallet(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.GetWalletResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
        {
            name: "InvalidID",
            req: &pb.GetWalletRequest{
				WalletId: "invalid-uuid",
			},
            setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, wallet.Username, time.Minute, token.TokenTypeAccessToken)
            },
            buildStubs : func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(t *testing.T, res *pb.GetWalletResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
        },
        {
            name: "InternalError",
            req: &pb.GetWalletRequest{
				WalletId: wallet.ID.String(),
			},
            setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, wallet.Username, time.Minute, token.TokenTypeAccessToken)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Eq(wallet.ID)).
                    Times(1).
                    Return(db.Wallet{}, sql.ErrConnDone)
            },
            checkResponse: func(t *testing.T, res *pb.GetWalletResponse, err error) {
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

			res, err := server.GetWallet(ctx, tc.req)
			tc.checkResponse(t, res, err)
        })
    }
}

func createRandomWallet() (db.CreateWalletParams, db.Wallet, db.UpdateWalletBalanceParams) {
    rand.New(rand.NewSource(int64(time.Now().UnixNano())))
	currencies := []string{"USD", "EUR", "BTC", "ETH", "LTC"}
	randomCurrency := currencies[rand.Intn(len(currencies))]

	randomEmail := "user" + uuid.New().String() + "@example.com"

	walletArgs := db.CreateWalletParams {
        Username: utils.RandomUser(),
		UserEmail: randomEmail,
		Currency: randomCurrency,
		Balance: decimal.NewFromFloat(0),
	}

	createWalletRows := db.Wallet {
		ID: uuid.New(),
        Username: walletArgs.Username,
		UserEmail: walletArgs.UserEmail,
		Currency: walletArgs.Currency,
		Balance: walletArgs.Balance,
		LockedBalance: decimal.NewFromFloat(0),
		CreatedAt: time.Now(),
	}

	updateWalletParams := db.UpdateWalletBalanceParams {
		Balance: decimal.NewFromFloat(100),
		LockedBalance: decimal.NewFromFloat(0),
	}

	return walletArgs, createWalletRows, updateWalletParams
}