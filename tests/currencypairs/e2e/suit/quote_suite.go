package suit

import (
	"PlataTest/tests/currencypairs/e2e/clients"
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

type QuoteSuite struct {
	suite.Suite
	Client *clients.QuoteClient
	Pool   *pgxpool.Pool
}

func (s *QuoteSuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	url := os.Getenv("SERVICE_URL")
	// s.Require().NotEmpty(url)
	if url == "" {
		url = "http://localhost:8080"
	}

	s.Client = clients.NewQuoteClient(clients.NewHttpClient(url))
	dsn := os.Getenv("POSTGRES_CONN")
	// s.Require().NotEmpty(dsn)
	if dsn == "" {
		dsn = "postgres://admin:secret@postgres:5432/quotes?sslmode=disable"
	}
	pool, err := pgxpool.New(ctx, dsn)
	s.Require().NoError(err)
	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		s.Require().NoError(err)
	}
	s.Pool = pool
}

func (s *QuoteSuite) TearDownSuite() {
	if s.Pool != nil {
		s.Pool.Close()
	}
}
func (s *QuoteSuite) DeleteQuote(id string) {
	_, err := s.Pool.Exec(context.Background(), `DELETE FROM quotes WHERE id = $1`, id)
	if err != nil {
		s.T().Logf("failed to delete quote %s: %v", id, err)
	}
}
