package db

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
)

type Store struct {
	*Queries
	db *sql.DB
}

func (store *Store) UpdateOrderTx(context context.Context, param any) any {
	panic("unimplemented")
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()

}

type TransactionsParams struct {
	UserID   uuid.UUID       `json:"user_id"`
	Type     TransactionType `json:"type"`
	Currency string          `json:"currency"`
	Amount   string          `json:"amount"`
	Status   string          `json:"status"`
	Address  sql.NullString  `json:"address"`
	TxHash   sql.NullString  `json:"tx_hash"`
}

type FeeParams struct {
	MarketID uuid.UUID      `json:"market_id"`
	Amount   sql.NullString `json:"amount"`
	TakerFee sql.NullString `json:"taker_fee"`
}

type UpdatedOrderParams struct {
	Status       OrderStatus    `json:"status"`
	FilledAmount sql.NullString `json:"filled_amount"`
	ID           uuid.UUID      `json:"id"`
}

type returnAmountParams struct {
	Amount sql.NullString `json:"amount"`
}


func (store *Store) CreateTransactionTx(ctx context.Context, arg TransactionsParams, feeArgs FeeParams) error {
	return store.execTx(ctx, func(q *Queries) error {
		_, err := q.CreateTransaction(ctx, CreateTransactionParams{
			UserID:   arg.UserID,
			Type:     arg.Type,
			Currency: arg.Currency,
			Amount:   arg.Amount,
			Address:  arg.Address,
			TxHash:   arg.TxHash,
		})
		if err != nil {
			return err
		}

		_, err = q.CreateFee(ctx, CreateFeeParams{
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

func (store *Store) UpdatedOrderTx(ctx context.Context, UpdatedOrderArgs UpdatedOrderParams) (returnAmountParams, error) {
	var returnAmount returnAmountParams
	var err error
	err = store.execTx(ctx, func(q *Queries) error {
		err = q.UpdateOrderStatusAndFilledAmount(ctx, UpdateOrderStatusAndFilledAmountParams{
			Status:       UpdatedOrderArgs.Status,
			FilledAmount: UpdatedOrderArgs.FilledAmount,
			ID:           UpdatedOrderArgs.ID,
		})

		returnAmount = returnAmountParams{
			Amount: UpdatedOrderArgs.FilledAmount,
		}

		if err != nil {
			return err
		}

		return nil

	})

	return returnAmount, err

}
