package api

import (
    db "crypto-system/db/sqlc"
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
)

type auditLogRequest struct {
    UserEmail string `json:"user_email"`
    Action    string `json:"action"`
    IPAddress sql.NullString `json:"ip_address"`
}

func (server *server) createAuditLog(ctx *gin.Context) {
    var req auditLogRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    arg := db.CreateAuditLogParams{
        UserEmail: req.UserEmail,
        Action:    req.Action,
        IpAddress: req.IPAddress,
    }

    auditLog, err := server.store.CreateAuditLog(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"email": auditLog.UserEmail})
}

func (server *server) getAuditLog(ctx *gin.Context) {
    email := ctx.Param("email")

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

func (server *server) listUserAuditLogs(ctx *gin.Context) {
    email := ctx.Param("user_email")

    auditLogs, err := server.store.GetAuditLogsByUserEmail(ctx, email)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, auditLogs)
}