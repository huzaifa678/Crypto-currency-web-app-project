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

func createRandomFee() (db.CreateFeeParams, db.Fee) {
    marketID := uuid.New()

    feeArgs := db.CreateFeeParams{
        MarketID: marketID,
        MakerFee: sql.NullString{String: "0.01", Valid: true},
        TakerFee: sql.NullString{String: "0.02", Valid: true},
    }

    fee := db.Fee{
        ID:        uuid.New(),
        MarketID:  marketID,
        MakerFee:  feeArgs.MakerFee,
        TakerFee:  feeArgs.TakerFee,
        CreatedAt: sql.NullTime{Time: time.Now(), Valid: true},
    }

    return feeArgs, fee
}

func TestCreateFeeAPI(t *testing.T) {
    feeArgs, fee := createRandomFee()

    testCases := []struct {
        name          string
        body          gin.H
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            body: gin.H{
                "market_id": feeArgs.MarketID,
                "maker_fee": feeArgs.MakerFee.String,
                "taker_fee": feeArgs.TakerFee.String,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateFee(gomock.Any(), gomock.Eq(feeArgs)).
                    Times(1).
                    Return(fee, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchFee(t, recorder.Body, fee)
            },
        },
        {
            name: "InternalError",
            body: gin.H{
                "market_id": feeArgs.MarketID,
                "maker_fee": feeArgs.MakerFee.String,
                "taker_fee": feeArgs.TakerFee.String,
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateFee(gomock.Any(), gomock.Eq(feeArgs)).
                    Times(1).
                    Return(db.Fee{}, sql.ErrConnDone)
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

            server := NewServer(store)
            recorder := httptest.NewRecorder()

            data, err := json.Marshal(tc.body)
            require.NoError(t, err)

            url := "/fees"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestGetFeeAPI(t *testing.T) {
    _, fee := createRandomFee()

    testCases := []struct {
        name          string
        marketID      uuid.UUID
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            marketID: fee.MarketID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetFeeByMarketID(gomock.Any(), gomock.Eq(fee.MarketID)).
                    Times(1).
                    Return(fee, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchFeeForGet(t, recorder.Body, fee)
            },
        },
        {
            name:     "NotFound",
            marketID: fee.MarketID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetFeeByMarketID(gomock.Any(), gomock.Eq(fee.MarketID)).
                    Times(1).
                    Return(db.Fee{}, sql.ErrNoRows)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name:     "InternalError",
            marketID: fee.MarketID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetFeeByMarketID(gomock.Any(), gomock.Eq(fee.MarketID)).
                    Times(1).
                    Return(db.Fee{}, sql.ErrConnDone)
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

            server := NewServer(store)
            recorder := httptest.NewRecorder()

            url := fmt.Sprintf("/fees/%s", tc.marketID)
            request, err := http.NewRequest(http.MethodGet, url, nil)
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestDeleteFeeAPI(t *testing.T) {
    _, fee := createRandomFee()

    testCases := []struct {
        name          string
        feeID        uuid.UUID
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name:   "OK",
            feeID: fee.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    DeleteFee(gomock.Any(), gomock.Eq(fee.ID)).
                    Times(1).
                    Return(nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name:   "NotFound",
            feeID: fee.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    DeleteFee(gomock.Any(), gomock.Eq(fee.ID)).
                    Times(1).
                    Return(sql.ErrNoRows)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name:   "InternalError",
            feeID: fee.ID,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    DeleteFee(gomock.Any(), gomock.Eq(fee.ID)).
                    Times(1).
                    Return(sql.ErrConnDone)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
		{
			name: "Invalid ID",
			feeID: uuid.Nil,
			buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    DeleteFee(gomock.Any(), gomock.Any()).
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

            server := NewServer(store)
            recorder := httptest.NewRecorder()

            url := fmt.Sprintf("/fees/%s", tc.feeID)
            request, err := http.NewRequest(http.MethodDelete, url, nil)
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func requireBodyMatchFee(t *testing.T, body *bytes.Buffer, fee db.Fee) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotFee db.Fee
    err = json.Unmarshal(data, &gotFee)
    require.NoError(t, err)

    require.Equal(t, fee.ID, gotFee.ID)
}

func requireBodyMatchFeeForGet(t *testing.T, body *bytes.Buffer, fee db.Fee) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotFee db.Fee
    err = json.Unmarshal(data, &gotFee)
    require.NoError(t, err)

    require.Equal(t, fee.ID, gotFee.ID)
    require.Equal(t, fee.MarketID, gotFee.MarketID)
    require.Equal(t, fee.MakerFee, gotFee.MakerFee)
    require.Equal(t, fee.TakerFee, gotFee.TakerFee)
    require.WithinDuration(t, fee.CreatedAt.Time, gotFee.CreatedAt.Time, time.Second)
}
