package gapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/huzaifa678/Crypto-currency-web-app-project/pb"
)

type binanceTradeEvent struct {
    EventType     string  `json:"e"`
    EventTime     int64   `json:"E"`
    Symbol        string  `json:"s"`
    TradeID       int64   `json:"t"`
    Price         string  `json:"p"`
    Quantity      string  `json:"q"`
    BuyerOrderID  int64   `json:"b"`
    SellerOrderID int64   `json:"a"`
    TradeTime     int64   `json:"T"`
    IsBuyerMaker  bool    `json:"m"`
    Ignore        bool    `json:"M"`
}

func parseStringToFloat64(str string) float64 {
    f, err := strconv.ParseFloat(str, 64)
    if err != nil {
        return 0
    }
    return f
}

func (server *server) StreamTrades(req *pb.TradeStreamRequest, stream pb.CryptoWebApp_StreamTradesServer) error {
    symbols := req.GetSymbols()
    if len(symbols) == 0 {
        return fmt.Errorf("no symbols provided")
    }

    streams := make([]string, len(symbols))
    for i, sym := range symbols {
        streams[i] = strings.ToLower(sym) + "@trade"
    }

    log.Printf("symbols: %v", symbols)
    log.Printf("streams: %v", streams)
    url := "wss://stream.binance.com:9443/stream?streams=" + strings.Join(streams, "/")

    log.Printf("Connecting to Binance WS: %s\n", url)

    c, _, err := websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        return fmt.Errorf("failed to connect to binance websocket: %v", err)
    }
    defer c.Close()

    for {
        select {
        case <-stream.Context().Done(): 
            log.Println("stream context canceled, closing trade stream")
            return nil
        default:
            _, message, err := c.ReadMessage()
            if err != nil {
                if stream.Context().Err() == context.Canceled {
                    log.Println("client canceled stream")
                    return nil
                }
                return fmt.Errorf("error reading from websocket: %v", err)
            }

            var wrapped struct {
                Stream string          `json:"stream"`
                Data   json.RawMessage `json:"data"`
            }
            if err := json.Unmarshal(message, &wrapped); err != nil {
                log.Printf("failed to unmarshal wrapped message: %v", err)
                continue
            }

            var event binanceTradeEvent
            if err := json.Unmarshal(wrapped.Data, &event); err != nil {
                log.Printf("failed to unmarshal trade event: %v", err)
                continue
            }

            trade := &pb.Trade{
                Symbol:       event.Symbol,
                Price:        parseStringToFloat64(event.Price),
                Quantity:     parseStringToFloat64(event.Quantity),
                TradeId:      event.TradeID,
                EventTime:    event.EventTime,
                TradeTime:    event.TradeTime,
                IsBuyerMaker: event.IsBuyerMaker,
                RawJson:      string(wrapped.Data),
            }

            if err := stream.Send(trade); err != nil {
                if stream.Context().Err() == context.Canceled {
                    log.Println("client canceled stream")
                    return nil
                }
                log.Printf("error sending trade: %v", err)
                return err
            }
        }
    }
}
