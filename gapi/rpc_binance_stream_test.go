package gapi

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type mockTradeStream struct {
	ctx     context.Context
	sent    []*pb.Trade
	sendErr error
}

func (m *mockTradeStream) Context() context.Context { return m.ctx }
func (m *mockTradeStream) Send(trade *pb.Trade) error {
	if m.sendErr != nil {
		return m.sendErr
	}
	m.sent = append(m.sent, trade)
	return nil
}

func (m *mockTradeStream) SetHeader(md metadata.MD) error { return nil }
func (m *mockTradeStream) SendHeader(md metadata.MD) error { return nil }
func (m *mockTradeStream) SetTrailer(md metadata.MD)       {}
func (m *mockTradeStream) RecvMsg(_ interface{}) error    { return nil }
func (m *mockTradeStream) SendMsg(_ interface{}) error    { return nil }


type fakeWSConn struct {
    messages chan []byte
	readErr  error
}

func (f *fakeWSConn) ReadMessage() (int, []byte, error) {
	if f.readErr != nil {
		return 0, nil, f.readErr
	}
    msg, ok := <-f.messages
    if !ok {
		return 0, nil, errors.New("closed")
}
	return 1, msg, nil 
}
func (f *fakeWSConn) Close() error { return nil }

func TestStreamTrades(t *testing.T) {
	event := binanceTradeEvent{
		EventType:    "trade",
		EventTime:    time.Now().Unix(),
		Symbol:       "BTCUSDT",
		TradeID:      42,
		Price:        "50000.50",
		Quantity:     "0.01",
		TradeTime:    time.Now().Unix(),
		IsBuyerMaker: true,
	}
	data, _ := json.Marshal(event)
	wrapped, _ := json.Marshal(map[string]interface{}{
		"stream": "btcusdt@trade",
		"data":   json.RawMessage(data),
	})

	testCases := []struct {
		name          string
		req           *pb.TradeStreamRequest
		buildStub     func() wsConn
		checkResponse func(t *testing.T, trades []*pb.Trade, err error)
	}{
		{
			name: "OK",
			req:  &pb.TradeStreamRequest{Symbols: []string{"BTCUSDT"}},
			buildStub: func() wsConn {
				ch := make(chan []byte, 1)
				ch <- wrapped
				close(ch)
				return &fakeWSConn{messages: ch}
			},
			checkResponse: func(t *testing.T, trades []*pb.Trade, err error) {
				require.NoError(t, err)
				require.Len(t, trades, 1)
				require.Equal(t, "BTCUSDT", trades[0].Symbol)
				require.Equal(t, float64(50000.50), trades[0].Price)
			},
		},
		{
			name: "NoSymbols",
			req:  &pb.TradeStreamRequest{Symbols: []string{}},
			buildStub: func() wsConn {
				return &fakeWSConn{readErr: errors.New("should not be called")}
			},
			checkResponse: func(t *testing.T, trades []*pb.Trade, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "WebSocketError",
			req:  &pb.TradeStreamRequest{Symbols: []string{"BTCUSDT"}},
			buildStub: func() wsConn {
				return &fakeWSConn{readErr: errors.New("websocket failure")}
			},
			checkResponse: func(t *testing.T, trades []*pb.Trade, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			stream := &mockTradeStream{ctx: ctx}

			go func() {
    			time.Sleep(10 * time.Millisecond)
    			cancel()
			}()

			oldDial := websocketDial
			websocketDial = func(_ string) (wsConn, error) {
				return tc.buildStub(), nil
			}
			defer func() { websocketDial = oldDial }()

			srv := server{}
			err := srv.StreamTrades(tc.req, stream)
			tc.checkResponse(t, stream.sent, err)
		})
	}
}