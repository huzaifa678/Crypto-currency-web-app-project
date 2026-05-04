package gapi

import (
	"context"
	"fmt"
	"log"

	"github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (server *server) GetTransactionsByUserEmail(ctx context.Context, req *pb.GetTransactionsByUserEmailRequest) (*pb.GetTransactionsByUserEmailResponse, error) {
	violations := validateGetTransactionByUserEmailRequest(req)

	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	transactions, err := server.store.GetTransactionsByUserEmail(ctx, req.GetUserEmail())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list the transactions for the user")
	}

	// 🔹 If no transactions exist
	if len(transactions) == 0 {
		return &pb.GetTransactionsByUserEmailResponse{
			Transactions: []*pb.Transaction{}, // empty list
		}, nil
	}

	// Use the first transaction record ever made by the user
	if transactions[0].Username != authPayload.Username {
		log.Println("USERNAME", transactions[0].Username)
		log.Println("AUTHUSERNAME", authPayload.Username)
		return nil, status.Errorf(codes.Unknown, "Not authorized")
	}


	log.Println("FIRST TX", transactions[0])


	res := &pb.GetTransactionsByUserEmailResponse{
		Transactions: convertTransactionList(transactions),
	}

	fmt.Println(res.Transactions[0].Status) 

	log.Println("PB TX", res)

	return res, nil
}


func validateGetTransactionByUserEmailRequest(req *pb.GetTransactionsByUserEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmail(req.GetUserEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}