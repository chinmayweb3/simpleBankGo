package sqlc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(TestDB)

	fAcc := RandomAccount(t)
	tAcc := RandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			transferTxParams := TransferTxParams{
				FromAccountId: fAcc.ID,
				ToAccountId:   tAcc.ID,
				Amount:        amount,
			}
			txResult, err := store.TransferTx(context.Background(), transferTxParams)
			errs <- err
			results <- txResult
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

		res := <-results
		require.NotEmpty(t, res.Transfer)

		transfer := res.Transfer
		require.Equal(t, transfer.FromAccountID, fAcc.ID)
		require.Equal(t, transfer.ToAccountID, tAcc.ID)

		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.queries.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := res.FromEntry
		require.NotZero(t, fromEntry.ID)
		require.NotEmpty(t, fromEntry.ID)
		require.Equal(t, fAcc.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)

		_, err = store.queries.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := res.ToEntry
		require.NotZero(t, toEntry.ID)
		require.NotEmpty(t, toEntry.ID)
		require.Equal(t, tAcc.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)

		_, err = store.queries.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check Account update
		fromAcc := res.FromAccount
		require.NotEmpty(t, fromAcc)
		require.Equal(t, fAcc.ID, fromAcc.ID)

		toAcc := res.ToAccount
		require.NotEmpty(t, toAcc)
		require.Equal(t, tAcc.ID, toAcc.ID)

		// check account balance
		diff1 := fAcc.Balance - fromAcc.Balance
		diff2 := toAcc.Balance - tAcc.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 >= 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)

	}

	ufAcc, err := store.queries.GetAccount(context.Background(), fAcc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, ufAcc)

	utAcc, err := store.queries.GetAccount(context.Background(), tAcc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, utAcc)

	require.Equal(t, fAcc.Balance-int64(n)*amount, ufAcc.Balance)
	require.Equal(t, tAcc.Balance+int64(n)*amount, utAcc.Balance)

}
