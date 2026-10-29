package main

import (
	dbconnection "PlataTest/internal/db_connection"
	"PlataTest/internal/repository"
	"PlataTest/internal/service"
	"PlataTest/internal/transport"
	"PlataTest/internal/worker"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// ПРИВЕСТИ ВСЕ ПАРЫ К НИЖНЕМУ РЕГИСТРУ И ЗАКИНУТЬ В БД НА УНИК КЛЮЧ
func main() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://admin:secret@localhost:5432/quotes?sslmode=disable"
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "84343d3b62c279fb5c47f52fbdf034ae"

	}
	ctx := context.Background()
	pool, err := dbconnection.CreateConnection(ctx, url)
	if err != nil {
		log.Printf("failed to connect postgres: %v", err)
		return
	}
	m, err := migrate.New("file://migrations", url)
	if err != nil {
		log.Fatalf("failed to create migration, err %v", err)
	}
	if err = m.Up(); err != nil {
		log.Fatalf("failed to migrate db, err %v", err)
	}
	m.Close()
	defer pool.Close()
	rep := repository.CreateRepository(pool)
	ser := service.CreateService(rep)
	handlers := transport.NewHandler(ser)
	r := transport.SetupRouter(handlers)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	go func() {
		worker.RunPriceUpdater(
			ctx,
			rep,
			client,
			apiKey,
		)
	}()
	log.Print("server is ready to accept conns")
	if err = http.ListenAndServe(":8080", r); err != nil {
		fmt.Println(err)
	}
}
