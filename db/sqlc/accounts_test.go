package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/rohanchauhan02/simplebank/util"
	"github.com/stretchr/testify/require"
)

func createRandomAccounts(t *testing.T) Account {
	args := CreateAccountsParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccounts(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.NotZero(t, account.ID)
	return account
}
func TestCreateAccounts(t *testing.T) {
	createRandomAccounts(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccounts(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccounts(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, arg.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccounts(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.Empty(t, account2)
	require.EqualError(t, err, sql.ErrNoRows.Error())
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccounts(t)
	}

	arg := ListAccountParams{
		Limit:  10,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, int(arg.Limit))

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}
