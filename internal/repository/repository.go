package repository

import (
	"PlataTest/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func CreateRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
func (r *Repository) CreateQuote(ctx context.Context, pair string) (int, error) {
	var id int
	query := `
	INSERT INTO quotes (pair)
	VALUES ($1)
	RETURNING id
	`
	err := r.pool.QueryRow(ctx, query, pair).Scan(&id)
	return id, err
}
func (r *Repository) GetLatestQuote(ctx context.Context, pair string) (*domain.Quote, error) {
	var quote domain.Quote
	query := `
	SELECT id, pair, price, updated_at, status FROM quotes
	WHERE pair = $1 and status = 'completed'
	ORDER BY updated_at DESC
	LIMIT 1
	`
	err := r.pool.QueryRow(ctx, query, pair).Scan(
		&quote.ID,
		&quote.Pair,
		&quote.Price,
		&quote.UpdatedAt,
		&quote.Status)
	return &quote, err
}

func (r *Repository) GetQuoteById(ctx context.Context, id int) (*domain.Quote, error) {
	var quote domain.Quote
	query := `
	SELECT id, pair, price, updated_at, status FROM quotes
	WHERE id = $1
	`
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&quote.ID,
		&quote.Pair,
		&quote.Price,
		&quote.UpdatedAt,
		&quote.Status)
	return &quote, err
}
func (r *Repository) UpdateQuote(ctx context.Context, q *domain.Quote) error {
	query := `
	UPDATE quotes 
	SET price = $1, updated_at = $2, status = 'completed'
	WHERE id = $3
	`
	res, err := r.pool.Exec(ctx, query, q.Price, q.UpdatedAt, q.ID)
	if err != nil {
		return fmt.Errorf("failed to update quote, %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("quote %d not found", q.ID)
	}
	return nil
}
func (r *Repository) UpdateStatus(ctx context.Context, pair string, status string) error {
	query := `
	UPDATE quotes 
	SET status = $1
	WHERE pair = $2
	`
	res, err := r.pool.Exec(ctx, query, status, pair)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("quote %v not found", pair)
	}
	return nil
}

func (r *Repository) GetToUpdateQuotes(ctx context.Context) ([]domain.Quote, error) {
	query := `
	SELECT pair FROM quotes
	WHERE status = 'pending' OR 'failed'
	ORDER BY created_at
	LIMIT 10
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending quotes: %w", err)
	}
	defer rows.Close()
	quotes := []domain.Quote{}
	for rows.Next() {
		quote := domain.Quote{}
		if err := rows.Scan(&quote.Pair); err != nil {
			return nil, fmt.Errorf("failed to parse quote to struct, err %w", err)
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}
