package repository

import (
	"database/sql"

	"github.com/google/uuid"
)

type ClientEntity struct {
	ID       string
	AgencyID string
	Name     string
	Status   string
}

type ClientRepository interface {
	Create(agencyID string, name string) (*ClientEntity, error)
	GetAllActive(agencyID string) ([]ClientEntity, error)
}

type postgresClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) ClientRepository {
	return &postgresClientRepository{db: db}
}

func (r *postgresClientRepository) Create(agencyID string, name string) (*ClientEntity, error) {
	id := uuid.New().String()
	_, err := r.db.Exec(`
		INSERT INTO clients (id, agency_id, name, status)
		VALUES ($1, $2, $3, 'active')
	`, id, agencyID, name)
	if err != nil {
		return nil, err
	}
	return &ClientEntity{ID: id, AgencyID: agencyID, Name: name, Status: "active"}, nil
}

func (r *postgresClientRepository) GetAllActive(agencyID string) ([]ClientEntity, error) {
	rows, err := r.db.Query(`
		SELECT id, name, status FROM clients 
		WHERE agency_id = $1 AND status = 'active'
		ORDER BY name
	`, agencyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []ClientEntity
	for rows.Next() {
		var c ClientEntity
		c.AgencyID = agencyID
		if err := rows.Scan(&c.ID, &c.Name, &c.Status); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}
