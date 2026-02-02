package handlers

import (
	"database/sql"
	"net/http"

	"github.com/agency-finance-reality/server/internal/db"
	"github.com/gin-gonic/gin"
)

type CreateClientRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateRetainerRequest struct {
	ClientID      string  `json:"client_id" binding:"required"`
	MonthlyAmount float64 `json:"monthly_amount" binding:"required,gt=0"`
}

func CreateClient(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		agency, err := db.GetAgencyByUserID(database, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		var req CreateClientRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		client, err := db.CreateClient(database, agency.ID, req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create client"})
			return
		}

		c.JSON(http.StatusCreated, client)
	}
}

func GetClients(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		agency, err := db.GetAgencyByUserID(database, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		clients, err := db.GetClients(database, agency.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, clients)
	}
}

func CreateRetainer(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		agency, err := db.GetAgencyByUserID(database, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		var req CreateRetainerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		err = db.CreateRetainer(database, agency.ID, req.ClientID, req.MonthlyAmount)
		if err != nil {
			if err.Error() == "client already has active retainer" {
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create retainer"})
			return
		}

		c.Status(http.StatusCreated)
	}
}

func GetRetainerSummary(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		agency, err := db.GetAgencyByUserID(database, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		summary, err := db.GetRetainerSummary(database, agency.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.JSON(http.StatusOK, summary)
	}
}
