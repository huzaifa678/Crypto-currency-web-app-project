package gapi

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	violations := validateCreateOrderRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	marketID, err := uuid.Parse(req.GetMarketId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid market ID: %v", err)
	}

	arg := db.CreateOrderParams{
		Username:  authPayload.Username,
		UserEmail: req.GetUserEmail(),
		MarketID:  marketID,
		Type:      db.OrderType(req.GetType()),
		Status:    db.OrderStatus(req.GetStatus()),
		Price:     sql.NullString{String: req.GetPrice(), Valid: req.GetPrice() != ""},
		Amount:    req.GetAmount(),
	}

	order, err := server.store.CreateOrder(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	res := &pb.CreateOrderResponse {
		OrderId: order.ID.String(),
	}

	return res, nil
}


func validateCreateOrderRequest(req *pb.CreateOrderRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateCreateOrderRequest(req.GetUserEmail(), req.GetMarketId(), req.GetPrice(), req.GetAmount(), req.GetType()); err != nil {
		violations = append(violations, fieldViolation("id", err))
		violations = append(violations, fieldViolation("price", err))
		violations = append(violations, fieldViolation("amount", err))
		violations = append(violations, fieldViolation("ordertype", err))
	}

	return violations
}
