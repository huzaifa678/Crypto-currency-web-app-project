package api

import (
    db "crypto-system/db/sqlc"
    "database/sql"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type orderRequest struct {
    UserEmail string          `json:"user_email"`
    MarketID  uuid.UUID       `json:"market_id"`
    Type      db.OrderType    `json:"type"`
    Status    db.OrderStatus  `json:"status"`
    Price     sql.NullString  `json:"price"`
    Amount    string          `json:"amount"`
}

func (server *server) createOrder(ctx *gin.Context) {
    var req orderRequest

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, errorResponse(err))
        return
    }

    arg := db.CreateOrderParams{
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

    order, err := server.store.GetOrderByID(ctx, orderID)
    if err != nil {
        if err == sql.ErrNoRows {
            ctx.JSON(http.StatusNotFound, errorResponse(err))
            return
        }
        ctx.JSON(http.StatusInternalServerError, errorResponse(err))
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

    err = server.store.DeleteOrder(ctx, orderID)
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