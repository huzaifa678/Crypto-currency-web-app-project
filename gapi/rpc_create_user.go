package gapi

import (
	"context"
	"log"
	"time"

	"github.com/hibiken/asynq"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/huzaifa678/Crypto-currency-web-app-project/val"
	"github.com/huzaifa678/Crypto-currency-web-app-project/worker"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Println("CreateUser called with request:", req)
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:    req.GetUsername(),
			Email:       req.GetEmail(),
			PasswordHash: hashedPassword,
			Role:        db.UserRole(req.GetRole().String()),
			IsVerified:  true,
		},
		AfterCreate: func(createUserRow db.CreateUserRow) error {
			taskPayload := &worker.PayloadSendVerifyEmail {
				Email: createUserRow.Email,
			}

			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	log.Println("ARGS:", arg)

	user, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	pbUser := convertCreateUser(user.CreateUserRow)

	res := &pb.CreateUserResponse{
		UserId: pbUser.Id,
	}

	return res, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		log.Println("ERROR: ", err)
		violations = append(violations, fieldViolation("username", err))
	}

	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		log.Println("ERROR: ", err)
		violations = append(violations, fieldViolation("email", err))
	}

	if err := val.ValidateString(req.GetPassword(), 6, 100); err != nil {
		log.Println("ERROR: ", err)
		violations = append(violations, fieldViolation("password", err))
	}

	if err := val.ValidateUserRole(req.GetRole()); err != nil {
		log.Println("ERROR: ", err)
		violations = append(violations, fieldViolation("role", err))
	}
	
	for _, v := range violations {
		log.Printf("- field: %s, description: %s", v.Field, v.Description)
	}

	return violations
}