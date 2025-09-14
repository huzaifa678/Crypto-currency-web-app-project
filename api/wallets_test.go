package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

func TestCreateWalletAPI(t *testing.T) {

    walletArgs, wallet, _ := createRandomWallet()

	log.Println("EMAIL: ", walletArgs.UserEmail)

    testCases := []struct {
        name          string
        body          gin.H
        buildStubs    func(store *mockdb.MockStore_interface)
        setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            body: gin.H{
                "user_email": walletArgs.UserEmail,
                "currency":   walletArgs.Currency,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateWallet(gomock.Any(), gomock.Eq(walletArgs)).
                    Times(1).
                    Return(wallet, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchWallet(t, recorder.Body, wallet)
            },
        },
        {
            name: "InternalError",
            body: gin.H{
                "user_email": walletArgs.UserEmail,
                "currency":   walletArgs.Currency,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateWallet(gomock.Any(), gomock.Eq(walletArgs)).
                    Times(1).
                    Return(db.Wallet{}, sql.ErrConnDone)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "InvalidEmail",
            body: gin.H{
                "user_email": "invalid-email",
                "currency":   walletArgs.Currency,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateWallet(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "DuplicateEmail",
            body: gin.H{
                "user_email": walletArgs.UserEmail,
                "currency":   walletArgs.Currency,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateWallet(gomock.Any(), gomock.Eq(walletArgs)).
                    Times(1).
                    Return(db.Wallet{}, &pq.Error{Code: "23505"})
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

            data, err := json.Marshal(tc.body)
            require.NoError(t, err)

            url := "/wallets"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            tc.setupAuth(t, request, server.tokenMaker)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func requireBodyMatchWallet(t *testing.T, body *bytes.Buffer, wallet db.Wallet) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotWallet db.Wallet
    err = json.Unmarshal(data, &gotWallet)
    require.NoError(t, err)

    require.Equal(t, wallet.ID, gotWallet.ID)
}

func TestGetWalletAPI(t *testing.T) {

    _, wallet, _ := createRandomWallet()

    testCases := []struct {
        name string
        WalletID uuid.UUID
        buildStubs func(store *mockdb.MockStore_interface)
        setupAuth  func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            WalletID: wallet.ID,
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Eq(wallet.ID)).
                    Times(1).
                    Return(wallet, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchWallet(t, recorder.Body, wallet)
            },
        },
        {
            
            name: "NotFound",
            WalletID: wallet.ID,
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Eq(wallet.ID)).
                    Times(1).
                    Return(db.Wallet{}, sql.ErrNoRows)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },

        },
        {
            name: "InvalidID",
            WalletID: uuid.Nil,
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            buildStubs : func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "InternalError",
            WalletID: wallet.ID,
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Eq(wallet.ID)).
                    Times(1).
                    Return(db.Wallet{}, sql.ErrConnDone)
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

			url := fmt.Sprintf("/wallets/%s", tc.WalletID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

            tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
        })
    }
}


func createRandomWallet() (db.CreateWalletParams, db.Wallet, db.UpdateWalletBalanceParams) {
    rand.Seed(uint64(time.Now().UnixNano()))
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
        ID: createWalletRows.ID,
	}

	return walletArgs, createWalletRows, updateWalletParams
}


func TestUpdateWalletAPI(t *testing.T) {
    _, wallet, updateWalletParams := createRandomWallet()
    ID := wallet.ID

    testCases := []struct {
        name          string
        body          gin.H
        buildStubs    func(store *mockdb.MockStore_interface)
        setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            body: gin.H{
                "balance":        updateWalletParams.Balance,
                "locked_balance": updateWalletParams.LockedBalance,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Eq(ID)).
                    Times(1).
                    Return(wallet, nil)

                store.EXPECT().
                    UpdateWalletBalance(gomock.Any(), gomock.Any()).
                        DoAndReturn(func(ctx context.Context, arg db.UpdateWalletBalanceParams) error {
                            require.True(t, arg.Balance.Equal(updateWalletParams.Balance))
                            require.True(t, arg.LockedBalance.Equal(updateWalletParams.LockedBalance))
                            require.Equal(t, updateWalletParams.ID, arg.ID)
                        return nil
                    }).
                    Times(1)
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "InternalError",
            body: gin.H{
                "balance":        updateWalletParams.Balance,
                "locked_balance": updateWalletParams.LockedBalance,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {

                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Eq(ID)).
                    Times(1).
                    Return(wallet, nil)

                store.EXPECT().
                    UpdateWalletBalance(gomock.Any(), gomock.Any()).
                    DoAndReturn(func(ctx context.Context, arg db.UpdateWalletBalanceParams) error {
                        require.True(t, arg.Balance.Equal(updateWalletParams.Balance))
                        require.True(t, arg.LockedBalance.Equal(updateWalletParams.LockedBalance))
                        require.Equal(t, updateWalletParams.ID, arg.ID)
                    return sql.ErrConnDone
                }).
                Times(1)
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
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

            data, err := json.Marshal(tc.body)
            require.NoError(t, err)

            url := fmt.Sprintf("/wallets/%s", wallet.ID)
            request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
            require.NoError(t, err)

            tc.setupAuth(t, request, server.tokenMaker)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestDeleteWalletAPI(t *testing.T) {
    _, wallet, _ := createRandomWallet()
    ID := wallet.ID

    testCases := []struct {
        name          string
        WalletID      uuid.UUID
        buildStubs    func(store *mockdb.MockStore_interface)
        setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            WalletID: wallet.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Eq(ID)).
                    Times(1).
                    Return(wallet, nil)

                store.EXPECT().
                    DeleteWallet(gomock.Any(), gomock.Eq(wallet.ID)).
                    Times(1).
                    Return(nil)
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "NotFound",
            WalletID: wallet.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Eq(ID)).
                    Times(1).
                    Return(wallet, sql.ErrNoRows)

                store.EXPECT().
                    DeleteWallet(gomock.Any(), gomock.Eq(wallet.ID)).
                    Times(0)
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name: "InvalidID",
            WalletID: uuid.Nil,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Eq(ID)).
                    Times(0)

                store.EXPECT().
                    DeleteWallet(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "InternalError",
            WalletID: wallet.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetWalletByID(gomock.Any(), gomock.Eq(ID)).
                    Times(1).
                    Return(wallet, nil)

                store.EXPECT().
                    DeleteWallet(gomock.Any(), gomock.Eq(wallet.ID)).
                    Times(1).
                    Return(sql.ErrConnDone)
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, wallet.Username, time.Minute)
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

            url := fmt.Sprintf("/wallets/%s", tc.WalletID)
            request, err := http.NewRequest(http.MethodDelete, url, nil)
            require.NoError(t, err)

            tc.setupAuth(t, request, server.tokenMaker)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}