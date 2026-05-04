package gapi

import (
	"context"
	"log"

	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *server) ListMarkets(ctx context.Context, req *pb.MarketListRequest) (*pb.MarketListResponse, error) {

	username := req.GetUsername()

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	log.Println("RECEIVED ListMarkets request for user:", username)
	log.Println("Authenticated user:", authPayload.Username)

	if authPayload.Username != username {
		return nil, status.Errorf(codes.Unknown, "Not authorized")
	}

	markets, err := server.store.ListMarketsByUsername(ctx, username)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to find the list of markets: %v", err)
	}

	pbMarkets := convertListMarkets(markets)

	res := &pb.MarketListResponse {
		Markets: pbMarkets.Markets,
	}

	return res, nil
}

