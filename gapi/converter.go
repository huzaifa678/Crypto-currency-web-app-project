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
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Role:      convertUserRole(user.Role),
		IsVerified: user.IsVerified,
	}
}


func convertGetByEmailUser(user db.GetUserByEmailRow) *pb.User {
	return &pb.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Role:      convertUserRole(user.Role),
		IsVerified: user.IsVerified,
	}
}

func convertGetUser(user db.GetUserByIDRow) *pb.User {
	return &pb.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		Role:      convertUserRole(user.Role),
		IsVerified: user.IsVerified,
	}
}

func convertMarket(market db.Market) *pb.Market {
	return &pb.Market{
		MarketId:     market.ID.String(),
		UserName:     market.Username,
		BaseCurrency: market.BaseCurrency,
		QuoteCurrency: market.QuoteCurrency,
		MinOrderAmount: market.MinOrderAmount,
		PricePrecision: int32(market.PricePrecision),
		CreatedAt:     timestamppb.New(market.CreatedAt),
	}
}

func convertListMarkets(markets []db.ListMarketsRow) *pb.MarketListResponse {
	listMarkets := make([]*pb.ListMarket, len(markets))

	for i, market := range markets {
		listMarkets[i] = &pb.ListMarket{
			MarketId:      market.ID.String(),
			BaseCurrency:  market.BaseCurrency,
			QuoteCurrency: market.QuoteCurrency,
			MinOrderAmount: market.MinOrderAmount,
			PricePrecision: int32(market.PricePrecision),
			CreatedAt:      timestamppb.New(market.CreatedAt),
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
		Price:        order.Price,
		Amount:       order.Amount,
		FilledAmount: order.FilledAmount,
		CreatedAt:    timestamppb.New(order.CreatedAt),
		UpdatedAt:    timestamppb.New(order.UpdatedAt),
	}
}

func convertWallet(wallet db.Wallet) *pb.Wallet {
	return &pb.Wallet{
		Id:            wallet.ID.String(),
		UserEmail:     wallet.UserEmail,
		Currency:      wallet.Currency,
		Balance:       wallet.Balance,
		LockedBalance: wallet.LockedBalance,
		CreatedAt:     timestamppb.New(wallet.CreatedAt),
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

func convertTransactionType(TransactionType db.TransactionType) pb.TransactionType {
	switch TransactionType {
		case db.TransactionTypeDeposit:
			return pb.TransactionType_DEPOSIT
		case db.TransactionTypeWithdrawal:
			return pb.TransactionType_WITHDRAWAL
		default:
			return pb.TransactionType_NONE
	}
}

func convertTransactionStatus(TransactionStatus db.TransactionStatus) pb.TransactionStatus {
	switch TransactionStatus {
		case db.TransactionStatusPending:
			return pb.TransactionStatus_PENDING
		case db.TransactionStatusCompleted:
			return pb.TransactionStatus_COMPLETED
		case db.TransactionStatusFailed:
			return pb.TransactionStatus_FAILED
		default:
			return pb.TransactionStatus_PENDING
	}
}

func convertTrade(trade db.Trade) (*pb.Trades) {
	return &pb.Trades{
		TradeId:     trade.ID.String(),
		Username:   trade.Username,
		BuyOrderId: trade.BuyOrderID.String(),
		SellOrderId: trade.SellOrderID.String(),
		MarketId:   trade.MarketID.String(),
		Price:      trade.Price,
		Amount:     trade.Amount,
		Fee:        trade.Fee,
		CreatedAt:  timestamppb.New(trade.CreatedAt), 
	}
}

func convertTransaction(transaction db.Transaction) (*pb.Transaction) {
	return &pb.Transaction{
		TransactionId: transaction.ID.String(),
		Username: transaction.Username,
		UserEmail: transaction.UserEmail,
		Type: convertTransactionType(transaction.Type),
		Currency: transaction.Currency,
		Amount: transaction.Amount,
		Status: convertTransactionStatus(transaction.Status),
		Address: transaction.Address,
		TxHash: transaction.TxHash,
		CreatedAt: timestamppb.New(transaction.CreatedAt),
	}
}

func convertCreateTransaction(userName string, transaction db.CreateTransactionRow) (*pb.Transaction) {
	return &pb.Transaction{
		TransactionId: transaction.ID.String(),
		Username: userName,
		UserEmail: transaction.UserEmail,
		Type: convertTransactionType(transaction.Type),
		Currency: transaction.Currency,
		Amount: transaction.Amount,
		Status: convertTransactionStatus(transaction.Status),
		Address: transaction.Address,
		TxHash: transaction.TxHash,
		CreatedAt: timestamppb.New(transaction.CreatedAt),
	}
}

