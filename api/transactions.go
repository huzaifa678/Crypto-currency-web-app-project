package api

import (
	db "crypto-system/db/sqlc"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionType string

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"
)


type transactionRequest struct {
	UserID   uuid.UUID       `json:"user_id"`
	Type     TransactionType `json:"type"`
	Currency string          `json:"currency"`
	Amount   string          `json:"amount"`
	Address  sql.NullString  `json:"address"`
	TxHash   sql.NullString  `json:"tx_hash"`
}

func (server *server) createTransaction(ctx *gin.Context) {
	var req transactionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTransactionParams {
		UserID: req.UserID,
		Type: db.TransactionType(req.Type),
		Currency: req.Currency,
		Amount: req.Amount,
		Address: req.Address,
		TxHash: req.TxHash,
	}

	transaction, err := server.store.CreateTransaction(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": transaction.ID})
}