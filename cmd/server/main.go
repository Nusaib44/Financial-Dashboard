package main

import (
	"log"
	"os"

	"github.com/agency-finance-reality/server/internal/db"
	internalHttp "github.com/agency-finance-reality/server/internal/http"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is required")
	}

	conn, err := db.Connect(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	if err := db.RunMigrations(conn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	defer conn.Close()

	router := internalHttp.NewRouter(conn)

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
