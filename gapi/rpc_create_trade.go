package gapi

import (
	"context"
	"log"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (server *server) CreateTrade(ctx context.Context, req *pb.CreateTradeRequest) (*pb.CreateTradeResponse, error) {
	violations := validateCreateTradeRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	buyOrderId, err := uuid.Parse(req.GetBuyOrderId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid buy order ID: %v", err)
	}

	sellOrderId, err := uuid.Parse(req.GetSellOrderId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid sell order ID: %v", err)
	}

	marketID, err := uuid.Parse(req.GetMarketId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid market ID: %v", err)
	}

	args := db.CreateTradeParams {
		Username: authPayload.Username,
		BuyOrderID: buyOrderId,
		SellOrderID: sellOrderId,
		MarketID: marketID,
		Price: req.GetPrice(),
		Amount: req.GetAmount(),
		Fee: req.GetFee(),
	}

	trade, err := server.store.CreateTrade(ctx, args)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create trade: %v", err)
	}

	convertToPb := convertTrade(trade)

	res := &pb.CreateTradeResponse{
		Trade: convertToPb,
	}

	return res, nil
}

func validateCreateTradeRequest(req *pb.CreateTradeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateCreateTradeRequest(req.GetBuyOrderId(), req.GetSellOrderId(), req.GetMarketId(), req.GetPrice(), req.GetAmount(), req.GetFee()); err != nil {
		log.Println("ERROR", err)
		violations = append(violations, fieldViolation("buy_order_id", err))
		violations = append(violations, fieldViolation("sell_order_id", err))
		violations = append(violations, fieldViolation("market_id", err))
		violations = append(violations, fieldViolation("price", err))
		violations = append(violations, fieldViolation("amount", err))
		violations = append(violations, fieldViolation("fee", err))
	}

	return violations
}