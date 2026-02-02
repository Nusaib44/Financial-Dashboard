package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = 'agencies'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Schema for 'users' table:")
	for rows.Next() {
		var name, dtype string
		if err := rows.Scan(&name, &dtype); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- %s: %s\n", name, dtype)
	}
}
