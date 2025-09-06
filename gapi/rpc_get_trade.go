package gapi

import (
	"context"
	"errors"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (server *server) GetTrade(ctx context.Context, req *pb.GetTradeByIDRequest) (*pb.GetTradeByIDResponse, error) {
	violations := validateGetTradeRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	tradeID, err := uuid.Parse(req.GetTradeId())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ID not parsed")
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	trade, err := server.store.GetTradeByID(ctx, tradeID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "trade not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find the trade for the user")
	}

	if authPayload.Username != trade.Username {
		return nil, status.Errorf(codes.Unknown, "unknown")
	}

	convertToPb := convertTrade(trade)

	res := &pb.GetTradeByIDResponse{
		Trade: convertToPb,
	}

	return res, nil
}

func validateGetTradeRequest(req *pb.GetTradeByIDRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateGetRequest(req.GetTradeId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return violations
}