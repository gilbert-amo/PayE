package payroll

import (
	"fmt"
	"strings"

	"github.com/gilbert-amo/PayE/pension"
	"github.com/gilbert-amo/PayE/types"
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

// CalculateTaxRecursive implements progressive tax bracket calculation
// where tax is calculated progressively on each bracket range
// The minimumWage parameter ensures that minimum wage is always taxed at 0%
func CalculateTaxRecursive(originalSalary float64, brackets []types.TaxBracket, minimumWage float64) (float64, []TaxCalculationStep) {
	var totalTax float64
	var steps []TaxCalculationStep

	// Create a copy of brackets and add minimum wage as 0% tax threshold
	// This ensures minimum wage is always tax-free
	enhancedBrackets := make([]types.TaxBracket, 0, len(brackets)+1)

	// Add minimum wage as 0% tax bracket at the beginning
	enhancedBrackets = append(enhancedBrackets, types.TaxBracket{
		Threshold: minimumWage,
		Rate:      0.0, // 0% tax on minimum wage
	})

	// Add all other brackets
	enhancedBrackets = append(enhancedBrackets, brackets...)

	// Progressive subtraction: each threshold amount is subtracted from basic salary
	remainingSalary := originalSalary
	for _, bracket := range enhancedBrackets {
		var taxAmount float64
		var taxableAmount float64

		if remainingSalary >= bracket.Threshold {
			// Tax at the threshold rate for the full threshold amount
			taxAmount = bracket.Threshold * (bracket.Rate / 100)
			taxableAmount = bracket.Threshold
			remainingSalary -= bracket.Threshold
		} else if remainingSalary > 0 {
			// Apply the same rate to the remaining salary and end calculation
			taxAmount = remainingSalary * (bracket.Rate / 100)
			taxableAmount = remainingSalary
			remainingSalary = 0 // Set to 0 since we're taxing the remaining amount
		} else {
			// No remaining salary to tax
			break
		}

		totalTax += taxAmount

		// Record this calculation step
		step := TaxCalculationStep{
			RemainingSalary:  remainingSalary + taxableAmount, // Show salary before this step
			TaxRate:          bracket.Rate,
			TaxAmount:        taxAmount,
			BracketThreshold: bracket.Threshold,
			TaxableAmount:    taxableAmount,
			BracketRange:     fmt.Sprintf("Tax at %.1f%%", bracket.Rate),
		}
		steps = append(steps, step)

		// If we've taxed the remaining salary (less than threshold), end the session
		if remainingSalary == 0 {
			break
		}

		// If salary becomes negative, stop
		if remainingSalary < 0 {
			break
		}
	}

	return totalTax, steps
}

// TaxCalculationStep represents a single step in the progressive tax calculation
type TaxCalculationStep struct {
	RemainingSalary  float64
	TaxRate          float64
	TaxAmount        float64
	BracketThreshold float64
	TaxableAmount    float64
	BracketRange     string
}

// findApplicableBracket finds the highest applicable tax bracket for a given salary
func findApplicableBracket(salary float64, brackets []types.TaxBracket) *types.TaxBracket {
	var applicableBracket *types.TaxBracket

	for i := len(brackets) - 1; i >= 0; i-- {
		if salary >= brackets[i].Threshold {
			applicableBracket = &brackets[i]
			break
		}
	}

	return applicableBracket
}

// CalculateTax is kept for backward compatibility but now uses the recursive method
func CalculateTax(salary float64, brackets []types.TaxBracket, minimumWage float64) float64 {
	tax, _ := CalculateTaxRecursive(salary, brackets, minimumWage)
	return tax
}

// GetTaxCalculationSummary returns a summary of the recursive tax calculation
func GetTaxCalculationSummary(steps []TaxCalculationStep) string {
	if len(steps) == 0 {
		return "No tax calculation steps available"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tax Calculation Summary:\n"))
	sb.WriteString(fmt.Sprintf("Total Steps: %d\n", len(steps)))
	sb.WriteString(fmt.Sprintf("Total Tax: %.2f\n", calculateTotalTax(steps)))
	sb.WriteString(fmt.Sprintf("Average Tax Rate: %.2f%%\n", calculateAverageTaxRate(steps)))

	return sb.String()
}

// calculateTotalTax calculates the total tax from all steps
func calculateTotalTax(steps []TaxCalculationStep) float64 {
	total := 0.0
	for _, step := range steps {
		total += step.TaxAmount
	}
	return total
}

// calculateAverageTaxRate calculates the average tax rate across all steps
func calculateAverageTaxRate(steps []TaxCalculationStep) float64 {
	if len(steps) == 0 {
		return 0.0
	}

	totalRate := 0.0
	for _, step := range steps {
		totalRate += step.TaxRate
	}
	return totalRate / float64(len(steps))
}
func CalculatePieceEarnings(pieces []types.PieceRateAggregation) float64 {
	total := 0.0
	for _, pw := range pieces {
		total += pw.Rate * pw.Quantity
	}
	return total
}

func PrintPayrollReport(emp types.Employee, country types.Country, originalBasic, gross, tax float64, pensionCalc *pension.Calculator, net float64, taxSteps []TaxCalculationStep) {
	fmt.Printf("\n=== %s ===\n", emp.Name)
	fmt.Printf("Country: %s (Min Wage: %.2f)\n", country.Name, country.MinimumWage)

	if originalBasic > 0 {
		fmt.Printf("\nOriginal Basic Salary: %.2f\n", originalBasic)
	}
	fmt.Printf("Current Basic Salary: %.2f\n", emp.BasicSalary)

	if emp.Allowance > 0 {
		fmt.Printf("Allowance: %.2f\n", emp.Allowance)
	}

	if len(emp.PieceRate) > 0 {
		fmt.Println("\nPiece-Rate Details:")
		for _, pw := range emp.PieceRate {
			fmt.Printf("- %s: %.0f × %.2f = %.2f\n",
				pw.Item, pw.Quantity, pw.Rate, pw.Rate*pw.Quantity)
		}
		fmt.Printf("Total Piece-Rate Earnings: %.2f\n", CalculatePieceEarnings(emp.PieceRate))
	}

	fmt.Println("\nTax Calculation:")
	for _, bracket := range country.TaxBrackets {
		fmt.Printf("- Threshold %.2f: %.1f%%\n", bracket.Threshold, bracket.Rate)
	}

	// Display progressive subtraction tax calculation steps if available
	if len(taxSteps) > 0 {
		fmt.Println("\nProgressive Subtraction Tax Calculation Steps:")
		fmt.Println("Step | Remaining Salary | Subtract Amount | Tax Rate | Tax Amount")
		fmt.Println("-----|------------------|----------------|----------|-----------")
		for i, step := range taxSteps {
			fmt.Printf("%4d | %16.2f | %14.2f | %8.1f%% | %10.2f\n",
				i+1, step.RemainingSalary, step.TaxableAmount, step.TaxRate, step.TaxAmount)
		}
		fmt.Printf("Total Tax: %.2f\n", tax)
	}

	// Print pension information
	fmt.Println("\nPension Contributions:")
	for k, v := range pensionCalc.GetContributionBreakdown() {
		fmt.Printf("%-20s: %.2f\n", k, v)
	}

	fmt.Println("\nPension Tier Allocations:")
	for k, v := range pensionCalc.GetTierBreakdown() {
		fmt.Printf("%-20s: %.2f\n", k, v)
	}

	fmt.Printf("\nGROSS SALARY: %.2f\n", gross)
	fmt.Printf("TAX DEDUCTION: %.2f\n", tax)
	fmt.Printf("PENSION DEDUCTION: %.2f\n", pensionCalc.EmployeeContribution)
	fmt.Printf("TOTAL DEDUCTIONS: %.2f\n", tax+pensionCalc.EmployeeContribution)
	fmt.Printf("NET SALARY: %.2f\n", net)
	fmt.Println(strings.Repeat("=", 30))
}
