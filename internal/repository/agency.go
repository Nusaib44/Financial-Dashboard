package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AgencyEntity struct {
	ID           string
	OwnerUserID  string
	Name         string
	BaseCurrency string
	CreatedAt    time.Time
}

type AgencyRepository interface {
	Create(userID string, name string, currency string, startingCash float64) error
	GetByUserID(userID string) (*AgencyEntity, error)
}

type postgresAgencyRepository struct {
	db *sql.DB
}

func NewAgencyRepository(db *sql.DB) AgencyRepository {
	return &postgresAgencyRepository{db: db}
}

func (r *postgresAgencyRepository) Create(userID string, name string, currency string, startingCash float64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	agencyID := uuid.New().String()
	_, err = tx.Exec(`
		INSERT INTO agencies (id, owner_user_id, name, base_currency, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, agencyID, userID, name, currency, time.Now())

	if err != nil {
		return fmt.Errorf("failed to insert agency: %v", err)
	}

	snapshotID := uuid.New().String()
	_, err = tx.Exec(`
		INSERT INTO cash_snapshots (id, agency_id, date, cash_balance, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, snapshotID, agencyID, time.Now(), startingCash, time.Now())

	if err != nil {
		return fmt.Errorf("failed to insert cash snapshot: %v", err)
	}

	return tx.Commit()
}

func (r *postgresAgencyRepository) GetByUserID(userID string) (*AgencyEntity, error) {
	row := r.db.QueryRow(`
		SELECT id, name, base_currency, created_at 
		FROM agencies 
		WHERE owner_user_id = $1
	`, userID)

	var a AgencyEntity
	err := row.Scan(&a.ID, &a.Name, &a.BaseCurrency, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}
