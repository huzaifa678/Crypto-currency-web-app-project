package gapi

import (
	db "github.com/huzaifa678/Crypto-currency-web-app-project/db/sqlc"
	pb "github.com/huzaifa678/Crypto-currency-web-app-project/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertCreateUser(user db.CreateUserRow) *pb.User {
	return &pb.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
		UpdatedAt: timestamppb.New(user.UpdatedAt.Time),
		Role:      convertUserRole(user.Role),
		IsVerified: user.IsVerified.Bool,
	}
}

func convertGetByEmailUser(user db.GetUserByEmailRow) *pb.User {
	return &pb.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
		UpdatedAt: timestamppb.New(user.UpdatedAt.Time),
		Role:      convertUserRole(user.Role),
		IsVerified: user.IsVerified.Bool,
	}
}

func convertGetUser(user db.GetUserByIDRow) *pb.User {
	return &pb.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
		UpdatedAt: timestamppb.New(user.UpdatedAt.Time),
		Role:      convertUserRole(user.Role),
		IsVerified: user.IsVerified.Bool,
	}
}

func convertMarket(market db.Market) *pb.Market {
	return &pb.Market{
		MarketId:     market.ID.String(),
		UserName:     market.Username,
		BaseCurrency: market.BaseCurrency,
		QuoteCurrency: market.QuoteCurrency,
		MinOrderAmount: market.MinOrderAmount.String,
		PricePrecision: int32(market.PricePrecision.Int32),
		CreatedAt:     timestamppb.New(market.CreatedAt.Time),
	}
}

func convertListMarkets(markets []db.ListMarketsRow) *pb.MarketListResponse {
	listMarkets := make([]*pb.ListMarket, len(markets))

	for i, market := range markets {
		listMarkets[i] = &pb.ListMarket{
			MarketId:      market.ID.String(),
			BaseCurrency:  market.BaseCurrency,
			QuoteCurrency: market.QuoteCurrency,
			MinOrderAmount: market.MinOrderAmount.String,
			PricePrecision: int32(market.PricePrecision.Int32),
			CreatedAt:      timestamppb.New(market.CreatedAt.Time),
		}
	}

	return &pb.MarketListResponse{
		Markets: listMarkets,
	}
}

func convertOrder(order db.Order) *pb.Order {
	return &pb.Order{
		Id:           order.ID.String(),
		UserName:     order.Username,
		UserEmail:    order.UserEmail,
		MarketId:     order.MarketID.String(),
		Type:         convertOrderType(order.Type),
		Status:       convertOrderStatus(order.Status),
		Price:        order.Price.String,
		Amount:       order.Amount,
		FilledAmount: order.FilledAmount.String,
		CreatedAt:    timestamppb.New(order.CreatedAt.Time),
		UpdatedAt:    timestamppb.New(order.UpdatedAt.Time),
	}
}

func convertWallet(wallet db.Wallet) *pb.Wallet {
	return &pb.Wallet{
		Id:            wallet.ID.String(),
		UserEmail:     wallet.UserEmail,
		Currency:      wallet.Currency,
		Balance:       wallet.Balance.String,
		LockedBalance: wallet.LockedBalance.String,
		CreatedAt:     timestamppb.New(wallet.CreatedAt.Time),
	}
}

func convertUserRole(role db.UserRole) pb.UserRole {
    switch role {
    	case db.UserRoleAdmin:
        	return pb.UserRole_USER_ROLE_ADMIN
   		case db.UserRoleUser:
        	return pb.UserRole_USER_ROLE_USER
    	default:
        	return pb.UserRole_USER_ROLE_USER 
    }
}

func convertOrderType(orderType db.OrderType) pb.OrderType {
	switch orderType {
		case db.OrderTypeBuy:
			return pb.OrderType_BUY
		case db.OrderTypeSell:
			return pb.OrderType_SELL
		default:
			return pb.OrderType_BUY
	}
}

func convertOrderStatus(status db.OrderStatus) pb.Status {
	switch status {
	case db.OrderStatusOpen:
		return pb.Status_OPEN
	case db.OrderStatusPartiallyFilled:
		return pb.Status_PARTIALLY_FILLED
	case db.OrderStatusFilled:
		return pb.Status_PARTIALLY_FILLED
	case db.OrderStatusCancelled:
		return pb.Status_CANCELLED
	default:
		return pb.Status_OPEN
	}
}