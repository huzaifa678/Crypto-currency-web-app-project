package gapi

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUpdateUserRPC(t *testing.T) {
	_, user, _, _, _, _ := createRandomUser()
	newPassword := "newpassword123"

	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				UserId:   user.ID.String(),
				Password: newPassword,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
		
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, "successfully updated the user information", res.Message)
			},
		},
		{
			name: "InvalidUUID",
			req: &pb.UpdateUserRequest{
				UserId:   "invalid-uuid",
				Password: newPassword,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InvalidPassword",
			req: &pb.UpdateUserRequest{
				UserId:   user.ID.String(),
				Password: "short", 
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.UpdateUserRequest{
				UserId:   user.ID.String(),
				Password: newPassword,
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
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

			res, err := server.UpdateUser(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
} 