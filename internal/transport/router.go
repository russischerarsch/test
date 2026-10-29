package transport

import (
	"net/http"
)

func SetupRouter(handlers *Handlers) http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/quotes", handlers.CreateQuote)
	mux.HandleFunc("GET /api/quotes/latest", handlers.GetLatestQuote)
	mux.HandleFunc("GET /api/quotes/{id}", handlers.GetById)

	return mux
}
