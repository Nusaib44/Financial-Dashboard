package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func Connect(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return db, nil
}

func EnsureUser(db *sql.DB, id string, email string) error {
	// Simple upsert-like check.
	// We trust ID is UUID from Cloudflare (or we treat it as text in DB? schema said UUID).
	// If CF sub is not uuid, we might have issues.
	// But usually CF Access 'sub' is a UUID.

	// Check if user exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM founders WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err = db.Exec("INSERT INTO founders (id, email, created_at) VALUES ($1, $2, $3)", id, email, time.Now())
		if err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}
	}
	return nil
}

type Agency struct {
	ID           string    `json:"id"`
	OwnerUserID  string    `json:"-"`
	Name         string    `json:"name"`
	BaseCurrency string    `json:"base_currency"`
	CreatedAt    time.Time `json:"created_at"`
}

func CreateAgency(db *sql.DB, userID string, name string, currency string, startingCash float64) error {
	// Transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Create Agency
	agencyID := uuid.New().String()
	_, err = tx.Exec(`
		INSERT INTO agencies (id, owner_user_id, name, base_currency, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, agencyID, userID, name, currency, time.Now())

	if err != nil {
		return fmt.Errorf("failed to insert agency: %v", err)
	}

	// 2. Create Cash Snapshot
	snapshotID := uuid.New().String()
	today := time.Now() // Date only used in DB, passed as timestamp usually or string
	_, err = tx.Exec(`
		INSERT INTO cash_snapshots (id, agency_id, date, cash_balance, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, snapshotID, agencyID, today, startingCash, time.Now())

	if err != nil {
		return fmt.Errorf("failed to insert cash snapshot: %v", err)
	}

	return tx.Commit()
}

func GetAgencyByUserID(db *sql.DB, userID string) (*Agency, error) {
	row := db.QueryRow(`
		SELECT id, name, base_currency, created_at 
		FROM agencies 
		WHERE owner_user_id = $1
	`, userID)

	var a Agency
	err := row.Scan(&a.ID, &a.Name, &a.BaseCurrency, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
