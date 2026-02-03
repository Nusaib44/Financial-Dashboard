package repository

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CashSnapshotEntity struct {
	Date        string
	CashBalance float64
}

type CashSnapshotRepository interface {
	CreateDaily(agencyID string, cashBalance float64) error
	GetToday(agencyID string) (*CashSnapshotEntity, error)
	GetLatestBefore(agencyID string, date string) (*float64, error)
	GetLatest(agencyID string) (*float64, error)
}

type postgresCashSnapshotRepository struct {
	db *sql.DB
}

func NewCashSnapshotRepository(db *sql.DB) CashSnapshotRepository {
	return &postgresCashSnapshotRepository{db: db}
}

func (r *postgresCashSnapshotRepository) CreateDaily(agencyID string, cashBalance float64) error {
	id := uuid.New().String()
	date := time.Now().Format("2006-01-02")

	_, err := r.db.Exec(`
		INSERT INTO daily_cash_snapshots (id, agency_id, date, cash_balance)
		VALUES ($1, $2, $3, $4)
	`, id, agencyID, date, cashBalance)

	return err
}

func (r *postgresCashSnapshotRepository) GetToday(agencyID string) (*CashSnapshotEntity, error) {
	today := time.Now().Format("2006-01-02")
	var snap CashSnapshotEntity
	err := r.db.QueryRow(`
		SELECT cash_balance 
		FROM daily_cash_snapshots 
		WHERE agency_id = $1 AND date = $2
	`, agencyID, today).Scan(&snap.CashBalance)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	snap.Date = today
	return &snap, nil
}

func (r *postgresCashSnapshotRepository) GetLatestBefore(agencyID string, date string) (*float64, error) {
	var balance float64
	err := r.db.QueryRow(`
		SELECT cash_balance 
		FROM daily_cash_snapshots 
		WHERE agency_id = $1 AND date < $2 
		ORDER BY date DESC 
		LIMIT 1
	`, agencyID, date).Scan(&balance)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (r *postgresCashSnapshotRepository) GetLatest(agencyID string) (*float64, error) {
	var balance float64
	err := r.db.QueryRow(`
		SELECT cash_balance FROM daily_cash_snapshots 
		WHERE agency_id = $1 
		ORDER BY date DESC 
		LIMIT 1
	`, agencyID).Scan(&balance)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &balance, nil
}
