package api

import (
    db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type marketRequest struct {
    BaseCurrency  string `json:"base_currency"`
    QuoteCurrency string `json:"quote_currency"`
}

func (server *server) createMarket(ctx *gin.Context) {
    var req marketRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    arg := db.CreateMarketParams{
        BaseCurrency:  req.BaseCurrency,
        QuoteCurrency: req.QuoteCurrency,
    }

    market, err := server.store.CreateMarket(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"id": market.ID})
}

func (server *server) getMarket(ctx *gin.Context) {
    id := ctx.Param("id")
    marketID, err := uuid.Parse(id)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    market, err := server.store.GetMarketByID(ctx, marketID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, market)
}

func (server *server) listMarkets(ctx *gin.Context) {
    markets, err := server.store.ListMarkets(ctx)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, markets)
}