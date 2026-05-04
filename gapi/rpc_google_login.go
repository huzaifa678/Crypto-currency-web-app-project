package gapi

import (
	"context"
	"errors"
	"log"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"github.com/huzaifa678/Crypto-currency-web-app-project/oauth2"
	"github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var oauth2VerifyGoogleIDToken = oauth2.VerifyGoogleIDToken

func (server *server) GoogleLogin(ctx context.Context, req *pb.GoogleLoginRequest) (*pb.GoogleLoginResponse, error) {
    payload, err := oauth2VerifyGoogleIDToken(ctx, req.IdToken, server.config.GoogleClientID)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "invalid google id_token: %s", err)
    }

    email, _ := payload.Claims["email"].(string)
    sub, _ := payload.Claims["sub"].(string)
    name, _ := payload.Claims["name"].(string)
    password, _ := payload.Claims["password"].(string)

    user, err := server.store.GetGoogleUserByProviderID(ctx, sub)
    if err != nil {
        if errors.Is(err, db.ErrRecordNotFound) {
            user, err = server.store.CreateGoogleUser(ctx, db.CreateGoogleUserParams{
                Email: email, 
                Username: name, 
                ProviderID: sub,
            })

            if err != nil {
                return nil, status.Errorf(codes.Internal, "failed to create google user: %s", err)
            }
        }
    }

    accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration, token.TokenTypeAccessToken)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "cannot create access token: %s", err)
    }

    refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username,server.config.RefreshTokenDuration, token.TokenTypeRefreshToken)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "cannot create refresh token: %s", err)
    }

	_, err = server.store.GetUserByEmail(ctx, email)
	if err != nil  {
        if errors.Is(err, db.ErrRecordNotFound) {
            _, err = server.store.CreateUser(ctx, db.CreateUserParams{
        	    Username:  user.Username,
        	    Email:     user.Email,
        	    PasswordHash: password,
        	    Role: db.UserRole(user.Role.String),
                IsVerified: user.CreatedAt.Valid,
    	    })
        }


		if err != nil {
        	return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
    	}
	}

    mtdt := server.extractMetadata(ctx)

    _, err = server.store.CreateSession(ctx, db.CreateSessionParams{
        ID:           refreshPayload.ID,
        Username:     user.Username,
        RefreshToken: refreshToken,
        UserAgent:    mtdt.UserAgent,
        ClientIp:     mtdt.ClientIP,
        IsBlocked:    false,
        ExpiresAt:    refreshPayload.ExpiredAt,
    })

    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to create session: %s", err)
    }

    log.Println("token duration", timestamppb.New(accessPayload.ExpiredAt), timestamppb.New(refreshPayload.ExpiredAt))

    return &pb.GoogleLoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        Client: &pb.LoginClient {
            Id:       user.ID.String(),
            Email:    user.Email,
            Username: user.Username,
            Role:     user.Role.String,
            CreatedAt: timestamppb.New(user.CreatedAt.Time),
        },
        AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
        RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
    }, nil
}