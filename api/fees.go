package api

import (
	"database/sql"
	"net/http"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FeeRequest struct {
    MarketID uuid.UUID `json:"market_id"`
    MakerFee decimal.Decimal `json:"maker_fee"`
    TakerFee decimal.Decimal `json:"taker_fee"`
}

func (server *server) createFee(ctx *gin.Context) {
    var req FeeRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }


    authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

    arg := db.CreateFeeParams{
        Username: authPayload.Username,
        MarketID: req.MarketID,
        MakerFee: req.MakerFee,
        TakerFee: req.TakerFee,
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

    fee, err := server.store.GetFeeByMarketID(ctx, feeID)

    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

    if authPayload.Username != fee.Username {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
        return
    }

    err = server.store.DeleteFee(ctx, feeID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}