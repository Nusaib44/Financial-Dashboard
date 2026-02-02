package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/agency-finance-reality/server/internal/db"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func AuthMiddleware(database *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("CF-Access-Jwt-Assertion")
			if tokenString == "" {
				// Local dev bypass if needed, but per requirements "If token missing -> 401"
				// I'll stick to strict requirements, but maybe allow a specific dev header/logic if I was asked.
				// Requirements say "Auth = gate, not logic".
				// But implementation plan mentioned mock auth.
				// Let's support a "Authorization: Bearer mock_sub:email" for local dev if CF header missing?
				// Or better, just manually set the CF header in curl/frontend for local dev.
				http.Error(w, "Missing authentication token", http.StatusUnauthorized)
				return
			}

			// Parse without verification (trusted header from CF)
			token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			sub, _ := claims["sub"].(string)
			email, _ := claims["email"].(string)

			if sub == "" || email == "" {
				http.Error(w, "Missing user info in token", http.StatusUnauthorized)
				return
			}

			// User Auto-Provisioning
			if err := db.EnsureUser(database, sub, email); err != nil {
				// Log error?
				fmt.Printf("Failed to provision user: %v\n", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
