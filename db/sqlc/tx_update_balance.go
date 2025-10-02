package db

import (
	"context"
	"sort"
)


type UpdateBalanceForTransactionTypeTxParams struct {
	CreateTransactionParams CreateTransactionParams
}

type UpdateBalanceForTransactionTypeTxResult struct {
	CreateTransactionRow CreateTransactionRow
}


func (store *SQLStore) CreateTransactionForTransactionTypeTx(ctx context.Context, args UpdateBalanceForTransactionTypeTxParams) (UpdateBalanceForTransactionTypeTxResult, error) {
	var result UpdateBalanceForTransactionTypeTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		amount := args.CreateTransactionParams.Amount

		wallets, err := q.GetWallets(ctx)

		sort.Slice(wallets, func(i, j int) bool {
    		return wallets[i].Balance.GreaterThan(wallets[j].Balance)
		})

		maxBalance := wallets[0].Balance

		switch args.CreateTransactionParams.Type {
			case TransactionTypeDeposit:
				result.CreateTransactionRow.Amount = amount
			case TransactionTypeWithdrawal:
				result.CreateTransactionRow.Amount = maxBalance.Sub(amount) 
			default:
				result.CreateTransactionRow.Amount = amount
		}

		result.CreateTransactionRow, err = q.CreateTransaction(ctx, args.CreateTransactionParams)

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}