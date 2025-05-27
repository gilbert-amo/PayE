# PayE - Payroll and Pension Management System

PayE is a comprehensive payroll and pension management system written in Go. It provides functionality for calculating salaries, taxes, and pension contributions for employees across different countries.

## Features

- Multi-country support with configurable tax brackets
- Flexible salary calculation supporting both basic salary and piece-rate work
- Pension contribution calculations with tiered allocation
- Detailed payroll reporting
- Configurable salary splitting between basic salary and allowances

## Package Structure

### Main Package
The main package provides the command-line interface for the PayE system. It handles:
- Country setup and configuration
- Employee data input
- Tax bracket configuration
- Salary splitting configuration
- Payroll report generation

### Types Package
Contains the core data structures:
- `Employee`: Represents an employee with basic salary, piece-rate work, and country information
- `PieceRateAggregation`: Represents piece-rate work items with rate and quantity
- `TaxBracket`: Defines tax thresholds and rates
- `Tier`: Represents pension contribution tiers
- `Country`: Contains country-specific information including minimum wage and tax brackets

### Payroll Package
Handles salary and tax calculations:
- Basic salary and piece-rate calculations
- Tax calculations based on country-specific brackets
- Salary breakdown generation
- Detailed payroll reporting

### Pension Package
Manages pension calculations:
- Employee and employer contribution calculations
- Tier-based contribution allocation
- Contribution breakdown reporting

## Usage

1. Set up countries with their respective tax brackets and minimum wages
2. Configure salary splitting if needed
3. Add employees with their basic salary and/or piece-rate work
4. Generate detailed payroll reports including:
   - Basic salary breakdown
   - Piece-rate earnings
   - Tax calculations
   - Pension contributions
   - Net salary calculations

## Configuration

### Country Setup
- Country code (3 letters)
- Country name
- Minimum wage
- Tax brackets with thresholds and rates

### Salary Configuration
- Basic salary
- Piece-rate work items
- Salary splitting (optional)
  - Basic salary ratio
  - Allowance ratio

### Pension Configuration
- Tier 1: 13.5%
- Tier 2: 55%
- Tier 3: 31.5%

## Example

```go
// Create a new employee
employee := types.Employee{
    Name:        "John Doe",
    BasicSalary: 5000,
    CountryCode: "USA",
    PieceRate: []types.PieceRateAggregation{
        {
            Item:     "Project A",
            Rate:     100,
            Quantity: 5,
        },
    },
}

// Calculate salary
salary := payroll.CalculateSalary(&employee)

// Calculate pension
pensionCalc := pension.NewCalculator(&employee)
pensionCalc.Calculate(tiers)

// Generate report
payroll.PrintPayrollReport(employee, country, originalBasic, gross, tax, pensionCalc, net)
```

## Dependencies

- Go 1.23.5 or higher

## License

[Add your license information here] 