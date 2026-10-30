package service

import (
	"PlataTest/internal/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	CreateQuote(ctx context.Context, pair string) (int, error)
	GetLatestQuote(ctx context.Context, pair string) (*domain.Quote, error)
	GetQuoteById(ctx context.Context, id int) (*domain.Quote, error)
	UpdateQuote(ctx context.Context, q *domain.Quote) error
	UpdateStatus(ctx context.Context, pair string, status string) error
	GetToUpdateQuote(ctx context.Context) (domain.Quote, error)
}
type QuoteService struct {
	repo Repository
}

func NewService(repo Repository) *QuoteService {
	return &QuoteService{
		repo: repo}
}
func (q *QuoteService) CreateQuote(ctx context.Context, pair string) (int, error) {
	pairLC := strings.ToLower(pair)
	if err := validatePair(pairLC); err != nil {
		return 0, domain.ErrInvalidPairFormat
	}

	id, err := q.repo.CreateQuote(ctx, pairLC)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (q *QuoteService) GetQuote(ctx context.Context, id int) (*domain.Quote, error) {
	quote, err := q.repo.GetQuoteById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	return quote, nil
}

func (q *QuoteService) GetLatestQuote(ctx context.Context, pair string) (*domain.Quote, error) {
	pairLC := strings.ToLower(pair)
	if err := validatePair(pairLC); err != nil {
		return nil, domain.ErrInvalidPairFormat
	}
	quote, err := q.repo.GetLatestQuote(ctx, pairLC)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return quote, nil
}
func validatePair(pair string) error {
	if len(pair) != 6 {
		return fmt.Errorf("pair must have 6 letters: %w", domain.ErrInvalidPairFormat)
	}

	for _, r := range pair {
		if r < 'a' || r > 'z' {
			return fmt.Errorf("pair must contain letters: %w", domain.ErrInvalidPairFormat)
		}
	}

	return nil
}
