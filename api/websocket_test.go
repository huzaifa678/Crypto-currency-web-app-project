package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

var testupgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func TestWebSocket(t *testing.T) {
	mockBinance := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := testupgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer ws.Close()

		for {
			messageType, msg, err := ws.ReadMessage()
			if err != nil {
				break 
			}
			err = ws.WriteMessage(messageType, msg)
			if err != nil {
				break 
			}
		}
	}))
	defer mockBinance.Close()

	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		WebSocket(c, "ws"+strings.TrimPrefix(mockBinance.URL, "http"))
	})

	server := httptest.NewServer(router)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(mockBinance.URL, "http")

	fmt.Println(wsURL)

	testCases := []struct {
		name          string
		setupMockData func() []byte
		checkResponse func(t *testing.T, msg []byte)
		expectError   bool
	}{
		{
			name: "ValidTradeData",
			setupMockData: func() []byte {
				mockTrade := TradeData{
					Symbol:   "BTCUSDT",
					Price:    "50000.00",
					Quantity: "0.5",
					Time:     time.Now().Unix(),
				}
				data, _ := json.Marshal(mockTrade)
				return data
			},
			checkResponse: func(t *testing.T, msg []byte) {
				var trade TradeData
				err := json.Unmarshal(msg, &trade)
				require.NoError(t, err)
				require.Equal(t, "BTCUSDT", trade.Symbol)
				require.Equal(t, "50000.00", trade.Price)
				require.Equal(t, "0.5", trade.Quantity)
			},
			expectError: false,
		},
		{
			name: "InvalidTradeData",
			setupMockData: func() []byte {
				
				return []byte(`{"symbol": "BTCUSDT", "price": "50000.00", "quantity": "0.5", "Time": "invalid"}`)
			},
			checkResponse: func(t *testing.T, msg []byte) {
				
			},
			expectError: true,
		},
		{
			name: "EmptyTradeData",
			setupMockData: func() []byte {
				
				return []byte(`{}`)
			},
			checkResponse: func(t *testing.T, msg []byte) {
				var trade TradeData
				err := json.Unmarshal(msg, &trade)
				require.NoError(t, err)
				require.Empty(t, trade.Symbol)
				require.Empty(t, trade.Price)
				require.Empty(t, trade.Quantity)
				require.Zero(t, trade.Time)
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			require.NoError(t, err)
			defer ws.Close()

			
			mockData := tc.setupMockData()
			err = ws.WriteMessage(websocket.TextMessage, mockData)
			require.NoError(t, err)

			
			var wg sync.WaitGroup
			wg.Add(1)

			
			done := make(chan struct{})
			var response []byte
			go func() {
				defer wg.Done()
				defer close(done)
				_, msg, err := ws.ReadMessage()
				if err != nil {
					if tc.expectError {
						return
					}
					t.Logf("read error: %v", err)
					return
				}
				response = msg
			}()

			select {
			case <-done:
				if response != nil {
					tc.checkResponse(t, response)
				}
			case <-time.After(5 * time.Second):
				if !tc.expectError {
					t.Fatal("test timed out")
				}
			}

			
			wg.Wait()
		})
	}
}