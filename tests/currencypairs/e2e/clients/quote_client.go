package clients

import (
	"PlataTest/internal/domain"
	"PlataTest/tests/currencypairs/models"
	"context"
	"net/http"
)

type QuoteClient struct {
	client *httpClient
}

func NewQuoteClient(client *httpClient) *QuoteClient {
	return &QuoteClient{
		client: client,
	}
}

func (c *QuoteClient) CreateQuote(ctx context.Context, req models.CreatePairReq) (*Response[models.CreatePairResp], error) {
	return Post[models.CreatePairResp](ctx, c.client, "/quotes", http.MethodPost, req)
}
func (c *QuoteClient) GetQuoteByID(ctx context.Context, id string) (*Response[domain.Quote], error) {
	return Get[domain.Quote](ctx, c.client, "/quotes/"+id, http.MethodGet, nil)
}
