package api

import (
	db "crypto-system/db/sqlc"

	"github.com/gin-gonic/gin"
)

type server struct {
	store *db.Store
	router *gin.Engine
}

func NewServer (Store *db.Store) *server {
	server := &server{store: Store}
	router := gin.Default()

	router.POST("/users", server.createUser)
    router.GET("/users/:id", server.getUser)
	router.PUT("/users/:id", server.updateUser)
    router.DELETE("/users/:id", server.deleteUser)

    router.POST("/transactions", server.createTransaction)
    router.GET("/transactions/:id", server.getTransaction)
	router.GET("/transactions/user/:email", server.listUserTransactions)

    router.POST("/markets", server.createMarket)
    router.GET("/markets/:id", server.getMarket)
	router.GET("/markets", server.listMarkets)

    router.POST("/fees", server.createFee)
    router.GET("/fees/:id", server.getFee)
	router.DELETE("/fees/:id", server.deleteFee)
	router.GET("/fees/market/:marketID", server.getFee)

    router.POST("/orders", server.createOrder)
    router.GET("/orders/:id", server.getOrder)

    router.POST("/trades", server.createTrade)
    router.GET("/trades/:id", server.getTrade)
	router.DELETE("/trades/:id", server.deleteTrade)

    router.POST("/wallets", server.createWallet)
    router.GET("/wallets/:id", server.getWallet)
	router.PUT("/wallets/:id", server.updateWallet)
	router.DELETE("/wallets/:id", server.deleteWallet)

	router.GET("/audit-logs/user/:email", server.listUserAuditLogs)
	router.GET("/audit-logs/:id", server.getAuditLog)

	server.router = router
	return server
}


func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}