package integration

import (
	"PlataTest/internal/repository"
	"PlataTest/internal/service"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var Pool *pgxpool.Pool
var URL string
var ServiceTest *service.QuoteService

func setUpEnvironment(ctx context.Context) (*pgxpool.Pool, testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "postgres:16-alpine",
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_USER":     "test",
		},
		Cmd:          []string{"postgres", "-c", "log_statement=all"},
		ExposedPorts: []string{"5432"},
		WaitingFor:   wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, pgContainer, err
	}
	host, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, pgContainer, err
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, pgContainer, err
	}

	URL = "postgres://test:test@" + host + ":" + port.Port() + "/testdb?sslmode=disable"

	pool, err := pgxpool.New(ctx, URL)
	if err != nil {
		return nil, pgContainer, err
	}

	if err := pool.Ping(ctx); err != nil {
		return pool, pgContainer, err
	}

	mig, err := migrate.New("file://../../../migrations", URL)
	if err != nil {
		return pool, pgContainer, err
	}
	defer mig.Close()
	if err := mig.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return pool, pgContainer, err
	}
	return pool, pgContainer, nil
}
func TestMain(m *testing.M) {
	ctx := context.Background()
	poolConns, container, err := setUpEnvironment(ctx)
	if err != nil {
		if poolConns != nil {
			poolConns.Close()
		}
		if termErr := testcontainers.TerminateContainer(container); termErr != nil {
			log.Printf("failed to terminate container, err %v", termErr)
		}
		log.Fatalf("failed to set up environment,err %v", err)
	}

	Pool = poolConns
	repo := repository.NewRepository(Pool)
	serv := service.NewService(repo)
	ServiceTest = serv
	m.Run()
	reader, _ := container.Logs(ctx)
	io.Copy(os.Stdout, reader)
	Pool.Close()
	if err := testcontainers.TerminateContainer(container); err != nil {
		log.Printf("failed to terminate container, err %v", err)
	}
}
