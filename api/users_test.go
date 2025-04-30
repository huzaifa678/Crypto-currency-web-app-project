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
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/rand"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

type eqUpdateUserParamsMatcher struct {
	arg      db.UpdateUserParams
	password string
}

func (e eqUpdateUserParamsMatcher) Matches(x interface{}) bool {
	panic("unimplemented")
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.ComparePasswords(arg.PasswordHash, e.password)
	if err != nil {
		return false
	}

	e.arg.PasswordHash = arg.PasswordHash
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func (e eqUpdateUserParamsMatcher) UpdateMatches(x interface{}) bool {
	arg, ok := x.(db.UpdateUserParams)
	if !ok {
		return false
	}

	err := utils.ComparePasswords(arg.PasswordHash, e.password)
	if err != nil {
		return false
	}

	e.arg.PasswordHash = arg.PasswordHash
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqUpdateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqUpdateUserParams(arg db.UpdateUserParams, password string) gomock.Matcher {
	return eqUpdateUserParamsMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	userArgs, user, userRows, _, userEmailGetArgs := createRandomUser()

	log.Println("userArgs.email: ", userArgs.Email)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":		 userArgs.Username,
				"email":         userArgs.Email,
				"password_hash": userArgs.PasswordHash,
				"role":          userArgs.Role,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
                    GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
                    Times(1).
                    Return(userEmailGetArgs, sql.ErrNoRows)

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(userArgs, userArgs.PasswordHash)).
					Times(1).
					Return(userRows, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":		 userArgs.Username,
				"email":         userArgs.Email,
				"password_hash": userArgs.PasswordHash,
				"role":          userArgs.Role,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
                    GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
                    Times(1).
                    Return(userEmailGetArgs, sql.ErrNoRows)

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(userArgs, userArgs.PasswordHash)).
					Times(1).
					Return(userRows, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateEmail",
			body: gin.H{
				"username":		 userArgs.Username,
				"email":         userArgs.Email,
				"password_hash": userArgs.PasswordHash,
				"role":          user.Role,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
                    GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
                    Times(1).
                    Return(userEmailGetArgs, sql.ErrNoRows)

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(userArgs, userArgs.PasswordHash)).
					Times(1).
					Return(userRows, &pq.Error{Code: "23505"})
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":	   userArgs.Username,
				"email":       "invalid-email",
				"password":    userArgs.PasswordHash,
				"role":        user.Role,
				"is_verified": user.IsVerified,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "TooShortPassword",
			body: gin.H{
				"username":	   userArgs.Username,
				"email":       user.Email,
				"password":    "123",
				"role":        user.Role,
				"is_verified": user.IsVerified,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
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

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestLoginUserAPI(t *testing.T) {
	userArgs, user, _, getUserArgs, getUserByEmailArgs := createRandomUser()

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    	 userArgs.Email,
				"password_hash": userArgs.PasswordHash,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				user.PasswordHash = getUserArgs.PasswordHash

				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(userArgs.Email)).
					Times(1).
					Return(getUserByEmailArgs, nil)

				HashedPassword, err := utils.HashPassword(getUserByEmailArgs.PasswordHash)
				getUserArgs.PasswordHash = HashedPassword
				err = utils.ComparePasswords(HashedPassword, getUserByEmailArgs.PasswordHash)
				require.NoError(t, err)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{
						ID:        uuid.New(),
						Username:  user.Username,
						ExpiresAt: time.Now().Add(time.Minute * 15),
					}, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "EmailNotFound",
			body: gin.H{
				"email":    userArgs.Email,
				"password_hash": userArgs.PasswordHash,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GetUserByEmailRow{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "WrongPassword",
			body: gin.H{
				"email":    	  userArgs.Email,
				"password_hash": "wrongpass12345678",
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				user.PasswordHash = getUserArgs.PasswordHash

				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(userArgs.Email)).
					Times(1).
					Return(getUserByEmailArgs, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "DBError",
			body: gin.H{
				"email":    user.Email,
				"password_hash": getUserArgs.PasswordHash,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(db.GetUserByEmailRow{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name: "InvalidBody",
			body: gin.H{
				"email": 12345,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
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

			request, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}


func TestGetUserAPI(t *testing.T) {
	_, user, _, getUserParams, _ := createRandomUser()

	testCases := []struct {
		name          string
		userID        uuid.UUID
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(getUserParams, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name:   "NotFound",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(getUserParams, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetUserByID(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(getUserParams, sql.ErrConnDone)
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

			url := fmt.Sprintf("/users/%s", tc.userID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateUserAPI(t *testing.T) {
	_, user, _, _, _ := createRandomUser()

	log.Println("UPDATE USER ARGS: ", user)
	testCases := []struct {
		name          string
		userID        uuid.UUID
		body          gin.H
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID,
			body: gin.H{
				"username": user.Username,
				"email": 	user.Email,
				"password_hash": "password123",
				"role":          user.Role,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				
				
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "UserNotFound",
			userID: user.ID,
			body: gin.H{
				"password_hash": user.PasswordHash,
				"role":          user.Role,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InternalErrorOnDelete",
			userID: user.ID,
			body: gin.H{
				"password_hash": user.PasswordHash,
				"role":          user.Role,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			userID: uuid.Nil,
			body: gin.H{
				"password_hash": "password123",
				"role":          user.Role,
			},

			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
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

			data, err := json.Marshal(tc.body)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/users/%s", tc.userID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteUserAPI(t *testing.T) {
	_, user, _, _, _ := createRandomUser()

	testCases := []struct {
		name          string
		userID        uuid.UUID
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "OK",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:   "NotFound",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InternalError",
			userID: user.ID,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:   "InvalidID",
			userID: uuid.Nil,
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
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

			url := fmt.Sprintf("/users/%s", tc.userID)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func createRandomUser() (db.CreateUserParams, db.User, db.CreateUserRow, db.GetUserByIDRow, db.GetUserByEmailRow) {
	randomEmail := fmt.Sprintf("testing%d@example.com", rand.Intn(1000))
	randomPassword := fmt.Sprintf("password%d", rand.Intn(1000))
	hashedPassword, _ := utils.HashPassword(randomPassword)

	userArgs := db.CreateUserParams{
		Username: 	  utils.RandomString(32),
		Email:        randomEmail,
		PasswordHash: randomPassword,
		Role:         db.UserRole("user"),
		IsVerified:   sql.NullBool{Bool: true, Valid: true},
	}

	user := db.User{
		ID: 		  uuid.New(),
		Username: 	  userArgs.Username,
		Email:        userArgs.Email,
		PasswordHash: hashedPassword,
		Role:         userArgs.Role,
		IsVerified:   userArgs.IsVerified,
	}

	userRow := db.CreateUserRow{
		ID:         user.ID,
		Username: 	user.Username,
		Email:      user.Email,
		CreatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:  sql.NullTime{Time: time.Now(), Valid: true},
		Role:       userArgs.Role,
		IsVerified: userArgs.IsVerified,
	}

	getUserRow := db.GetUserByIDRow{
		ID:           user.ID,
		Username: 	  user.Username,	
		Email:        user.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		Role:         userArgs.Role,
		IsVerified:   userArgs.IsVerified,
	}

	userEmailGetArgs := db.GetUserByEmailRow {
		ID:           user.ID,
		Username: 	  user.Username,	
		Email:        user.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		Role:         userArgs.Role,
		IsVerified:   userArgs.IsVerified,
	}

	return userArgs, user, userRow, getUserRow, userEmailGetArgs
}


func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	log.Println("Response body:", string(data))
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)

	require.NotZero(t, user.ID, gotUser.ID)
}
