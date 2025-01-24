package api

import (
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type RequestTradeParams struct {
	Price       string         `json:"price" binding:"required,price"`
	Amount      string         `json:"amount" binding:"required,amount"`
	Fee         sql.NullString `json:"fee" binding:"required,fee"`
}


func (server *server) createTrade(ctx *gin.Context) {
	var req RequestTradeParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateTradeParams{
		Price:     req.Price,
		Amount:    req.Amount,
		Fee:       req.Fee,
	}

	trade, err := server.store.CreateTrade(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": trade.ID})
}

func (server *server) getTrade(ctx *gin.Context) {
	id := ctx.Param("id")
	tradeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	trade, err := server.store.GetTradeByID(ctx, tradeID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, trade)
}

func (server *server) deleteTrade(ctx *gin.Context) {
	id := ctx.Param("id")
	tradeID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.DeleteTrade(ctx, tradeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

