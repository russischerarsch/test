package http

import "github.com/gin-gonic/gin"

func SetupRouter(handlers *Handlers) *gin.Engine {
	r := gin.Default()
	r.POST("/api/quotes", handlers.CreateQuoteHandler)
	r.GET("/api/quotes/latest", handlers.GetLatestQuote)
	r.GET("/api/quotes/:id", handlers.GetByIdHandler)
	return r
}
