package http

import (
	"database/sql"

	"github.com/agency-finance-reality/server/internal/handlers"
	"github.com/agency-finance-reality/server/internal/repository"
	"github.com/agency-finance-reality/server/internal/services"
	"github.com/gin-gonic/gin"
)

func NewRouter(db *sql.DB) *gin.Engine {
	// Repositories
	founderRepo := repository.NewFounderRepository(db)
	agencyRepo := repository.NewAgencyRepository(db)
	cashRepo := repository.NewCashSnapshotRepository(db)
	financeRepo := repository.NewFinanceRepository(db)
	clientRepo := repository.NewClientRepository(db)
	retainerRepo := repository.NewRetainerRepository(db)
	timeRepo := repository.NewTimeEntryRepository(db)

	// Services
	authService := services.NewAuthService(founderRepo)
	agencyService := services.NewAgencyService(agencyRepo)
	financeService := services.NewFinanceService(cashRepo, financeRepo, retainerRepo, timeRepo)
	clientService := services.NewClientService(clientRepo, retainerRepo, financeRepo)
	utilizationService := services.NewUtilizationService(timeRepo)

	// Handlers
	agencyHandler := handlers.NewAgencyHandler(agencyService)
	cashHandler := handlers.NewCashSnapshotHandler(agencyService, financeService)
	financeHandler := handlers.NewDailyFinanceHandler(agencyService, financeService)
	realityScoreHandler := handlers.NewRealityScoreHandler(agencyService, financeService)
	clientHandler := handlers.NewClientHandler(agencyService, clientService)
	retainerHandler := handlers.NewRetainerHandler(agencyService, clientService)
	utilizationHandler := handlers.NewUtilizationHandler(agencyService, utilizationService)
	survivalHandler := handlers.NewSurvivalHandler(agencyService, financeService)

	r := gin.New()
	r.Use(gin.Recovery())

	// Public
	r.GET("/health", handlers.Health)

	// Private
	api := r.Group("/")
	api.Use(AuthMiddleware(authService))

	api.GET("/agency", agencyHandler.GetAgency)
	api.POST("/agency", agencyHandler.CreateAgency)

	api.GET("/cash-snapshot/today", cashHandler.GetTodaysCash)
	api.POST("/cash-snapshot", cashHandler.RecordDailyCash)

	api.POST("/revenue", financeHandler.AddRevenue)
	api.POST("/cost", financeHandler.AddCost)
	api.GET("/daily-summary/today", financeHandler.GetDailySummary)
	api.GET("/cost-breakdown", financeHandler.GetCostBreakdown)

	api.GET("/burn-runway", survivalHandler.GetBurnRunway)

	api.POST("/clients", clientHandler.CreateClient)
	api.GET("/clients", clientHandler.GetClients)
	api.POST("/retainers", retainerHandler.CreateRetainer)
	api.GET("/retainer-summary", retainerHandler.GetRetainerSummary)

	api.POST("/time-entry", utilizationHandler.AddTimeEntry)
	api.GET("/utilization", utilizationHandler.GetUtilization)

	api.GET("/agency-reality-score", realityScoreHandler.GetRealityScore)

	return r
}
