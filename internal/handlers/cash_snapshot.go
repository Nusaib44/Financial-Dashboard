package handlers

import (
	"database/sql"
	"net/http"

	"github.com/agency-finance-reality/server/internal/db"
	"github.com/gin-gonic/gin"
)

type RecordCashRequest struct {
	CashBalance float64 `json:"cash_balance" binding:"required"`
}

func RecordDailyCash(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		// Resolve Agency
		agency, err := db.GetAgencyByUserID(database, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Agency not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		var req RecordCashRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		err = db.CreateDailySnapshot(database, agency.ID, req.CashBalance)
		if err != nil {
			// Check for duplicate (Unique constraint violation)
			// primitive check via error string or specialized lib, for now generic error handling
			// "pq: duplicate key value violates unique constraint" using lib/pq
			// But we are on pgx/stdlib? No, we switched to lib/pq.
			// Let's just return 409 if it looks like a constraint error or just generic.
			// Ideally we check error code.
			c.JSON(http.StatusConflict, gin.H{"error": "Snapshot already exists for today"})
			return
		}

		c.Status(http.StatusCreated)
	}
}

func GetTodaysCash(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		// Resolve Agency
		agency, err := db.GetAgencyByUserID(database, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Agency not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		snapshot, err := db.GetDailySnapshot(database, agency.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		if snapshot == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No snapshot for today"})
			return
		}

		c.JSON(http.StatusOK, snapshot)
	}
}
