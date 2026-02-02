package handlers

import (
	"database/sql"
	"net/http"

	"github.com/agency-finance-reality/server/internal/db"
	"github.com/gin-gonic/gin"
)

type AddRevenueRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Source string  `json:"source" binding:"required"`
}

type AddCostRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
	Type   string  `json:"type" binding:"required,oneof=fixed variable"`
	Label  string  `json:"label" binding:"required"`
}

func AddRevenue(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		agency, err := db.GetAgencyByUserID(database, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Agency not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		var req AddRevenueRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}

		if err := db.AddRevenue(database, agency.ID, req.Amount, req.Source); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add revenue"})
			return
		}

		c.Status(http.StatusCreated)
	}
}

func AddCost(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		agency, err := db.GetAgencyByUserID(database, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Agency not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		var req AddCostRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}

		if err := db.AddCost(database, agency.ID, req.Amount, req.Type, req.Label); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add cost"})
			return
		}

		c.Status(http.StatusCreated)
	}
}

func GetDailySummary(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		agency, err := db.GetAgencyByUserID(database, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Agency not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		summary, err := db.GetDailySummary(database, agency.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, summary)
	}
}
