package db

import (
	"context"
	"fmt"
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

		wallets, err := q.GetWalletsByUserEmail(ctx, args.CreateTransactionParams.UserEmail)

		if err != nil {
			fmt.Println("Error fetching wallets:", err)
		}

		if len(wallets) == 0 {
			return fmt.Errorf("no wallets found for user email: %s", args.CreateTransactionParams.UserEmail)
		}

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

		q.UpdateWalletBalance(ctx, UpdateWalletBalanceParams{
			Balance: result.CreateTransactionRow.Amount,
			LockedBalance: wallets[0].LockedBalance,
			ID: wallets[0].ID,

		})

		if err != nil {
			return err
		}

		return nil
	})

	return result, err
}