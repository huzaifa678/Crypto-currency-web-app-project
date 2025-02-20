package api

import (
    db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type MarketRequest struct {
    BaseCurrency  string `json:"base_currency"`
    QuoteCurrency string `json:"quote_currency"`
    MinOrderAmount sql.NullString `json:"min_order_amount"`
    PricePrecision sql.NullInt32 `json:"price_precision"`
}

func (server *server) createMarket(ctx *gin.Context) {
    var req MarketRequest

    var err error

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    arg := db.CreateMarketParams{
        BaseCurrency:  req.BaseCurrency,
        QuoteCurrency: req.QuoteCurrency,
        MinOrderAmount: req.MinOrderAmount,
        PricePrecision: req.PricePrecision,
    }

    if req.BaseCurrency == "" || req.QuoteCurrency == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "no currency added"})
        return 
    }

    currencies := []string{"USD", "EUR", "BTC", "ETH", "JPY"}

    if !isValidCurrency(req.BaseCurrency, currencies) || !isValidCurrency(req.QuoteCurrency, currencies) {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid currency"})
        return
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

func (server *server) deleteMarket(ctx *gin.Context) {
    id := ctx.Param("id")
    marketID, err := uuid.Parse(id)

    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    if id == uuid.Nil.String() {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
        return
    }

    err = server.store.DeleteMarket(ctx, marketID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (server *server) listMarkets(ctx *gin.Context) {
    markets, err := server.store.ListMarkets(ctx)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, markets)
}

func isValidCurrency(currency string, currencies []string) bool {
    for _, c := range currencies {
        if currency == c {
            return true
        }
    }
    return false
}
