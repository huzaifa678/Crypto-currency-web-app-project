package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/huzaifa678/Crypto-currency-web-app-project/worker"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

const (
	AuthorizationTypeBearer = "bearer"
)

type TestServerBuilder struct{
	server *server
}

func (b *TestServerBuilder) setStore(store db.Store_interface) (*TestServerBuilder) {
	b.server.store = store
	return b
}

func (b *TestServerBuilder) setTaskDistributor(taskDistributor worker.TaskDistributor) (*TestServerBuilder) {
	b.server.taskDistributor = taskDistributor
	return b
}

func NewTestServerBuilder() *TestServerBuilder {
	config := utils.Config{
		PasetoSymmetricKey: utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	tokenMaker, err := token.NewPasetoMaker(config.PasetoSymmetricKey)
	if err != nil {
		return nil
	}

	return &TestServerBuilder{
		server: &server{
			config: config,
			tokenMaker: tokenMaker,
		},
	}
}

func (b *TestServerBuilder) NewTestServer2(t *testing.T) *server {
	return b.server
}

func NewTestServer(t *testing.T, store db.Store_interface, taskDistributor worker.TaskDistributor) *server {
	config := utils.Config{
		PasetoSymmetricKey: utils.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(store, config, taskDistributor)

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