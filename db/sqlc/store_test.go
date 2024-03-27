package sqlc

import (
	"context"
	"github/chinmayweb3/simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {

	store := NewStore(TestDB)

	fAcc := RandomAccount(t)
	tAcc := RandomAccount(t)

	transferTxParams := TransferTxParams{
		FromAccountId: fAcc.ID,
		ToAccountId:   tAcc.ID,
		Amount:        util.RandomMoney(),
	}
	const n = 5

	errs := make(chan error, n)
	results := make(chan TransferTxResult, n)

	for range [n]int{} {
		go func() {
			txResult, err := store.TransferTx(context.Background(), transferTxParams)
			errs <- err
			results <- txResult
		}()
	}

	for range [n]int{} {
		err := <-errs
		require.NoError(t, err)

		res := <-results

		// require.NotEmpty(t, res.FromAccount)
		// require.NotEmpty(t, res.ToAccount)
		require.NotEmpty(t, res.Transfer)
		require.Equal(t, res.Transfer.FromAccountID, transferTxParams.FromAccountId)
		require.Equal(t, res.Transfer.ToAccountID, transferTxParams.ToAccountId)

		require.NotZero(t, res.Transfer.ID)
		require.NotZero(t, res.Transfer.CreatedAt)

		_, err = store.queries.GetTransfer(context.Background(), res.Transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := res.FromEntry
		require.NotZero(t, fromEntry.ID)
		require.NotEmpty(t, fromEntry.ID)
		require.Equal(t, transferTxParams.FromAccountId, fromEntry.AccountID)
		require.Equal(t, -transferTxParams.Amount, fromEntry.Amount)

		_, err = store.queries.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := res.ToEntry
		require.NotZero(t, toEntry.ID)
		require.NotEmpty(t, toEntry.ID)
		require.Equal(t, transferTxParams.FromAccountId, toEntry.AccountID)
		require.Equal(t, transferTxParams.Amount, toEntry.Amount)

		_, err = store.queries.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check Account update

	}

}
