package gapi

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteUserRPC(t *testing.T) {
	_, user, _, _, _, _ := createRandomUser()

	testCases := []struct {
		name          string
		req           *pb.DeleteUserRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(t *testing.T, res *pb.DeleteUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.DeleteUserRequest{
				UserId: user.ID.String(),
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "successfully deleted the user", res.Message)
			},
		},
		{
			name: "InvalidUUID",
			req: &pb.DeleteUserRequest{
				UserId: "invalid-uuid",
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "UserNotFound",
			req: &pb.DeleteUserRequest{
				UserId: user.ID.String(),
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(db.ErrRecordNotFound)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.NotFound, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.DeleteUserRequest{
				UserId: user.ID.String(),
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					DeleteUser(gomock.Any(), gomock.Eq(user.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.DeleteUserResponse, err error) {
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

			res, err := server.DeleteUser(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
} 