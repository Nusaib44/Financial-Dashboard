package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func RunMigrations(db *sql.DB) error {
	// Simple migration runner: reads all .sql files in migrations/ and executes them.
	// In a real app we'd track versions in a table. For this "skeleton", re-running valid idempotent SQL (CREATE IF NOT EXISTS) is fine.

	// Assuming migrations dir is passed or we find it relative to cwd
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %v", err)
	}

	// Sort by name
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, f := range files {
		if filepath.Ext(f.Name()) == ".sql" {
			log.Printf("Applying migration: %s", f.Name())
			content, err := os.ReadFile(filepath.Join("migrations", f.Name()))
			if err != nil {
				return fmt.Errorf("failed to read migration %s: %v", f.Name(), err)
			}

			if _, err := db.Exec(string(content)); err != nil {
				return fmt.Errorf("failed to execute migration %s: %v", f.Name(), err)
			}
		}
	}
	return nil
}
