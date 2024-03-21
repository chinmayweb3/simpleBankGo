package sqlc

import (
	"context"
	"database/sql"
	"github/chinmayweb3/simplebank/util"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func RandomAccount(t *testing.T) Account {
	account := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	resAcc, err := TestQueries.CreateAccount(context.Background(), account)
	require.NoError(t, err)
	require.NotEmpty(t, resAcc)

	require.Equal(t, resAcc.Balance, account.Balance)
	require.Equal(t, resAcc.Currency, account.Currency)
	require.Equal(t, resAcc.Owner, account.Owner)

	require.NotZero(t, resAcc.ID)
	require.NotZero(t, resAcc.CreatedAt)
	return resAcc

}
func TestCreateAccount(t *testing.T) {
	RandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc1 := RandomAccount(t)

	resAcc, err := TestQueries.GetAccount(context.Background(), acc1.ID)
	require.NoError(t, err)
	require.Equal(t, acc1.ID, resAcc.ID)
	require.Equal(t, acc1.Owner, resAcc.Owner)
	require.Equal(t, acc1.Balance, resAcc.Balance)
	require.Equal(t, acc1.Currency, resAcc.Currency)
	require.WithinDuration(t, acc1.CreatedAt, resAcc.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	acc := RandomAccount(t)
	upAcc := UpdateAccountParams{ID: acc.ID, Balance: util.RandomMoney()}
	resAcc, err := TestQueries.UpdateAccount(context.Background(), upAcc)

	require.NoError(t, err)
	require.NotEmpty(t, resAcc)

	require.Equal(t, acc.ID, resAcc.ID)
	require.Equal(t, acc.Currency, resAcc.Currency)
	require.Equal(t, acc.Owner, resAcc.Owner)
	require.Equal(t, upAcc.Balance, resAcc.Balance)
}

func TestDeleteAccount(t *testing.T) {
	acc := RandomAccount(t)

	err := TestQueries.DeleteAccount(context.Background(), acc.ID)
	require.NoError(t, err)

	acc2, err := TestQueries.GetAccount(context.Background(), acc.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, acc2)
}

func TestListAccount(t *testing.T) {

	for range [10]int{} {
		RandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accs, err := TestQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accs, int(arg.Limit))

	for _, acc := range accs {
		require.NotEmpty(t, acc.Owner)
	}
}
