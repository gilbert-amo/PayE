package payroll

import (
	"fmt"
	"github.com/gilbert-amo/PayE/pension"
	"github.com/gilbert-amo/PayE/types"
	"sort"
	"strings"
)

type Config struct {
	SplitEnabled     bool
	BasicSalaryRatio float64
	AllowanceRatio   float64
}

func CalculateSalary(e *types.Employee) float64 {
	totalPieceRate := calculatePieceRate(e)

	if e.BasicSalary > 0 {
		// Piece work is a bonus on top of basic salary
		return e.BasicSalary + totalPieceRate
	}
	// Piece work is the entire salary
	return totalPieceRate
}

// calculatePieceRate computes the total from piece-rate work
func calculatePieceRate(e *types.Employee) float64 {
	total := 0.0
	for _, item := range e.PieceRate {
		total += item.Rate * item.Quantity
	}
	return total
}

// AddPieceRate adds a piece-rate work item to the employee
func AddPieceRate(e *types.Employee, item string, rate, quantity float64) {
	e.PieceRate = append(e.PieceRate, types.PieceRateAggregation{
		Item:     item,
		Rate:     rate,
		Quantity: quantity,
	})
}

func GetSalaryBreakdownWithPension(e *types.Employee, pensionCalc *pension.Calculator) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\nSalary Breakdown for %s:\n", e.Name))
	sb.WriteString(strings.Repeat("-", 30) + "\n")

	if e.BasicSalary > 0 {
		sb.WriteString(fmt.Sprintf("Basic Salary: $%.2f\n", e.BasicSalary))
	}

	if len(e.PieceRate) > 0 {
		fmt.Println("\nPiece Rate Earnings:")
		for _, item := range e.PieceRate {
			earnings := item.Rate * item.Quantity
			fmt.Printf("- %s: %.0f units @ $%.2f = $%.2f\n",
				item.Item, item.Quantity, item.Rate, earnings)
		}
		fmt.Printf("Total Piece Rate: $%.2f\n", calculatePieceRate(e))
	}

	total := CalculateSalary(e)
	if e.BasicSalary > 0 && len(e.PieceRate) > 0 {
		fmt.Printf("\nTotal Salary (Basic + Piece Rate): $%.2f\n", total)
	} else {
		fmt.Printf("\nTotal Earnings: $%.2f\n", total)
	}

	pensionBreakdown := pensionCalc.GetContributionBreakdown()
	sb.WriteString("\nPension Contributions:\n")
	for k, v := range pensionBreakdown {
		sb.WriteString(fmt.Sprintf("%-20s: %.2f\n", k, v))
	}

	return sb.String()
}

func CalculateTax(salary float64, brackets []types.TaxBracket) float64 {
	tax := 0.0

	// Sort brackets just in case (from lowest to highest threshold)
	sort.Slice(brackets, func(i, j int) bool {
		return brackets[i].Threshold < brackets[j].Threshold
	})

	// Find which bracket the salary falls into
	for i := len(brackets) - 1; i >= 0; i-- {
		if salary >= brackets[i].Threshold {
			// Apply the rate for this bracket to the entire salary
			tax = salary * (brackets[i].Rate / 100)
			break
		}
	}

	// If salary is below all thresholds, no tax
	return tax
}
