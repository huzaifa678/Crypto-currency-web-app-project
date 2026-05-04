package gapi

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (server *server) DeleteTrade(ctx context.Context, req *pb.DeleteTradeRequest) (*pb.DeleteTradeResponse, error) {
	violations := validateDeleteTradeRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	tradeID, err := uuid.Parse(req.GetTradeId())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ID not parsed")
	}

	trade, err := server.store.GetTradeByID(ctx, tradeID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "trade not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find the trade for the user")
	}

	if authPayload.Username != trade.Username {
		log.Println(err)
		return nil, status.Errorf(codes.PermissionDenied, "Not authorized")
	}

	err = server.store.DeleteTrade(ctx, tradeID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "trade not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete the trade for the user")
	}

	res := &pb.DeleteTradeResponse{
		Success: "Succesfully deleted the trade",
	}

	return res, nil
}

func validateDeleteTradeRequest(req *pb.DeleteTradeRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateDeleteRequest(req.GetTradeId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return violations
}