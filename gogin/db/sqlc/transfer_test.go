package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        1000,
	}

	createdTransferId, _ := testQueries.CreateTransfer(context.Background(), arg)
	transferId, _ := createdTransferId.LastInsertId()
	transfer, err := testQueries.GetTransfer(context.Background(), transferId)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}
