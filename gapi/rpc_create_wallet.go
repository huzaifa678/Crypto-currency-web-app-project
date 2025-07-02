package gapi

import (
	"context"
	"database/sql"
	"log"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *server) CreateWallet(ctx context.Context, req *pb.CreateWalletRequest) (*pb.CreateWalletResponse, error) {
	violations := validateCreateWalletRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		log.Println("ERROR:", err.Error())
		return nil, unauthenticatedError(err)
	}

	arg := db.CreateWalletParams{
		Username:  authPayload.Username,
		UserEmail: req.GetUserEmail(),
		Currency:  req.GetCurrency(),
		Balance:   sql.NullString{String: "0", Valid: true},
	}

	wallet, err := server.store.CreateWallet(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create wallet: %v", err)
	}

	return &pb.CreateWalletResponse{
		WalletId: wallet.ID.String(),
	}, nil
}

func validateCreateWalletRequest(req *pb.CreateWalletRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(req.GetUserEmail()); err != nil {
		violations = append(violations, fieldViolation("user_email", err))
	}
	if err := val.ValidateCurrency(req.GetCurrency()); err != nil {
		violations = append(violations, fieldViolation("currency", err))
	}
	return violations
}
