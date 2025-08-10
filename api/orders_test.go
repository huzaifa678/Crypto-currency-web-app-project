package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)


func TestCreateOrderAPI(t *testing.T) {
    createOrderParams, order, _, createOrderRow := createRandomOrder()

    testCases := []struct {
        name          string
        body          gin.H
        setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            body: gin.H{
                "user_name":  createOrderParams.Username,
                "user_email": createOrderParams.UserEmail,
                "market_id":  createOrderParams.MarketID,
                "type":       createOrderParams.Type,
                "status":     createOrderParams.Status,
                "price":      createOrderParams.Price,
                "amount":     createOrderParams.Amount,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, createOrderParams.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateOrder(gomock.Any(), gomock.Eq(createOrderParams)).
                    Times(1).
                    Return(createOrderRow, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchOrder(t, recorder.Body, order)
            },
        },
        {
            name: "InternalError",
            body: gin.H{
                "user_name":  createOrderParams.Username,
                "user_email": createOrderParams.UserEmail,
                "market_id":  createOrderParams.MarketID,
                "type":       createOrderParams.Type,
                "status":     createOrderParams.Status,
                "price":      createOrderParams.Price,
                "amount":     createOrderParams.Amount,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, createOrderParams.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateOrder(gomock.Any(), gomock.Eq(createOrderParams)).
                    Times(1).
                    Return(createOrderRow, sql.ErrConnDone)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "InvalidMarketID",
            body: gin.H{
                "user_name":  createOrderParams.Username,
                "user_email": createOrderParams.UserEmail,
                "market_id":  "invalid-uuid",
                "type":       createOrderParams.Type,
                "status":     createOrderParams.Status,
                "price":      createOrderParams.Price,
                "amount":     createOrderParams.Amount,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, createOrderParams.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateOrder(gomock.Any(), gomock.Any()).
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

            url := "/orders"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            tc.setupAuth(t, request, server.tokenMaker)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestGetOrderAPI(t *testing.T) {

    _, order, _, _ := createRandomOrder()

	testCases := []struct {
		name          string
		orderID       string
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			orderID: order.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, order.Username, time.Minute)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetOrderByID(gomock.Any(), gomock.Eq(order.ID)).
					Times(1).
					Return(order, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchOrderForGet(t, recorder.Body, order)
			},
		},
		{
			name:    "NotFound",
			orderID: order.ID.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, order.Username, time.Minute)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetOrderByID(gomock.Any(), gomock.Eq(order.ID)).
					Times(1).
					Return(order, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InvalidUUID",
			orderID: "invalid-uuid",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, order.Username, time.Minute)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetOrderByID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			orderID: uuid.Nil.String(),
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, order.Username, time.Minute)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetOrderByID(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/orders/%s", tc.orderID)
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(nil))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteOrderAPI(t *testing.T) {
    trade, _ := createRandomTrade()

    testCases := []struct {
        name          string
        tradeID       uuid.UUID
        buildStubs    func(store *mockdb.MockStore_interface)
        setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name:    "OK",
            tradeID: trade.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
                    Times(1).
                    Return(trade, nil)

				store.EXPECT().
                    DeleteTrade(gomock.Any(), gomock.Eq(trade.ID)).
                    Times(1).
                    Return(nil)
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, trade.Username, time.Minute)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name:    "NotFound",
            tradeID: trade.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
                    Times(1).
                    Return(trade, sql.ErrNoRows)

				store.EXPECT().
                    DeleteTrade(gomock.Any(), gomock.Eq(trade.ID)).
                    Times(0)
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, trade.Username, time.Minute)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name:    "InvalidUUID",
            tradeID: uuid.Nil,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetTradeByID(gomock.Any(), gomock.Any()).
                    Times(0)

				store.EXPECT().
                    DeleteTrade(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, trade.Username, time.Minute)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name:    "InternalError",
            tradeID: trade.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
                    Times(1).
                    Return(trade, nil)

				store.EXPECT().
                    DeleteTrade(gomock.Any(), gomock.Eq(trade.ID)).
                    Times(1).
                    Return(sql.ErrConnDone)
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, trade.Username, time.Minute)
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

            url := fmt.Sprintf("/trades/%s", tc.tradeID)
            request, err := http.NewRequest(http.MethodDelete, url, nil)
            require.NoError(t, err)

            tc.setupAuth(t, request, server.tokenMaker)
            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}


func createRandomOrder() (createOrderParams db.CreateOrderParams, order db.Order, updatedOrderParams db.UpdateOrderStatusAndFilledAmountParams, createOrderRow db.CreateOrderRow) {
	username := utils.RandomUser()
	email := "hello" + fmt.Sprint(rand.Intn(10000)) + "@example.com"
	marketID := uuid.New()
	orderType := db.OrderType(fmt.Sprint(rand.Intn(2))) 
	orderStatus := db.OrderStatus(fmt.Sprint(rand.Intn(3))) 
	price := "100.50"
	amount := "10"

	createOrderParams = db.CreateOrderParams{
		Username: username,
		UserEmail: email,
		MarketID:  marketID,
		Type:      orderType,
		Status:    orderStatus,
		Price:     price,
		Amount:    amount,
	}


	createdOrder := db.Order{
		ID:           uuid.New(),
		Username: 	  username,
		UserEmail:    email,
		MarketID:     marketID,
		Type:         orderType,
		Status:       orderStatus,
		Price:        price,
		Amount:       amount,
		FilledAmount: "5", 
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	updatedOrderParams = db.UpdateOrderStatusAndFilledAmountParams{
		Status:       db.OrderStatus(fmt.Sprint(1)), 
		FilledAmount: "10", 
		ID:           createdOrder.ID,
	}

	createOrderRow = db.CreateOrderRow {
		ID:           createdOrder.ID,
		UserEmail:    email,
		MarketID:     marketID,
		Type:         orderType,
		Status:       orderStatus,
		Price:        price,
		Amount:       amount,
		FilledAmount: "5", 
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return createOrderParams, createdOrder, updatedOrderParams, createOrderRow
}


func requireBodyMatchOrderForGet(t *testing.T, body *bytes.Buffer, order db.Order) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotOrder db.Order
    err = json.Unmarshal(data, &gotOrder)
    require.NoError(t, err)

    require.Equal(t, order.ID, gotOrder.ID)
    require.Equal(t, order.UserEmail, order.UserEmail)
    require.Equal(t, order.MarketID, gotOrder.MarketID)
    require.Equal(t, order.Type, gotOrder.Type)
    require.Equal(t, order.Status, gotOrder.Status)
    require.Equal(t, order.Price, order.Price)
    require.Equal(t, order.Amount, gotOrder.Amount)
    require.Equal(t, order.FilledAmount, gotOrder.FilledAmount)
    require.WithinDuration(t, order.CreatedAt, gotOrder.CreatedAt, time.Second)
}


func requireBodyMatchOrder(t *testing.T, body *bytes.Buffer, order db.Order) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotOrder db.Order
    err = json.Unmarshal(data, &gotOrder)
    require.NoError(t, err)

    require.Equal(t, order.ID, gotOrder.ID)
}