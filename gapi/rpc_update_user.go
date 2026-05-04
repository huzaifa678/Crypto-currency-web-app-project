package gapi

import (
	"context"
	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (server *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	id, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %v", err)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
	}

	arg := db.UpdateUserParams{
		PasswordHash: hashedPassword,
		IsVerified:   false,
		ID:           id,
	}

	err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}


	res := &pb.UpdateUserResponse {
		Message: "successfully updated the user information",
	}

	return res, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUpdateUser(req.GetUserId(), req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("id", err))
		violations = append(violations, fieldViolation("password", err))
	}
	return violations
}
