package db

import (
	"context"
	"github.com/yongjeong-kim/go/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t * testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	// account, err := testQueries.CreateAccount(context.Background(), arg)
	err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	//require.NotEmpty(t, account)
	//
	//require.Equal(t, arg.Owner, account.Owner)
	//require.Equal(t, arg.Balance, account.Balance)
	//require.Equal(t, arg.Currency, account.Currency)
	//
	//require.NotZero(t, account.ID)
	//require.NotZero(t, account.CreatedAt)
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccouont(context.Background(), account1.ID)
}