package dbconnection

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateConnection(ctx context.Context, url string) (*pgxpool.Pool, error) {

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
	return pool, nil
}
