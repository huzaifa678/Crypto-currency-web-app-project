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
	"github.com/huzaifa678/Crypto-currency-web-app-project/worker"
	mockwk "github.com/huzaifa678/Crypto-currency-web-app-project/worker/mock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user db.CreateUserRow
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	if err := utils.ComparePasswords(actualArg.PasswordHash, expected.password); err != nil {
		log.Printf("error in comparing crypto/bcrypt: %v", err)
		return false
	}

	expected.arg.PasswordHash = actualArg.PasswordHash

	log.Println("expected.arg: ", expected.arg.CreateUserParams)
	log.Println("actualArg: ", actualArg.CreateUserParams)

	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	err := actualArg.AfterCreate(expected.user)
	return err == nil
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.CreateUserRow) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

func TestCreateUserAPI(t *testing.T) {
	userArgs, user, userRows, _, _, password := createRandomUser()

	log.Println("userArgs.email: ", userArgs.Email)
	log.Println("password: ", password)
	log.Println("userArgs", userArgs)

	testCases := []struct {
		name          string
		body          *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore_interface, taskDistributor *mockwk.MockTaskDistributor)
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
			buildStubs: func(store *mockdb.MockStore_interface, taskDistributor *mockwk.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						Email:    user.Email,
						PasswordHash: password,
						Role:     db.UserRole(pb.UserRole_USER_ROLE_USER.String()),
						IsVerified: userArgs.IsVerified,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, userRows)).
					Times(1).
					Return(db.CreateUserTxResult{CreateUserRow: userRows}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{
					Email: userRows.Email,
				}
				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createdUserId := res.GetUserId()
				require.NotEmpty(t, createdUserId)
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
			buildStubs: func(store *mockdb.MockStore_interface, taskDistributor *mockwk.MockTaskDistributor) {
				

				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, sql.ErrConnDone)

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
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
			buildStubs: func(store *mockdb.MockStore_interface, taskDistributor *mockwk.MockTaskDistributor) {		
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, &pq.Error{Code: "23505"})

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
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
			buildStubs: func(store *mockdb.MockStore_interface, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0)

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
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
			buildStubs: func(store *mockdb.MockStore_interface, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(0)

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
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
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore_interface(storeCtrl)

			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()
			taskDistributor := mockwk.NewMockTaskDistributor(taskCtrl)

			tc.buildStubs(store, taskDistributor)
			server := NewTestServer(t, store, taskDistributor)

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
	password := fmt.Sprintf("password%d", rand.Intn(20))
	hashedPassword, _ := utils.HashPassword(password)

	username := fmt.Sprintf("usertesting%d", rand.Intn(100))

	userArgs := db.CreateUserParams{
		Username:    username,
		Email:       randomEmail,
		PasswordHash: hashedPassword, 
		Role:        db.UserRole(pb.UserRole_USER_ROLE_USER.String()),
		IsVerified:  true,
	}

	user := db.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        randomEmail,
		PasswordHash: hashedPassword,
		Role:         db.UserRole(pb.UserRole_USER_ROLE_USER.String()), 
		IsVerified:   userArgs.IsVerified,
	}

	now := time.Now()

	userRow := db.CreateUserRow{
		ID:         user.ID,
		Username:   user.Username,
		Email:      user.Email,
		CreatedAt:  now,
		UpdatedAt:  now,
		Role:       user.Role,
		IsVerified: user.IsVerified,
	}

	getUserRow := db.GetUserByIDRow{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
		Role:         user.Role,
		IsVerified:   user.IsVerified,
	}

	userEmailGetArgs := db.GetUserByEmailRow{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
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
