package gapi

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (server *server) GetOrderByMarketID(ctx context.Context, req *pb.GetOrderByMarketIDRequest) (*pb.GetOrderByMarketIDResponse, error) {

	marketID, err := uuid.Parse(req.GetMarketId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ID not parsed")
	}

	orders, err := server.store.ListOrdersByMarketID(ctx, marketID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find the orders for the market")
	}

	pbOrders := convertListOrders(orders)

	res := &pb.GetOrderByMarketIDResponse {
		Orders: pbOrders.Orders,
	}

	return res, nil
}