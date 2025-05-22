package types

import (
	"bufio"
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

type ProcessEntry struct {
	ID          string
	ProcessName string
	InputItems  []*Item
	OutputItem  *Item
	StartTime   time.Time
	EndTime     time.Time
	Status      string
	Error       error
}

type Process struct {
	Name          string
	Description   string
	InputTypes    []string
	OutputType    string
	TransformFunc func([]*Item) (*Item, error)
	CycleTime     time.Duration
	FailureRate   float64
}

type ProcessFlow struct {
	Processes     []*Process
	Inventory     map[string][]*Item
	ProcessQueues map[string][]*Item
	Completed     []*Item
	Defects       []*Item
	History       []*ProcessEntry
	reader        *bufio.Reader
}

type Item struct {
	ID         string
	Name       string
	Type       string
	Properties map[string]interface{}
	CreatedAt  time.Time
	ModifiedAt time.Time
}
