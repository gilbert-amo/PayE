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
		if salary == 0 {
			fmt.Println("\nAdding piece-rate work (required for employees with no basic salary)")
		} else {
			fmt.Println("\nAdding optional piece-rate bonus work")
		}

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

		// Validate at least some piece-rate work if no basic salary
		if emp.BasicSalary == 0 && len(emp.PieceRate) == 0 {
			fmt.Println("Error: Piece-rate employees must have at least one piece-rate item")
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

	// Process and print results
	fmt.Println("\n=== Payroll Results ===")
	for _, emp := range employees {
		country := countries[emp.CountryCode]
		totalEarnings := emp.BasicSalary

		// Calculate piece-rate earnings
		pieceEarnings := 0.0
		if len(emp.PieceRate) > 0 {
			fmt.Printf("\n%s's Piece-Rate Aggregation:\n", emp.Name)
			for _, pw := range emp.PieceRate {
				earning := pw.Rate * pw.Quantity
				fmt.Printf("- %s: %.0f @ %.2f = %.2f\n", pw.Item, pw.Quantity, pw.Rate, earning)
				pieceEarnings += earning
			}
			fmt.Printf("Total Piece-Rate Earnings: %.2f\n", pieceEarnings)
		}

		// Apply minimum wage check to basic salary
		if emp.BasicSalary > 0 && emp.BasicSalary < country.MinimumWage {
			fmt.Printf("Adjusting basic salary to meet minimum wage (%.2f -> %.2f)\n",
				emp.BasicSalary, country.MinimumWage)
			emp.BasicSalary = country.MinimumWage

			//fmt.Printf("Basic Salary for", "%s: %.2f\n", emp.Name, country.MinimumWage)
		}

		totalEarnings = country.MinimumWage + pieceEarnings

		// For piece-rate only employees, ensure total meets minimum wage
		if emp.BasicSalary == 0 && totalEarnings < country.MinimumWage {
			shortfall := country.MinimumWage - totalEarnings
			fmt.Printf("Warning: Total earnings (%.2f) below minimum wage (%.2f)\n",
				totalEarnings, country.MinimumWage)
			fmt.Printf("Shortfall: %.2f (must be paid to employee)\n", shortfall)
			totalEarnings = country.MinimumWage
		}

		fmt.Printf("\n%s's Total Earnings: %.2f\n", emp.Name, totalEarnings)
		fmt.Println(strings.Repeat("-", 40))
	}
}
