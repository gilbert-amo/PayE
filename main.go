package main

import (
	"bufio"
	"fmt"
	"github.com/gilbert-amo/PayE/payroll"
	"github.com/gilbert-amo/PayE/pension"
	"github.com/gilbert-amo/PayE/types"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Country struct {
	Name        string
	MinimumWage float64
	TaxBrackets []types.TaxBracket
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	countries := make(map[string]Country)
	config := payroll.Config{SplitEnabled: false}

	pensionTiers := []types.Tier{
		{Name: "Tier 1", Percentage: 0.135},
		{Name: "Tier 2", Percentage: 0.55},
		{Name: "Tier 3", Percentage: 0.315},
	}

	// Country setup (unchanged)
	fmt.Println("=== Country Setup ===")
	for {
		fmt.Print("\nEnter country code (3 letters, or 'done' to finish): ")
		code, _ := reader.ReadString('\n')
		code = strings.ToUpper(strings.TrimSpace(code))

		if code == "DONE" {
			break
		}

		if len(code) != 3 {
			fmt.Println("Country code must be exactly 3 letters")
			continue
		}

		fmt.Print("Enter country name: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		fmt.Print("Enter minimum wage: ")
		minWageStr, _ := reader.ReadString('\n')
		minWage, err := strconv.ParseFloat(strings.TrimSpace(minWageStr), 64)
		if err != nil {
			fmt.Println("Invalid wage. Please enter a valid number.")
			continue
		}

		// Tax bracket setup
		var taxBrackets []types.TaxBracket
		// When setting up country tax brackets:
		fmt.Println("\nSetting up peak tax brackets:")
		fmt.Println("Enter thresholds where specific rates should apply")
		fmt.Println("Example: If salary reaches 500, apply 5% to entire amount")
		//fmt.Println("         If salary reaches 1000, apply 10% to entire amount")

		for {
			fmt.Print("Enter threshold amount (or 'done'): ")
			thresholdStr, _ := reader.ReadString('\n')
			thresholdStr = strings.TrimSpace(thresholdStr)
			if strings.ToLower(thresholdStr) == "done" {
				break
			}

			threshold, err := strconv.ParseFloat(thresholdStr, 64)
			if err != nil {
				fmt.Println("Invalid threshold. Please enter a valid number.")
				continue
			}

			fmt.Print("Enter tax rate to apply when threshold is reached (e.g., 5 for 5%): ")
			rateStr, _ := reader.ReadString('\n')
			rate, err := strconv.ParseFloat(strings.TrimSpace(rateStr), 64)
			if err != nil {
				fmt.Println("Invalid rate. Please enter a valid number.")
				continue
			}

			taxBrackets = append(taxBrackets, types.TaxBracket{
				Threshold: threshold,
				Rate:      rate,
			})
		}

		// Sort brackets by threshold
		sort.Slice(taxBrackets, func(i, j int) bool {
			return taxBrackets[i].Threshold < taxBrackets[j].Threshold
		})

		countries[code] = Country{
			Name:        name,
			MinimumWage: minWage,
			TaxBrackets: taxBrackets,
		}
	}

	if len(countries) == 0 {
		fmt.Println("No countries entered. Exiting.")
		return
	}

	// Configure splitting (unchanged)
	fmt.Print("\nEnable salary splitting when piece-rate ≥ basic? (y/n): ")
	splitInput, _ := reader.ReadString('\n')
	if strings.ToLower(strings.TrimSpace(splitInput)) == "y" {
		config.SplitEnabled = true

		fmt.Print("Enter basic salary ratio (e.g., 0.7 for 70%): ")
		basicRatioStr, _ := reader.ReadString('\n')
		basicRatio, err := strconv.ParseFloat(strings.TrimSpace(basicRatioStr), 64)
		if err != nil || basicRatio <= 0 || basicRatio >= 1 {
			fmt.Println("Invalid ratio. Using default 0.7")
			basicRatio = 0.7
		}
		config.BasicSalaryRatio = basicRatio
		config.AllowanceRatio = 1 - basicRatio
	}

	// Employee setup
	fmt.Println("\n=== Employee Setup ===")
	var employees []types.Employee
	for {
		emp := types.Employee{}

		fmt.Print("\nEnter employee name (or 'done' to finish): ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if strings.ToLower(name) == "done" {
			break
		}
		emp.Name = name

		fmt.Print("Enter basic salary (0 for piece-rate only): ")
		salaryStr, _ := reader.ReadString('\n')
		salary, err := strconv.ParseFloat(strings.TrimSpace(salaryStr), 64)
		if err != nil {
			fmt.Println("Invalid salary. Please enter a valid number.")
			continue
		}
		emp.BasicSalary = salary

		// Add piece-rate work
		for {
			fmt.Print("Add piece-rate item (name or 'done'): ")
			item, _ := reader.ReadString('\n')
			item = strings.TrimSpace(item)
			if strings.ToLower(item) == "done" {
				break
			}

			fmt.Print("Enter unit price: ")
			rateStr, _ := reader.ReadString('\n')
			rate, err := strconv.ParseFloat(strings.TrimSpace(rateStr), 64)
			if err != nil {
				fmt.Println("Invalid rate. Please try again.")
				continue
			}

			fmt.Print("Enter quantity: ")
			qtyStr, _ := reader.ReadString('\n')
			qty, err := strconv.ParseFloat(strings.TrimSpace(qtyStr), 64)
			if err != nil {
				fmt.Println("Invalid quantity. Please try again.")
				continue
			}

			emp.PieceRate = append(emp.PieceRate, types.PieceRateAggregation{
				Item:     item,
				Rate:     rate,
				Quantity: qty,
			})
		}

		if len(emp.PieceRate) == 0 && emp.BasicSalary == 0 {
			fmt.Println("Error: Employee must have either basic salary or piece-rate work")
			continue
		}

		fmt.Print("Enter employee's country code: ")
		countryCode, _ := reader.ReadString('\n')
		countryCode = strings.ToUpper(strings.TrimSpace(countryCode))

		if _, exists := countries[countryCode]; !exists {
			fmt.Println("Invalid country code. Please try again.")
			continue
		}
		emp.CountryCode = countryCode

		employees = append(employees, emp)
	}

	if len(employees) == 0 {
		fmt.Println("No employees entered. Exiting.")
		return
	}

	// Process payroll
	fmt.Println("\n=== Payroll Results ===")
	for _, emp := range employees {
		country := countries[emp.CountryCode]
		pieceEarnings := calculatePieceEarnings(emp.PieceRate)
		originalBasic := emp.BasicSalary
		emp.Allowance = 0

		// Handle employees with no basic salary
		if originalBasic == 0 && len(emp.PieceRate) > 0 {
			fmt.Printf("\n%s has no basic salary - using piece-rate as basic\n", emp.Name)
			emp.BasicSalary = pieceEarnings
			pieceEarnings = 0 // Reset since we've converted to basic salary
		}

		// Conditional splitting for employees with both
		if originalBasic > 0 && config.SplitEnabled && pieceEarnings >= originalBasic {
			fmt.Printf("\nPiece-rate (%.2f) ≥ basic salary (%.2f)\n", pieceEarnings, originalBasic)
			fmt.Println("Converting piece-rate to basic salary + allowance")

			emp.BasicSalary = pieceEarnings * config.BasicSalaryRatio
			emp.Allowance = pieceEarnings * config.AllowanceRatio
			fmt.Printf("- New Basic: %.2f\n- Allowance: %.2f\n", emp.BasicSalary, emp.Allowance)
		} else if originalBasic > 0 && len(emp.PieceRate) > 0 {
			// Keep original basic and add piece-rate as bonus
			fmt.Printf("\nAdding piece-rate earnings (%.2f) as bonus\n", pieceEarnings)
			emp.Allowance = pieceEarnings
		}

		// Minimum wage enforcement
		if emp.BasicSalary < country.MinimumWage {
			fmt.Printf("Adjusting basic salary to meet minimum wage (%.2f → %.2f)\n",
				emp.BasicSalary, country.MinimumWage)
			emp.BasicSalary = country.MinimumWage
		}

		totalEarnings := emp.BasicSalary + emp.Allowance

		taxAmount := 0.0
		if len(country.TaxBrackets) > 0 {
			taxAmount = payroll.CalculateTax(totalEarnings, country.TaxBrackets)
		}

		// Calculate pension
		pensionCalc := pension.NewCalculator(&emp)
		pensionCalc.Calculate(pensionTiers)
		pensionDeduction := pensionCalc.EmployeeContribution
		totalDeductions := taxAmount + pensionDeduction

		netSalary := totalEarnings - totalDeductions

		printPayrollReport(emp, country, originalBasic, totalEarnings, taxAmount, pensionCalc, netSalary)
	}
}

func calculatePieceEarnings(pieces []types.PieceRateAggregation) float64 {
	total := 0.0
	for _, pw := range pieces {
		total += pw.Rate * pw.Quantity
	}
	return total
}

func printPayrollReport(emp types.Employee, country Country, originalBasic, gross, tax float64, pensionCalc *pension.Calculator, net float64) {
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
		fmt.Printf("Total Piece-Rate Earnings: %.2f\n", calculatePieceEarnings(emp.PieceRate))
	}

	fmt.Println("\nTax Calculation:")
	for i, bracket := range country.TaxBrackets {
		if i == 0 {
			fmt.Printf("- Up to %.2f: %.1f%%\n", bracket.Threshold, bracket.Rate)
		} else {
			fmt.Printf("- %.2f to %.2f: %.1f%%\n",
				country.TaxBrackets[i-1].Threshold, bracket.Threshold, bracket.Rate)
		}
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
