package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Store_interface interface {
	Querier
	CreateTransactionTx(ctx context.Context, arg TransactionsParams, feeArgs FeeParams) error
	UpdatedOrderTx(ctx context.Context, UpdatedOrderArgs UpdatedOrderParams) (ReturnAmountParams, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
	CreateTransactionForTransactionTypeTx(ctx context.Context, arg UpdateBalanceForTransactionTypeTxParams) (UpdateBalanceForTransactionTypeTxResult, error)
	CreateTradeTx(ctx context.Context, arg CreateTradeTxParams) (CreateTradeTxResult, error)
}

// Defining the SQLStore struct to execute the queries and transactions defined
type SQLStore struct {
	*Queries
	connPool *pgxpool.Pool
}

// func (store *SQLStore) UpdateOrderTx(context context.Context, param any) any {
// 	panic("unimplemented")
// }

func NewStore(connPool *pgxpool.Pool) Store_interface {
	return &SQLStore{
		Queries:  New(connPool),
		connPool: connPool,
	}
}

type TransactionsParams struct {
	Username  string          `json:"username"`
	UserEmail string          `json:"user_email"`
	Type      TransactionType `json:"type"`
	Currency  string          `json:"currency"`
	Amount    decimal.Decimal          `json:"amount"`
	Status    string          `json:"status"`
	Address   string          `json:"address"`
	TxHash    string          `json:"tx_hash"`
}

type FeeParams struct {
	MarketID uuid.UUID `json:"market_id"`
	Amount   decimal.Decimal    `json:"amount"`
	TakerFee decimal.Decimal    `json:"taker_fee"`
}

type UpdatedOrderParams struct {
	Status       OrderStatus `json:"status"`
	FilledAmount decimal.Decimal    `json:"filled_amount"`
	ID           uuid.UUID   `json:"id"`
}

type ReturnAmountParams struct {
	Amount decimal.Decimal `json:"amount"`
}

func (store *SQLStore) CreateTransactionTx(ctx context.Context, arg TransactionsParams, feeArgs FeeParams) error {
	return store.execTx(ctx, func(q *Queries) error {
		_, err := q.CreateTransaction(ctx, CreateTransactionParams{
			Username:  arg.Username,
			UserEmail: arg.UserEmail,
			Type:      arg.Type,
			Currency:  arg.Currency,
			Amount:    arg.Amount,
			Address:   arg.Address,
			TxHash:    arg.TxHash,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateFee(ctx, CreateFeeParams{
			Username: arg.Username,
			MarketID: feeArgs.MarketID,
			MakerFee: feeArgs.Amount,
			TakerFee: feeArgs.TakerFee,
		})

		if err != nil {
			return err
		}

		return nil
	})
}

func (store *SQLStore) UpdatedOrderTx(ctx context.Context, UpdatedOrderArgs UpdatedOrderParams) (ReturnAmountParams, error) {
	var returnAmount ReturnAmountParams
	var err error
	err = store.execTx(ctx, func(q *Queries) error {
		err = q.UpdateOrderStatusAndFilledAmount(ctx, UpdateOrderStatusAndFilledAmountParams{
			Status:       UpdatedOrderArgs.Status,
			FilledAmount: UpdatedOrderArgs.FilledAmount,
			ID:           UpdatedOrderArgs.ID,
		})

		returnAmount = ReturnAmountParams{
			Amount: UpdatedOrderArgs.FilledAmount,
		}

		if err != nil {
			return err
		}

		return nil

	})

	return returnAmount, err

}
