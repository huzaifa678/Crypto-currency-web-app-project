package gapi

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGoogleLogin(t *testing.T) {
	userID := uuid.New()
	sub := "google-sub-id"
	email := "test2443@example.com"
	username := "testuser"
	password := "go234w"

	originalVerify := oauth2VerifyGoogleIDToken
	defer func() { oauth2VerifyGoogleIDToken = originalVerify }()

	testCases := []struct {
		name          string
		req           *pb.GoogleLoginRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		stubVerify    func()
		checkResponse func(t *testing.T, res *pb.GoogleLoginResponse, err error)
	}{
		{
			name: "OK_NewUserCreated",
			req:  &pb.GoogleLoginRequest{IdToken: "valid-id-token"},
			stubVerify: func() {
				oauth2VerifyGoogleIDToken = func(ctx context.Context, rawIDToken, clientID string) (*idtoken.Payload, error) {
					return &idtoken.Payload{
						Claims: map[string]interface{}{
							"email":    email,
							"sub":      sub,
							"name":     username,
							"password": password,
						},
					}, nil
				}
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetGoogleUserByProviderID(gomock.Any(), sub).
					Times(1).
					Return(db.GoogleAuth{}, db.ErrRecordNotFound)

				store.EXPECT().
					CreateGoogleUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.GoogleAuth{
						ID:         userID,
						Email:      email,
						Username:   username,
						ProviderID: sub,
						Role:       pgtype.Text{String: "user", Valid: true},
					}, nil)

				store.EXPECT().
					GetUserByEmail(gomock.Any(), email).
					Times(1).
					Return(db.GetUserByEmailRow{}, db.ErrRecordNotFound)

				store.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserRow{
						ID:           userID,
						Username:     username,
						Email:        email,
						CreatedAt: 	  time.Now(),
						UpdatedAt: 	  time.Now(),
						Role:         db.UserRole("user"),
					}, nil)
				
				store.EXPECT().
        			CreateSession(gomock.Any(), gomock.Any()).
        			Times(1).
        			Return(db.Session{}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.GoogleLoginResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, email, res.Client.Email)
				require.Equal(t, username, res.Client.Username)
				require.NotEmpty(t, res.AccessToken)
				require.NotEmpty(t, res.RefreshToken)
			},
		},
		{
			name: "InvalidToken",
			req:  &pb.GoogleLoginRequest{IdToken: "invalid"},
			stubVerify: func() {
				oauth2VerifyGoogleIDToken = func(ctx context.Context, rawIDToken, clientID string) (*idtoken.Payload, error) {
					return nil, errors.New("invalid token")
				}
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
			},
			checkResponse: func(t *testing.T, res *pb.GoogleLoginResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "DBErrorOnCreateGoogleUser",
			req:  &pb.GoogleLoginRequest{IdToken: "valid"},
			stubVerify: func() {
				oauth2VerifyGoogleIDToken = func(ctx context.Context, rawIDToken, clientID string) (*idtoken.Payload, error) {
					return &idtoken.Payload{
						Claims: map[string]interface{}{
							"email": email,
							"sub":   sub,
							"name":  username,
						},
					}, nil
				}
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					GetGoogleUserByProviderID(gomock.Any(), sub).
					Times(1).
					Return(db.GoogleAuth{}, db.ErrRecordNotFound)

				store.EXPECT().
    			CreateGoogleUser(gomock.Any(), gomock.Any()).
    			Times(1).
    			Return(db.GoogleAuth{}, errors.New("db failure"))
			},
			checkResponse: func(t *testing.T, res *pb.GoogleLoginResponse, err error) {
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
			tc.stubVerify()

			server := NewTestServer(t, store, nil)

			res, err := server.GoogleLogin(context.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
