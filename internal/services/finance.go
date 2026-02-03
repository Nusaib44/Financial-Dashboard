package services

import (
	"time"

	"github.com/agency-finance-reality/server/internal/models"
	"github.com/agency-finance-reality/server/internal/repository"
)

type FinanceService interface {
	RecordCashSnapshot(agencyID string, cashBalance float64) error
	GetDailySnapshot(agencyID string) (*models.DailySnapshotView, error)
	AddRevenue(agencyID string, amount float64, source string) error
	AddCost(agencyID string, amount float64, costType string, label string, category string) error
	GetDailySummary(agencyID string) (*models.DailySummaryView, error)
	GetSurvivalMetrics(agencyID string) (*models.SurvivalMetricsView, error)
	GetRealityScore(agencyID string) (*models.RealityScoreView, error)
	GetCostBreakdown(agencyID string) (*models.CostBreakdownView, error)
}

type financeService struct {
	cashRepo     repository.CashSnapshotRepository
	financeRepo  repository.FinanceRepository
	retainerRepo repository.RetainerRepository
	timeRepo     repository.TimeEntryRepository
}

func NewFinanceService(
	cashRepo repository.CashSnapshotRepository,
	financeRepo repository.FinanceRepository,
	retainerRepo repository.RetainerRepository,
	timeRepo repository.TimeEntryRepository,
) FinanceService {
	return &financeService{
		cashRepo:     cashRepo,
		financeRepo:  financeRepo,
		retainerRepo: retainerRepo,
		timeRepo:     timeRepo,
	}
}

func (s *financeService) RecordCashSnapshot(agencyID string, cashBalance float64) error {
	return s.cashRepo.CreateDaily(agencyID, cashBalance)
}

func (s *financeService) GetDailySnapshot(agencyID string) (*models.DailySnapshotView, error) {
	snap, err := s.cashRepo.GetToday(agencyID)
	if err != nil || snap == nil {
		return nil, err
	}

	view := &models.DailySnapshotView{
		Date:        snap.Date,
		CashBalance: snap.CashBalance,
	}

	prev, err := s.cashRepo.GetLatestBefore(agencyID, snap.Date)
	if err == nil && prev != nil {
		view.PreviousCashBalance = prev
		d := snap.CashBalance - *prev
		view.Delta = &d
	}

	return view, nil
}

func (s *financeService) AddRevenue(agencyID string, amount float64, source string) error {
	return s.financeRepo.AddRevenue(agencyID, amount, source)
}

func (s *financeService) AddCost(agencyID string, amount float64, costType string, label string, category string) error {
	return s.financeRepo.AddCost(agencyID, amount, costType, label, category)
}

func (s *financeService) GetDailySummary(agencyID string) (*models.DailySummaryView, error) {
	today := time.Now().Format("2006-01-02")
	rev, err := s.financeRepo.SumRevenues(agencyID, today)
	if err != nil {
		return nil, err
	}
	cost, err := s.financeRepo.SumCosts(agencyID, today)
	if err != nil {
		return nil, err
	}

	return &models.DailySummaryView{
		Date:    today,
		Revenue: rev,
		Costs:   cost,
		Net:     rev - cost,
	}, nil
}

func (s *financeService) GetSurvivalMetrics(agencyID string) (*models.SurvivalMetricsView, error) {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	cash, err := s.cashRepo.GetLatest(agencyID)
	if err != nil || cash == nil {
		return nil, err
	}

	burn, err := s.financeRepo.SumFixedCostsInRange(agencyID, thirtyDaysAgo)
	if err != nil {
		return nil, err
	}

	retainers, err := s.retainerRepo.SumActiveRetainers(agencyID)
	if err != nil {
		return nil, err
	}

	view := &models.SurvivalMetricsView{
		CashBalance:     *cash,
		MonthlyBurn:     burn,
		TotalRetainers:  retainers,
		OperatingMargin: retainers - burn,
	}

	if burn > 0 {
		runway := *cash / burn
		runway = float64(int(runway*10)) / 10
		view.RunwayMonths = &runway
	}

	return view, nil
}

