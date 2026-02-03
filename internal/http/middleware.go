package http

import (
	"net/http"

	"github.com/agency-finance-reality/server/internal/auth"
	"github.com/agency-finance-reality/server/internal/services"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("CF-Access-Jwt-Assertion")

		// Validate using auth package
		sub, email, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// User Auto-Provisioning
		if err := authService.EnsureUser(sub, email); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		c.Set("user_id", sub)
		c.Set("email", email)
		c.Next()
	}
}
