package gapi

import (
	"context"
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestLoginUserRPC(t *testing.T) {
	userArgs, user, _, _, GetUserByEmailRow, _ := createRandomUser()

	testCases := []struct {
		name          string
		req           *pb.LoginUserRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(t *testing.T, res *pb.LoginUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.LoginUserRequest{
				Email:    userArgs.Email,
				Password: userArgs.PasswordHash,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(GetUserByEmailRow, nil)

				store.EXPECT().
					CreateSession(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Session{
						ID:           uuid.New(),
						Username:     user.Username,
						RefreshToken: "refresh_token",
						UserAgent:    "user_agent",
						ClientIp:     "client_ip",
						IsBlocked:    false,
						ExpiresAt:    time.Now().Add(time.Hour),
					}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				log.Println("ERROR: ", err)
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.AccessToken)
				require.NotEmpty(t, res.RefreshToken)
				require.NotNil(t, res.AccessTokenExpiration)
				require.NotNil(t, res.RefreshTokenExpiration)
				require.NotNil(t, res.User)
				require.Equal(t, user.Username, res.User.Username)
				require.Equal(t, user.Email, res.User.Email)
			},
		},
		{
			name: "UserNotFound",
			req: &pb.LoginUserRequest{
				Email:    user.Email,
				Password: user.PasswordHash,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(GetUserByEmailRow, db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				log.Println("ERROR: ", err)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "IncorrectPassword",
			req: &pb.LoginUserRequest{
				Email:    user.Email,
				Password: "wrongpassword",
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(db.GetUserByEmailRow{
						ID:           user.ID,
						Username:     user.Username,
						Email:        user.Email,
						PasswordHash: user.PasswordHash,
						Role:         user.Role,
					}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InvalidEmail",
			req: &pb.LoginUserRequest{
				Email:    "invalid-email",
				Password: user.PasswordHash,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InvalidPassword",
			req: &pb.LoginUserRequest{
				Email:    user.Email,
				Password: "123", 
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.LoginUserRequest{
				Email:    user.Email,
				Password: user.PasswordHash,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetUserByEmail(gomock.Any(), gomock.Eq(user.Email)).
					Times(1).
					Return(db.GetUserByEmailRow{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.LoginUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
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
			res, err := server.LoginUser(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
} 