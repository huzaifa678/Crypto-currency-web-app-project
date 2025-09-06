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


func (server *server) DeleteWallet(ctx context.Context, req *pb.DeleteWalletRequest) (*pb.DeleteWalletResponse, error) {
	violations := validateDeleteWalletRequest(req)
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

	wallet, _ := server.store.GetWalletByID(ctx, walletID)

	if authPayload.Username != wallet.Username {
		return nil, status.Errorf(codes.Unknown, "Not authorized")
	}

	err = server.store.DeleteWallet(ctx, walletID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "wallet not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete the wallet for the user")
	}


	res := &pb.DeleteWalletResponse {
		Message: "successfully delete the wallet",
	}

	return res, nil
}

func validateDeleteWalletRequest(req *pb.DeleteWalletRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateDeleteRequest(req.GetWalletId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return violations
}


