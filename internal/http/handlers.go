package http

import (
	"PlataTest/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *Handlers) CreateQuote(w http.ResponseWriter, r *http.Request) {
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
	r.Header.Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(id); err != nil {
		log.Printf("failed to send id, err %v", err)
	}
}

func (h *Handlers) CreateQuoteHandler(c *gin.Context) {

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quoteId, err := h.service.CreateQuote(c.Request.Context(), req.Pair)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPairFormat) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": quoteId})
}

func (h *Handlers) GetByIdHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	quote, err := h.service.GetQuote(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "quote not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, quote)
}

func (h *Handlers) GetLatestQuote(c *gin.Context) {
	pair := c.Query("pair")
	if pair == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pair is required"})
		return
	}
	quote, err := h.service.GetLatestQuote(c.Request.Context(), pair)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no quote avaiable"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, quote)
}
