package transport

import (
	"PlataTest/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type Service interface {
	CreateQuote(ctx context.Context, pair string) (int, error)
	GetLatestQuote(ctx context.Context, pair string) (*domain.Quote, error)
	GetQuote(ctx context.Context, id int) (*domain.Quote, error)
}

type Handlers struct {
	service Service
}
type IDJSON struct {
	ID int `json:"id"`
}

func NewHandler(service Service) *Handlers {
	return &Handlers{service: service}
}

type errorResponse struct {
	Message string `json:"error"`
}
type Request struct {
	Pair string `json:"pair"`
}

func (h *Handlers) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Pair string `json:"pair"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	id, err := h.service.CreateQuote(r.Context(), req.Pair)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPairFormat) {
			writeJSONError(w, http.StatusBadRequest, "invalid format")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		log.Printf("failed to create quote, err %v", err)
		return
	}

	writeJSON(w, http.StatusCreated, IDJSON{ID: id})
}

func (h *Handlers) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	quote, err := h.service.GetQuote(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "quote not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		log.Printf("failed to get quote, err %v", err)
		return
	}
	writeJSON(w, http.StatusOK, quote)
}

func (h *Handlers) GetLatestQuote(w http.ResponseWriter, r *http.Request) {
	pair := r.URL.Query().Get("pair")
	if pair == "" {
		writeJSONError(w, http.StatusBadRequest, "pair is empty")
		return
	}
	quote, err := h.service.GetLatestQuote(r.Context(), pair)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "quotes not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "internal error")
		log.Printf("failed to find quote, err %v", err)
		return
	}
	writeJSON(w, http.StatusOK, quote)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to write response: %v", err)
	}
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Message: msg})
}
