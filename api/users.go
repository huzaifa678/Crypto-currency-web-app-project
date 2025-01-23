package api

import (
	db "crypto-system/db/sqlc"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gin-contrib/sessions"
)

type UserRole string

const (
	Admin  UserRole = "admin"
	User   UserRole = "user"
	Guest  UserRole = "guest"
)


type userRequest struct {
	Email        string       `json:"email" binding:"required,email"` 
	PasswordHash string       `json:"password_hash" binding:"required"`
	Role         UserRole     `json:"role" binding:"required"`
}


func (server *server) createUser(ctx *gin.Context) {
	var req userRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateUserParams {
		Email: req.Email,
		PasswordHash: req.PasswordHash,
		Role: db.UserRole(req.Role),
		IsVerified: sql.NullBool{Bool: true, Valid: true},
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session := sessions.Default(ctx)
    session.Set("user_id", user.ID.String())
    
	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	log.Println("User created with ID:", user.ID)
	ctx.JSON(http.StatusOK, gin.H{"id": user.ID})
}

func (server *server) getUser(ctx *gin.Context) {

	var err error

	session := sessions.Default(ctx)
    id := session.Get("user_id")

	if id == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No session found"})
		return
	}

	if id == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

	

	parsedID, err := uuid.Parse(id.(string))

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

	var req userRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	id := ctx.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserParams {
		PasswordHash: req.PasswordHash,
		IsVerified: sql.NullBool{Bool: true, Valid: true},
		ID: userID,
	}

	err = server.store.UpdateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (server *server) deleteUser(ctx *gin.Context) {

	session := sessions.Default(ctx)
    id := session.Get("user_id")

	if id == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No session found"})
		return
	}

	if id == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

	parsedID, err := uuid.Parse(id.(string))

	err = server.store.DeleteUser(ctx, parsedID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}


