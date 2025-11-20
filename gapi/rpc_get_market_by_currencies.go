package gapi

import (
	"context"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *server) GetMarketByCurrencies(ctx context.Context, req *pb.GetMarketByCurrenciesRequest) (*pb.GetMarketByCurrenciesResponse, error) {

	username := req.GetUsername()

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}
	
	if authPayload.Username != username {
		return nil, status.Errorf(codes.Unknown, "Not authorized")
	}
	
	baseCurrency := req.GetBaseCurrency()
	quoteCurrency := req.GetQuoteCurrency()

	market, err := server.store.GetMarketByCurrencies(ctx, db.GetMarketByCurrenciesParams{
		BaseCurrency:  baseCurrency,
		QuoteCurrency: quoteCurrency,
	})

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to find the list of markets: %v", err)
	}

	pbMarket := convertMarket(market)
	
	res := &pb.GetMarketByCurrenciesResponse {
		Market: pbMarket,
	}

	return res, nil
}