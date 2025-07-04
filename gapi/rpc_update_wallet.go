package gapi

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *server) UpdateWallet(ctx context.Context, req *pb.UpdateWalletRequest) (*pb.UpdateWalletResponse, error) {
	violations := validateUpdateWalletRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	walletID, err := uuid.Parse(req.GetWalletId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %v", err)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	arg := db.UpdateWalletBalanceParams{
		Balance:       sql.NullString{String: req.GetBalance(), Valid: req.GetBalance() != ""},
		LockedBalance: sql.NullString{String: req.GetLockedBalance(), Valid: req.GetLockedBalance() != ""},
		ID:            walletID,
	}

	err = server.store.UpdateWalletBalance(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update wallet: %v", err)
	}

	wallet, _ := server.store.GetWalletByID(ctx, walletID)

	if authPayload.Username != wallet.Username {
		return nil, status.Errorf(codes.Unknown, "unknown")
	}

	res := &pb.UpdateWalletResponse {
		Message: "successfully updated the wallet",
	}

	return res, nil
}


func validateUpdateWalletRequest(req *pb.UpdateWalletRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUpdateWalletRequest(req.GetWalletId(), req.GetBalance(), req.GetLockedBalance()); err != nil {
		violations = append(violations, fieldViolation("id", err))
		violations = append(violations, fieldViolation("balance", err))
		violations = append(violations, fieldViolation("lockedbalance", err))
	}

	return violations
}
