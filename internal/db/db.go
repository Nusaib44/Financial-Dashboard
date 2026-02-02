package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func Connect(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	return db, nil
}

func EnsureUser(db *sql.DB, id string, email string) error {
	// Simple upsert-like check.
	// Check if user exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM founders WHERE id=$1)", id).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		_, err = db.Exec("INSERT INTO founders (id, email, created_at) VALUES ($1, $2, $3)", id, email, time.Now())
		if err != nil {
			return fmt.Errorf("failed to create user: %v", err)
		}
	}
	return nil
}

type Agency struct {
	ID           string    `json:"id"`
	OwnerUserID  string    `json:"-"`
	Name         string    `json:"name"`
	BaseCurrency string    `json:"base_currency"`
	CreatedAt    time.Time `json:"created_at"`
}

func CreateAgency(db *sql.DB, userID string, name string, currency string, startingCash float64) error {
	// Transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Create Agency
	agencyID := uuid.New().String()
	_, err = tx.Exec(`
		INSERT INTO agencies (id, owner_user_id, name, base_currency, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, agencyID, userID, name, currency, time.Now())

	if err != nil {
		return fmt.Errorf("failed to insert agency: %v", err)
	}

	// 2. Create Cash Snapshot (First one is also a daily snapshot technically, but Ticket 02 used cash_snapshots)
	// Ticket 03 introduces daily_cash_snapshots and says "One snapshot per agency per date".
	// The requirement implies we should write to daily_cash_snapshots for the "Starting Cash" too?
	// Ticket 02 said "Create initial cash snapshot (today)".
	// Ticket 03 says "daily_cash_snapshots" table.
	// I'll leave Ticket 02 logic as is (cash_snapshots table) to not break it,
	// but I will ALSO write to daily_cash_snapshots for the initial day so the graph/delta works from day 0?
	// Actually, strict scope says "no backfilling".
	// But getting a 404 on day 1 after creating agency feels wrong.
	// I will just let the user "Record Today's Reality" even if they just created the agency.
	// Or should CreateAgency also populate daily_cash_snapshots?
	// Scope says "No backfilling".
	// I'll stick to Ticket 02 logic for CreateAgency (writing to `cash_snapshots`).
	// Ticket 03 logic is separate. The user will record today's cash manually.

	snapshotID := uuid.New().String()
	today := time.Now()
	_, err = tx.Exec(`
		INSERT INTO cash_snapshots (id, agency_id, date, cash_balance, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, snapshotID, agencyID, today, startingCash, time.Now())

	if err != nil {
		return fmt.Errorf("failed to insert cash snapshot: %v", err)
	}

	return tx.Commit()
}

func GetAgencyByUserID(db *sql.DB, userID string) (*Agency, error) {
	row := db.QueryRow(`
		SELECT id, name, base_currency, created_at 
		FROM agencies 
		WHERE owner_user_id = $1
	`, userID)

	var a Agency
	err := row.Scan(&a.ID, &a.Name, &a.BaseCurrency, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

type DailySnapshot struct {
	Date                string   `json:"date"` // stored as YYYY-MM-DD
	CashBalance         float64  `json:"cash_balance"`
	PreviousCashBalance *float64 `json:"previous_cash_balance"`
	Delta               *float64 `json:"delta"`
}

func CreateDailySnapshot(db *sql.DB, agencyID string, cashBalance float64) error {
	id := uuid.New().String()
	date := time.Now().Format("2006-01-02")

	_, err := db.Exec(`
		INSERT INTO daily_cash_snapshots (id, agency_id, date, cash_balance)
		VALUES ($1, $2, $3, $4)
	`, id, agencyID, date, cashBalance)

	return err
}

func GetDailySnapshot(db *sql.DB, agencyID string) (*DailySnapshot, error) {
	// 1. Get Today's Snapshot
	today := time.Now().Format("2006-01-02")

	var snap DailySnapshot
	snap.Date = today

	err := db.QueryRow(`
		SELECT cash_balance 
		FROM daily_cash_snapshots 
		WHERE agency_id = $1 AND date = $2
	`, agencyID, today).Scan(&snap.CashBalance)

	if err == sql.ErrNoRows {
		return nil, nil // Not found today
	} else if err != nil {
		return nil, err
	}

	// 2. Get Previous Snapshot (latest before today)
	var prevBalance float64
	err = db.QueryRow(`
		SELECT cash_balance 
		FROM daily_cash_snapshots 
		WHERE agency_id = $1 AND date < $2 
		ORDER BY date DESC 
		LIMIT 1
	`, agencyID, today).Scan(&prevBalance)

	if err == nil {
		snap.PreviousCashBalance = &prevBalance
		d := snap.CashBalance - prevBalance
		snap.Delta = &d
	} else if err != sql.ErrNoRows {
		return nil, err
	}

	return &snap, nil
}

// Revenue & Cost methods (Ticket 04)

func AddRevenue(db *sql.DB, agencyID string, amount float64, source string) error {
	id := uuid.New().String()
	date := time.Now().Format("2006-01-02")

	_, err := db.Exec(`
		INSERT INTO daily_revenues (id, agency_id, date, amount, source)
		VALUES ($1, $2, $3, $4, $5)
	`, id, agencyID, date, amount, source)

	return err
}

func AddCost(db *sql.DB, agencyID string, amount float64, costType string, label string) error {
	id := uuid.New().String()
	date := time.Now().Format("2006-01-02")

	_, err := db.Exec(`
		INSERT INTO daily_costs (id, agency_id, date, amount, type, label)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, id, agencyID, date, amount, costType, label)

	return err
}

type DailySummary struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Costs   float64 `json:"costs"`
	Net     float64 `json:"net"`
}

func GetDailySummary(db *sql.DB, agencyID string) (*DailySummary, error) {
	today := time.Now().Format("2006-01-02")

	var summary DailySummary
	summary.Date = today

	// Sum revenues for today
	err := db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) FROM daily_revenues 
		WHERE agency_id = $1 AND date = $2
	`, agencyID, today).Scan(&summary.Revenue)
	if err != nil {
		return nil, err
	}

	// Sum costs for today
	err = db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) FROM daily_costs 
		WHERE agency_id = $1 AND date = $2
	`, agencyID, today).Scan(&summary.Costs)
	if err != nil {
		return nil, err
	}

	summary.Net = summary.Revenue - summary.Costs
	return &summary, nil
}

// Burn Rate & Runway (Ticket 05) - Updated Ticket 09

type BurnRunway struct {
	CashBalance     float64  `json:"cash_balance"`
	MonthlyBurn     float64  `json:"monthly_burn"`
	RunwayMonths    *float64 `json:"runway_months"`    // null if burn=0
	OperatingMargin float64  `json:"operating_margin"` // Retainers - Fixed Costs
	TotalRetainers  float64  `json:"total_retainers"`
}

func GetBurnRunway(db *sql.DB, agencyID string) (*BurnRunway, error) {
	var result BurnRunway
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	// 1. Get latest cash balance from daily_cash_snapshots
	err := db.QueryRow(`
		SELECT cash_balance FROM daily_cash_snapshots 
		WHERE agency_id = $1 
		ORDER BY date DESC 
		LIMIT 1
	`, agencyID).Scan(&result.CashBalance)

	if err == sql.ErrNoRows {
		return nil, nil // No cash snapshot
	} else if err != nil {
		return nil, err
	}

	// 2. Sum fixed costs in last 30 days
	err = db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) FROM daily_costs 
		WHERE agency_id = $1 AND type = 'fixed' AND date >= $2
	`, agencyID, thirtyDaysAgo).Scan(&result.MonthlyBurn)
	if err != nil {
		return nil, err
	}

	// 3. Get total active retainers
	db.QueryRow(`SELECT COALESCE(SUM(monthly_amount), 0) FROM retainers WHERE agency_id = $1 AND active = true`, agencyID).Scan(&result.TotalRetainers)

	// 4. Calculate Panic Runway
	if result.MonthlyBurn > 0 {
		runway := result.CashBalance / result.MonthlyBurn
		runway = float64(int(runway*10)) / 10
		result.RunwayMonths = &runway
	}

	// 5. Calculate Operating Margin (Retainers - Fixed Costs)
	result.OperatingMargin = result.TotalRetainers - result.MonthlyBurn

	return &result, nil
}

// Clients & Retainers (Ticket 06)

type Client struct {
	ID       string `json:"id"`
	AgencyID string `json:"-"`
	Name     string `json:"name"`
	Status   string `json:"status"`
}

func CreateClient(db *sql.DB, agencyID string, name string) (*Client, error) {
	id := uuid.New().String()
	_, err := db.Exec(`
		INSERT INTO clients (id, agency_id, name, status)
		VALUES ($1, $2, $3, 'active')
	`, id, agencyID, name)
	if err != nil {
		return nil, err
	}
	return &Client{ID: id, AgencyID: agencyID, Name: name, Status: "active"}, nil
}

func GetClients(db *sql.DB, agencyID string) ([]Client, error) {
	rows, err := db.Query(`
		SELECT id, name, status FROM clients 
		WHERE agency_id = $1 AND status = 'active'
		ORDER BY name
	`, agencyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var c Client
		c.AgencyID = agencyID
		if err := rows.Scan(&c.ID, &c.Name, &c.Status); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

func CreateRetainer(db *sql.DB, agencyID string, clientID string, amount float64) error {
	// Check if client already has active retainer
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM retainers WHERE client_id = $1 AND active = true)
	`, clientID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("client already has active retainer")
	}

	id := uuid.New().String()
	_, err = db.Exec(`
		INSERT INTO retainers (id, agency_id, client_id, monthly_amount, active)
		VALUES ($1, $2, $3, $4, true)
	`, id, agencyID, clientID, amount)
	return err
}

type RetainerSummary struct {
	TotalRetainerRevenue float64 `json:"total_retainer_revenue"`
	FixedCosts           float64 `json:"fixed_costs"`
	CoverageRatio        float64 `json:"coverage_ratio"`
	TopClientPercentage  float64 `json:"top_client_percentage"`
}

func GetRetainerSummary(db *sql.DB, agencyID string) (*RetainerSummary, error) {
	var summary RetainerSummary

	// 1. Sum all active retainers
	err := db.QueryRow(`
		SELECT COALESCE(SUM(monthly_amount), 0) FROM retainers 
		WHERE agency_id = $1 AND active = true
	`, agencyID).Scan(&summary.TotalRetainerRevenue)
	if err != nil {
		return nil, err
	}

	// 2. Get fixed costs (last 30 days) - reuse burn rate logic
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	err = db.QueryRow(`
		SELECT COALESCE(SUM(amount), 0) FROM daily_costs 
		WHERE agency_id = $1 AND type = 'fixed' AND date >= $2
	`, agencyID, thirtyDaysAgo).Scan(&summary.FixedCosts)
	if err != nil {
		return nil, err
	}

	// 3. Calculate coverage ratio
	if summary.FixedCosts > 0 {
		summary.CoverageRatio = summary.TotalRetainerRevenue / summary.FixedCosts
		// Round to 2 decimals
		summary.CoverageRatio = float64(int(summary.CoverageRatio*100)) / 100
	}

	// 4. Get top client percentage
	if summary.TotalRetainerRevenue > 0 {
		var maxRetainer float64
		err = db.QueryRow(`
			SELECT COALESCE(MAX(monthly_amount), 0) FROM retainers 
			WHERE agency_id = $1 AND active = true
		`, agencyID).Scan(&maxRetainer)
		if err != nil {
			return nil, err
		}
		summary.TopClientPercentage = maxRetainer / summary.TotalRetainerRevenue
		// Round to 2 decimals
		summary.TopClientPercentage = float64(int(summary.TopClientPercentage*100)) / 100
	}

	return &summary, nil
}

// Time Entries & Utilization (Ticket 07)

const MonthlyCapacityHours = 160.0 // 1 person × 8 hours × 20 days

func AddTimeEntry(db *sql.DB, agencyID string, clientID *string, hours float64) error {
	id := uuid.New().String()
	date := time.Now().Format("2006-01-02")

	_, err := db.Exec(`
		INSERT INTO time_entries (id, agency_id, client_id, date, hours)
		VALUES ($1, $2, $3, $4, $5)
	`, id, agencyID, clientID, date, hours)

	return err
}

type Utilization struct {
	UsedHours          float64 `json:"used_hours"`
	CapacityHours      float64 `json:"capacity_hours"`
	UtilizationPercent float64 `json:"utilization_percent"`
}

func GetUtilization(db *sql.DB, agencyID string) (*Utilization, error) {
	var result Utilization
	result.CapacityHours = MonthlyCapacityHours

	// Sum hours in last 30 days
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	err := db.QueryRow(`
		SELECT COALESCE(SUM(hours), 0) FROM time_entries 
		WHERE agency_id = $1 AND date >= $2
	`, agencyID, thirtyDaysAgo).Scan(&result.UsedHours)
	if err != nil {
		return nil, err
	}

	// Calculate utilization %
	if result.CapacityHours > 0 {
		result.UtilizationPercent = (result.UsedHours / result.CapacityHours) * 100
		// Round to 1 decimal
		result.UtilizationPercent = float64(int(result.UtilizationPercent*10)) / 10
	}

	return &result, nil
}

// Agency Reality Score (Ticket 08) - Updated Ticket 09

type ScoreBreakdown struct {
	RetainerSafety      int `json:"retainer_safety"`
	Runway              int `json:"runway"`
	ClientConcentration int `json:"client_concentration"`
	Profitability       int `json:"profitability"`
	CapacityPressure    int `json:"capacity_pressure"`
}

type AgencyRealityScore struct {
	Score              int            `json:"score"`
	Breakdown          ScoreBreakdown `json:"breakdown"`
	Status             string         `json:"status"`
	CashOnHand         float64        `json:"cash_on_hand"`
	CommittedRetainers float64        `json:"committed_retainers"`
}

func GetAgencyRealityScore(db *sql.DB, agencyID string) (*AgencyRealityScore, error) {
	var result AgencyRealityScore
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	// A. Retainer Safety (25 pts) - Coverage = Retainer / Fixed Costs
	var totalRetainer, fixedCosts float64
	db.QueryRow(`SELECT COALESCE(SUM(monthly_amount), 0) FROM retainers WHERE agency_id = $1 AND active = true`, agencyID).Scan(&totalRetainer)
	db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM daily_costs WHERE agency_id = $1 AND type = 'fixed' AND date >= $2`, agencyID, thirtyDaysAgo).Scan(&fixedCosts)

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
	var cashBalance float64
	err := db.QueryRow(`SELECT cash_balance FROM daily_cash_snapshots WHERE agency_id = $1 ORDER BY date DESC LIMIT 1`, agencyID).Scan(&cashBalance)
	if err == nil && fixedCosts > 0 {
		runway := cashBalance / fixedCosts
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

	// C. Client Concentration (20 pts) - Top Client % of Retainers
	if totalRetainer > 0 {
		var maxRetainer float64
		db.QueryRow(`SELECT COALESCE(MAX(monthly_amount), 0) FROM retainers WHERE agency_id = $1 AND active = true`, agencyID).Scan(&maxRetainer)
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

	// D. Profitability (20 pts) - Net Profit Margin (last 30 days)
	var revenue, costs float64
	db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM daily_revenues WHERE agency_id = $1 AND date >= $2`, agencyID, thirtyDaysAgo).Scan(&revenue)
	db.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM daily_costs WHERE agency_id = $1 AND date >= $2`, agencyID, thirtyDaysAgo).Scan(&costs)

	if revenue > 0 {
		margin := ((revenue - costs) / revenue) * 100
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

	// E. Capacity Pressure (15 pts) - Utilization %
	var usedHours float64
	db.QueryRow(`SELECT COALESCE(SUM(hours), 0) FROM time_entries WHERE agency_id = $1 AND date >= $2`, agencyID, thirtyDaysAgo).Scan(&usedHours)
	utilization := (usedHours / MonthlyCapacityHours) * 100

	if utilization >= 60 && utilization <= 85 {
		result.Breakdown.CapacityPressure = 15
	} else if utilization >= 50 && utilization < 60 {
		result.Breakdown.CapacityPressure = 10
	} else if utilization >= 40 && utilization < 50 {
		result.Breakdown.CapacityPressure = 6
	} else if utilization > 85 && utilization <= 100 {
		result.Breakdown.CapacityPressure = 6
	}

	// Calculate total score
	result.Score = result.Breakdown.RetainerSafety +
		result.Breakdown.Runway +
		result.Breakdown.ClientConcentration +
		result.Breakdown.Profitability +
		result.Breakdown.CapacityPressure

	// Status mapping
	if result.Score >= 80 {
		result.Status = "Healthy"
	} else if result.Score >= 60 {
		result.Status = "Watch"
	} else if result.Score >= 40 {
		result.Status = "At Risk"
	} else {
		result.Status = "Danger"
	}

	// Add header metrics (Ticket 09)
	result.CashOnHand = cashBalance
	result.CommittedRetainers = totalRetainer

	return &result, nil
}
