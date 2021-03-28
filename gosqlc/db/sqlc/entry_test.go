package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yongjeong-kim/go/util"
)

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	randomAmount := util.RandomMoney()

	createdEntryId, _ := testQueries.CreateEntry(context.Background(), CreateEntryParams{
		AccountID: account.ID,
		Amount:    randomAmount,
	})

	entryId, _ := createdEntryId.LastInsertId()
	entry, err := testQueries.GetEntry(context.Background(), entryId)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, entry.Amount, randomAmount)

	require.NotZero(t, entry.CreatedAt)
	require.NotZero(t, account.CreatedAt)
}
