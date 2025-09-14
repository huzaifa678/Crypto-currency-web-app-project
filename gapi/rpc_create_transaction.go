package gapi

import (
	"context"
	"log"
	"strings"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"github.com/shopspring/decimal"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *server) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	log.Println("req", req)

	amount := decimal.NewFromInt(req.GetAmount()).Div(decimal.New(1, scale))

	violations := validateCreateTransactionRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	pbType := strings.ToLower(req.GetType().String())

	log.Println("TYPE", pbType)

	args := db.CreateTransactionParams {
		Username: authPayload.Username,
		UserEmail: req.GetUserEmail(),
		Type: db.TransactionType(pbType),
		Currency: req.GetCurrency(),
		Amount: amount,
		Address: req.GetAddress(),
		TxHash: req.TxHash,
	}
	
	transaction, err := server.store.CreateTransaction(ctx, args)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create transaction: %v", err)
	}

	transaction.Status = db.TransactionStatus("completed")

	log.Println("STATUS", transaction.Status)

	res := &pb.CreateTransactionResponse{
		Transaction: convertCreateTransaction(req.GetUsername(), transaction),
	}

	log.Println("PB STATUS", res.Transaction.Status)

	return res, nil
}

func validateCreateTransactionRequest(req *pb.CreateTransactionRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateCreateTransactionRequest(req.UserEmail, decimal.NewFromFloat(float64(req.Amount)), req.Type); err != nil {
		violations = append(violations, fieldViolation("user_email", err))
		violations = append(violations, fieldViolation("amount", err))
		violations = append(violations, fieldViolation("type", err))
	}

	return violations
}