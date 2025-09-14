package gapi

import (
	"context"

	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (server *server) GetWallets(ctx context.Context, req *emptypb.Empty) (*pb.GetWalletsResponse, error) {
	wallets, err := server.store.GetWallets(ctx)
	if err != nil {
		return nil, err
	}

	pbWallets := convertGetWallets(wallets)

	res := &pb.GetWalletsResponse{
		Wallets: pbWallets.Wallets,
	}

	return res, nil
}