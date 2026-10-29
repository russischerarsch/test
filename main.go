package main

import (
	dbconnection "PlataTest/internal/db_connection"
	"PlataTest/internal/repository"
	"PlataTest/internal/service"
	"PlataTest/internal/transport"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate"
)

// ПРИВЕСТИ ВСЕ ПАРЫ К НИЖНЕМУ РЕГИСТРУ И ЗАКИНУТЬ В БД НА УНИК КЛЮЧ
func main() {
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://admin:secret@localhost:5432/quotes?sslmode=disable"
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

	if err = http.ListenAndServe(":8080", r); err != nil {
		fmt.Println(err)
	}
}
