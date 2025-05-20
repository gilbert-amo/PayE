package main

import (
	"PayE/payroll"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Country struct {
	Name        string
	MinimumWage float64
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	countries := make(map[string]Country)
	config := payroll.PayrollConfig{SplitEnabled: false} // Default disabled

	// Country setup
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

		countries[code] = Country{
			Name:        name,
			MinimumWage: minWage,
		}
	}

	if len(countries) == 0 {
		fmt.Println("No countries entered. Exiting.")
		return
	}

	// Configure splitting
	fmt.Print("\nEnable salary splitting when piece-rate exceeds basic? (y/n): ")
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
	var employees []payroll.Employee
	for {
		emp := payroll.Employee{}

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

			emp.PieceRate = append(emp.PieceRate, payroll.PieceRateAggregation{
				Item:     item,
				Rate:     rate,
				Quantity: qty,
			})
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
		totalEarnings := emp.BasicSalary + pieceEarnings

		// Apply splitting if enabled and conditions met
		if config.SplitEnabled && emp.BasicSalary > 0 && pieceEarnings >= emp.BasicSalary {
			splitAmount := pieceEarnings
			emp.BasicSalary += splitAmount * config.BasicSalaryRatio
			emp.Allowance = splitAmount * config.AllowanceRatio
			pieceEarnings = 0 // Reset since we've converted to salary+allowance
			fmt.Printf("\nSplit %.2f piece-rate into:\n- Basic: %.2f\n- Allowance: %.2f\n",
				splitAmount, emp.BasicSalary, emp.Allowance)
		}

		// Minimum wage enforcement
		if emp.BasicSalary > 0 && emp.BasicSalary < country.MinimumWage {
			adjustment := country.MinimumWage - emp.BasicSalary
			emp.BasicSalary = country.MinimumWage
			fmt.Printf("Adjusted basic salary to meet minimum wage (+%.2f)\n", adjustment)
		}

		// Final calculation
		totalEarnings = emp.BasicSalary + emp.Allowance + pieceEarnings

		// Print results
		printPayrollReport(emp, country, pieceEarnings, totalEarnings)
	}
}

func calculatePieceEarnings(pieces []payroll.PieceRateAggregation) float64 {
	total := 0.0
	for _, pw := range pieces {
		total += pw.Rate * pw.Quantity
	}
	return total
}

func printPayrollReport(emp payroll.Employee, country Country, pieceEarnings, total float64) {
	fmt.Printf("\n=== %s ===\n", emp.Name)
	fmt.Printf("Country: %s (Min Wage: %.2f)\n", country.Name, country.MinimumWage)

	if emp.BasicSalary > 0 {
		fmt.Printf("Basic Salary: %.2f\n", emp.BasicSalary)
	}
	if emp.Allowance > 0 {
		fmt.Printf("Allowance: %.2f\n", emp.Allowance)
	}
	if len(emp.PieceRate) > 0 {
		fmt.Println("\nPiece-Rate Work:")
		for _, pw := range emp.PieceRate {
			fmt.Printf("- %s: %.0f × %.2f = %.2f\n",
				pw.Item, pw.Quantity, pw.Rate, pw.Rate*pw.Quantity)
		}
		fmt.Printf("Total Piece-Rate: %.2f\n", pieceEarnings)
	}

	fmt.Printf("\nTOTAL EARNINGS: %.2f\n", total)
	fmt.Println(strings.Repeat("=", 30))
}
