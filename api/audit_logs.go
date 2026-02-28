//nolint:revive
package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"

	"github.com/gin-gonic/gin"
)

type AuditLogRequest struct {
	UserEmail string `json:"user_email"`
	Action    string `json:"action"`
	IPAddress string `json:"ip_address"`
}

func (server *server) createAuditLog(ctx *gin.Context) {
	var req AuditLogRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

	arg := db.CreateAuditLogParams{
		Username:  authPayload.Username,
		UserEmail: req.UserEmail,
		Action:    req.Action,
		IpAddress: req.IPAddress,
	}

	auditLog, err := server.store.CreateAuditLog(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user_email": auditLog.UserEmail})
}

func (server *server) getAuditLog(ctx *gin.Context) {
	email := ctx.Param("user_email")

	log.Printf("Extracted userEmail: %s", email)

	auditLog, err := server.store.GetAuditLogsByUserEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

	if authPayload.Username != auditLog[1].Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	ctx.JSON(http.StatusOK, auditLog)
}

func (server *server) DeleteAuditLog(ctx *gin.Context) {
	id := ctx.Param("id")
	userEmail := ctx.Query("user_email")

	log.Printf("Extracted userEmail: %s", userEmail)

	if userEmail == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("user_email is required")))
		return
	}

	auditLogs, err := server.store.GetAuditLogsByUserEmail(ctx, userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	auditLogID, err := uuid.Parse(id)
	if err != nil || auditLogID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if len(auditLogs) == 0 {
		ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("no audit logs found for the given user email")))
		return
	}

	authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

	if authPayload.Username != auditLogs[0].Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	err = server.store.DeleteAuditLog(ctx, auditLogID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (server *server) listUserAuditLogs(ctx *gin.Context) {
	email := ctx.Param("user_email")

	auditLogs, err := server.store.GetAuditLogsByUserEmail(ctx, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, auditLogs)
}
