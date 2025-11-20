package db

import (
	"context"
	"fmt"
	"log"
)


type OrderForCurrencyTxParams struct {
	UserEmail     string
	BaseCurrency  string
	QuoteCurrency string
}


func (store *SQLStore) OrderForCurrencyTx(ctx context.Context, arg OrderForCurrencyTxParams) (error) {
	err := store.execTx(ctx, func(q *Queries) error {

		wallets, err := q.GetWalletsByUserEmail(ctx, arg.UserEmail)

		if err != nil {
			return fmt.Errorf("failed to get wallets for the user: %w", err)
		}

		log.Printf("MARKET currencies received in tx: Base=%q Quote=%q", arg.BaseCurrency, arg.QuoteCurrency)



		walletMap := make(map[string]bool)

		for _, wallet := range wallets {
    		log.Println("WALLET: ", wallet.Currency)
    		walletMap[wallet.Currency] = true
		}

		if !walletMap[arg.BaseCurrency] || !walletMap[arg.QuoteCurrency] {
    		return fmt.Errorf("user does not have wallet for the required currencies")
		}

		return err

	})

	return err
}