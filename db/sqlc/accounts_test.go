package db

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"github.com/longln/simplebank/utils"
	"github.com/stretchr/testify/require"
)




func createRandomAccount(t *testing.T) Account {
	user := createRandomUser(t)
	arg := CreateAccountParams{
		Owner: user.UserName,
		Balance: utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.NotZero(t, account.ID)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}


func TestGetAccount(t *testing.T) {
	// create random account
	randomAccount := createRandomAccount(t)

	// get account and compare
	account, err := testQueries.GetAccount(context.Background(), randomAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.ID, randomAccount.ID)
	require.Equal(t, account.Owner, randomAccount.Owner)
	require.Equal(t, account.Balance, randomAccount.Balance)
	require.Equal(t, account.Currency, randomAccount.Currency)
	require.WithinDuration(t, randomAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestGetAccountForUpdate(t *testing.T) {
	// create random account
	randomAccount := createRandomAccount(t)

	// get account and compare
	account, err := testQueries.GetAccountForUpdate(context.Background(), randomAccount.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.ID, randomAccount.ID)
	require.Equal(t, account.Owner, randomAccount.Owner)
	require.Equal(t, account.Balance, randomAccount.Balance)
	require.Equal(t, account.Currency, randomAccount.Currency)
	require.WithinDuration(t, randomAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	randomAccount := createRandomAccount(t)
	balance := utils.RandomMoney()
	arg := UpdateAccountParams {
		ID: randomAccount.ID,
		Balance: balance,
	}
	account, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account.ID, randomAccount.ID)
	require.Equal(t, account.Owner, randomAccount.Owner)
	require.Equal(t, account.Balance, balance)
	require.Equal(t, account.Currency, randomAccount.Currency)
	require.WithinDuration(t, randomAccount.CreatedAt, account.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	// create random account
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	accountTest, err := testQueries.GetAccount(context.Background(), account.ID)
	require.Empty(t, accountTest)
	require.Error(t, err, sql.ErrNoRows)
}


func TestListAccounts(t *testing.T) {
	var lastAccount Account
	// create random 10 account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Owner: lastAccount.Owner,
		Limit: 5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range(accounts) {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}