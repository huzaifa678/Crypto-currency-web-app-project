package gapi

import (
	"context"
	"errors"
	"log"
	"strings"

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

	log.Println("UPDATE REQ", req)

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
		return nil, status.Errorf(codes.Unknown, "Not authorized")
	}

	log.Println("FETCHED TRANSACTION STATUS", transaction.Status)
	log.Println("STATUS ENUM", req.GetStatus())
	log.Println("STRINGED STATUS ENUM", req.GetStatus().String())

	pbStatus := strings.ToLower(req.GetStatus().String())

	log.Println("PB STATUS", pbStatus)

	args := db.UpdateTransactionStatusParams {
		Status: db.TransactionStatus(pbStatus),
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
		violations = append(violations, fieldViolation("status", err))
	}

	return violations
}

