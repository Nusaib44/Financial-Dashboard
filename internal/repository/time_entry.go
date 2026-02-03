package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type TimeEntryRepository interface {
	Add(agencyID string, clientID *string, hours float64) error
	SumHoursInRange(agencyID string, startDate string) (float64, error)
}

type postgresTimeEntryRepository struct {
	db *sql.DB
}

func NewTimeEntryRepository(db *sql.DB) TimeEntryRepository {
	return &postgresTimeEntryRepository{db: db}
}

func (r *postgresTimeEntryRepository) Add(agencyID string, clientID *string, hours float64) error {
	id := uuid.New().String()
	date := time.Now().Format("2006-01-02")
	_, err := r.db.Exec(`
		INSERT INTO time_entries (id, agency_id, client_id, date, hours)
		VALUES ($1, $2, $3, $4, $5)
	`, id, agencyID, clientID, date, hours)
	return err
}

func (r *postgresTimeEntryRepository) SumHoursInRange(agencyID string, startDate string) (float64, error) {
	var total float64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(hours), 0) FROM time_entries 
		WHERE agency_id = $1 AND date >= $2
	`, agencyID, startDate).Scan(&total)
	return total, err
}
