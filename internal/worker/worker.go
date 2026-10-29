package worker

import (
	"PlataTest/internal/domain"
	"PlataTest/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPriceUpdater(ctx context.Context, pool *pgxpool.Pool) error {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			childCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			quotes, err := repository.GetPendingQuotes(childCtx, pool)
			if err != nil {
				log.Printf("failed to send query, err %v", err)
				continue
			}
			for _, v := range quotes {

			}
		}
	}
}
func FetchPrice(ctx context.Context, pair string) (float64, error) {
	first := pair[:3]
	second := pair[3:]
	var ApiKey string
	if ApiKey = os.Getenv("API_KEY"); ApiKey == "" {
		ApiKey = "84343d3b62c279fb5c47f52fbdf034ae"
		// подобный ключ нужно хранить в переменной окружения. В проде упал бы с ошибкой чтения
	}
	url := fmt.Sprintf("http://api.exchangeratesapi.io/v1/latest?access_key=%s&base=%s&symbols=%s", ApiKey, first, second)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)

	}
	resp, err := q.client.Do(req)
	if err != nil {
		return 0, domain.ErrAPINotAvailable
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API service returned %d status code", resp.StatusCode)
	}
	defer resp.Body.Close()

	var info struct {
		Success bool               `json:"success"`
		Rates   map[string]float64 `json:"rates"`
		Error   *struct {
			Code int    `json:"code"`
			Info string `json:"info"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return 0, err
	}
	if !info.Success {
		if info.Error != nil {
			return 0, fmt.Errorf(
				"API error %d: %s",
				info.Error.Code,
				info.Error.Info,
			)
		}
		return 0, errors.New("API request failed")
	}
	price, ok := info.Rates[second]
	if !ok {
		return 0, fmt.Errorf("rate %s was not found\n", second)
	}
	return price, nil
}
func validatePair(pair string) error {
	if len(pair) != 6 {
		return fmt.Errorf("pair must have 6 letters: %w", domain.ErrInvalidPairFormat)
	}

	for _, r := range pair {
		if r < 'A' || r > 'Z' {
			return fmt.Errorf("pair must contain uppercase letters: %w", domain.ErrInvalidPairFormat)
		}
	}

	return nil
}
