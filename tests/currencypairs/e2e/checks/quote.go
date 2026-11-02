package checks

import (
	"PlataTest/tests/currencypairs/models"
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func CreatePairResponseValid(t *testing.T, resp *models.CreatePairResp) {
	t.Helper()

	require.NotNil(t, resp)
	require.NotEmpty(t, resp.ID)
}

func QuoteResponseMatches(t *testing.T, expPair, expID string, got *models.Quote) {
	t.Helper()

	require.NotNil(t, got)
	require.Equal(t, expPair, got.Pair)
	require.Equal(t, expID, got.ID)

}

func CheckQuoteStatus(t *testing.T, ctx context.Context, pool *pgxpool.Pool, ID, expPair string) {
	t.Helper()
	query := `
	SELECT pair, status, price FROM quotes 
	WHERE id = $1
	`
	var quote *models.Quote
	err := pool.QueryRow(ctx, query, ID).Scan(&quote.Pair, &quote.Status, &quote.Price)
	require.NoError(t, err)

}

func PairInDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, expPair, ID string) {
	t.Helper()
	query := `
	SELECT pair FROM quotes 
	WHERE id = $1
	`
	var pair string
	err := pool.QueryRow(ctx, query, ID).Scan(&pair)
	require.NoError(t, err)
	require.Equal(t, expPair, pair)
}
