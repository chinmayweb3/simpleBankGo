package sqlc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {

	account := CreateAccountParams{
		Owner:    "chinmay",
		Balance:  1000,
		Currency: "USD",
	}

	resAcc, err := TestQueries.CreateAccount(context.Background(), account)
	require.NoError(t, err)
	require.NotEmpty(t, resAcc)

	require.Equal(t, resAcc.Balance, account.Balance)
	require.Equal(t, resAcc.Currency, account.Currency)
	require.Equal(t, resAcc.Owner, account.Owner)

	require.NotZero(t, resAcc.ID)
	require.NotZero(t, resAcc.CreatedAt)
}
