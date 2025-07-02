package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

const (
	AuthorizationTypeBearer = "bearer"
)

func NewTestServer(t *testing.T, store db.Store_interface) *server {
	config := utils.Config{
		PasetoSymmetricKey: utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(store, config)

	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, duration time.Duration, tokenType token.TokenType) context.Context {
	accessToken, _, err := tokenMaker.CreateToken(username, duration, tokenType)
	require.NoError(t, err)

	bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
	md := metadata.MD{
		authorizationHeader: []string{
			bearerToken,
		},
	}

	return metadata.NewIncomingContext(context.Background(), md)
}