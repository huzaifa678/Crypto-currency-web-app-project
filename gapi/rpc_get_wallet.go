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


func (server *server) GetWallet(ctx context.Context, req *pb.GetWalletRequest) (*pb.GetWalletResponse, error) {
	violations := validateGetWalletRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	walletID, err := uuid.Parse(req.GetWalletId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %v", err)
	}

	wallet, err := server.store.GetWalletByID(ctx, walletID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "wallet not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find the wallet for the user")
	}

	if authPayload.Username != wallet.Username {
		return nil, status.Errorf(codes.PermissionDenied, "Username does not match")
	}

	res := &pb.GetWalletResponse {
		Wallet: convertWallet(wallet),
	}
	return res, nil
}

func validateGetWalletRequest(req *pb.GetWalletRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateGetRequest(req.GetWalletId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}
	return violations
}