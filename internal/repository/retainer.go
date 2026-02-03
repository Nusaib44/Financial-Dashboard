package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

type RetainerRepository interface {
	Create(agencyID string, clientID string, amount float64) error
	SumActiveRetainers(agencyID string) (float64, error)
	GetMaxRetainer(agencyID string) (float64, error)
	HasActiveRetainer(clientID string) (bool, error)
}

type postgresRetainerRepository struct {
	db *sql.DB
}

func NewRetainerRepository(db *sql.DB) RetainerRepository {
	return &postgresRetainerRepository{db: db}
}

func (r *postgresRetainerRepository) Create(agencyID string, clientID string, amount float64) error {
	id := uuid.New().String()
	_, err := r.db.Exec(`
		INSERT INTO retainers (id, agency_id, client_id, monthly_amount, active)
		VALUES ($1, $2, $3, $4, true)
	`, id, agencyID, clientID, amount)
	return err
}

func (r *postgresRetainerRepository) SumActiveRetainers(agencyID string) (float64, error) {
	var total float64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(monthly_amount), 0) FROM retainers 
		WHERE agency_id = $1 AND active = true
	`, agencyID).Scan(&total)
	return total, err
}

func (r *postgresRetainerRepository) GetMaxRetainer(agencyID string) (float64, error) {
	var max float64
	err := r.db.QueryRow(`
		SELECT COALESCE(MAX(monthly_amount), 0) FROM retainers 
		WHERE agency_id = $1 AND active = true
	`, agencyID).Scan(&max)
	return max, err
}

func (r *postgresRetainerRepository) HasActiveRetainer(clientID string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM retainers WHERE client_id = $1 AND active = true)
	`, clientID).Scan(&exists)
	return exists, err
}
