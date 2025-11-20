package gapi

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *server) ListTrades(ctx context.Context, req *pb.TradeListRequest) (*pb.TradeListResponse, error) {
	
	marketId, err := uuid.Parse(req.GetMarketId())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %v", err)
	}

	trades, err := server.store.GetTradesByMarketID(ctx, marketId)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list trades: %v", err)
	}

	pbTrades := convertListTrades(trades)

	res := &pb.TradeListResponse{
		Trades: pbTrades,
	}

	return res, nil

}
