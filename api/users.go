package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserRole string

const (
	Admin  UserRole = "admin"
	User   UserRole = "user"
	Guest  UserRole = "guest"
)


type UserRequest struct {
	Username	 string	   `json:"username" binding:"required"`
	Email        string    `json:"email" binding:"required,email"` 
	Password 	 string    `json:"password_hash" binding:"required"`
	Role         UserRole  `json:"role" binding:"required"`
}

type UserLoginRequest struct {
	Email 	 string `json:"email" binding:"required,email"`
	Password string `json:"password_hash" binding:"required"`
}

type UserLoginResponse struct {
	SessionID				string  			`json:"session_id"`
	AccessToken 			string  			`json:"access_token"`
	AccessTokenExpiration 	time.Time 			`json:"access_token_expiration"`
	RefreshToken 			string 				`json:"refresh_token"`
	RefreshTokenExpiration 	time.Time 			`json:"refresh_token_expiration"`
	User					db.GetUserByEmailRow`json:"user"`
}


func (server *server) loginUser(ctx *gin.Context) {
	var req UserLoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		log.Println("Error in binding JSON", err)
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("error", err)
			ctx.JSON(http.StatusNotFound, gin.H{"Email not found": "Username not found with the given email"})
			return
		}
		log.Println("Error in getting user by email", err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	log.Println("user_email args", user)

	err = utils.ComparePasswords(user.PasswordHash, req.Password)
	log.Println("Password's Hash", user.PasswordHash)
	log.Println("Password", req.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	log.Println("User", user)

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration, token.TokenTypeAccessToken,)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenDuration, token.TokenTypeAccessToken,)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	args := db.CreateSessionParams {
		ID: refreshTokenPayload.ID,
		Username: user.Username,
		RefreshToken: refreshToken,
		UserAgent: ctx.Request.UserAgent(),
		ClientIp: ctx.ClientIP(),
		IsBlocked: false,
		ExpiresAt: refreshTokenPayload.ExpiredAt,
	}

	session, err := server.store.CreateSession(ctx, args)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}


	res := UserLoginResponse {
		SessionID: session.ID.String(),
		AccessToken: accessToken,
		AccessTokenExpiration: accessTokenPayload.ExpiredAt,
		RefreshToken: refreshToken,
		RefreshTokenExpiration: refreshTokenPayload.ExpiredAt,
		User: user,
	}

	ctx.JSON(http.StatusOK, res)

}


func (server *server) createUser(ctx *gin.Context) {
	var req UserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		log.Println("Error in binding JSON", err)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err = server.store.GetUserByEmail(ctx, req.Email)
	if err == nil { 
		ctx.JSON(http.StatusConflict, gin.H{"error": "User with this email already exists"})
		return
	}

	if err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams {
		Username: req.Username,
		Email: req.Email,
		PasswordHash: hashedPassword,
		Role: db.UserRole(req.Role),
		IsVerified: sql.NullBool{Bool: true, Valid: true},
	}


	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	log.Println("User created with ID:", user.ID)
	ctx.JSON(http.StatusOK, gin.H{"id": user.ID})
}

func (server *server) getUser(ctx *gin.Context) {

	var err error

	id := ctx.Param("id")

	parsedID, err := uuid.Parse(id)

	user, err := server.store.GetUserByID(ctx, parsedID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (server *server) updateUser(ctx *gin.Context) {

	var req UserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id := ctx.Param("id")

	if id == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

	parsedID, err := uuid.Parse(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams {
		PasswordHash: hashedPassword,
		IsVerified: sql.NullBool{Bool: true, Valid: true},
		ID: parsedID,
	}

	err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (server *server) deleteUser(ctx *gin.Context) {

	id := ctx.Param("id")

	if id == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

    userID, err := uuid.Parse(id)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

	if userID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

    err = server.store.DeleteUser(ctx, userID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
