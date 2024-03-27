package sqlc

import (
	"context"
	"github/chinmayweb3/simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T, toAcc Account, fromAcc Account) Transfer {
	args := CreateTransferParams{
		FromAccountID: fromAcc.ID,
		ToAccountID:   toAcc.ID,
		Amount:        util.RandomMoney(),
	}

	resTran, err := TestQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, resTran)

	require.Equal(t, args.FromAccountID, fromAcc.ID)
	require.Equal(t, args.ToAccountID, toAcc.ID)

	require.NotZero(t, resTran.ID)
	require.NotZero(t, resTran.CreatedAt)
	return resTran

}

func TestCreateTransfer(t *testing.T) {
	toAcc := RandomAccount(t)
	fromAcc := RandomAccount(t)
	CreateRandomTransfer(t, toAcc, fromAcc)

}
func TestGetTransfer(t *testing.T) {
	toAcc := RandomAccount(t)
	fromAcc := RandomAccount(t)
	trans := CreateRandomTransfer(t, toAcc, fromAcc)

	resTran, err := TestQueries.GetTransfer(context.Background(), trans.ID)
	require.NoError(t, err)
	require.NotEmpty(t, resTran)

	require.Equal(t, toAcc.ID, resTran.ToAccountID)
	require.Equal(t, fromAcc.ID, resTran.FromAccountID)

	require.NotZero(t, resTran.ID)
	require.NotZero(t, resTran.CreatedAt)
}
func TestListTransfer(t *testing.T) {
	toAcc := RandomAccount(t)
	fromAcc := RandomAccount(t)

	for range [10]int{} {
		CreateRandomTransfer(t, toAcc, fromAcc)
	}

	args := ListTransferParams{
		Limit:  5,
		Offset: 5,
	}
	resTrans, err := TestQueries.ListTransfer(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, resTrans, int(args.Limit))

	for _, tran := range resTrans {
		require.NotEmpty(t, tran.ID)
	}
}
