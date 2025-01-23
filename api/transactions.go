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
	UserEmail string         `json:"user_email"`
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
		UserEmail: req.UserEmail,
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

func (server *server) getTransaction(ctx *gin.Context) {
    id := ctx.Param("id")
    transactionID, err := uuid.Parse(id)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
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

    ctx.JSON(http.StatusOK, transaction)
}

func (server *server) listUserTransactions(ctx *gin.Context) {
    email := ctx.Param("user_email")

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

    err = server.store.DeleteTransaction(ctx, transactionID)
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