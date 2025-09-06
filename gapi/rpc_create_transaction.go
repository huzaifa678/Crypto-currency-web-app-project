package gapi

import (
	"context"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (server *server) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	violations := validateCreateTransactionRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	args := db.CreateTransactionParams {
		Username: authPayload.Username,
		UserEmail: req.GetUserEmail(),
		Type: db.TransactionType(req.GetType()),
		Currency: req.GetCurrency(),
		Amount: req.GetAmount(),
		Address: req.GetAddress(),
		TxHash: req.TxHash,
	}
	
	transaction, err := server.store.CreateTransaction(ctx, args)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create transaction: %v", err)
	}

	res := &pb.CreateTransactionResponse{
		Transaction: convertCreateTransaction(req.GetUsername(), transaction),
	}

	return res, nil
}

func validateCreateTransactionRequest(req *pb.CreateTransactionRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateCreateTransactionRequest(req.UserEmail, req.Amount, req.Type); err != nil {
		violations = append(violations, fieldViolation("user_email", err))
		violations = append(violations, fieldViolation("amount", err))
		violations = append(violations, fieldViolation("type", err))
	}

	return violations
}