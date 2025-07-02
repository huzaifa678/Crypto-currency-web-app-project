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


func (server *server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	violations := validateDeleteUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	id, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %v", err)
	}

	err = server.store.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	res := &pb.DeleteUserResponse {
		Message: "successfully deleted the user",
	}

	return res, nil
}

func validateDeleteUserRequest(req *pb.DeleteUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateDeleteRequest(req.GetUserId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}
	return violations
}
