package services

import (
	"fmt"
	"time"

	"github.com/agency-finance-reality/server/internal/models"
	"github.com/agency-finance-reality/server/internal/repository"
)

type ClientService interface {
	CreateClient(agencyID string, name string) (*models.ClientView, error)
	GetClients(agencyID string) ([]models.ClientView, error)
	CreateRetainer(agencyID string, clientID string, amount float64) error
	GetRetainerSummary(agencyID string) (*models.RetainerSummaryView, error)
}

type clientService struct {
	clientRepo   repository.ClientRepository
	retainerRepo repository.RetainerRepository
	financeRepo  repository.FinanceRepository
}

func NewClientService(
	clientRepo repository.ClientRepository,
	retainerRepo repository.RetainerRepository,
	financeRepo repository.FinanceRepository,
) ClientService {
	return &clientService{
		clientRepo:   clientRepo,
		retainerRepo: retainerRepo,
		financeRepo:  financeRepo,
	}
}

func (s *clientService) CreateClient(agencyID string, name string) (*models.ClientView, error) {
	client, err := s.clientRepo.Create(agencyID, name)
	if err != nil {
		return nil, err
	}
	return &models.ClientView{
		ID:     client.ID,
		Name:   client.Name,
		Status: client.Status,
	}, nil
}

func (s *clientService) GetClients(agencyID string) ([]models.ClientView, error) {
	entities, err := s.clientRepo.GetAllActive(agencyID)
	if err != nil {
		return nil, err
	}
	views := make([]models.ClientView, len(entities))
	for i, e := range entities {
		views[i] = models.ClientView{
			ID:     e.ID,
			Name:   e.Name,
			Status: e.Status,
		}
	}
	return views, nil
}

func (s *clientService) CreateRetainer(agencyID string, clientID string, amount float64) error {
	exists, err := s.retainerRepo.HasActiveRetainer(clientID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("client already has active retainer")
	}
	return s.retainerRepo.Create(agencyID, clientID, amount)
}

func (s *clientService) GetRetainerSummary(agencyID string) (*models.RetainerSummaryView, error) {
	total, err := s.retainerRepo.SumActiveRetainers(agencyID)
	if err != nil {
		return nil, err
	}

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	fixed, err := s.financeRepo.SumFixedCostsInRange(agencyID, thirtyDaysAgo)
	if err != nil {
		return nil, err
	}

	view := &models.RetainerSummaryView{
		TotalRetainerRevenue: total,
		FixedCosts:           fixed,
	}

	if fixed > 0 {
		view.CoverageRatio = float64(int((total/fixed)*100)) / 100
	}

	if total > 0 {
		max, err := s.retainerRepo.GetMaxRetainer(agencyID)
		if err == nil {
			view.TopClientPercentage = float64(int((max/total)*100)) / 100
		}
	}

	return view, nil
}
