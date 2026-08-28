package repository

import (
	"PlataTest/domain"
	"context"
	"fmt"
	"time"

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
func (r *Repository) UpdateQuote(ctx context.Context, id int, price float64) error {
	query := `
	UPDATE quotes 
	SET price = $1, updated_at = $2, status = 'completed'
	WHERE id = $3
	`
	res, err := r.pool.Exec(ctx, query, price, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update quote, %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("quote %d not found", id)
	}
	return nil
}
func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `
	UPDATE quotes 
	SET status = $1
	WHERE id = $2
	`
	res, err := r.pool.Exec(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("quote %d not found", id)
	}
	return nil
}
