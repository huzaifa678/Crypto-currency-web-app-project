package gapi

import (
	"context"
	"log"

	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)


func (server *server) ListOrder (ctx context.Context, req *emptypb.Empty) (*pb.OrderListResponse, error) {
	log.Println("RECEIVED ListOrder request")

	orders, err := server.store.ListOrders(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	pbOrders := convertListOrders(orders)

	res := &pb.OrderListResponse{
		Orders: pbOrders.Orders,
	}

	return res, nil
}