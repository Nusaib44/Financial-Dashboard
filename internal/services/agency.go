package services

import (
	"github.com/agency-finance-reality/server/internal/models"
	"github.com/agency-finance-reality/server/internal/repository"
)

type AgencyService interface {
	CreateAgency(userID string, name string, currency string, startingCash float64) error
	GetAgencyByUserID(userID string) (*models.AgencyView, error)
}

type agencyService struct {
	agencyRepo repository.AgencyRepository
}

func NewAgencyService(agencyRepo repository.AgencyRepository) AgencyService {
	return &agencyService{agencyRepo: agencyRepo}
}

func (s *agencyService) CreateAgency(userID string, name string, currency string, startingCash float64) error {
	return s.agencyRepo.Create(userID, name, currency, startingCash)
}

func (s *agencyService) GetAgencyByUserID(userID string) (*models.AgencyView, error) {
	entity, err := s.agencyRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return &models.AgencyView{
		ID:           entity.ID,
		Name:         entity.Name,
		BaseCurrency: entity.BaseCurrency,
		CreatedAt:    entity.CreatedAt,
	}, nil
}
