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
)

func TestCreateTradeAPI(t *testing.T) {
	trade, createTradeParams := createRandomTrade()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"buy_order_id":  createTradeParams.BuyOrderID,
				"sell_order_id": createTradeParams.SellOrderID,
				"market_id":     createTradeParams.MarketID,
				"price":         createTradeParams.Price,
				"amount":        createTradeParams.Amount,
				"fee":           createTradeParams.Fee.String,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateTrade(gomock.Any(), gomock.Eq(createTradeParams)).
					Times(1).
					Return(trade, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTrade(t, recorder.Body, trade)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"buy_order_id":  trade.BuyOrderID,
				"sell_order_id": trade.SellOrderID,
				"market_id":     trade.MarketID,
				"price":         trade.Price,
				"amount":        trade.Amount,
				"fee":           trade.Fee.String,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateTrade(gomock.Any(), gomock.Eq(createTradeParams)).
					Times(1).
					Return(trade, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidFee",
			body: gin.H{
				"buy_order_id":  trade.BuyOrderID,
				"sell_order_id": trade.SellOrderID,
				"market_id":     trade.MarketID,
				"price":         trade.Price,
				"amount":        trade.Amount,
				"fee":           "invalid-fee",
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateTrade(gomock.Any(), gomock.Any()).
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

			url := "/trades"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetTradeAPI(t *testing.T) {
	trade, _ := createRandomTrade()

	testCases := []struct {
		name          string
		tradeID       uuid.UUID
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			tradeID: trade.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(trade, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTradeForGet(t, recorder.Body, trade)
			},
		},
		{
			name: "NotFound",
			tradeID: trade.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(trade, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			tradeID: uuid.Nil,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			tradeID: trade.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradeByID(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(trade, sql.ErrConnDone)
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
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteTradeAPI(t *testing.T) {
	trade, _ := createRandomTrade()

	testCases := []struct {
		name          string
		tradeID       uuid.UUID
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:    "OK",
			tradeID: trade.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteTrade(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(nil)
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
					DeleteTrade(gomock.Any(), gomock.Eq(trade.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:    "InvalidID",
			tradeID: uuid.Nil,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteTrade(gomock.Any(), gomock.Any()).
					Times(0)
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
					DeleteTrade(gomock.Any(), gomock.Eq(trade.ID)).
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

			url := fmt.Sprintf("/trades/%s", tc.tradeID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestListTradesByMarketIDAPI(t *testing.T) {
	marketID := uuid.New()
	trade1 := createRandomTradeWithMarketID(marketID)
	trade2 := createRandomTradeWithMarketID(marketID)
	trades := []db.Trade{trade1, trade2}

	testCases := []struct {
		name          string
		marketID      uuid.UUID
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			marketID: marketID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradesByMarketID(gomock.Any(), gomock.Eq(marketID)).
					Times(1).
					Return(trades, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTradesList(t, recorder.Body, trades)
			},
		},
		{
			name:     "InternalError",
			marketID: marketID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradesByMarketID(gomock.Any(), gomock.Eq(marketID)).
					Times(1).
					Return(nil, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "EmptyList",
			marketID: marketID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetTradesByMarketID(gomock.Any(), gomock.Eq(marketID)).
					Times(1).
					Return([]db.Trade{}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchTradesList(t, recorder.Body, []db.Trade{})
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

			url := fmt.Sprintf("/trades/market/%s", tc.marketID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func requireBodyMatchTrade(t *testing.T, body *bytes.Buffer, trade db.Trade) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTrade db.Trade
	err = json.Unmarshal(data, &gotTrade)
	require.NoError(t, err)

	require.Equal(t, trade.ID, gotTrade.ID)
}

func requireBodyMatchTradeForGet(t *testing.T, body *bytes.Buffer, trade db.Trade) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTrade db.Trade
	err = json.Unmarshal(data, &gotTrade)
	require.NoError(t, err)

	require.Equal(t, trade.ID, gotTrade.ID)
	require.Equal(t, trade.BuyOrderID, gotTrade.BuyOrderID)
	require.Equal(t, trade.SellOrderID, gotTrade.SellOrderID)
	require.Equal(t, trade.MarketID, gotTrade.MarketID)
	require.Equal(t, trade.Price, gotTrade.Price)
	require.Equal(t, trade.Amount, gotTrade.Amount)
	require.Equal(t, trade.Fee, gotTrade.Fee)
	require.WithinDuration(t, trade.CreatedAt.Time, gotTrade.CreatedAt.Time, time.Second)
}

func requireBodyMatchTradesList(t *testing.T, body *bytes.Buffer, trades []db.Trade) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotTrades []db.Trade
	err = json.Unmarshal(data, &gotTrades)
	require.NoError(t, err)

	require.Equal(t, len(trades), len(gotTrades))
	for i := range trades {
		require.Equal(t, trades[i].ID, gotTrades[i].ID)
		require.Equal(t, trades[i].BuyOrderID, gotTrades[i].BuyOrderID)
		require.Equal(t, trades[i].SellOrderID, gotTrades[i].SellOrderID)
		require.Equal(t, trades[i].MarketID, gotTrades[i].MarketID)
		require.Equal(t, trades[i].Price, gotTrades[i].Price)
		require.Equal(t, trades[i].Amount, gotTrades[i].Amount)
		require.Equal(t, trades[i].Fee, gotTrades[i].Fee)
		require.WithinDuration(t, trades[i].CreatedAt.Time, gotTrades[i].CreatedAt.Time, time.Second)
	}
}

func createRandomTrade() (trade db.Trade, createTradeParams db.CreateTradeParams) {

	_, sellOrder, _ := createRandomOrder()
	_, BuyOrder, _ := createRandomOrder()
	_, market, _ := createRandomMarket()

	createTradeParams = db.CreateTradeParams {
		BuyOrderID: BuyOrder.ID,
    	SellOrderID: sellOrder.ID,   
    	MarketID:    market.ID,      
    	Price:       "0.0",   
    	Amount:      "0.0",         
    	Fee:         sql.NullString{String: "5", Valid: true},
	}

	Trade := db.Trade {
		ID: uuid.New(),
		BuyOrderID: BuyOrder.ID,
    	SellOrderID: sellOrder.ID,   
    	MarketID:    market.ID,      
    	Price:       createTradeParams.Price,   
    	Amount:      createTradeParams.Amount,         
    	Fee:         createTradeParams.Fee,
		CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
	}

	return Trade, createTradeParams
}

func createRandomTradeWithMarketID(marketID uuid.UUID) db.Trade {
	trade, _ := createRandomTrade()
	trade.MarketID = marketID
	return trade
}


