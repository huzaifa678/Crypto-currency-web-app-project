package gapi

import (
	"context"
	"log"

	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *server) ListOrder (ctx context.Context, req *pb.OrderListRequest) (*pb.OrderListResponse, error) {
	log.Println("RECEIVED ListOrder request")

	username := req.GetUsername()

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	if authPayload.Username != username {
		return nil, status.Errorf(codes.Unknown, "Not authorized")
	}

	orders, err := server.store.ListOrdersByUsername(ctx, username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list orders: %v", err)
	}

	pbOrders := convertListOrders(orders)

	res := &pb.OrderListResponse{
		Orders: pbOrders.Orders,
	}

	return res, nil
}