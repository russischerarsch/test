package main

import (
	dbconnection "PlataTest/internal/db_connection"
	"PlataTest/internal/http"
	"PlataTest/internal/repository"
	"PlataTest/internal/service"
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	pool, err := dbconnection.CreateConnection(ctx)
	if err != nil {
		fmt.Println("failed to connect postgres", err)
		return
	}
	defer pool.Close()
	rep := repository.CreateRepository(pool)
	ser := service.CreateService(rep)
	handlers := http.CreateHandler(ser)
	r := http.SetupRouter(handlers)
	if err = r.Run(":8080"); err != nil {
		fmt.Println(err)
	}
}
