package gapi

import (
	"context"
	"fmt"
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

func TestCreateTransactionRPC(t *testing.T) {
	transaction, transactionRow, createTxParams := createRandomTransaction()

	testCases := []struct {
		name          string
		req           *pb.CreateTransactionRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.CreateTransactionResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateTransactionRequest{
				UserEmail: createTxParams.UserEmail,
				Type:      pb.TransactionType_DEPOSIT,
				Currency:  transaction.Currency,
				Amount:    transaction.Amount.Mul(decimal.New(1, scale)).IntPart(),
				Address:   transaction.Address,
				TxHash:    transaction.TxHash,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, transaction.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, arg db.CreateTransactionParams) (db.CreateTransactionRow, error) {
						require.True(t, createTxParams.Amount.Equal(arg.Amount))
						require.Equal(t, createTxParams.UserEmail, arg.UserEmail)
						require.Equal(t, createTxParams.Type, arg.Type)
						require.Equal(t, createTxParams.Currency, arg.Currency)
						require.Equal(t, createTxParams.Address, arg.Address)
						require.Equal(t, createTxParams.TxHash, arg.TxHash)
						return transactionRow, nil
					}).
					Times(1)
			},
			checkResponse: func(t *testing.T, res *pb.CreateTransactionResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotNil(t, res.Transaction)
				require.Equal(t, transactionRow.ID.String(), res.Transaction.GetTransactionId())
			},
		},
		{
			name: "Unauthorized",
			req: &pb.CreateTransactionRequest{
				UserEmail: transaction.UserEmail,
				Type:      pb.TransactionType_DEPOSIT,
				Currency:  transaction.Currency,
				Amount:    transaction.Amount.IntPart(),
				Address:   transaction.Address,
				TxHash:    transaction.TxHash,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return context.Background()
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateTransactionResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidEmail",
			req: &pb.CreateTransactionRequest{
				UserEmail: "invalid-email",
				Type:      pb.TransactionType_DEPOSIT,
				Currency:  "BTC",
				Amount:    10,
				Address:   "1BitcoinAddress",
				TxHash:    uuid.New().String(),
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, transaction.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateTransactionResponse, err error) {
				log.Println("ERROR: ", err)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateTransactionRequest{
				UserEmail: transaction.UserEmail,
				Type:      pb.TransactionType_DEPOSIT,
				Currency:  transaction.Currency,
				Amount:    transaction.Amount.IntPart(),
				Address:   transaction.Address,
				TxHash:    transaction.TxHash,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(t, tokenMaker, transaction.Username, time.Minute, token.TokenTypeAccessToken)
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateTransaction(gomock.Any(), gomock.Any()).
					Times(1).
					Return(transactionRow, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.CreateTransactionResponse, err error) {
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

			res, err := server.CreateTransaction(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}

func createRandomTransaction() (db.Transaction, db.CreateTransactionRow, db.CreateTransactionParams) {
	username := fmt.Sprintf("example-%s", uuid.New().String())
	email := fmt.Sprintf("example-%s@example.com", uuid.New().String())
	txType := db.TransactionType("deposit")
	txStatus := db.TransactionStatus("pending")
	currency := "BTC"
	amount := decimal.NewFromFloat(10.5)
	address := "1BitcoinAddress"
	txHash := uuid.New().String()

	createTxParams := db.CreateTransactionParams{
		Username:  username,
		UserEmail: email,
		Type:      txType,
		Currency:  currency,
		Amount:    amount,
		Address:   address,
		TxHash:    txHash,
	}

	transactionRow := db.CreateTransactionRow{
		ID:        uuid.New(),
		UserEmail: email,
		Type:      txType,
		Currency:  currency,
		Amount:    amount,
		Status:    txStatus,
		Address:   address,
		TxHash:    txHash,
		CreatedAt: time.Now(),
	}

	transaction := db.Transaction{
		ID:        uuid.New(),
		Username:  username,
		UserEmail: email,
		Type:      txType,
		Currency:  currency,
		Amount:    amount,
		Status:    txStatus,
		Address:   address,
		TxHash:    txHash,
		CreatedAt: time.Now(),
	}

	return transaction, transactionRow, createTxParams
}
