package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	mockdb "github.com/huzaifa678/Crypto-currency-web-app-project/db/mock"
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"github.com/huzaifa678/Crypto-currency-web-app-project/token"
	"github.com/huzaifa678/Crypto-currency-web-app-project/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateOrderRPC(t *testing.T) {
	createOrderParams, order, _, createOrderRow := createRandomOrder()

	testCases := []struct {
		name          string
		req           *pb.CreateOrderRequest
		buildStubs    func(store *mockdb.MockStore_interface)
		setupAuth     func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.CreateOrderResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateOrderRequest{
				UserEmail: order.UserEmail,
				MarketId:  order.MarketID.String(),
				Type:      pb.OrderType_BUY,
				Status:    pb.Status_OPEN,
				Price:     order.Price.String,
				Amount:    order.Amount,
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, order.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				arg := createOrderParams

				store.EXPECT().
					CreateOrder(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(createOrderRow, nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateOrderResponse, err error) {
				log.Println("ERROR: ", err)
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.OrderId)
			},
		},
		{
			name: "Unauthorized",
			req: &pb.CreateOrderRequest{
				UserEmail: order.UserEmail,
				MarketId:  order.MarketID.String(),
				Type:      pb.OrderType_BUY,
				Status:    pb.Status_OPEN,
				Price:     "100.50",
				Amount:    "1.5",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
				return context.Background()
			},
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateOrderResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidEmail",
			req: &pb.CreateOrderRequest{
				UserEmail: "invalid-email",
				MarketId:  order.ID.String(),
				Type:      pb.OrderType_BUY,
				Status:    pb.Status_OPEN,
				Price:     "100.50",
				Amount:    "1.5",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, order.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateOrderResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InvalidMarketID",
			req: &pb.CreateOrderRequest{
				UserEmail: order.UserEmail,
				MarketId:  "invalid-uuid",
				Type:      pb.OrderType_BUY,
				Status:    pb.Status_OPEN,
				Price:     "100.50",
				Amount:    "1.5",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, order.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateOrderResponse, err error) {
				log.Println("ERROR: ", err)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "InternalError",
			req: &pb.CreateOrderRequest{
				UserEmail: order.UserEmail,
				MarketId:  order.MarketID.String(),
				Type:      pb.OrderType_BUY,
				Status:    pb.Status_OPEN,
				Price:     "100.50",
				Amount:    "1.5",
			},
			setupAuth: func(t *testing.T, tokenMaker token.Maker) context.Context{
                return newContextWithBearerToken(t, tokenMaker, order.Username, time.Minute, token.TokenTypeAccessToken)
            },
			buildStubs: func(store *mockdb.MockStore_interface) {
				store.EXPECT().
					CreateOrder(gomock.Any(), gomock.Any()).
					Times(1).
					Return(createOrderRow, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, res *pb.CreateOrderResponse, err error) {
				log.Println("ERROR: ", err)
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore_interface(ctrl)
			tc.buildStubs(store)

			server := NewTestServer(t, store)
			ctx := tc.setupAuth(t, server.tokenMaker)

			res, err := server.CreateOrder(ctx, tc.req)
			tc.checkResponse(t, res, err)
		})
	}
} 


func createRandomOrder() (createOrderParams db.CreateOrderParams, order db.Order, updatedOrderParams db.UpdateOrderStatusAndFilledAmountParams, createOrderRow db.CreateOrderRow) {
	username := utils.RandomUser()
	email := "hello" + fmt.Sprint(rand.Intn(10000)) + "@example.com"
	marketID := uuid.New()
	orderType := db.OrderType(pb.OrderType_BUY) 
	orderStatus := db.OrderStatus(pb.Status_OPEN) 
	price := sql.NullString{String: "100.50", Valid: true}
	amount := "10"

	createOrderParams = db.CreateOrderParams{
		Username: username,
		UserEmail: email,
		MarketID:  marketID,
		Type:      orderType,
		Status:    orderStatus,
		Price:     price,
		Amount:    amount,
	}


	createdOrder := db.Order{
		ID:           uuid.New(),
		Username: 	  username,
		UserEmail:    email,
		MarketID:     marketID,
		Type:         orderType,
		Status:       orderStatus,
		Price:        price,
		Amount:       amount,
		FilledAmount: sql.NullString{String: "5", Valid: true}, 
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	updatedOrderParams = db.UpdateOrderStatusAndFilledAmountParams{
		Status:       db.OrderStatus(fmt.Sprint(1)), 
		FilledAmount: sql.NullString{String: "10", Valid: true}, 
		ID:           createdOrder.ID,
	}

	createOrderRow = db.CreateOrderRow {
		ID:           createdOrder.ID,
		UserEmail:    email,
		MarketID:     marketID,
		Type:         orderType,
		Status:       orderStatus,
		Price:        price,
		Amount:       amount,
		FilledAmount: sql.NullString{String: "5", Valid: true}, 
		CreatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	return createOrderParams, createdOrder, updatedOrderParams, createOrderRow
}