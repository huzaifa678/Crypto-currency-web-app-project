package gapi

import (
	"context"
	"database/sql"
	"errors"
	"log"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	_, err := server.store.GetUserByEmail(ctx, req.GetEmail())

	if err != nil {
		if !errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.Internal, "failed to check existing user")
		}
	} else {
		return nil, status.Errorf(codes.AlreadyExists, "user with this email already exists")
	}


	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	arg := db.CreateUserParams{
		Username:    req.GetUsername(),
		Email:       req.GetEmail(),
		PasswordHash: hashedPassword,
		Role:        db.UserRole(req.GetRole().String()),
		IsVerified:  sql.NullBool{Bool: true, Valid: true},
	}

	log.Println("ARGS:", arg)

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	pbUser := convertCreateUser(user)

	res := &pb.CreateUserResponse{
		UserId: pbUser.Id,
	}

	return res, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	if err := val.ValidateString(req.GetPassword(), 6, 100); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateUserRole(req.GetRole()); err != nil {
		violations = append(violations, fieldViolation("role", err))
	}

	log.Println("Validation errors:")
	for _, v := range violations {
		log.Printf("- field: %s, description: %s", v.Field, v.Description)
	}

	return violations
}
