package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type TradeData struct {
	Symbol      string `json:"symbol"` 
	Price       string `json:"price"` 
	Quantity    string `json:"quantity"` 
	Time        int64  `json:"Time"` 
}


var upgrader = websocket.Upgrader{
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

func WebSocket(ctx *gin.Context, url string) {
    // Upgrade HTTP connection to WebSocket
    conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
    if err != nil {
        fmt.Println("WebSocket upgrade failed:", err)
        return
    }
    defer conn.Close()

    // Connect to Binance WebSocket
    binanceConn, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        log.Fatal("Binance WebSocket connection failed:", err)
        return
    }
    defer binanceConn.Close()

    fmt.Println("Connected to Binance WebSocket...")

    // Forward messages from Binance to client
    for {
        _, message, err := binanceConn.ReadMessage()
        if err != nil {
            log.Println("Binance read error:", err)
            break
        }

        var trade TradeData
        err = json.Unmarshal(message, &trade)
        if err != nil {
            log.Println("JSON Unmarshal error:", err)
            continue
        }

        // Send formatted trade data to client
        err = conn.WriteJSON(trade)
        if err != nil {
            log.Println("Client write error:", err)
            break
        }
    }
}