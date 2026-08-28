package http

import (
	apperrors "PlataTest/app_errors"
	"PlataTest/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *service.QuoteService
}

func CreateHandler(service *service.QuoteService) *Handlers {
	return &Handlers{service: service}
}

func (h *Handlers) CreateQuoteHandler(c *gin.Context) {
	var req struct {
		Pair string `json:"pair"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quoteId, err := h.service.CreateQuote(c.Request.Context(), req.Pair)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidPairFormat) {
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
		if errors.Is(err, apperrors.ErrNotFound) {
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
		if errors.Is(err, apperrors.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no quote avaiable"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, quote)
}
