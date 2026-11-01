package integration

import (
	"PlataTest/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateQuote_success(t *testing.T) {
	cleanupDB(t)
	id, err := ServiceTest.CreateQuote(t.Context(), "EURUSD")

	require.NoError(t, err)
	var ReqID int
	err = Pool.QueryRow(t.Context(), `SELECT id FROM quotes WHERE pair = 'eurusd'`).Scan(&ReqID)
	assert.NoError(t, err)
	assert.Equal(t, id, ReqID)
}

func TestCreateQuote_duplicate(t *testing.T) {
	cleanupDB(t)
	_, err := ServiceTest.CreateQuote(t.Context(), "EURUSD")
	require.NoError(t, err)
	_, err = ServiceTest.CreateQuote(t.Context(), "EURUSD")
	require.ErrorIs(t, err, domain.ErrAlreadyExists)
}

func TestCreateQuote_LongPair(t *testing.T) {
	cleanupDB(t)
	_, err := ServiceTest.CreateQuote(t.Context(), "EURBIMXN")
	require.ErrorIs(t, err, domain.ErrInvalidPairFormat)
	var amount int
	err = Pool.QueryRow(t.Context(), `SELECT COUNT(*) FROM quotes`).Scan(&amount)
	assert.NoError(t, err)
	assert.Equal(t, 0, amount)
}

func TestGetQuoteByID_success(t *testing.T) {
	cleanupDB(t)
	id, err := ServiceTest.CreateQuote(t.Context(), "EURUSD")
	require.NoError(t, err)

	quote, err := ServiceTest.GetQuote(t.Context(), id)
	require.NoError(t, err)
	assert.Equal(t, "EURUSD", quote.Pair)
}

func cleanupDB(t *testing.T) {
	t.Helper()
	_, err := Pool.Exec(t.Context(), `TRUNCATE quotes`)
	require.NoError(t, err)
}
