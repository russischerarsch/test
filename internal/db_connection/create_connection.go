package dbconnection

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateConnection(ctx context.Context, url string) (*pgxpool.Pool, error) {

	var pool *pgxpool.Pool
	var err error
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
			pool, err = pgxpool.New(ctx, url)
			if err != nil {
				continue
			}
			if err = pool.Ping(ctx); err == nil {
				return pool, nil
			}
			pool.Close()
		}
	}
	return nil, err
}
