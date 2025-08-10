package api

import (
	"database/sql"
	"log"
	"net/http"

	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderRequest struct {
    UserEmail string          `json:"user_email"`
    MarketID  uuid.UUID       `json:"market_id"`
    Type      db.OrderType    `json:"type"`
    Status    db.OrderStatus  `json:"status"`
    Price     string          `json:"price"`
    Amount    string          `json:"amount"`
}

func (server *server) createOrder(ctx *gin.Context) {
    var req OrderRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        log.Println("JSON Binding Error:", err)
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

    arg := db.CreateOrderParams{
        Username:  authPayload.Username,
        UserEmail: req.UserEmail,
        MarketID:  req.MarketID,
        Type:      req.Type,
        Status:    req.Status,
        Price:     req.Price,
        Amount:    req.Amount,
    }

    order, err := server.store.CreateOrder(ctx, arg)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"id": order.ID})
}

func (server *server) getOrder(ctx *gin.Context) {
    id := ctx.Param("id")
    orderID, err := uuid.Parse(id)

    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    if orderID == uuid.Nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
        return
    }

    order, err := server.store.GetOrderByID(ctx, orderID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

    if authPayload.Username != order.Username {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
        return
    }

    ctx.JSON(http.StatusOK, order)
}

func (server *server) deleteOrder(ctx *gin.Context) {
    id := ctx.Param("id")
    orderID, err := uuid.Parse(id)

    if err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    if orderID == uuid.Nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
        return
    }

    order, err := server.store.GetOrderByID(ctx, orderID)

    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
        return
    }

    authPayload := ctx.MustGet(AuthorizationPayloadKey).(*token.Payload)

    if authPayload.Username != order.Username {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
        return
    }

    err = server.store.DeleteOrder(ctx, orderID)
    if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

    ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}