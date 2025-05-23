package types

import (
	"time"
)

type Employee struct {
	Name        string
	BasicSalary float64
	CountryCode string
	PieceRate   []PieceRateAggregation
	Allowance   float64
}

type PieceRateAggregation struct {
	Item     string
	Rate     float64
	Quantity float64
}

type TaxBracket struct {
	Threshold float64
	Rate      float64
}

type Tier struct {
	Name       string
	Percentage float64
}

type Process struct {
	ID           int
	Name         string
	Description  string
	Supervisor   string
	StaffWorker  string
	StartTime    time.Time
	EndTime      time.Time
	Status       string // "Pending", "In Progress", "Completed"
	Quantity     int
	QualityCheck bool
	Notes        string
}

// Product represents the product being manufactured
type Product struct {
	ID            int
	Name          string
	Description   string
	Processes     []Process
	TotalQuantity int
	StartDate     time.Time
	TargetDate    time.Time
}

// Worker represents a staff member
type Worker struct {
	ID      int
	Name    string
	Role    string // "Supervisor" or "Staff"
	Contact string
	Shift   string
}
