package gapi

import (
	"context"
	"log"
	"strings"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var scale int32 = 8


func (server *server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log.Println("RECEIVED CreateOrder request:", req)

	userEmail := req.GetUserEmail()
	baseCurrency := req.GetBaseCurrency()
	quoteCurrency := req.GetQuoteCurrency()

	price := decimal.NewFromInt(req.GetPrice()).Div(decimal.New(1, scale))

	amount := decimal.NewFromInt(req.GetAmount()).Div(decimal.New(1, scale))

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

	err = server.store.OrderForCurrencyTx(ctx, db.OrderForCurrencyTxParams{
		UserEmail:     userEmail,
		BaseCurrency:  baseCurrency,
		QuoteCurrency: quoteCurrency,
	})

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user cannot place order for the given currencies: %v", err)
	}

	pbType := strings.ToLower(req.GetType().String())

	arg := db.CreateOrderParams{
		Username:  authPayload.Username,
		UserEmail: req.GetUserEmail(),
		MarketID:  marketID,
		Type:      db.OrderType(pbType),
		Status:    db.OrderStatus("open"), // default status
		Price:     price,
		Amount:    amount,
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
	if err := val.ValidateCreateOrderRequest(req.GetUserEmail(), req.GetMarketId(), decimal.NewFromFloat(float64(req.GetPrice())), decimal.NewFromFloat(float64(req.GetAmount())), req.GetType()); err != nil {
		violations = append(violations, fieldViolation("user_email", err))
		violations = append(violations, fieldViolation("id", err))
		violations = append(violations, fieldViolation("price", err))
		violations = append(violations, fieldViolation("amount", err))
		violations = append(violations, fieldViolation("ordertype", err))
	}

	return violations
}
