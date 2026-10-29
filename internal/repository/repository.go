package repository

import (
	"PlataTest/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
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
	INSERT INTO quotes (pair, created_at)
	VALUES ($1, $2)
	RETURNING id
	`
	err := r.pool.QueryRow(ctx, query, pair, time.Now()).Scan(&id)
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) && pgError.Code == "23505" { // номер ошибки дубликата
		return 0, fmt.Errorf("%w, pair %s", domain.ErrAlreadyExists, pair)
	}
	return id, nil
}
func (r *Repository) GetLatestQuote(ctx context.Context, pair string) (*domain.Quote, error) {
	var quote domain.Quote
	query := `
	SELECT id, pair, price, updated_at, status, created_at FROM quotes
	WHERE pair = $1
	`
	err := r.pool.QueryRow(ctx, query, pair).Scan(
		&quote.ID,
		&quote.Pair,
		&quote.Price,
		&quote.UpdatedAt,
		&quote.Status,
		&quote.CreatedAt,
	)
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
		return fmt.Errorf("quote %s not found", q.Pair)
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

func (r *Repository) GetToUpdateQuote(ctx context.Context) (domain.Quote, error) {
	query := `
	SELECT id, pair FROM quotes
	WHERE status = 'pending'
	ORDER BY created_at
	`
	quote := domain.Quote{}
	err := r.pool.QueryRow(ctx, query).Scan(&quote.ID, &quote.Pair)
	if err != nil {
		return domain.Quote{}, fmt.Errorf("failed to get quote to update: %w", err)
	}
	return quote, nil
}
