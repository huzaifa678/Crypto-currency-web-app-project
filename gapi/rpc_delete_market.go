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


func (server *server) DeleteMarket(ctx context.Context, req *pb.DeleteMarketRequest) (*pb.DeleteMarketResponse, error) {
	violations := validateDeleteMarketRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	marketID, err := uuid.Parse(req.GetMarketId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %v", err)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	market, err := server.store.GetMarketByID(ctx, marketID)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Not found")
	}

	if authPayload.Username != market.Username {
		return nil, status.Errorf(codes.Unknown, "unknown")
	}

	err = server.store.DeleteMarket(ctx, marketID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "market not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find market for the user")
	}

	res := &pb.DeleteMarketResponse {
		Message: "Market deleted successfully",
	}

	return res, nil
}

func validateDeleteMarketRequest(req *pb.DeleteMarketRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateDeleteRequest(req.GetMarketId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}
	return violations
}

