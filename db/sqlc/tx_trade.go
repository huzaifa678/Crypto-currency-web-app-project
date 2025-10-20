package db

import (
	"context"
	"fmt"

)

type CreateTradeTxParams struct {
	TradeParams CreateTradeParams
}

type CreateTradeTxResult struct {
	Trade Trade
}

func (store *SQLStore) CreateTradeTx(ctx context.Context, arg CreateTradeTxParams) (CreateTradeTxResult, error) {
	var result CreateTradeTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Trade, err = q.CreateTrade(ctx, arg.TradeParams)
		if err != nil {
			return fmt.Errorf("failed to create trade: %w", err)
		}

		trade := result.Trade
		amount := trade.Amount
		price := trade.Price
		fee := trade.Fee
		total := amount.Mul(price)

		buyerWallets, err := q.GetWalletsByUserEmail(ctx, trade.Username)
		if err != nil || len(buyerWallets) == 0 {
			return fmt.Errorf("buyer wallets not found for %s: %w", trade.Username, err)
		}

		sellerWallets, err := q.GetWalletsByUserEmail(ctx, trade.Username) // replace later with seller's email
		if err != nil || len(sellerWallets) == 0 {
			return fmt.Errorf("seller wallets not found: %w", err)
		}

		// Picking the first wallet in buyer/seller lists for now
		// logic might be extended later
		buyerWallet := buyerWallets[0]
		sellerWallet := sellerWallets[0]

		// Buyer pays total + fee, Seller receives total - fee
		newBuyerBalance := buyerWallet.Balance.Sub(total.Add(fee))
		newSellerBalance := sellerWallet.Balance.Add(total.Sub(fee))

		// Keeping locked balances unchanged for now
		newBuyerLocked := buyerWallet.LockedBalance
		newSellerLocked := sellerWallet.LockedBalance

		err = q.UpdateWalletBalance(ctx, UpdateWalletBalanceParams{
			ID:            buyerWallet.ID,
			Balance:       newBuyerBalance,
			LockedBalance: newBuyerLocked,
		})
		if err != nil {
			return fmt.Errorf("failed to update buyer wallet: %w", err)
		}

		err = q.UpdateWalletBalance(ctx, UpdateWalletBalanceParams{
			ID:            sellerWallet.ID,
			Balance:       newSellerBalance,
			LockedBalance: newSellerLocked,
		})
		if err != nil {
			return fmt.Errorf("failed to update seller wallet: %w", err)
		}

		return nil
	})

	return result, err
}
