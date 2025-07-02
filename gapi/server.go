package gapi

import (
	"fmt"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
)


type server struct {
	store       db.Store_interface
	tokenMaker  token.Maker
	config 	    utils.Config
	pb.UnimplementedCryptoWebAppServer
}

func NewServer(store db.Store_interface, config utils.Config) (*server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.PasetoSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	return server, nil
}
