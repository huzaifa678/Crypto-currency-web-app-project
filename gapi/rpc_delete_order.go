package gapi

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *server) DeleteOrder(ctx context.Context, req *pb.DeleteOrderRequest) (*pb.DeleteOrderResponse, error) {

	violations := validateDeleteOrderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %v", err)
	}

	order, err := server.store.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find the order")
	}

	log.Println("AUTH USERNAME", authPayload.Username)
	log.Println("ORDER USERNAME", order.Username)

	if authPayload.Username != order.Username {
		log.Println(err)
		return nil, status.Errorf(codes.PermissionDenied, "not authorized to delete this order")
	}

	err = server.store.DeleteOrder(ctx, orderID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find the order for the user")
	}

	res := &pb.DeleteOrderResponse {
		Message: "Order deleted successfully",
	}

	return res, nil
}


func validateDeleteOrderRequest(req *pb.DeleteOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateDeleteRequest(req.GetOrderId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return violations
}