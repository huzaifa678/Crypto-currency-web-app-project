package gapi

import (
	"context"
	"database/sql"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *server) CreateMarket(ctx context.Context, req *pb.CreateMarketRequest) (*pb.CreateMarketResponse, error) {
	violations := validateCreateMarketRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	arg := db.CreateMarketParams{
		Username: authPayload.Username,
		BaseCurrency:  req.GetBaseCurrency(),
		QuoteCurrency: req.GetQuoteCurrency(),
		MinOrderAmount: sql.NullString{String: req.GetMinOrderAmount(), Valid: req.GetMinOrderAmount() != ""},
		PricePrecision: sql.NullInt32{Int32: req.GetPricePrecision(), Valid: req.GetPricePrecision() != 0},
	}

	market, err := server.store.CreateMarket(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create market: %v", err)
	}

	return &pb.CreateMarketResponse{MarketId: market.ID.String()}, nil
}

func validateCreateMarketRequest(req *pb.CreateMarketRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateCreateMarketRequest(req.GetBaseCurrency(), req.GetQuoteCurrency(), req.GetMinOrderAmount(), req.GetPricePrecision()); err != nil {
		violations = append(violations, fieldViolation("base_currency", err))
		violations = append(violations, fieldViolation("quote_currency", err))
		violations = append(violations, fieldViolation("min_order_amount", err))
		violations = append(violations, fieldViolation("price_precision", err))
	}

	return violations
}


