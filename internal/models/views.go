package models

import "time"

// Agency models
type AgencyView struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	BaseCurrency string    `json:"base_currency"`
	CreatedAt    time.Time `json:"created_at"`
}

// Finance models
type DailySnapshotView struct {
	Date                string   `json:"date"`
	CashBalance         float64  `json:"cash_balance"`
	PreviousCashBalance *float64 `json:"previous_cash_balance"`
	Delta               *float64 `json:"delta"`
}

type DailySummaryView struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Costs   float64 `json:"costs"`
	Net     float64 `json:"net"`
}

type CostDriver struct {
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

type CostBreakdownView struct {
	TotalFixedCosts float64            `json:"total_fixed_costs"`
	Breakdown       map[string]float64 `json:"breakdown"`
	PrimaryDriver   CostDriver         `json:"primary_driver"`
}

type SurvivalMetricsView struct {
	CashBalance     float64  `json:"cash_balance"`
	MonthlyBurn     float64  `json:"monthly_burn"`
	RunwayMonths    *float64 `json:"runway_months"`
	OperatingMargin float64  `json:"operating_margin"`
	TotalRetainers  float64  `json:"total_retainers"`
}

// Client models
type ClientView struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type RetainerSummaryView struct {
	TotalRetainerRevenue float64 `json:"total_retainer_revenue"`
	FixedCosts           float64 `json:"fixed_costs"`
	CoverageRatio        float64 `json:"coverage_ratio"`
	TopClientPercentage  float64 `json:"top_client_percentage"`
}

// Utilization models
type UtilizationView struct {
	UsedHours          float64 `json:"used_hours"`
	CapacityHours      float64 `json:"capacity_hours"`
	UtilizationPercent float64 `json:"utilization_percent"`
}

// Reality Score models
type ScoreBreakdownView struct {
	RetainerSafety      int `json:"retainer_safety"`
	Runway              int `json:"runway"`
	ClientConcentration int `json:"client_concentration"`
	Profitability       int `json:"profitability"`
	CapacityPressure    int `json:"capacity_pressure"`
}

type RealityScoreView struct {
	Score              int                `json:"score"`
	Breakdown          ScoreBreakdownView `json:"breakdown"`
	Status             string             `json:"status"`
	CashOnHand         float64            `json:"cash_on_hand"`
	CommittedRetainers float64            `json:"committed_retainers"`
	PrimaryRisk        string             `json:"primary_risk"`
}
