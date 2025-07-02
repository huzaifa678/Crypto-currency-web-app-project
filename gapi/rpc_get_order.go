package gapi

import (
	"context"
	"errors"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (server *server) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	violations := validateGetOrderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ID not parsed")
	}

	order, err := server.store.GetOrderByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find the order for the user")
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	if authPayload.Username != order.Username {
		return nil, status.Errorf(codes.Unknown, "unknown")
	}

	res := &pb.GetOrderResponse {
		Order: convertOrder(order),
	}

	return res, nil
}

func validateGetOrderRequest(req *pb.GetOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateGetRequest(req.GetOrderId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return violations
}
