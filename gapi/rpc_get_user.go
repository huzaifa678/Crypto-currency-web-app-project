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


func (server *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	violations := validateGetUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	id, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid UUID: %v", err)
	}

	user, err := server.store.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	res := &pb.GetUserResponse{
		User: convertGetUser(user),
	}

	return res, nil
}

func validateGetUserRequest(req *pb.GetUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateGetRequest(req.GetUserId()); err != nil {
		violations = append(violations, fieldViolation("id", err))
	}
	return violations
}
