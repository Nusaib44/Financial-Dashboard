package http

import (
	"database/sql"

	"github.com/agency-finance-reality/server/internal/handlers"
	"github.com/gin-gonic/gin"
)

func NewRouter(database *sql.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// Public
	r.GET("/health", handlers.Health)

	// Private
	api := r.Group("/")
	api.Use(AuthMiddleware(database))

	api.GET("/agency", handlers.GetAgency(database))
	api.POST("/agency", handlers.CreateAgency(database))

	return r
}
