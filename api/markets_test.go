package api

import (
	"bytes"
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
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)


func TestCreateMarketAPI(t *testing.T) {
	marketArgs, market, marketRow := createRandomMarket()

    testCases := []struct {
        name          string
        body          gin.H
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            body: gin.H{
                "base_currency":  market.BaseCurrency,
                "quote_currency": market.QuoteCurrency,
                "min_order_amount": marketArgs.MinOrderAmount,
                "price_precision": marketArgs.PricePrecision,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateMarket(gomock.Any(), gomock.Eq(marketArgs)).
                    Times(1).
                    Return(marketRow, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchMarket(t, recorder.Body, market)
            },
        },
        {
            name: "InternalError",
            body: gin.H{
                "base_currency":  marketArgs.BaseCurrency,
                "quote_currency": marketArgs.QuoteCurrency,
                "min_order_amount": marketArgs.MinOrderAmount,
                "price_precision": marketArgs.PricePrecision,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateMarket(gomock.Any(), gomock.Eq(marketArgs)).
                    Times(1).
                    Return(db.CreateMarketRow{}, sql.ErrConnDone)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "InvalidCurrency",
            body: gin.H{
                "base_currency":  nil,
                "quote_currency": marketArgs.QuoteCurrency,
                "min_order_amount": marketArgs.MinOrderAmount,
                "price_precision": marketArgs.PricePrecision,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateMarket(gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusBadRequest, recorder.Code)
            },
        },
        {
            name: "DuplicateMarket",
            body: gin.H{
                "base_currency":  marketArgs.BaseCurrency,
                "quote_currency": marketArgs.QuoteCurrency,
                "min_order_amount": marketArgs.MinOrderAmount,
                "price_precision": marketArgs.PricePrecision,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateMarket(gomock.Any(), gomock.Eq(marketArgs)).
                    Times(1).
                    Return(db.CreateMarketRow{}, &pq.Error{Code: "23505"})
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

            url := "/markets"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestGetMarketByIDAPI(t *testing.T) {
    _, market, _ := createRandomMarket()

    testCases := []struct {
        name          string
        MarketID      uuid.UUID
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            MarketID: market.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetMarketByID(gomock.Any(), gomock.Eq(market.ID)).
                    Times(1).
                    Return(market, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchMarkets(t, recorder.Body, market)
            },
        },
        {
            name: "NotFound",
            MarketID: market.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetMarketByID(gomock.Any(), gomock.Eq(market.ID)).
                    Times(1).
                    Return(market, sql.ErrNoRows)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name: "InternalError",
            MarketID: market.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetMarketByID(gomock.Any(), gomock.Eq(market.ID)).
                    Times(1).
                    Return(market, sql.ErrConnDone)
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

            url := fmt.Sprintf("/markets/%s", tc.MarketID)
            request, err := http.NewRequest(http.MethodGet, url, nil)
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestDeleteMarketAPI(t *testing.T) {
    _, market, _ := createRandomMarket()

    testCases := []struct {
        name          string
        MarketID      uuid.UUID
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            MarketID: market.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    DeleteMarket(gomock.Any(), gomock.Eq(market.ID)).
                    Times(1).
                    Return(nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "NotFound",
            MarketID: market.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    DeleteMarket(gomock.Any(), gomock.Eq(market.ID)).
                    Times(1).
                    Return(sql.ErrNoRows)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
        {
            name: "InternalError",
            MarketID: market.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    DeleteMarket(gomock.Any(), gomock.Eq(market.ID)).
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

            url := fmt.Sprintf("/markets/%s", tc.MarketID)
            request, err := http.NewRequest(http.MethodDelete, url, nil)
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}


func TestListMarketsAPI(t *testing.T) {
    _, market1, _ := createRandomMarket()
    _, market2, _ := createRandomMarket()
    _, market3, _ := createRandomMarket()

    markets := []db.Market{
        market1,
        market2,
        market3,
    }

    testCases := []struct {
        name          string
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    ListMarkets(gomock.Any()).
                    Times(1).
                    Return(markets, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchMarketsForLists(t, recorder.Body, markets)
            },
        },
        {
            name: "InternalError",
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    ListMarkets(gomock.Any()).
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

            url := "/markets"
            request, err := http.NewRequest(http.MethodGet, url, nil)
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func requireBodyMatchMarket(t *testing.T, body *bytes.Buffer, market db.Market) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotMarket db.CreateMarketRow
    err = json.Unmarshal(data, &gotMarket)
    require.NoError(t, err)

    log.Println("DATA: ", string(data))

    require.Equal(t, market.ID, gotMarket.ID)
}

func requireBodyMatchMarkets(t *testing.T, body *bytes.Buffer, market db.Market) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotMarket db.Market
    err = json.Unmarshal(data, &gotMarket)
    require.NoError(t, err)

    require.Equal(t, market.ID, gotMarket.ID)
    require.Equal(t, market.BaseCurrency, gotMarket.BaseCurrency)
    require.Equal(t, market.QuoteCurrency, gotMarket.QuoteCurrency)
    require.Equal(t, market.MinOrderAmount, gotMarket.MinOrderAmount)
    require.Equal(t, market.PricePrecision, gotMarket.PricePrecision)
    require.Equal(t, market.CreatedAt.Time.Unix(), gotMarket.CreatedAt.Time.Unix())
}

func requireBodyMatchMarketsForLists(t *testing.T, body *bytes.Buffer, markets []db.Market) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotMarkets []db.Market
    err = json.Unmarshal(data, &gotMarkets)
    require.NoError(t, err)

    require.Equal(t, len(markets), len(gotMarkets))
    for i := range markets {
        require.Equal(t, markets[i].ID, gotMarkets[i].ID)
        require.Equal(t, markets[i].BaseCurrency, gotMarkets[i].BaseCurrency)
        require.Equal(t, markets[i].QuoteCurrency, gotMarkets[i].QuoteCurrency)
        require.Equal(t, markets[i].MinOrderAmount, gotMarkets[i].MinOrderAmount)
        require.Equal(t, markets[i].PricePrecision, gotMarkets[i].PricePrecision)
        require.WithinDuration(t, markets[i].CreatedAt.Time, gotMarkets[i].CreatedAt.Time, time.Second)
    }
}

func createRandomMarket() (db.CreateMarketParams, db.Market, db.CreateMarketRow) {
    rand.Seed(uint64(time.Now().UnixNano()))
    currencies := []string{"USD", "EUR", "BTC", "ETH", "JPY"}
    baseCurrency := currencies[rand.Intn(len(currencies))]
    quoteCurrency := currencies[rand.Intn(len(currencies))]

    for baseCurrency == quoteCurrency {
        quoteCurrency = currencies[rand.Intn(len(currencies))]
    }

    marketArgs := db.CreateMarketParams{
        BaseCurrency:  baseCurrency,
        QuoteCurrency: quoteCurrency,
        MinOrderAmount: sql.NullString{
            String: "0.1",
            Valid:  true,
        },
        PricePrecision: sql.NullInt32{
            Int32: 8,
            Valid: true,
        },
    }

    market := db.Market{
        ID:            uuid.New(),
        BaseCurrency:  marketArgs.BaseCurrency,
        QuoteCurrency: marketArgs.QuoteCurrency,
        MinOrderAmount: marketArgs.MinOrderAmount,
        PricePrecision: marketArgs.PricePrecision,
        CreatedAt:     sql.NullTime{Time: time.Now(), Valid: true},
    }

	marketRow := db.CreateMarketRow {
		ID: market.ID,
		BaseCurrency: market.BaseCurrency,
		QuoteCurrency: market.QuoteCurrency,
		CreatedAt: market.CreatedAt,
	}

    return marketArgs, market, marketRow
}