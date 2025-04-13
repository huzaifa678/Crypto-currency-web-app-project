package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/stretchr/testify/require"
)


func addAuthMiddleware(t *testing.T, request *http.Request, tokenMaker token.Maker, authorizationType string, username string, duration time.Duration) {
	token, err := tokenMaker.CreateToken(username, duration)

	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	request.Header.Set("Authorization", authorizationHeader)
}

func TestAuthMiddleWare(t *testing.T) {

	testCases := []struct {
		name 		  string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, "user1", time.Minute)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recoder.Code)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "Invalid Token",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthMiddleware(t, request, tokenMaker, "bearer token access", "user2", time.Minute)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "Unsupported Type",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthMiddleware(t, request, tokenMaker, "Basic", "user3", time.Minute)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
		{
			name: "Expired Token",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthMiddleware(t, request, tokenMaker, AuthorizationTypeBearer, "user4", -time.Minute)
			},
			checkResponse: func(t *testing.T, recoder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recoder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {


            server := NewTestServer(t, nil)

			url := "/auth"
			server.router.GET(url, authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, req, server.tokenMaker)
			server.router.ServeHTTP(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
	
}