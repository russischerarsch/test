package worker

import (
	"PlataTest/internal/domain"
	"PlataTest/internal/service"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

func RunPriceUpdater(
	ctx context.Context,
	repository service.Repository,
	client http.Client,
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
	client http.Client,
	apiKey string,
) {
	queryCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	quote, err := repository.GetToUpdateQuote(queryCtx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return
		}
		log.Printf("failed to get pending quotes, err %v", err)
		return
	}

	if err := updateQuote(ctx, repository, client, apiKey, &quote); err != nil {
		log.Printf("pair %s: %v", quote.Pair, err)
	}

}

func updateQuote(
	ctx context.Context,
	repository service.Repository,
	client http.Client,
	apiKey string,
	q *domain.Quote,
) error {

	childCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	URL := fmt.Sprintf(
		"http://api.exchangeratesapi.io/v1/latest?access_key=%s&base=%s&symbols=%s",
		apiKey, q.Pair[:3], q.Pair[3:],
	)

	req, err := http.NewRequestWithContext(childCtx, http.MethodGet, URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create req: %w", err)
	}
	r, err := client.Do(req)
	if err != nil {
		return domain.ErrAPINotAvailable
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("API service returned %d status code", r.StatusCode)
	}
	defer r.Body.Close()
	var response domain.Info
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}
	if !response.Success {
		if response.Error != nil {
			return fmt.Errorf(
				"API error %d: %s",
				response.Error.Code,
				response.Error.Info,
			)
		}
		return errors.New("API request failed")
	}
	if !response.Success {
		if response.Error != nil {
			return fmt.Errorf("api error %d: %s", response.Error.Code, response.Error.Info)
		}
		return fmt.Errorf("api returned success=false")
	}
	if len(response.Rates) == 0 {
		return fmt.Errorf("empty rates")
	}

	now := time.Now()
	var rate domain.Quote
	for currency, value := range response.Rates {
		rate = domain.Quote{
			ID:        q.ID,
			Pair:      currency,
			Price:     &value,
			UpdatedAt: &now,
		}
	}

	if err := repository.UpdateQuote(childCtx, &rate); err != nil {
		return fmt.Errorf("update quote: %w", err)
	}
	return nil
}
