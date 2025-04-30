package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	token "github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
)

type server struct {
	store       db.Store_interface
	router 	   *gin.Engine
	tokenMaker  token.Maker
	config 	    utils.Config
}

const url = "wss://stream.binance.com:9443/ws/btcusdt@trade/ethusdt@trade/eurusdt@trade/jpyusdt@trade"



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


	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	router.POST("/token/renew_token", server.renewAccessToken)

	router.POST("/users/login", server.loginUser)
	router.POST("/users", server.createUser)
	router.GET("/users/:id", server.getUser)
	router.PUT("/users/:id", server.updateUser)
	router.DELETE("/users/:id", server.deleteUser)

	authRoutes.POST("/transactions", server.createTransaction)
	authRoutes.GET("/transactions/:id", server.getTransaction)
	router.GET("/transactions/user/:user_email", server.listUserTransactions)
	authRoutes.DELETE("/transactions/:id", server.deleteTransaction)

	authRoutes.POST("/markets", server.createMarket)
	authRoutes.GET("/markets/:id", server.getMarket)
	authRoutes.DELETE("/markets/:id", server.deleteMarket)
	router.GET("/markets", server.listMarkets)

	authRoutes.POST("/fees", server.createFee)
	authRoutes.GET("/fees/:market_id", server.getFee)
	authRoutes.DELETE("/fees/:id", server.deleteFee)

	authRoutes.POST("/orders", server.createOrder)
	authRoutes.GET("/orders/:id", server.getOrder)
	authRoutes.DELETE("/orders/:id", server.deleteOrder)

	authRoutes.POST("/trades", server.createTrade)
	authRoutes.GET("/trades/:id", server.getTrade)
	router.GET("/trades/market/:market_id", server.listTrades)
	authRoutes.DELETE("/trades/:id", server.deleteTrade)

	authRoutes.POST("/wallets", server.createWallet)
	authRoutes.GET("/wallets/:id", server.getWallet)
	authRoutes.PUT("/wallets/:id", server.updateWallet)
	authRoutes.DELETE("/wallets/:id", server.deleteWallet)

	authRoutes.POST("/audit-logs", server.createAuditLog)
	router.GET("/audit-logs/user/:user_email", server.listUserAuditLogs)
	authRoutes.GET("/audit-logs/:user_email", server.getAuditLog)
	authRoutes.DELETE("/audit-logs/:id", server.DeleteAuditLog)

	router.GET("/ws", func(ctx *gin.Context) {
			WebSocket(ctx, url)
		})



	server.router = router
	return server, nil
}

func (server *server) Start (address string) error {
    return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
