package services

import (
	"time"

	"github.com/agency-finance-reality/server/internal/models"
	"github.com/agency-finance-reality/server/internal/repository"
)

type UtilizationService interface {
	AddTimeEntry(agencyID string, clientID *string, hours float64) error
	GetUtilization(agencyID string) (*models.UtilizationView, error)
}

type utilizationService struct {
	timeRepo repository.TimeEntryRepository
}

func NewUtilizationService(timeRepo repository.TimeEntryRepository) UtilizationService {
	return &utilizationService{timeRepo: timeRepo}
}

func (s *utilizationService) AddTimeEntry(agencyID string, clientID *string, hours float64) error {
	return s.timeRepo.Add(agencyID, clientID, hours)
}

func (s *utilizationService) GetUtilization(agencyID string) (*models.UtilizationView, error) {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	used, err := s.timeRepo.SumHoursInRange(agencyID, thirtyDaysAgo)
	if err != nil {
		return nil, err
	}

	capacity := 160.0 // Hardcoded capacity
	view := &models.UtilizationView{
		UsedHours:     used,
		CapacityHours: capacity,
	}

	if capacity > 0 {
		percent := (used / capacity) * 100
		view.UtilizationPercent = float64(int(percent*10)) / 10
	}

	return view, nil
}
