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

	api.GET("/cash-snapshot/today", handlers.GetTodaysCash(database))
	api.POST("/cash-snapshot", handlers.RecordDailyCash(database))

	api.POST("/revenue", handlers.AddRevenue(database))
	api.POST("/cost", handlers.AddCost(database))
	api.GET("/daily-summary/today", handlers.GetDailySummary(database))

	api.GET("/burn-runway", handlers.GetBurnRunway(database))

	api.POST("/clients", handlers.CreateClient(database))
	api.GET("/clients", handlers.GetClients(database))
	api.POST("/retainers", handlers.CreateRetainer(database))
	api.GET("/retainer-summary", handlers.GetRetainerSummary(database))

	api.POST("/time-entry", handlers.AddTimeEntry(database))
	api.GET("/utilization", handlers.GetUtilization(database))

	api.GET("/agency-reality-score", handlers.GetAgencyRealityScore(database))

	return r
}
