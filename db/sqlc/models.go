// Code generated by sqlc. DO NOT EDIT.

package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusOpen            OrderStatus = "open"
	OrderStatusPartiallyFilled OrderStatus = "partially_filled"
	OrderStatusFilled          OrderStatus = "filled"
	OrderStatusCancelled       OrderStatus = "cancelled"
)

func (e *OrderStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderStatus(s)
	case string:
		*e = OrderStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderStatus: %T", src)
	}
	return nil
}

type OrderType string

const (
	OrderTypeBuy  OrderType = "buy"
	OrderTypeSell OrderType = "sell"
)

func (e *OrderType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OrderType(s)
	case string:
		*e = OrderType(s)
	default:
		return fmt.Errorf("unsupported scan type for OrderType: %T", src)
	}
	return nil
}

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
)

func (e *TransactionStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TransactionStatus(s)
	case string:
		*e = TransactionStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for TransactionStatus: %T", src)
	}
	return nil
}

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
)

func (e *TransactionType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TransactionType(s)
	case string:
		*e = TransactionType(s)
	default:
		return fmt.Errorf("unsupported scan type for TransactionType: %T", src)
	}
	return nil
}

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

func (e *UserRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserRole(s)
	case string:
		*e = UserRole(s)
	default:
		return fmt.Errorf("unsupported scan type for UserRole: %T", src)
	}
	return nil
}

type AuditLog struct {
	ID        uuid.UUID      `json:"id"`
	Username  string         `json:"username"`
	UserEmail string         `json:"user_email"`
	Action    string         `json:"action"`
	IpAddress sql.NullString `json:"ip_address"`
	CreatedAt sql.NullTime   `json:"created_at"`
}

type Fee struct {
	ID        uuid.UUID      `json:"id"`
	Username  string         `json:"username"`
	MarketID  uuid.UUID      `json:"market_id"`
	MakerFee  sql.NullString `json:"maker_fee"`
	TakerFee  sql.NullString `json:"taker_fee"`
	CreatedAt sql.NullTime   `json:"created_at"`
}

type Market struct {
	ID             uuid.UUID      `json:"id"`
	Username       string         `json:"username"`
	BaseCurrency   string         `json:"base_currency"`
	QuoteCurrency  string         `json:"quote_currency"`
	MinOrderAmount sql.NullString `json:"min_order_amount"`
	PricePrecision sql.NullInt32  `json:"price_precision"`
	CreatedAt      sql.NullTime   `json:"created_at"`
}

type Order struct {
	ID           uuid.UUID      `json:"id"`
	Username     string         `json:"username"`
	UserEmail    string         `json:"user_email"`
	MarketID     uuid.UUID      `json:"market_id"`
	Type         OrderType      `json:"type"`
	Status       OrderStatus    `json:"status"`
	Price        sql.NullString `json:"price"`
	Amount       string         `json:"amount"`
	FilledAmount sql.NullString `json:"filled_amount"`
	CreatedAt    sql.NullTime   `json:"created_at"`
	UpdatedAt    sql.NullTime   `json:"updated_at"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type Trade struct {
	ID          uuid.UUID      `json:"id"`
	Username    string         `json:"username"`
	BuyOrderID  uuid.UUID      `json:"buy_order_id"`
	SellOrderID uuid.UUID      `json:"sell_order_id"`
	MarketID    uuid.UUID      `json:"market_id"`
	Price       string         `json:"price"`
	Amount      string         `json:"amount"`
	Fee         sql.NullString `json:"fee"`
	CreatedAt   sql.NullTime   `json:"created_at"`
}

type Transaction struct {
	ID        uuid.UUID         `json:"id"`
	Username  string            `json:"username"`
	UserEmail string            `json:"user_email"`
	Type      TransactionType   `json:"type"`
	Currency  string            `json:"currency"`
	Amount    string            `json:"amount"`
	Status    TransactionStatus `json:"status"`
	Address   sql.NullString    `json:"address"`
	TxHash    sql.NullString    `json:"tx_hash"`
	CreatedAt sql.NullTime      `json:"created_at"`
}

type User struct {
	ID           uuid.UUID    `json:"id"`
	Username     string       `json:"username"`
	Email        string       `json:"email"`
	PasswordHash string       `json:"password_hash"`
	CreatedAt    sql.NullTime `json:"created_at"`
	UpdatedAt    sql.NullTime `json:"updated_at"`
	IsVerified   sql.NullBool `json:"is_verified"`
	Role         UserRole     `json:"role"`
}

type Wallet struct {
	ID            uuid.UUID      `json:"id"`
	Username      string         `json:"username"`
	UserEmail     string         `json:"user_email"`
	Currency      string         `json:"currency"`
	Balance       sql.NullString `json:"balance"`
	LockedBalance sql.NullString `json:"locked_balance"`
	CreatedAt     sql.NullTime   `json:"created_at"`
}
