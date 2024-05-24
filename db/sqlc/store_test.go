package db

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(account1.Balance)
	fmt.Println(account2.Balance)

	arg := TransferTxParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        10,
	}

	n := 100
	var results = make(chan TransferTxResult)
	var errs = make(chan error)
	wg := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        10,
			})
			results <- result
			errs <- err

		}()
	}
	wg.Wait()

	var fromAccount Account
	var toAccount Account

	for i := 0; i < n; i++ {
		// chu y thu tu don kenh vao phan kenh *****
		result := <-results
		require.NotEmpty(t, result)
		err := <-errs
		require.NoError(t, err)

		// check transfer
		transfer := result.Transfer
		require.Equal(t, transfer.FromAccountID, arg.FromAccountID)
		require.Equal(t, transfer.ToAccountID, arg.ToAccountID)
		require.Equal(t, transfer.Amount, arg.Amount)
		require.WithinDuration(t, time.Now(), transfer.CreatedAt, time.Second)
		//check entry
		fromEntry := result.FromEntry
		require.Equal(t, fromEntry.AccountID, arg.FromAccountID)
		require.Equal(t, fromEntry.Amount, -arg.Amount)
		require.WithinDuration(t, time.Now(), fromEntry.CreatedAt, time.Second)

		toEntry := result.ToEntry
		require.Equal(t, toEntry.AccountID, arg.ToAccountID)
		require.Equal(t, toEntry.Amount, arg.Amount)
		require.WithinDuration(t, time.Now(), toEntry.CreatedAt, time.Second)

		// check account
		fromAccount = result.FromAccount
		require.Equal(t, fromAccount.Owner, account1.Owner)
		require.Equal(t, fromAccount.Currency, account1.Currency)

		toAccount = result.ToAccount
		require.Equal(t, toAccount.Owner, account2.Owner)
		require.Equal(t, toAccount.Currency, account2.Currency)

	}

	fmt.Printf("after From : %d , To : %d \n", fromAccount.Balance, toAccount.Balance)
	require.Equal(t, account1.Balance-int64(n)*arg.Amount, fromAccount.Balance)
	require.Equal(t, account2.Balance+int64(n)*arg.Amount, toAccount.Balance)

}
