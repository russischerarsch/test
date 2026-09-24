package service

import (
	apperrors "PlataTest/app_errors"
	"PlataTest/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type Quote interface {
	CreateQuote(ctx context.Context, pair string) (int, error)
	GetLatestQuote(ctx context.Context, pair string) (*domain.Quote, error)
	GetQuoteById(ctx context.Context, id int) (*domain.Quote, error)
	UpdateQuote(ctx context.Context, id int, price float64) error
	UpdateStatus(ctx context.Context, id int, status string) error
}
type QuoteService struct {
	repo   Quote
	client *http.Client
}

func CreateService(repo Quote) *QuoteService {
	return &QuoteService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		}}
}
func (q *QuoteService) CreateQuote(ctx context.Context, pair string) (int, error) {
	if err := validatePair(pair); err != nil {
		return 0, apperrors.ErrInvalidPairFormat
	}

	id, err := q.repo.CreateQuote(ctx, pair)
	if err != nil {
		return 0, err
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := q.updatePrice(ctx, id, pair); err != nil {
			log.Printf("failed to update quote %d: %v", id, err)
		}
	}()

	return id, nil
}

func (q *QuoteService) GetQuote(ctx context.Context, id int) (*domain.Quote, error) {
	quote, err := q.repo.GetQuoteById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, err
	}

	return quote, nil
}

func (q *QuoteService) GetLatestQuote(ctx context.Context, pair string) (*domain.Quote, error) {
	if err := validatePair(pair); err != nil {
		return nil, apperrors.ErrInvalidPairFormat
	}
	quote, err := q.repo.GetLatestQuote(ctx, pair)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return quote, nil
}

func (q *QuoteService) updatePrice(ctx context.Context, id int, pair string) error {
	price, err := q.FetchPrice(ctx, pair)
	if err != nil {
		if updateErr := q.repo.UpdateStatus(ctx, id, "failed"); updateErr != nil {
			return fmt.Errorf("update quote: %w; update status: %v", err, updateErr)
		}
		return fmt.Errorf("failed to fetch price, %w", err)
	}
	if err := q.repo.UpdateQuote(ctx, id, price); err != nil {
		if updateErr := q.repo.UpdateStatus(ctx, id, "failed"); updateErr != nil {
			return fmt.Errorf("fetch price: %w; update status: %v", err, updateErr)
		}
		return fmt.Errorf("failed to update price, %w", err)
	}
	return nil
}

func (q *QuoteService) FetchPrice(ctx context.Context, pair string) (float64, error) {
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
		return 0, apperrors.ErrAPINotAvailable
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
		return fmt.Errorf("pair must have 6 letters: %w", apperrors.ErrInvalidPairFormat)
	}

	for _, r := range pair {
		if r < 'A' || r > 'Z' {
			return fmt.Errorf("pair must contain uppercase letters: %w", apperrors.ErrInvalidPairFormat)
		}
	}

	return nil
}
