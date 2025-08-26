package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gilbert-amo/PayE/payroll"
	"github.com/gilbert-amo/PayE/pension"
	"github.com/gilbert-amo/PayE/types"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("=== Country Setup ===")
		countries := make(map[string]types.Country)
		config := payroll.Config{SplitEnabled: false}

		// Pension tier setup
		fmt.Println("\n=== Pension Tier Setup ===")
		fmt.Println("Enter pension tiers (total must equal 100%)")
		fmt.Println("Example: Tier 1: 13.5%, Tier 2: 55%, Tier 3: 31.5%")
		fmt.Println("Type 'done' for tier name when finished")

		var pensionTiers []types.Tier
		totalPercentage := 0.0

		for {
			fmt.Printf("\nCurrent total: %.1f%%\n", totalPercentage*100)

			fmt.Print("Enter tier name (or 'done' to finish): ")
			tierName, _ := reader.ReadString('\n')
			tierName = strings.TrimSpace(tierName)
			if strings.ToLower(tierName) == "done" {
				break
			}

			if tierName == "" {
				fmt.Println("Tier name cannot be empty. Please try again.")
				continue
			}

			remainingPercentage := (1.0 - totalPercentage) * 100
			if remainingPercentage <= 0 {
				fmt.Println("Total already equals 100%. Tier setup complete!")
				break
			}

			fmt.Printf("Enter tier percentage (0.1 to %.1f): ", remainingPercentage)
			percentageStr, _ := reader.ReadString('\n')
			percentage, err := strconv.ParseFloat(strings.TrimSpace(percentageStr), 64)
			if err != nil || percentage <= 0 {
				fmt.Println("Invalid percentage. Please enter a valid positive number.")
				continue
			}

			// Convert percentage to decimal
			percentageDecimal := percentage / 100

			if totalPercentage+percentageDecimal > 1.001 { // Allow small floating point precision
				fmt.Printf("Total would exceed 100%%. Maximum allowed: %.1f%%\n", remainingPercentage)
				continue
			}

			pensionTiers = append(pensionTiers, types.Tier{
				Name:       tierName,
				Percentage: percentageDecimal,
			})

			totalPercentage += percentageDecimal
			fmt.Printf("Added %s: %.1f%%\n", tierName, percentage)

			// Check if we've reached 100% (with small tolerance)
			if totalPercentage >= 0.999 {
				fmt.Printf("Total reached %.1f%%. Tier setup complete!\n", totalPercentage*100)
				break
			}
		}

		// Validate total equals 100% (with tolerance)
		if totalPercentage < 0.999 {
			fmt.Printf("Warning: Total percentage is %.1f%%, not 100%%\n", totalPercentage*100)
			fmt.Println("Using default tiers instead")
			pensionTiers = []types.Tier{
				{Name: "Tier 1", Percentage: 0.135},
				{Name: "Tier 2", Percentage: 0.55},
				{Name: "Tier 3", Percentage: 0.315},
			}
		} else {
			fmt.Printf("✓ Pension tiers configured successfully! Total: %.1f%%\n", totalPercentage*100)
		}

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
			fmt.Println("\nSetting up tax brackets:")
			fmt.Printf("Note: Minimum wage (%.2f) will automatically be taxed at 0%%\n", minWage)
			fmt.Println("Enter additional threshold amounts that will be taxed at specific rates")
			fmt.Println("Each threshold amount will be taxed at its rate and subtracted from salary")
			fmt.Println("Example: Threshold 500 with 5% rate = 500 × 5% = 25 tax")
			fmt.Println("Note: 0% tax rates are valid (e.g., for tax-free allowances)")

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

				// Allow 0% tax rates for progressive tax systems (e.g., tax-free allowance)
				// 0% rates are valid and commonly used in progressive tax systems

				taxBrackets = append(taxBrackets, types.TaxBracket{
					Threshold: threshold,
					Rate:      rate,
				})
			}

			// Tax brackets are processed in the order they were entered (no sorting)

			countries[code] = types.Country{
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
			pieceEarnings := payroll.CalculatePieceEarnings(emp.PieceRate)
			originalBasic := emp.BasicSalary // Track the original basic salary for recursive tax calculation
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
			var taxCalculationSteps []payroll.TaxCalculationStep
			if len(country.TaxBrackets) > 0 {
				// Use recursive PAYE tax calculation on the ORIGINAL basic salary entered by user
				// This ensures we track the tax calculation from the initial salary amount
				// The minimum wage is automatically applied as a 0% tax threshold
				taxAmount, taxCalculationSteps = payroll.CalculateTaxRecursive(originalBasic, country.TaxBrackets, country.MinimumWage)
			}

			// Calculate pension
			pensionCalc := pension.NewCalculator(&emp)

			// Option 1: Use the tiers entered during setup
			pensionCalc.Calculate(pensionTiers)

			// Option 2: Use internal tiers (uncomment to use instead)
			// pensionCalc.SetTiers(pensionTiers)
			// pensionCalc.CalculateWithInternalTiers()

			pensionDeduction := pensionCalc.EmployeeContribution
			totalDeductions := taxAmount + pensionDeduction

			netSalary := totalEarnings - totalDeductions

			payroll.PrintPayrollReport(emp, country, originalBasic, totalEarnings, taxAmount, pensionCalc, netSalary, taxCalculationSteps)
		}
	}

}
