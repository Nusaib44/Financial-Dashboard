package handlers

import (
	"database/sql"
	"net/http"

	"github.com/agency-finance-reality/server/internal/db"
	"github.com/gin-gonic/gin"
)

type AddTimeEntryRequest struct {
	ClientID *string `json:"client_id"`
	Hours    float64 `json:"hours" binding:"required,gt=0"`
}

func AddTimeEntry(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		agency, err := db.GetAgencyByUserID(database, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		var req AddTimeEntryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
			return
		}

		err = db.AddTimeEntry(database, agency.ID, req.ClientID, req.Hours)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add time entry"})
			return
		}

		c.Status(http.StatusCreated)
	}
}

func GetUtilization(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		agency, err := db.GetAgencyByUserID(database, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		utilization, err := db.GetUtilization(database, agency.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, utilization)
	}
}
