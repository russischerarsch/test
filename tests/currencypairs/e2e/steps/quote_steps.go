package steps

import (
	"PlataTest/tests/currencypairs/e2e/clients"
	"PlataTest/tests/currencypairs/models"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type httpClient struct {
	client *httpClient
}

func NewClient(cl *httpClient) *httpClient {
	return &httpClient{client: cl}
}

func (c *httpClient) CreateQuote(
	t *testing.T,
	ctx context.Context,
	cli *clients.QuoteClient,
	req models.CreatePairReq,
) *models.CreatePairResp {
	t.Helper()
	resp, err := cli.CreateQuote(ctx, models.CreatePairReq{Pair: "EURUSD"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.NotNil(t, resp.Body)
	return resp.Data
}

func (c *httpClient) GetQuoteByID(
	t *testing.T,
	ctx context.Context,
	cli *clients.QuoteClient,
	id string,
) *models.Quote {
	t.Helper()
	resp, err := cli.GetQuoteByID(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.NotNil(t, resp.Body)
	var res *models.Quote
	json.Unmarshal(resp.Body, &res)
	return res
}
