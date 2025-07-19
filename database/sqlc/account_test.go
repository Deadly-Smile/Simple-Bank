package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/Deadly-Smile/simple-bank/database/utils"
	"github.com/stretchr/testify/require"
)

// CreateRandomAccount creates a random account for testing purposes.
func CreateRandomAccount(t *testing.T, minBalence ...int64) Account {
	ctx := context.Background()

	min := int64(0)
	if len(minBalence) > 0 {
		min = minBalence[0]
	}

	// Create a user first
	createdUser := CreateRandomUser(t)

	currencies := []string{"USD", "EUR", "GBP"}
	arg := CreateAccountParams{
		Owner:    createdUser.ID,
		Balance:  utils.RandomInt(min, min+1000),
		Currency: currencies[utils.RandomInt(0, 2)],
	}

	account, err := testQueries.CreateAccount(ctx, arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.WithinDuration(t, time.Now(), account.CreatedAt, time.Second)

	return account
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	createdAccount := CreateRandomAccount(t)
	fetchedAccount, err := testQueries.GetAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedAccount)

	require.Equal(t, createdAccount.ID, fetchedAccount.ID)
	require.Equal(t, createdAccount.Owner, fetchedAccount.Owner)
	require.Equal(t, createdAccount.Balance, fetchedAccount.Balance)
	require.Equal(t, createdAccount.Currency, fetchedAccount.Currency)
	require.WithinDuration(t, createdAccount.CreatedAt, fetchedAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	createdAccount := CreateRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), createdAccount.ID)
	require.NoError(t, err)

	fetchedAccount, err := testQueries.GetAccount(context.Background(), createdAccount.ID)
	require.Error(t, err)
	require.Empty(t, fetchedAccount)
}

func TestUpdateAccount(t *testing.T) {
	createdAccount := CreateRandomAccount(t)

	newBalance := utils.RandomInt(0, 1000)
	arg := UpdateAccountParams{
		ID:      createdAccount.ID,
		Balance: newBalance,
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	require.Equal(t, createdAccount.ID, updatedAccount.ID)
	require.Equal(t, createdAccount.Owner, updatedAccount.Owner)
	require.Equal(t, newBalance, updatedAccount.Balance)
	require.Equal(t, createdAccount.Currency, updatedAccount.Currency)
	require.WithinDuration(t, createdAccount.CreatedAt, updatedAccount.CreatedAt, time.Second)
}

func TestListAccount(t *testing.T) {
	createdAccount1 := CreateRandomAccount(t)
	createdAccount2 := CreateRandomAccount(t)

	arg := ListAccountsParams{
		Limit:  2,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 2)

	require.Contains(t, accounts, createdAccount1)
	require.Contains(t, accounts, createdAccount2)
}
