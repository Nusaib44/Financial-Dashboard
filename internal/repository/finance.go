package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type FinanceRepository interface {
	AddRevenue(agencyID string, amount float64, source string) error
	AddCost(agencyID string, amount float64, costType string, label string, category string) error
	SumRevenues(agencyID string, date string) (float64, error)
	SumCosts(agencyID string, date string) (float64, error)
	SumFixedCostsInRange(agencyID string, startDate string) (float64, error)
	SumAllRevenuesInRange(agencyID string, startDate string) (float64, error)
	SumAllCostsInRange(agencyID string, startDate string) (float64, error)
	GetGroupedFixedCosts(agencyID string, startDate string) (map[string]float64, error)
}

type postgresFinanceRepository struct {
	db *sql.DB
}

func NewFinanceRepository(db *sql.DB) FinanceRepository {
	return &postgresFinanceRepository{db: db}
}

func (r *postgresFinanceRepository) AddRevenue(agencyID string, amount float64, source string) error {
	id := uuid.New().String()
	date := time.Now().Format("2006-01-02")
	_, err := r.db.Exec(`
		INSERT INTO daily_revenues (id, agency_id, date, amount, source)
		VALUES ($1, $2, $3, $4, $5)
	`, id, agencyID, date, amount, source)
	return err
}

func (r *postgresFinanceRepository) AddCost(agencyID string, amount float64, costType string, label string, category string) error {
	id := uuid.New().String()
	date := time.Now().Format("2006-01-02")
	_, err := r.db.Exec(`
		INSERT INTO daily_costs (id, agency_id, date, amount, type, label, category)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, id, agencyID, date, amount, costType, label, category)
	return err
}

func (r *postgresFinanceRepository) SumRevenues(agencyID string, date string) (float64, error) {
	var total float64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) FROM daily_revenues 
		WHERE agency_id = $1 AND date = $2
	`, agencyID, date).Scan(&total)
	return total, err
}

func (r *postgresFinanceRepository) SumCosts(agencyID string, date string) (float64, error) {
	var total float64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) FROM daily_costs 
		WHERE agency_id = $1 AND date = $2
	`, agencyID, date).Scan(&total)
	return total, err
}

func (r *postgresFinanceRepository) SumFixedCostsInRange(agencyID string, startDate string) (float64, error) {
	var total float64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) FROM daily_costs 
		WHERE agency_id = $1 AND type = 'fixed' AND date >= $2
	`, agencyID, startDate).Scan(&total)
	return total, err
}

func (r *postgresFinanceRepository) SumAllRevenuesInRange(agencyID string, startDate string) (float64, error) {
	var total float64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) FROM daily_revenues 
		WHERE agency_id = $1 AND date >= $2
	`, agencyID, startDate).Scan(&total)
	return total, err
}

func (r *postgresFinanceRepository) SumAllCostsInRange(agencyID string, startDate string) (float64, error) {
	var total float64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) FROM daily_costs 
		WHERE agency_id = $1 AND date >= $2
	`, agencyID, startDate).Scan(&total)
	return total, err
}

func (r *postgresFinanceRepository) GetGroupedFixedCosts(agencyID string, startDate string) (map[string]float64, error) {
	rows, err := r.db.Query(`
		SELECT category, COALESCE(SUM(amount), 0) FROM daily_costs 
		WHERE agency_id = $1 AND type = 'fixed' AND date >= $2
		GROUP BY category
	`, agencyID, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]float64)
	for rows.Next() {
		var cat string
		var amt float64
		if err := rows.Scan(&cat, &amt); err != nil {
			return nil, err
		}
		result[cat] = amt
	}
	return result, nil
}
