package api

import (
    db "crypto-system/db/sqlc"
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type walletRequest struct {
    UserEmail string         `json:"user_email" binding:"required,email"`
    Currency  string         `json:"currency"`
}

type updateWalletRequest struct {
	Balance       sql.NullString `json:"balance"`
	LockedBalance sql.NullString `json:"locked_balance"`
	UserEmail     string         `json:"user_email"`
	Currency      string         `json:"currency"`
}

func (server *server) createWallet(ctx *gin.Context) {
    var req walletRequest

    arg := db.CreateWalletParams{
        UserEmail: req.UserEmail,
        Currency:  req.Currency,
        Balance:   sql.NullString{String: "0", Valid: true},
    }

    wallet, err := server.store.CreateWallet(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"email": wallet.UserEmail, "currency": wallet.Currency})
}

func (server *server) getWallet(ctx *gin.Context) {
	userEmail := ctx.Query("user_email")
    currency := ctx.Query("currency")

    if userEmail == "" || currency == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_email and currency are required"})
        return
    }

    arg := db.GetWalletByUserEmailAndCurrencyParams{
        UserEmail: userEmail,
        Currency:  currency,
    }

    wallet, err := server.store.GetWalletByUserEmailAndCurrency(ctx, arg)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, wallet)
}

func (server *server) updateWallet(ctx *gin.Context) {

	var req updateWalletRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

	arg := db.UpdateWalletBalanceParams{
		Balance:       req.Balance,
		LockedBalance: req.LockedBalance,
		UserEmail:     req.UserEmail,
		Currency:      req.Currency,
	}

    err := server.store.UpdateWalletBalance(ctx, arg)
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

    err = server.store.DeleteWallet(ctx, walletID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}