func (s *financeService) GetRealityScore(agencyID string) (*models.RealityScoreView, error) {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	var result models.RealityScoreView

	// A. Retainer Safety (25 pts)
	var totalRetainer, fixedCosts float64
	totalRetainer, _ = s.retainerRepo.SumActiveRetainers(agencyID)
	fixedCosts, _ = s.financeRepo.SumFixedCostsInRange(agencyID, thirtyDaysAgo)
	if fixedCosts > 0 {
		coverage := totalRetainer / fixedCosts
		if coverage >= 1.5 {
			result.Breakdown.RetainerSafety = 25
		} else if coverage >= 1.2 {
			result.Breakdown.RetainerSafety = 20
		} else if coverage >= 1.0 {
			result.Breakdown.RetainerSafety = 15
		} else if coverage >= 0.8 {
			result.Breakdown.RetainerSafety = 10
		}
	}

	// B. Runway Health (20 pts)
	cash, _ := s.cashRepo.GetLatest(agencyID)
	if cash != nil && fixedCosts > 0 {
		runway := *cash / fixedCosts
		if runway >= 6 {
			result.Breakdown.Runway = 20
		} else if runway >= 4 {
			result.Breakdown.Runway = 15
		} else if runway >= 2 {
			result.Breakdown.Runway = 8
		} else if runway >= 1 {
			result.Breakdown.Runway = 4
		}
	}

	// C. Client Concentration (20 pts)
	if totalRetainer > 0 {
		maxRetainer, _ := s.retainerRepo.GetMaxRetainer(agencyID)
		topPct := (maxRetainer / totalRetainer) * 100
		if topPct < 30 {
			result.Breakdown.ClientConcentration = 20
		} else if topPct < 40 {
			result.Breakdown.ClientConcentration = 15
		} else if topPct < 50 {
			result.Breakdown.ClientConcentration = 8
		} else if topPct < 60 {
			result.Breakdown.ClientConcentration = 4
		}
	}

	// D. Profitability (20 pts) - 30-day rolling window (Ticket 09 fix)
	rev, _ := s.financeRepo.SumAllRevenuesInRange(agencyID, thirtyDaysAgo)
	costs, _ := s.financeRepo.SumAllCostsInRange(agencyID, thirtyDaysAgo)
	if rev > 0 {
		margin := ((rev - costs) / rev) * 100
		if margin >= 20 {
			result.Breakdown.Profitability = 20
		} else if margin >= 10 {
			result.Breakdown.Profitability = 15
		} else if margin >= 0 {
			result.Breakdown.Profitability = 8
		} else if margin >= -10 {
			result.Breakdown.Profitability = 4
		}
	}

	// E. Capacity Pressure (15 pts)
	usedHours, _ := s.timeRepo.SumHoursInRange(agencyID, thirtyDaysAgo)
	utilization := (usedHours / 160.0) * 100 // Hardcoded capacity 160
	if utilization >= 60 && utilization <= 85 {
		result.Breakdown.CapacityPressure = 15
	} else if utilization >= 50 && utilization < 60 {
		result.Breakdown.CapacityPressure = 10
	} else if utilization >= 40 && utilization < 50 {
		result.Breakdown.CapacityPressure = 6
	} else if utilization > 85 && utilization <= 100 {
		result.Breakdown.CapacityPressure = 6
	}

	result.Score = result.Breakdown.RetainerSafety + result.Breakdown.Runway + result.Breakdown.ClientConcentration + result.Breakdown.Profitability + result.Breakdown.CapacityPressure
	if result.Score >= 80 {
		result.Status = "Healthy"
	} else if result.Score >= 60 {
		result.Status = "Watch"
	} else if result.Score >= 40 {
		result.Status = "At Risk"
	} else {
		result.Status = "Danger"
	}

	if cash != nil {
		result.CashOnHand = *cash
	}
	result.CommittedRetainers = totalRetainer

	// Primary Risk Attribution (Ticket 11)
	if result.Score >= 80 {
		result.PrimaryRisk = "Healthy"
	} else {
		// Priority order
		if fixedCosts > totalRetainer && fixedCosts > 0 {
			result.PrimaryRisk = "High Fixed Costs"
		} else if totalRetainer/fixedCosts < 1.0 {
			result.PrimaryRisk = "Low Retainer Base"
		} else {
			maxRetainer, _ := s.retainerRepo.GetMaxRetainer(agencyID)
			topPct := 0.0
			if totalRetainer > 0 {
				topPct = (maxRetainer / totalRetainer) * 100
			}
			if topPct > 60 {
				result.PrimaryRisk = "Client Concentration"
			} else if cash != nil && fixedCosts > 0 && (*cash/fixedCosts) < 2 {
				result.PrimaryRisk = "Low Runway"
			} else {
				result.PrimaryRisk = "Healthy"
			}
		}
	}

	return &result, nil
}

func (s *financeService) GetCostBreakdown(agencyID string) (*models.CostBreakdownView, error) {
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	breakdown, err := s.financeRepo.GetGroupedFixedCosts(agencyID, thirtyDaysAgo)
	if err != nil {
		return nil, err
	}

	total := 0.0
	for _, amt := range breakdown {
		total += amt
	}

	view := &models.CostBreakdownView{
		TotalFixedCosts: total,
		Breakdown:       breakdown,
	}

	if total > 0 {
		maxAmt := -1.0
		maxCat := "other"
		for cat, amt := range breakdown {
			if amt > maxAmt {
				maxAmt = amt
				maxCat = cat
			}
		}
		view.PrimaryDriver = models.CostDriver{
			Category:   maxCat,
			Amount:     maxAmt,
			Percentage: float64(int((maxAmt/total)*1000)) / 10,
		}
	} else {
		view.PrimaryDriver = models.CostDriver{Category: "other", Amount: 0, Percentage: 0}
	}

	return view, nil
}
