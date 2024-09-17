package db

import (
	"context"
	"testing"
	"time"

	"github.com/longln/simplebank/utils"
	"github.com/stretchr/testify/require"
)


func createRandomEntry(t *testing.T, id int64, amount int64) Entry {
	arg := CreateEntryParams{
		AccountID: id,
		Amount: amount,
	}

	entry, err := Query.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	return entry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	amount := utils.RandomMoney()
	createRandomEntry(t, account.ID, int64(amount))
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t)
	amount := utils.RandomMoney()
	entry := createRandomEntry(t, account.ID, int64(amount))

	newEntry, err := Query.GetEntry(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, newEntry)

	require.Equal(t, newEntry.AccountID, entry.AccountID)
	require.Equal(t, newEntry.ID, entry.ID)
	require.Equal(t, newEntry.Amount, entry.Amount)
	require.WithinDuration(t, newEntry.CreatedAt, entry.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {
	// create random account
	account := createRandomAccount(t)

	// create 10 entries from account
	for i := 0; i < 10; i++ {
		amount := utils.RandomMoney()
		createRandomEntry(t, account.ID, amount)
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit: 5,
		Offset: 5,
	}

	entries, err := Query.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	require.Equal(t, 5, len(entries))

	for _, entry := range(entries) {
		require.NotEmpty(t, entry)
		require.Equal(t, account.ID, entry.AccountID)
	}
}