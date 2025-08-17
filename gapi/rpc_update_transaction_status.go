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


func (server *server) UpdateTransactionStatus(ctx context.Context, req *pb.UpdateTransactionStatusRequest) (*pb.UpdateTransactionStatusResponse, error) {
	violations := validateUpdateTransactionStatus(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	transactionID, err := uuid.Parse(req.GetTransactionId())

	
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ID not parsed")
	}

	transaction, err := server.store.GetTransactionByID(ctx, transactionID)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "trade not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find the trade for the user")
	}

	if authPayload.Username != transaction.Username {
		return nil, status.Errorf(codes.Unknown, "unknown")
	}

	args := db.UpdateTransactionStatusParams {
		Status: db.TransactionStatus(req.GetStatus()),
		ID: 	transactionID,
	}

	err = server.store.UpdateTransactionStatus(ctx, args)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update the trade status for the user")
	}

	res := &pb.UpdateTransactionStatusResponse{
		Success: "Transaction status updated successfully",
	}

	return res, nil
}

func validateUpdateTransactionStatus(req *pb.UpdateTransactionStatusRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUpdateTransactionStatusRequest(req.GetTransactionId(), req.GetStatus()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}

	return violations
}

