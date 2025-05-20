package payroll

import (
	"fmt"
	"strings"
)

type Employee struct {
	Name        string
	BasicSalary float64
	CountryCode string
	PieceRate   []PieceRateAggregation
}

type PieceRateAggregation struct {
	Item     string
	Rate     float64 // unit price
	Quantity float64
}

func (e *Employee) CalculateSalary() float64 {
	totalPieceRate := e.calculatePieceRate()

	if e.BasicSalary > 0 {
		// Piece work is a bonus on top of basic salary
		return e.BasicSalary + totalPieceRate
	}
	// Piece work is the entire salary
	return totalPieceRate
}

// calculatePieceRate computes the total from piece-rate work
func (e *Employee) calculatePieceRate() float64 {
	total := 0.0
	for _, item := range e.PieceRate {
		total += item.Rate * item.Quantity
	}
	return total
}

// AddPieceRate adds a piece-rate work item to the employee
func (e *Employee) AddPieceRate(item string, rate, quantity float64) {
	e.PieceRate = append(e.PieceRate, PieceRateAggregation{
		Item:     item,
		Rate:     rate,
		Quantity: quantity,
	})
}

// PrintSalaryBreakdown displays the compensation details
func (e *Employee) PrintSalaryBreakdown() {
	fmt.Printf("\nSalary Breakdown for %s:\n", e.Name)
	fmt.Println(strings.Repeat("-", 30))

	if e.BasicSalary > 0 {
		fmt.Printf("Basic Salary: $%.2f\n", e.BasicSalary)
	}

	if len(e.PieceRate) > 0 {
		fmt.Println("\nPiece Rate Earnings:")
		for _, item := range e.PieceRate {
			earnings := item.Rate * item.Quantity
			fmt.Printf("- %s: %.0f units @ $%.2f = $%.2f\n",
				item.Item, item.Quantity, item.Rate, earnings)
		}
		fmt.Printf("Total Piece Rate: $%.2f\n", e.calculatePieceRate())
	}

	total := e.CalculateSalary()
	if e.BasicSalary > 0 && len(e.PieceRate) > 0 {
		fmt.Printf("\nTotal Salary (Basic + Piece Rate): $%.2f\n", total)
	} else {
		fmt.Printf("\nTotal Earnings: $%.2f\n", total)
	}
}
