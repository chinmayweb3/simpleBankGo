package sqlc

import (
	"context"
	"github/chinmayweb3/simplebank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomEntry(t *testing.T, acc Account) Entry {
	args := CreateEntryParams{
		AccountID: acc.ID,
		Amount:    util.RandomMoney(),
	}

	resEn, err := TestQueries.CreateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, resEn)

	require.Equal(t, args.AccountID, resEn.AccountID)
	require.Equal(t, args.Amount, resEn.Amount)
	return resEn
}

func TestCreateEntry(t *testing.T) {
	acc := RandomAccount(t)
	CreateRandomEntry(t, acc)
}

func TestGetEntry(t *testing.T) {
	acc := RandomAccount(t)
	en := CreateRandomEntry(t, acc)

	resEn, err := TestQueries.GetEntry(context.Background(), en.ID)
	require.NoError(t, err)
	require.NotEmpty(t, resEn)

	require.Equal(t, en.AccountID, resEn.AccountID)
	require.Equal(t, en.Amount, resEn.Amount)
	require.Equal(t, en.ID, resEn.ID)
	require.Equal(t, en.CreatedAt, resEn.CreatedAt)
}

func TestListEntries(t *testing.T) {
	acc := RandomAccount(t)
	for range [10]int{} {
		CreateRandomEntry(t, acc)
	}

	args := ListEntriesParams{
		Limit:     5,
		Offset:    5,
		AccountID: acc.ID,
	}
	resEntries, err := TestQueries.ListEntries(context.Background(), args)
	require.NoError(t, err)
	require.Len(t, resEntries, int(args.Limit))

	for _, en := range resEntries {
		require.NotEmpty(t, en.AccountID)
	}
}
