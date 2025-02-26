package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type TransactionRequest struct {
	UserEmail string         `json:"user_email"`
	Type     db.TransactionType `json:"type"`
	Currency string          `json:"currency"`
	Amount   string          `json:"amount"`
	Address  string  `json:"address"`
	TxHash   string  `json:"tx_hash"`
}

func (server *server) createTransaction(ctx *gin.Context) {
	var req TransactionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
        log.Printf("JSON Binding Error: %v", err)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

    currencies := []string{"USD", "EUR", "BTC", "ETH", "JPY"}

    if !isValidCurrency(req.Currency, currencies) {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid currency"})
        return
    }

    authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

	arg := db.CreateTransactionParams{
        Username: authPayload.Username,
		UserEmail: req.UserEmail,
		Type: db.TransactionType(req.Type),
		Currency: req.Currency,
		Amount: req.Amount,
		Address: sql.NullString{String: req.Address, Valid: req.Address != ""},
		TxHash: sql.NullString{String: req.TxHash, Valid: req.TxHash != ""},
	}

	transaction, err := server.store.CreateTransaction(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": transaction.ID})
}

func (server *server) getTransaction(ctx *gin.Context) {
    id := ctx.Param("id")
    transactionID, err := uuid.Parse(id)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    if transactionID == uuid.Nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error":  "invalid UUID"})
        return
    }

    transaction, err := server.store.GetTransactionByID(ctx, transactionID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

    if authPayload.Username != transaction.Username {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
        return
    }

    ctx.JSON(http.StatusOK, transaction)
}

func (server *server) listUserTransactions(ctx *gin.Context) {
    email := ctx.Param("user_email")
    fmt.Printf("Extracted Email: '%s'\n", email)

    transactions, err := server.store.GetTransactionsByUserEmail(ctx, email)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, transactions)
}

func (server *server) deleteTransaction(ctx *gin.Context) {
    id := ctx.Param("id")
    transactionID, err := uuid.Parse(id)

    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    if transactionID == uuid.Nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error":  "invalid UUID"})
        return
    }

    transaction, err := server.store.GetTransactionByID(ctx, transactionID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

    if authPayload.Username != transaction.Username {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
        return
    }

    err = server.store.DeleteTransaction(ctx, transactionID)
    if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

    ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}