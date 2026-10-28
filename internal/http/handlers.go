package http

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

func CreateHandler(service Service) *Handlers {
	return &Handlers{service: service}
}

type errorResponse struct {
	Message string `json:"error"`
}

func (h *Handlers) CreateQuote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method now allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Pair string `json:"pair"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to parse pair, err: "+err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.service.CreateQuote(r.Context(), req.Pair)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPairFormat) {
			http.Error(w, "invalid pair format, err: "+err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to create quote, err: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(id); err != nil {
		log.Printf("failed to send id, err %v", err)
	}
}

func (h *Handlers) GetById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method now allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "failed to convert path id into int, err "+err.Error(), http.StatusBadRequest)
		return
	}
	quote, err := h.service.GetQuote(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			http.Error(w, "quote not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to find quote, err "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(quote); err != nil {
		log.Printf("failed to send quote, err %v", err)
	}
}

func (h *Handlers) GetLatestQuote(w http.ResponseWriter, r *http.Request) {
	pair := r.URL.Query().Get("pair")
	if pair == "" {
		writeJSONError(w, http.StatusBadRequest, "path is empty")
		return
	}
	quote, err := h.service.GetLatestQuote(r.Context(), pair)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			writeJSONError(w, http.StatusNotFound, "quoutes not found")
			log.Printf("quote not found, err %v", err)
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "failed to find quote")
		log.Printf("failed to find quote, err %v", err)
		return
	}
	writeJSON(w, http.StatusOK)
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
