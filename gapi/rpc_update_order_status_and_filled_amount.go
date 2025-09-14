package gapi

import (
	"context"
	"strings"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	val "github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *server) UpdateOrderStatusAndFilledAmount(ctx context.Context, req *pb.UpdateOrderStatusAndFilledAmountRequest) (*pb.UpdateOrderStatusAndFilledAmountResponse, error) {
	violations := validateUpdateOrderStatusAndFilledAmount(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	orderId, err := uuid.Parse(req.GetOrderId())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ID not parsed")
	}

	order, err := server.store.GetOrderByID(ctx, orderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find the order for the user")
	}

	if authPayload.Username != order.Username {
		return nil, status.Errorf(codes.Unknown, "Not authorized")
	}

	pbStatus := strings.ToLower(req.GetOrderStatus().String())

	args := db.UpdateOrderStatusAndFilledAmountParams{
		Status:       db.OrderStatus(pbStatus),
		FilledAmount: decimal.NewFromFloat(float64(req.GetFilledAmount())),
		ID:           orderId,
	}

	err = server.store.UpdateOrderStatusAndFilledAmount(ctx, args)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update the order status and filled amount")
	}

	res := &pb.UpdateOrderStatusAndFilledAmountResponse{
		Success: "successfully updated the order status and filled amount",
	}

	return res, nil
}


func validateUpdateOrderStatusAndFilledAmount(req *pb.UpdateOrderStatusAndFilledAmountRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUpdateOrderStatusAndFilledAmount(req.GetOrderId(), req.GetOrderStatus(), decimal.NewFromFloat(float64(req.GetFilledAmount()))); err != nil {
		violations = append(violations, fieldViolation("id", err))
		violations = append(violations, fieldViolation("status", err))
		violations = append(violations, fieldViolation("filled_amount", err))
	}

	return violations
}