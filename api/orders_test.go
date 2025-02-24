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
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)


func TestCreateOrderAPI(t *testing.T) {
    createOrderParams, order, _ := createRandomOrder()

    testCases := []struct {
        name          string
        body          gin.H
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            body: gin.H{
                "user_email": createOrderParams.UserEmail,
                "market_id":  createOrderParams.MarketID,
                "type":       createOrderParams.Type,
                "status":     createOrderParams.Status,
                "price":      createOrderParams.Price.String,
                "amount":     createOrderParams.Amount,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateOrder(gomock.Any(), gomock.Eq(createOrderParams)).
                    Times(1).
                    Return(order, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchOrder(t, recorder.Body, order)
            },
        },
        {
            name: "InternalError",
            body: gin.H{
                "user_email": createOrderParams.UserEmail,
                "market_id":  createOrderParams.MarketID,
                "type":       createOrderParams.Type,
                "status":     createOrderParams.Status,
                "price":      createOrderParams.Price.String,
                "amount":     createOrderParams.Amount,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateOrder(gomock.Any(), gomock.Eq(createOrderParams)).
                    Times(1).
                    Return(order, sql.ErrConnDone)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "InvalidMarketID",
            body: gin.H{
                "user_email": createOrderParams.UserEmail,
                "market_id":  "invalid-uuid",
                "type":       createOrderParams.Type,
                "status":     createOrderParams.Status,
                "price":      createOrderParams.Price.String,
                "amount":     createOrderParams.Amount,
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

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestGetOrderAPI(t *testing.T) {

    _, order, _ := createRandomOrder()

	testCases := []struct {
		name          string
		orderID       string
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			orderID: order.ID.String(),
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetOrderByID(gomock.Any(), gomock.Eq(order.ID)).
					Times(1).
					Return(order, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchOrder(t, recorder.Body, order)
			},
		},
		{
			name:    "NotFound",
			orderID: order.ID.String(),
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetOrderByID(gomock.Any(), gomock.Eq(order.ID)).
					Times(1).
					Return(db.Order{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InvalidUUID",
			orderID: "invalid-uuid",
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

			url := "/orders/" + tc.orderID
			request, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(nil))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}


func TestDeleteOrderAPI(t *testing.T) {

    _, order, _ := createRandomOrder()

	testCases := []struct {
		name          string
		orderID       string
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			orderID: order.ID.String(),
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteOrder(gomock.Any(), gomock.Eq(order.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:    "NotFound",
			orderID: order.ID.String(),
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteOrder(gomock.Any(), gomock.Eq(order.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:    "InvalidUUID",
			orderID: "invalid-uuid",
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteOrder(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			orderID: uuid.Nil.String(),
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteOrder(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	ctrl := gomock.NewController(t)
    defer ctrl.Finish()

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore_interface(ctrl)
   			server := NewTestServer(t, store)

			tc.buildStubs(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/orders/%s", tc.orderID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}


func createRandomOrder() (createOrderParams db.CreateOrderParams, order db.Order, updatedOrderParams db.UpdateOrderStatusAndFilledAmountParams) {
	email := "hello" + fmt.Sprint(rand.Intn(10000)) + "@example.com"
	marketID := uuid.New()
	orderType := db.OrderType(fmt.Sprint(rand.Intn(2))) 
	orderStatus := db.OrderStatus(fmt.Sprint(rand.Intn(3))) 
	price := sql.NullString{String: "100.50", Valid: true}
	amount := "10"

	createOrderParams = db.CreateOrderParams{
		UserEmail: email,
		MarketID:  marketID,
		Type:      orderType,
		Status:    orderStatus,
		Price:     price,
		Amount:    amount,
	}


	createdOrder := db.Order{
		ID:           uuid.New(),
		UserEmail:    email,
		MarketID:     marketID,
		Type:         orderType,
		Status:       orderStatus,
		Price:        price,
		Amount:       amount,
		FilledAmount: sql.NullString{String: "5", Valid: true}, 
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	updatedOrderParams = db.UpdateOrderStatusAndFilledAmountParams{
		Status:       db.OrderStatus(fmt.Sprint(1)), 
		FilledAmount: sql.NullString{String: "10", Valid: true}, 
		ID:           createdOrder.ID,
	}

	return createOrderParams, createdOrder, updatedOrderParams
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
    require.WithinDuration(t, order.CreatedAt.Time, gotOrder.CreatedAt.Time, time.Second)
}


func requireBodyMatchOrder(t *testing.T, body *bytes.Buffer, order db.Order) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotOrder db.Order
    err = json.Unmarshal(data, &gotOrder)
    require.NoError(t, err)

    require.Equal(t, order.ID, gotOrder.ID)
}