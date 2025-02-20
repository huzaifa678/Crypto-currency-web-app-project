package api

import (
    db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type FeeRequest struct {
    MarketID uuid.UUID `json:"market_id"`
    MakerFee string `json:"maker_fee"`
    TakerFee string `json:"taker_fee"`
}

func (server *server) createFee(ctx *gin.Context) {
    var req FeeRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    arg := db.CreateFeeParams{
        MarketID: req.MarketID,
        MakerFee: sql.NullString{String: req.MakerFee, Valid: req.MakerFee != ""},
        TakerFee: sql.NullString{String: req.TakerFee, Valid: req.TakerFee != ""},
    }

    fee, err := server.store.CreateFee(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"market_id": fee.MarketID})
}

func (server *server) getFee(ctx *gin.Context) {
    id := ctx.Param("market_id")
    feeID, err := uuid.Parse(id)


    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    if feeID == uuid.Nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
    }

    fee, err := server.store.GetFeeByMarketID(ctx, feeID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, fee)
}

func (server *server) deleteFee(ctx *gin.Context) {
    id := ctx.Param("id")
    feeID, err := uuid.Parse(id)

    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    if feeID == uuid.Nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    err = server.store.DeleteFee(ctx, feeID)
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