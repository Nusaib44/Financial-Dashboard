package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type FounderRepository interface {
	EnsureUser(id string, email string) error
}

type postgresFounderRepository struct {
	db *sql.DB
}

func NewFounderRepository(db *sql.DB) FounderRepository {
	return &postgresFounderRepository{db: db}
}

func (r *postgresFounderRepository) EnsureUser(id string, email string) error {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM founders WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err = r.db.Exec("INSERT INTO founders (id, email, created_at) VALUES ($1, $2, $3)", id, email, time.Now())
		if err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}
	}
	return nil
}
