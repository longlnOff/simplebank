package db

import (
	"context"
	"testing"
	"time"

	"github.com/longln/simplebank/utils"
	"github.com/stretchr/testify/require"
)






func createRandomTransfer(t *testing.T, fromAccount Account, toAccount Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID: toAccount.ID,
		Amount: utils.RandomMoney(),
	}

	transfer, err := Query.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, fromAccount.ID, transfer.FromAccountID)
	require.Equal(t, toAccount.ID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	createRandomTransfer(t, fromAccount, toAccount)
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	randomTransfer := createRandomTransfer(t, fromAccount, toAccount)	

	transfer, err := Query.GetTransfer(context.Background(), randomTransfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, transfer.FromAccountID, fromAccount.ID)
	require.Equal(t, transfer.ToAccountID, toAccount.ID)
	require.Equal(t, randomTransfer.Amount, transfer.Amount)
	require.WithinDuration(t, randomTransfer.CreatedAt, transfer.CreatedAt, time.Second)
}


func TestListTransfers(t *testing.T) {
	// create a random account
	fromAccount := createRandomAccount(t)
	toAccount := createRandomAccount(t)
	// create random 10 transfers from account to account
	for i := 0; i < 10; i++ {
		createRandomTransfer(t, fromAccount, toAccount)	
	}

	arg := ListTransfersParams{
		FromAccountID: fromAccount.ID,
		ToAccountID: toAccount.ID,
		Limit: 5,
		Offset: 5,
	}

	transfers, err := Query.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, 5, len(transfers))
}