package dbconnection

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateConnection(ctx context.Context) (*pgxpool.Pool, error) {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://admin:secret@localhost:5432/quotes?sslmode=disable"
	}
	var pool *pgxpool.Pool
	var err error
	for i := 0; i < 10; i++ {
		pool, err = pgxpool.New(ctx, url)
		if err != nil {
			time.Sleep(2 * time.Second)
			continue
		}
		if err = pool.Ping(ctx); err == nil {
			break
		}
		pool.Close()
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	query := `
	CREATE TABLE IF NOT EXISTS quotes(
		id SERIAL PRIMARY KEY,
		pair VARCHAR(7) NOT NULL,
		price DECIMAL,
		updated_at TIMESTAMP,
		status VARCHAR(10) DEFAULT 'pending'
	)
	`
	_, err = pool.Exec(ctx, query)
	if err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}
