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


func (server *server) DeleteTransaction(ctx context.Context, req *pb.DeleteTransactionRequest) (*pb.DeleteTransactionResponse, error) {
	violations := validateDeleteTransactionRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	transactionID, err := uuid.Parse(req.GetTransactionId())

	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ID not parsed")
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	transaction, err := server.store.GetTransactionByID(ctx, transactionID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "transaction not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find the transaction for the user")
	}

	if authPayload.Username != transaction.Username {
		return nil, status.Errorf(codes.Unknown, "unknown")
	}

	err = server.store.DeleteTransaction(ctx, transactionID)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete the transaction for the user")
	}

	res := &pb.DeleteTransactionResponse{
		Success: "Successfully deleted the transaction",
	}

	return res, nil
}


func validateDeleteTransactionRequest(req *pb.DeleteTransactionRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateGetRequest(req.GetTransactionId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}
	return violations
}