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
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAuditLog() (db.CreateAuditLogParams, db.AuditLog) {
    id := uuid.New()
    email := "user@example.com"
    action := "login"
    ipAddress := "192.168.1.1"
    createdAt := sql.NullTime{Time: time.Now(), Valid: true}

    auditLogArgs := db.CreateAuditLogParams{
        UserEmail:  email,
        Action:     action,
        IpAddress:  sql.NullString{String: ipAddress, Valid: true},
    }

    auditLog := db.AuditLog{
        ID:         id,
        UserEmail:  email,
        Action:     action,
        IpAddress:  sql.NullString{String: ipAddress, Valid: true},
        CreatedAt:  createdAt,
    }

    return auditLogArgs, auditLog
}

func TestCreateAuditLogAPI(t *testing.T) {
    auditLogArgs, auditLog := createRandomAuditLog()

	log.Println("Audit log args: ", auditLogArgs)


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
                "user_email": auditLogArgs.UserEmail,
                "action":     auditLogArgs.Action,
                "ip_address": auditLogArgs.IpAddress.String,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, auditLog.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateAuditLog(gomock.Any(), gomock.Eq(auditLogArgs)).
                    Times(1).
                    Return(auditLog, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchAuditLog(t, recorder.Body, auditLog)
            },
        },
        {
            name: "InternalError",
            body: gin.H{
                "user_email": auditLogArgs.UserEmail,
                "action":     auditLogArgs.Action,
                "ip_address": auditLogArgs.IpAddress.String,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, auditLog.Username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    CreateAuditLog(gomock.Any(), gomock.Eq(auditLogArgs)).
                    Times(1).
                    Return(db.AuditLog{}, sql.ErrConnDone)
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

            url := "/audit-logs"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            tc.setupAuth(t, request, server.tokenMaker)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestGetAuditLogsByUserEmailAPI(t *testing.T) {
    numAuditLogs := 5
    userEmail := "user@example.com"
    username := utils.RandomUser()
    
    auditLogs := make([]db.AuditLog, numAuditLogs)
    for i := 0; i < numAuditLogs; i++ {
        _, auditLog := createRandomAuditLog()
        auditLog.Username = username
        auditLog.UserEmail = userEmail
        auditLogs[i] = auditLog
    }

    testCases := []struct {
        name          string
        userEmail     string
        buildStubs    func(store *mockdb.MockStore_interface)
        setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            userEmail: userEmail,
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetAuditLogsByUserEmail(gomock.Any(), gomock.Eq(userEmail)).
                    Times(1).
                    Return(auditLogs, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchAuditLogs(t, recorder.Body, auditLogs)
            },
        },
        {
            name: "InternalError",
            userEmail: userEmail,
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetAuditLogsByUserEmail(gomock.Any(), gomock.Eq(userEmail)).
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

            url := fmt.Sprintf("/audit-logs/%s", tc.userEmail)
            request, err := http.NewRequest(http.MethodGet, url, nil)
            require.NoError(t, err)

            tc.setupAuth(t, request, server.tokenMaker)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestListUserAuditLogsAPI(t *testing.T) {
    numAuditLogs := 5
    userEmail := "user@example.com"
    
    auditLogs := make([]db.AuditLog, numAuditLogs)
    for i := 0; i < numAuditLogs; i++ {
        _, auditLog := createRandomAuditLog()
        auditLog.UserEmail = userEmail
        auditLogs[i] = auditLog
    }

    testCases := []struct {
        name          string
        userEmail     string
        buildStubs    func(store *mockdb.MockStore_interface)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            userEmail: userEmail,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetAuditLogsByUserEmail(gomock.Any(), gomock.Eq(userEmail)).
                    Times(1).
                    Return(auditLogs, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchAuditLogs(t, recorder.Body, auditLogs)
            },
        },
        {
            name: "NoAuditLogs",
            userEmail: "nonexistent",
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetAuditLogsByUserEmail(gomock.Any(), gomock.Any()).
                    Times(1).
                    Return(auditLogs, nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
                requireBodyMatchAuditLogs(t, recorder.Body, auditLogs)
            },
        },
        {
            name: "InternalError",
            userEmail: userEmail,
            buildStubs: func(store *mockdb.MockStore_interface) {
                store.EXPECT().
                    GetAuditLogsByUserEmail(gomock.Any(), gomock.Eq(userEmail)).
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

            url := fmt.Sprintf("/audit-logs/user/%s", tc.userEmail)
            request, err := http.NewRequest(http.MethodGet, url, nil)
            require.NoError(t, err)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}

func TestDeleteAuditLogsAPI(t *testing.T) {
    _, auditLog := createRandomAuditLog()

    userEmail := "user@example.com"
    username := utils.RandomUser()
    
    auditLog.UserEmail = userEmail
    auditLog.Username = username

    testCases := []struct {
        name          string
        AuditLogId    uuid.UUID
        UserEmail     string
        buildStubs    func(store *mockdb.MockStore_interface)
        setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            AuditLogId: auditLog.ID,
            UserEmail: userEmail,
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {

                store.EXPECT().
                    GetAuditLogsByUserEmail(gomock.Any(), gomock.Eq(userEmail)).
                    Times(1).
                    Return([]db.AuditLog{auditLog}, nil)

                store.EXPECT().
                    DeleteAuditLog(gomock.Any(), gomock.Eq(auditLog.ID)).
                    Times(1).
                    Return(nil)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
        {
            name: "Not Found",
            AuditLogId: auditLog.ID,
            UserEmail: userEmail,
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {

                store.EXPECT().
                    GetAuditLogsByUserEmail(gomock.Any(), gomock.Eq(userEmail)).
                    Times(1).
                    Return([]db.AuditLog{auditLog}, sql.ErrNoRows)

                store.EXPECT().
                    DeleteAuditLog(gomock.Any(), gomock.Eq(auditLog.ID)).
                    Times(0)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusNotFound, recorder.Code)
            },
        },
        {
            name: "InternalError",
            AuditLogId: auditLog.ID,
            UserEmail: userEmail,
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {

                store.EXPECT().
                    GetAuditLogsByUserEmail(gomock.Any(), gomock.Eq(userEmail)).
                    Times(1).
                    Return([]db.AuditLog{auditLog}, nil)

                store.EXPECT().
                    DeleteAuditLog(gomock.Any(), gomock.Eq(auditLog.ID)).
                    Times(1).
                    Return(sql.ErrConnDone)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusInternalServerError, recorder.Code)
            },
        },
		{
			name: "InvalidID",
			AuditLogId: uuid.Nil,
            UserEmail: userEmail,
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, username, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore_interface) {

                store.EXPECT().
                    GetAuditLogsByUserEmail(gomock.Any(), gomock.Eq(userEmail)).
                    Times(1).
                    Return([]db.AuditLog{auditLog}, nil)

                store.EXPECT().
                    DeleteAuditLog(gomock.Any(), gomock.Any()).
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

            url := fmt.Sprintf("/audit-logs/%s?user_email=%s", tc.AuditLogId, tc.UserEmail)
            request, err := http.NewRequest(http.MethodDelete, url, nil)
            log.Printf("Request URL: %s", request.URL.String())
            require.NoError(t, err)
            
            tc.setupAuth(t, request, server.tokenMaker)

            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }
}



func requireBodyMatchAuditLog(t *testing.T, body *bytes.Buffer, auditLog db.AuditLog) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotAuditLog db.AuditLog
    err = json.Unmarshal(data, &gotAuditLog)
    require.NoError(t, err)

	log.Println("Data: ", string(data))

    require.Equal(t, auditLog.UserEmail, gotAuditLog.UserEmail)
}

func requireBodyMatchAuditLogForGet(t *testing.T, body *bytes.Buffer, auditLog db.AuditLog) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotAuditLog db.AuditLog
    err = json.Unmarshal(data, &gotAuditLog)
    require.NoError(t, err)

    require.Equal(t, auditLog.ID, gotAuditLog.ID)
    require.Equal(t, auditLog.UserEmail, gotAuditLog.UserEmail)
    require.Equal(t, auditLog.Action, gotAuditLog.Action)
    require.Equal(t, auditLog.IpAddress, gotAuditLog.IpAddress)
    require.WithinDuration(t, auditLog.CreatedAt.Time, gotAuditLog.CreatedAt.Time, time.Second)
}

func requireBodyMatchAuditLogs(t *testing.T, body *bytes.Buffer, auditLogs []db.AuditLog) {
    data, err := io.ReadAll(body)
    require.NoError(t, err)

    var gotAuditLogs []db.AuditLog
    err = json.Unmarshal(data, &gotAuditLogs)
    require.NoError(t, err)

    require.Equal(t, len(auditLogs), len(gotAuditLogs))
    for i := range auditLogs {
        require.Equal(t, auditLogs[i].ID, gotAuditLogs[i].ID)
        require.Equal(t, auditLogs[i].UserEmail, gotAuditLogs[i].UserEmail)
        require.Equal(t, auditLogs[i].Action, gotAuditLogs[i].Action)
        require.Equal(t, auditLogs[i].IpAddress, gotAuditLogs[i].IpAddress)
        require.WithinDuration(t, auditLogs[i].CreatedAt.Time, gotAuditLogs[i].CreatedAt.Time, time.Second)
    }
}


