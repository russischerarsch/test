package checks

import (
	"PlataTest/tests/currencypairs/models"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
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

func CheckQuoteStatus(t *testing.T, ctx context.Context, pool *pgxpool.Pool, ID, expPair, expStatus string) {
	t.Helper()
	var quote models.Quote
	require.Eventually(t, func() bool {
		query := `
		SELECT pair, status, price FROM quotes 
		WHERE id = $1
		`

		err := pool.QueryRow(ctx, query, ID).Scan(&quote.Pair, &quote.Status, &quote.Price)
		if err != nil {
			t.Logf("failed to get quote, err %v", err)
			return false
		}
		return quote.Status == expStatus

	}, 12*time.Second, 1*time.Second)
	assert.Equal(t, expPair, quote.Pair)
	assert.Greater(t, quote.Price, float64(0))
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
