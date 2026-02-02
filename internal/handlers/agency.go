package handlers

import (
	"database/sql"
	"net/http"

	"github.com/agency-finance-reality/server/internal/db"
	"github.com/gin-gonic/gin"
)

type CreateAgencyRequest struct {
	Name         string  `json:"name"`
	BaseCurrency string  `json:"base_currency"`
	StartingCash float64 `json:"starting_cash"`
}

func CreateAgency(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		// Check if agency exists
		existing, err := db.GetAgencyByUserID(database, userID)
		if err == nil && existing != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Agency already exists"})
			return
		} else if err != nil && err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		var req CreateAgencyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := db.CreateAgency(database, userID, req.Name, req.BaseCurrency, req.StartingCash); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.Status(http.StatusCreated)
	}
}

func GetAgency(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		agency, err := db.GetAgencyByUserID(database, userID)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agency not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, agency)
	}
}
