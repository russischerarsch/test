package worker

import (
	quotespb "PlataTest/contracts"
	"PlataTest/internal/domain"
	"PlataTest/internal/service"
	"context"
	"fmt"
	"log"
	"time"
)

func RunPriceUpdater(
	ctx context.Context,
	repository service.Repository,
	client quotespb.QuotesClient,
	apiKey string,
) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			updatePrices(ctx, repository, client, apiKey)
		}
	}
}

func updatePrices(
	ctx context.Context,
	repository service.Repository,
	client quotespb.QuotesClient,
	apiKey string,
) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	quotes, err := repository.GetToUpdateQuotes(queryCtx)
	if err != nil {
		log.Printf("failed to get pending quotes, err %v", err)
		return
	}

	for _, v := range quotes {
		if err := updateQuote(ctx, repository, client, apiKey, v.Pair); err != nil {
			log.Printf("pair %s: %v", v.Pair, err)
		}
	}
}

func updateQuote(
	ctx context.Context,
	repository service.Repository,
	client quotespb.QuotesClient,
	apiKey string,
	pair string,
) error {

	childCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	URL := fmt.Sprintf(
		"http://api.exchangeratesapi.io/v1/latest?access_key=%s&base=%s&symbols=%s",
		apiKey, pair[:3], pair[3:],
	)

	resp, err := client.FetchPrice(childCtx, &quotespb.QuoteRequest{URL: URL})
	if err != nil {
		return fmt.Errorf("fetch price: %w", err)
	}

	quote, err := toRates(resp)
	if err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	if err := repository.UpdateQuote(childCtx, quote); err != nil {
		return fmt.Errorf("update quote: %w", err)
	}
	return nil
}

func toRates(info *quotespb.QuoteResponse) (*domain.Quote, error) {
	if !info.Success {
		if info.Error != nil {
			return nil, fmt.Errorf("api error %d: %s", info.Error.Code, info.Error.Info)
		}
		return nil, fmt.Errorf("api returned success=false")
	}
	if len(info.Rates) == 0 {
		return nil, fmt.Errorf("empty rates")
	}

	now := time.Now()
	var rate domain.Quote
	for currency, value := range info.Rates {
		rate = domain.Quote{
			Pair:      currency,
			Price:     value,
			UpdatedAt: now,
		}
	}
	return &rate, nil
}
