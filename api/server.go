package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
)

type server struct {
	store      db.Store_interface
	router     *gin.Engine
	tokenMaker token.Maker
	config     utils.Config
}


func (server *server) Start(serverAddr string) error {
	return server.router.Run(serverAddr)
}

func NewServer(store db.Store_interface, config utils.Config) (*server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.PasetoSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	router := gin.Default()

	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.PUT("/users/:id", server.updateUser)
	router.DELETE("/users/:id", server.deleteUser)

	router.POST("/transactions", server.createTransaction)
	router.GET("/transactions/:id", server.getTransaction)
	router.GET("/transactions/user/:user_email", server.listUserTransactions)
	router.DELETE("/transactions/:id", server.deleteTransaction)

	router.POST("/markets", server.createMarket)
	router.GET("/markets/:id", server.getMarket)
	router.DELETE("/markets/:id", server.deleteMarket)
	router.GET("/markets", server.listMarkets)

	router.POST("/fees", server.createFee)
	router.GET("/fees/:market_id", server.getFee)
	router.DELETE("/fees/:id", server.deleteFee)

	router.POST("/orders", server.createOrder)
	router.GET("/orders/:id", server.getOrder)
	router.DELETE("/orders/:id", server.deleteOrder)

	router.POST("/trades", server.createTrade)
	router.GET("/trades/:id", server.getTrade)
	router.GET("/trades/market/:market_id", server.listTrades)
	router.DELETE("/trades/:id", server.deleteTrade)

	router.POST("/wallets", server.createWallet)
	router.GET("/wallets/:id", server.getWallet)
	router.PUT("/wallets/:id", server.updateWallet)
	router.DELETE("/wallets/:id", server.deleteWallet)

	router.POST("/audit-logs", server.createAuditLog)
	router.GET("/audit-logs/user/:user_email", server.listUserAuditLogs)
	router.GET("/audit-logs/:user_email", server.getAuditLog)
	router.DELETE("/audit-logs/:id", server.DeleteAuditLog)

	server.router = router
	return server, nil
}

func (server *server) start (address string) error {
    return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
