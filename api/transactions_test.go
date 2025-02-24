package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/stretchr/testify/require"
)


func TestCreateTransactionAPI(t *testing.T) {
    transaction, transactionArgs := createRandomTransaction()

	log.Println(transactionArgs)

    testCases := []struct {
        name          string
        body          gin.H
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            body: gin.H{
                "user_email":     transactionArgs.UserEmail,
                "type":           transactionArgs.Type,
                "currency":       transactionArgs.Currency,
                "amount":         transactionArgs.Amount,
                "address":        transactionArgs.Address.String,
                "tx_hash":        transactionArgs.TxHash.String,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateTransaction(gomock.Any(), gomock.Eq(transactionArgs)).
                    Times(1).
                    Return(transaction, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchTransaction(t, recorder.Body, transaction)
            },
        },
        {
            name: "InternalError",
            body: gin.H{
                "user_email":     transactionArgs.UserEmail,
                "type":           transactionArgs.Type,
                "currency":       transactionArgs.Currency,
                "amount":         transactionArgs.Amount,
                "address":        transactionArgs.Address.String,
                "tx_hash":        transactionArgs.TxHash.String,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateTransaction(gomock.Any(), gomock.Eq(transactionArgs)).
                    Times(1).
                    Return(transaction, sql.ErrConnDone)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "InvalidCurrency",
            body: gin.H{
                "user_email":     transactionArgs.UserEmail,
                "type":           transactionArgs.Type,
                "currency":       "INVALID",
                "amount":         transactionArgs.Amount,
                "address":        transactionArgs.Address.String,
                "tx_hash":        transactionArgs.TxHash.String,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateTransaction(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
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
            recorder := httptest.NewRecorder()

            data, err := json.Marshal(tc.body)
            require.NoError(t, err)

            url := "/transactions"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestGetTransactionAPI(t *testing.T) {

    transaction, _ := createRandomTransaction()
    
    testCases := []struct {
        name          string
        txID          uuid.UUID
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            txID: transaction.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                GetTransactionByID(gomock.Any(), gomock.Eq(transaction.ID)).
                Times(1).
                Return(transaction, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchTransactions(t, recorder.Body, transaction)
            },
        },
        {
            name:   "NotFound",
			txID: transaction.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(transaction, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
        },
        {
            name: "InvalidID",
            txID: uuid.Nil,
            buildStubs : func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetTransactionByID(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name:   "InternalError",
			txID: transaction.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTransactionByID(gomock.Any(), gomock.Eq(transaction.ID)).
					Times(1).
					Return(transaction, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
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
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/transactions/%s", tc.txID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListUserTransactionsAPI(t *testing.T) {
    userEmail := "mock123@example.com"

	transaction1 := createRandomTransactionWithEmail(userEmail)
	transaction2 := createRandomTransactionWithEmail(userEmail)
    transactions := []db.Transaction{
        transaction1,
        transaction2,
    }

    testCases := []struct {
        name          string
        userEmail     string
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name:      "OK",
            userEmail: userEmail,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetTransactionsByUserEmail(gomock.Any(), gomock.Eq(userEmail)).
                    Times(1).
                    Return(transactions, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name:      "InternalError",
            userEmail: userEmail,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetTransactionsByUserEmail(gomock.Any(), gomock.Eq(userEmail)).
                    Times(1).
                    Return(nil, sql.ErrConnDone)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
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
            recorder := httptest.NewRecorder()

            url := fmt.Sprintf("/transactions/user/%s", tc.userEmail)
            request, err := http.NewRequest(http.MethodGet, url, nil)
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestDeleteTransactionAPI(t *testing.T) {
    transaction, _ := createRandomTransaction()
    
    testCases := []struct {
        name          string
        txID          uuid.UUID
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            txID: transaction.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    DeleteTransaction(gomock.Any(), gomock.Eq(transaction.ID)).
                    Times(1).
                    Return(nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "NotFound",
            txID: transaction.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    DeleteTransaction(gomock.Any(), gomock.Eq(transaction.ID)).
                    Times(1).
                    Return(sql.ErrNoRows)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name: "InvalidID",
            txID: uuid.Nil,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    DeleteTransaction(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "InternalError",
            txID: transaction.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    DeleteTransaction(gomock.Any(), gomock.Eq(transaction.ID)).
                    Times(1).
                    Return(sql.ErrConnDone)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
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
            recorder := httptest.NewRecorder()

            url := fmt.Sprintf("/transactions/%s", tc.txID)
            request, err := http.NewRequest(http.MethodDelete, url, nil)
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}



func requireBodyMatchTransaction(t *testing.T, body *bytes.Buffer, transactionArgs db.Transaction) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotTransaction db.Transaction
    err = json.Unmarshal(data, &gotTransaction)
    require.NoError(t, err)

    require.Equal(t, transactionArgs.ID, gotTransaction.ID)
}

func requireBodyMatchTransactions(t *testing.T, body *bytes.Buffer, transactionArgs db.Transaction) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotTransaction db.Transaction
    err = json.Unmarshal(data, &gotTransaction)
    require.NoError(t, err)

    require.Equal(t, transactionArgs.ID, gotTransaction.ID)
    require.Equal(t, transactionArgs.UserEmail, gotTransaction.UserEmail)
    require.Equal(t, transactionArgs.Type, gotTransaction.Type)
    require.Equal(t, transactionArgs.Currency, gotTransaction.Currency)
    require.Equal(t, transactionArgs.Amount, gotTransaction.Amount)
    require.Equal(t, transactionArgs.Status, gotTransaction.Status)
    require.Equal(t, transactionArgs.Address, gotTransaction.Address)
    require.Equal(t, transactionArgs.TxHash, gotTransaction.TxHash)
    require.WithinDuration(t, transactionArgs.CreatedAt.Time, gotTransaction.CreatedAt.Time, time.Second)
}

func createRandomTransaction() (transaction db.Transaction, transactionArgs db.CreateTransactionParams) {
	rand.Seed(time.Now().UnixNano())

	id := uuid.New()

	email := "hello" + fmt.Sprint(rand.Intn(10000)) + "@example.com"

	txTypes := []string{"withdraw", "deposit"}

    txType := db.TransactionType(txTypes[rand.Intn(len(txTypes))])

	currencies := []string{"USD", "EUR", "BTC", "ETH", "JPY"}

	currency := currencies[rand.Intn(len(currencies))]

	amount := fmt.Sprint(rand.Intn(10000))

	status := []string{"completed", "pending"}

	txStatus := db.TransactionStatus(status[rand.Intn(2)])

	randomAddress := []string{"x", "y", "z"}

	address := sql.NullString{String: randomAddress[rand.Intn(len(randomAddress))], Valid: true}

	txHash := sql.NullString{String: RandomString(64), Valid: true}

	createdAt := sql.NullTime{Time: time.Now(), Valid: true}

	transactions := db.Transaction{
		ID: id,
		UserEmail: email,
		Type: txType,
		Currency: currency,
		Amount: amount,
		Status: txStatus,
		Address: address,
		TxHash: txHash,
		CreatedAt: createdAt,
	}

	args := db.CreateTransactionParams{
		UserEmail: email,
		Type: txType,
		Currency: currency,
		Amount: amount,
		Address: address,
		TxHash: txHash,
	}

	return transactions, args
}

func createRandomTransactionWithEmail(Email string) (transaction db.Transaction) {
	rand.Seed(time.Now().UnixNano())

	id := uuid.New()

	email := Email

	txTypes := []string{"withdraw", "deposit"}

    txType := db.TransactionType(txTypes[rand.Intn(len(txTypes))])

	currencies := []string{"USD", "EUR", "BTC", "ETH", "JPY"}

	currency := currencies[rand.Intn(len(currencies))]

	amount := fmt.Sprint(rand.Intn(10000))

	status := []string{"completed", "pending"}

	txStatus := db.TransactionStatus(status[rand.Intn(2)])

	randomAddress := []string{"x", "y", "z"}

	address := sql.NullString{String: randomAddress[rand.Intn(len(randomAddress))], Valid: true}

	txHash := sql.NullString{String: RandomString(64), Valid: true}

	createdAt := sql.NullTime{Time: time.Now(), Valid: true}

	transactions := db.Transaction{
		ID: id,
		UserEmail: email,
		Type: txType,
		Currency: currency,
		Amount: amount,
		Status: txStatus,
		Address: address,
		TxHash: txHash,
		CreatedAt: createdAt,
	}

	return transactions
}

func RandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTVWXYZ01234567890"

	var sb strings.Builder

	for i := 0; i < length; i++ {
		sb.WriteByte(letters[rand.Intn(len(letters))])
	}
	return sb.String()
}