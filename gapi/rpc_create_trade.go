package gapi

import (
	"context"
	"log"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"github.com/shopspring/decimal"
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

	buyerUserEmail := req.GetBuyerUserEmail()

	sellerUserEmail := req.GetSellerUserEmail()

	marketID, err := uuid.Parse(req.GetMarketId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid market ID: %v", err)
	}

	price, err := decimal.NewFromString(req.GetPrice())
	if err != nil {
    	return nil, status.Errorf(codes.InvalidArgument, "invalid price: %v", err)
	}

	amount, err := decimal.NewFromString(req.GetAmount())
	if err != nil {
    	return nil, status.Errorf(codes.InvalidArgument, "invalid amount: %v", err)
	}

	fee, err := decimal.NewFromString(req.GetFee())
	if err != nil {
    	return nil, status.Errorf(codes.InvalidArgument, "invalid fee: %v", err)
	}

	args := db.CreateTradeTxParams{
		TradeParams: db.CreateTradeParams{	
			Username: authPayload.Username,
			BuyerUserEmail: buyerUserEmail,
			SellerUserEmail: sellerUserEmail,
			BuyOrderID: buyOrderId,
			SellOrderID: sellOrderId,
			MarketID: marketID,
			Price: price,
			Amount: amount,
			Fee: fee,
		},
	}

	trade, err := server.store.CreateTradeTx(ctx, args)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create trade: %v", err)
	}

	convertToPb := convertTrade(trade.Trade)

	res := &pb.CreateTradeResponse{
		Trade: convertToPb,
	}

	return res, nil
}

func validateCreateTradeRequest(req *pb.CreateTradeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	price, err := decimal.NewFromString(req.GetPrice())
	if err != nil {
		violations = append(violations, fieldViolation("price", err))
	}

	amount, err := decimal.NewFromString(req.GetAmount())
	if err != nil {
		violations = append(violations, fieldViolation("amount", err))
	}

	fee, err := decimal.NewFromString(req.GetFee())
	if err != nil {
		violations = append(violations, fieldViolation("fee", err))
	}

	if len(violations) == 0 {
		if err := val.ValidateCreateTradeRequest(
			req.GetBuyerUserEmail(),
			req.GetSellerUserEmail(),
			req.GetBuyOrderId(),
			req.GetSellOrderId(),
			req.GetMarketId(),
			price,
			amount,
			fee,
		); err != nil {
			log.Println("ERROR", err)
			violations = append(violations, fieldViolation("buy_order_id", err))
			violations = append(violations, fieldViolation("sell_order_id", err))
			violations = append(violations, fieldViolation("market_id", err))
			violations = append(violations, fieldViolation("price", err))
			violations = append(violations, fieldViolation("amount", err))
			violations = append(violations, fieldViolation("fee", err))
		}
	}

	return violations
}
