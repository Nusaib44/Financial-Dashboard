package handlers

import (
	"database/sql"
	"net/http"

	"github.com/agency-finance-reality/server/internal/db"
	"github.com/gin-gonic/gin"
)

func GetBurnRunway(database *sql.DB) gin.HandlerFunc {
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

		burnRunway, err := db.GetBurnRunway(database, agency.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		if burnRunway == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No cash snapshot available"})
			return
		}

		c.JSON(http.StatusOK, burnRunway)
	}
}
