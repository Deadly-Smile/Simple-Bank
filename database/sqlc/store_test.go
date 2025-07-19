package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/Deadly-Smile/simple-bank/database/utils"
	"github.com/stretchr/testify/require"
)

func TestTransaction(t *testing.T) {
	store := NewStore(testDB)
	// number of transactions
	limit := int32(5)

	account1 := CreateRandomAccount(t, int64(1000*limit))
	account2 := CreateRandomAccount(t, int64(1000*limit))

	errs := make(chan error)
	results := make(chan TransfarResult)
	amounts := make(chan int64)
	fullAmount := 0

	for i := 0; i < int(limit); i++ {
		go func() {
			amount := utils.RandomInt(0, 1000)
			transfarParams := TransfarParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			}

			result, err := store.TransfarTx(context.Background(), transfarParams)
			errs <- err
			results <- result
			amounts <- amount
			fullAmount += int(amount)
		}()
	}

	for i := 0; i < int(limit); i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		amount := <-amounts

		transfar := result.Transfar
		require.NotEmpty(t, transfar)
		require.Equal(t, transfar.FromAccountID, account1.ID)
		require.Equal(t, transfar.ToAccountID, account2.ID)
		require.Equal(t, transfar.Amount, amount)
		require.NotZero(t, transfar.ID)

		_, err = store.GetTransfar(context.Background(), transfar.ID)
		require.NoError(t, err)
		require.NotEmpty(t, transfar.CreatedAt)
		require.WithinDuration(t, transfar.CreatedAt, time.Now(), time.Second)

		// Check both entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -transfar.Amount)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotEmpty(t, fromEntry.CreatedAt)
		require.WithinDuration(t, fromEntry.CreatedAt, transfar.CreatedAt, time.Second)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, transfar.Amount)
		require.Equal(t, toEntry.Amount, amount)
		require.NotEmpty(t, toEntry.CreatedAt)
		require.WithinDuration(t, toEntry.CreatedAt, transfar.CreatedAt, time.Second)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)
		require.Equal(t, fromAccount.Balance, int64(account1.Balance-int64(fullAmount)))

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)
		require.Equal(t, toAccount.Balance, int64(account2.Balance+int64(fullAmount)))
	}

	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount1)
	require.Equal(t, updateAccount1.Balance, int64(account1.Balance)-int64(fullAmount))

	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount2)
	require.Equal(t, updateAccount2.Balance, int64(account2.Balance)+int64(fullAmount))
}

func TestBothwayTransactionDeadlock(t *testing.T) {
	store := NewStore(testDB)
	// number of transactions
	limit := int32(10)

	account1 := CreateRandomAccount(t, int64(1000*limit))
	account2 := CreateRandomAccount(t, int64(1000*limit))

	errs := make(chan error)
	fullAmount := 0

	for i := 0; i < int(limit); i++ {
		go func() {
			amount := utils.RandomInt(0, 1000)
			if i&1 == 0 {
				transfarParams := TransfarParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        amount,
				}

				_, err := store.TransfarTx(context.Background(), transfarParams)
				errs <- err
				fullAmount += int(amount)
			} else {
				transfarParams := TransfarParams{
					FromAccountID: account2.ID,
					ToAccountID:   account1.ID,
					Amount:        amount,
				}
				_, err := store.TransfarTx(context.Background(), transfarParams)
				errs <- err
				fullAmount -= int(amount)
			}
		}()
	}

	for i := 0; i < int(limit); i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount1)
	require.Equal(t, updateAccount1.Balance, int64(account1.Balance)-int64(fullAmount))

	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount2)
	require.Equal(t, updateAccount2.Balance, int64(account2.Balance)+int64(fullAmount))
}
