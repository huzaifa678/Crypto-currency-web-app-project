package gapi

import (
	"fmt"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/huzaifa678/Crypto-currency-web-app-project/worker"
)


type server struct {
	store       db.Store_interface
	tokenMaker  token.Maker
	config 	    utils.Config
	pb.UnimplementedCryptoWebAppServer
	taskDistributor worker.TaskDistributor
}

func NewServer(store db.Store_interface, config utils.Config, taskDistributor worker.TaskDistributor) (*server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.PasetoSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
