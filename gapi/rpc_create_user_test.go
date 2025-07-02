package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	log.Println("PASSWORD", e.arg.PasswordHash)
	log.Println("ARGS", arg)

	

	if err := utils.ComparePasswords(arg.PasswordHash, e.arg.PasswordHash); err != nil {
		log.Println("password mismatch:", err)
		return false
	}

	e.arg.PasswordHash = arg.PasswordHash
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v", e.arg)
}

func EqCreateUserParams(arg db.CreateUserParams) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg}
}

func TestCreateUserAPI(t *testing.T) {
	userArgs, user, userRows, _, userEmailGetArgs, password := createRandomUser()

	log.Println("userArgs.email: ", userArgs.Email)
	log.Println("password: ", password)
	log.Println("userArgs", userArgs)

	testCases := []struct {
		name          string
		body          *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			body: &pb.CreateUserRequest{
				Username:	   userArgs.Username,
				Email:         userArgs.Email,
				Password: 	   password,
				Role:          pb.UserRole_USER_ROLE_USER,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
                    GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
                    Times(1).
                    Return(userEmailGetArgs, db.ErrRecordNotFound)

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(userArgs)).
					Times(1).
					Return(userRows, nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				log.Println("ERROR: ", err)
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, user.ID.String(), res.UserId)
			},
		},
		{
			name: "InternalError",
			body: &pb.CreateUserRequest{
				Username:	   userArgs.Username,
				Email:         user.Email,
				Password:      password,
				Role:          pb.UserRole_USER_ROLE_USER,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
                    GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
                    Times(1).
                    Return(userEmailGetArgs, db.ErrRecordNotFound)

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(userArgs)).
					Times(1).
					Return(userRows, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				log.Println("ERROR: ", err)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name: "DuplicateEmail",
			body: &pb.CreateUserRequest{
				Username:	   userArgs.Username,
				Email:         user.Email,
				Password:      password,
				Role:          pb.UserRole_USER_ROLE_USER,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
                    GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
                    Times(1).
                    Return(userEmailGetArgs, db.ErrRecordNotFound)

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(userArgs)).
					Times(1).
					Return(userRows, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				log.Println("ERROR: ", err)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
		{
			name: "InvalidEmail",
			body: &pb.CreateUserRequest{
				Username:	   userArgs.Username,
				Email:       "invalid-email",
				Password:    userArgs.PasswordHash,
				Role:        pb.UserRole_USER_ROLE_USER,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				log.Println("ERROR: ", err)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "TooShortPassword",
			body: &pb.CreateUserRequest{
				Username:	   userArgs.Username,
				Email:       user.Email,
				Password:    "123",
				Role:        pb.UserRole_USER_ROLE_USER,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				log.Println("ERROR: ", err)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
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

			res, err := server.CreateUser(context.Background(), tc.body)
			tc.checkResponse(t, res, err)
		})
	}
}

func createRandomUser() (
	db.CreateUserParams, db.User, db.CreateUserRow,
	db.GetUserByIDRow, db.GetUserByEmailRow, string, 
) {
	randomEmail := fmt.Sprintf("testing%d@example.com", rand.Intn(1000))
	password := utils.RandomString(20)
	hashedPassword, _ := utils.HashPassword(password)

	username := RandomUsername(10)

	userArgs := db.CreateUserParams{
		Username:    username,
		Email:       randomEmail,
		PasswordHash: password,
		Role:        db.UserRole(pb.UserRole_USER_ROLE_USER.String()),
		IsVerified:  sql.NullBool{Bool: true, Valid: true},
	}

	user := db.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        randomEmail,
		PasswordHash: hashedPassword,
		Role:         db.UserRole("user"),
		IsVerified:   userArgs.IsVerified,
	}

	now := time.Now()

	userRow := db.CreateUserRow{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		CreatedAt:  sql.NullTime{Time: now, Valid: true},
		UpdatedAt:  sql.NullTime{Time: now, Valid: true},
		Role:       user.Role,
		IsVerified: user.IsVerified,
	}

	getUserRow := db.GetUserByIDRow{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    sql.NullTime{Time: now, Valid: true},
		UpdatedAt:    sql.NullTime{Time: now, Valid: true},
		Role:         user.Role,
		IsVerified:   user.IsVerified,
	}

	userEmailGetArgs := db.GetUserByEmailRow{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    sql.NullTime{Time: now, Valid: true},
		UpdatedAt:    sql.NullTime{Time: now, Valid: true},
		Role:         user.Role,
		IsVerified:   user.IsVerified,
	}

	return userArgs, user, userRow, getUserRow, userEmailGetArgs, password
}

func RandomUsername(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789_")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
