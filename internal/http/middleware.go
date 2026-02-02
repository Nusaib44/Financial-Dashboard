package http

import (
	"database/sql"
	"net/http"

	"github.com/agency-finance-reality/server/internal/auth"
	"github.com/agency-finance-reality/server/internal/db"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(database *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("CF-Access-Jwt-Assertion")

		// Validate using auth package
		sub, email, err := auth.ValidateToken(tokenString)
		if err != nil {
			// dev bypass/mock is handled by injecting a valid-ish JWT in the header for now
			// If validation fails (missing or bad structure), we 401.
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// User Auto-Provisioning
		if err := db.EnsureUser(database, sub, email); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.Set("user_id", sub)
		c.Set("email", email)
		c.Next()
	}
}
