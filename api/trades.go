package api

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type RequestTradeParams struct {
	BuyOrderID  uuid.UUID 	   `json:"buy_order_id" binding:"required"`
	SellOrderID uuid.UUID	   `json:"sell_order_id" binding:"required"`
	MarketID    uuid.UUID 	   `json:"market_id" binding:"required"`
	Price       string         `json:"price" binding:"required"`
	Amount      string         `json:"amount" binding:"required"`
	Fee         string 		   `json:"fee" binding:"required"`
}


func (server *server) createTrade(ctx *gin.Context) {
	var req RequestTradeParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !isValidFee(req.Fee) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid fee"})
		return
	}

	arg := db.CreateTradeParams{
		BuyOrderID:  req.BuyOrderID,
		SellOrderID: req.SellOrderID,
		MarketID:    req.MarketID,
		Price:     req.Price,
		Amount:    req.Amount,
		Fee:       sql.NullString{String: req.Fee, Valid: req.Fee != ""},
	}


	trade, err := server.store.CreateTrade(ctx, arg)
	if err != nil {
		log.Println("ERROR:", err)
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

	if tradeID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
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


func (server *server) listTrades(ctx *gin.Context) {
    id := ctx.Param("market_id")
	marketID, err := uuid.Parse(id)

    trades, err := server.store.GetTradesByMarketID(ctx, marketID)

    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, trades)
}


func isValidFee(fee string) bool {
	_, err := strconv.ParseFloat(fee, 64)
	return err == nil
}

