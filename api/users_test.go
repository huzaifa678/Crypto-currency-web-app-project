package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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

func createRandomUser() (db.CreateUserParams, db.User, db.CreateUserRow, db.GetUserByIDRow) {
	randomEmail := fmt.Sprintf("test%d@example.com", rand.Intn(1000))
	randomPasswordHash := fmt.Sprintf("passwordhash%d", rand.Intn(1000))

	userArgs := db.CreateUserParams{
		Email:        randomEmail,
		PasswordHash: randomPasswordHash,
		Role:         db.UserRole("user"),
		IsVerified:   sql.NullBool{Bool: true, Valid: true},
	}


	user := db.User{
        ID: uuid.New(),
		Email:        userArgs.Email,
		PasswordHash: userArgs.PasswordHash,
		Role:         userArgs.Role,
		IsVerified:   userArgs.IsVerified,
	}

	userRow := db.CreateUserRow{
		ID:         user.ID,
		Email:      user.Email,
		CreatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		Role:       userArgs.Role,
		IsVerified: userArgs.IsVerified,
	}

    getUserRow := db.GetUserByIDRow{
        ID:         user.ID,
        Email:      user.Email,
        PasswordHash: user.PasswordHash,
        CreatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
        UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
        Role:       userArgs.Role,
        IsVerified: userArgs.IsVerified,
    }

	return userArgs, user, userRow, getUserRow
}

func initTestServer(t *testing.T, ctrl *gomock.Controller) (*server, *mockdb.MockStore_interface, *httptest.ResponseRecorder) {
	store := mockdb.NewMockStore_interface(ctrl)
	server := NewServer(store)
	recorder := httptest.NewRecorder()
	return server, store, recorder
}

func sendPostRequest(t *testing.T, server *server, endpoint string, body interface{}) *httptest.ResponseRecorder {
	data, err := json.Marshal(body)
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(data))
	require.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	server.router.ServeHTTP(recorder, request)

	return recorder
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server, store, recorder := initTestServer(t, ctrl)

	userArgs, _, userRow, _ := createRandomUser()

	store.EXPECT().
		CreateUser(gomock.Any(), gomock.Eq(userArgs)).
		Times(1).
		Return(userRow, nil)


    log.Println(userRow.ID)
    
	recorder = sendPostRequest(t, server, "/users", userArgs)

	require.Equal(t, http.StatusOK, recorder.Code)

	var got gin.H
	err := json.Unmarshal(recorder.Body.Bytes(), &got)
	require.NoError(t, err)

	require.Equal(t, userRow.ID.String(), got["id"])
}

func TestGetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server, store, recorder := initTestServer(t, ctrl)

	userArgs, _, userRow, getUserRow := createRandomUser()

    
    store.EXPECT().
		CreateUser(gomock.Any(), gomock.Eq(userArgs)).
		Times(1).
		Return(userRow, nil)

    log.Println(userRow.ID)

	
    recorder = sendPostRequest(t, server, "/users", userArgs)

    userId, err := uuid.Parse(getUserRow.ID.String())

    store.EXPECT().
		GetUserByID(gomock.Any(), gomock.Eq(userId)).
		Times(1).
		Return(getUserRow, nil)

    recorder = httptest.NewRecorder()

    endpoint := fmt.Sprintf("/users/%s", userId)

    request, err := http.NewRequest(http.MethodGet, endpoint, nil)
    server.router.ServeHTTP(recorder, request)

    log.Println(recorder.Body.String())

	require.NoError(t, err)

	require.Equal(t, http.StatusOK, recorder.Code)

	require.NoError(t, err)

    role := db.UserRole(getUserRow.Role)

    var got gin.H
    err = json.Unmarshal(recorder.Body.Bytes(), &got)
	require.Equal(t, getUserRow.ID.String(), got["id"])
    require.Equal(t, getUserRow.PasswordHash, got["password_hash"])
    require.Equal(t, getUserRow.Email, got["email"])
	require.Equal(t, role, db.UserRole(got["role"].(string)))
    isVerified, ok := got["is_verified"].(map[string]interface{})
    require.True(t, ok, "is_verified should be a map containing Bool and Valid")
	boolValue, ok := isVerified["Bool"].(bool)
    require.True(t, ok, "Bool should be a boolean")
    require.Equal(t, getUserRow.IsVerified.Bool, boolValue)

    validValue, ok := isVerified["Valid"].(bool)
    require.True(t, ok, "Valid should be a boolean")
    require.Equal(t, getUserRow.IsVerified.Valid, validValue)
}
