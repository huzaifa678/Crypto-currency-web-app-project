package gapi

import (
	"context"

	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *server) ListMarkets(ctx context.Context) (*pb.MarketListResponse, error) {
	markets, err := server.store.ListMarkets(ctx)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to find the list of markets: %v", err)
	}

	pbMarkets := convertListMarkets(markets)

	res := &pb.MarketListResponse {
		Markets: pbMarkets.Markets,
	}

	return res, nil
}

