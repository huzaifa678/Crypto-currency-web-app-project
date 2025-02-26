package api

import (
	"database/sql"
	"net/http"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WalletRequest struct {
    UserEmail string         `json:"user_email" binding:"required,email"`
    Currency  string         `json:"currency"`
}

type UpdateWalletRequest struct {
	Balance       sql.NullString `json:"balance"`
	LockedBalance sql.NullString `json:"locked_balance"`
	UserEmail     string         `json:"user_email"`
	Currency      string         `json:"currency"`
}

func (server *server) createWallet(ctx *gin.Context) {
    var req WalletRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    authPayload, _ := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

    arg := db.CreateWalletParams{
        Username:  authPayload.Username,
        UserEmail: req.UserEmail,
        Currency:  req.Currency,
        Balance:   sql.NullString{String: "0", Valid: true},
    }

    wallet, err := server.store.CreateWallet(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"id": wallet.ID})
}

func (server *server) getWallet(ctx *gin.Context) {

    id := ctx.Param("id")


    walletID, err := uuid.Parse(id)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    if walletID == uuid.Nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
        return
    }


    wallet, err := server.store.GetWalletByID(ctx, walletID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }


    authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

    if authPayload.Username != wallet.Username {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    ctx.JSON(http.StatusOK, wallet)
}

func (server *server) updateWallet(ctx *gin.Context) {
	var req UpdateWalletRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    id := ctx.Param("id")

	if id == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

	walletID, err := uuid.Parse(id)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

    wallet, err := server.store.GetWalletByID(ctx, walletID)

    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }
    authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

    if authPayload.Username != wallet.Username {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

	arg := db.UpdateWalletBalanceParams{
		Balance:       req.Balance,
		LockedBalance: req.LockedBalance,
	}

    err = server.store.UpdateWalletBalance(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (server *server) deleteWallet(ctx *gin.Context) {
    id := ctx.Param("id")
    walletID, err := uuid.Parse(id)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    if walletID == uuid.Nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
        return
    }

    wallet, err := server.store.GetWalletByID(ctx, walletID)

    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

    if authPayload.Username != wallet.Username {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    err = server.store.DeleteWallet(ctx, walletID)

    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}