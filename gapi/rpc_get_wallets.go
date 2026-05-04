package gapi

import (
	"context"

	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
)

func (server *server) GetWallets(ctx context.Context, req *pb.GetWalletsRequest) (*pb.GetWalletsResponse, error){
	
	email := req.GetUserEmail()

	wallets, err := server.store.GetWalletsByUserEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	pbWallets := convertGetWallets(wallets)

	res := &pb.GetWalletsResponse{
		Wallets: pbWallets.Wallets,
	}

	return res, nil
}