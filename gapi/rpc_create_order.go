package gapi

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log.Println("RECEIVED CreateOrder request:", req)
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

	pbType := strings.ToLower(req.GetType().String())

	log.Printf("UserEmail bytes: %q\n", []byte(req.GetUserEmail()))
	log.Printf("Price bytes: %q\n", []byte(req.GetPrice()))
	log.Printf("Amount bytes: %q\n", []byte(req.GetAmount()))
	log.Printf("Username bytes: %q\n", []byte(authPayload.Username))

	arg := db.CreateOrderParams{
		Username:  authPayload.Username,
		UserEmail: req.GetUserEmail(),
		MarketID:  marketID,
		Type:      db.OrderType(pbType),
		Status:    db.OrderStatus("open"), // default status
		Price:     req.GetPrice(),
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
		violations = append(violations, fieldViolation("user_email", err))
		violations = append(violations, fieldViolation("id", err))
		violations = append(violations, fieldViolation("price", err))
		violations = append(violations, fieldViolation("amount", err))
		violations = append(violations, fieldViolation("ordertype", err))
	}

	return violations
}
