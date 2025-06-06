package types

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

type Country struct {
	Name        string
	MinimumWage float64
	TaxBrackets []TaxBracket
}
