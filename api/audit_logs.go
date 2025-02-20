package api

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"

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

    arg := db.CreateAuditLogParams{
        UserEmail: req.UserEmail,
        Action:    req.Action,
        IpAddress: sql.NullString{String: req.IPAddress, Valid: req.IPAddress != ""},
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

    auditLog, err := server.store.GetAuditLogsByUserEmail(ctx, email)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, auditLog)
}

func (server *server) DeleteAuditLog(ctx *gin.Context) {
    id := ctx.Param("id")

    auditLogId, err := uuid.Parse(id)

    if auditLogId == uuid.Nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    err = server.store.DeleteAuditLog(ctx, auditLogId)

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