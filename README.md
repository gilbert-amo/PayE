# PayE - Payroll and Pension Management System

PayE is a comprehensive payroll and pension management system written in Go. It provides functionality for calculating salaries, taxes, and pension contributions for employees across different countries.

## Features

- Multi-country support with configurable tax brackets
- Flexible salary calculation supporting both basic salary and piece-rate work
- **Recursive PAYE Tax Calculation**: Tax is calculated as a percentage of chargeable income and subtracted from the basic salary recursively until salary reaches zero
- **Automatic Minimum Wage Protection**: Minimum wage is automatically taxed at 0%
- **Dynamic Pension Tier Management**: Add, remove, and update pension tiers at runtime
- Pension contribution calculations with tiered allocation
- Detailed payroll reporting with step-by-step tax calculation breakdown
- Configurable salary splitting between basic salary and allowances

## Quick Start

1. **Clone and Build**
   ```bash
   git clone <repository>
   cd PayE
   go build
   ```

2. **Run the Program**
   ```bash
   ./PayE
   ```

3. **Follow the Interactive Setup**
   - Configure pension tiers (e.g., Conservative 40%, Balanced 35%, Growth 25%)
   - Set up countries with tax brackets
   - Add employees and their salary information
   - Generate detailed payroll reports

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
- **Recursive PAYE Tax Calculation**: Implements the recursive tax calculation system where tax is calculated as a percentage of chargeable income and subtracted from the basic salary until it reaches zero
- Tax calculation step tracking and reporting
- Salary breakdown generation
- Detailed payroll reporting with tax calculation breakdown

### Pension Package
Manages pension calculations with dynamic tier management:
- Employee and employer contribution calculations
- **Dynamic Tier Management**: Add, remove, and update pension tiers at runtime
- Tier-based contribution allocation with validation
- Contribution breakdown reporting
- **Tier Rate Entry**: Interactive tier setup with percentage validation
- Support for both external and internal tier management

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
- **Dynamic Tier Setup**: Configure custom pension tiers during runtime
- **Percentage Validation**: Ensures total tier percentages equal 100%
- **Example Default Tiers**:
  - Tier 1: 13.5%
  - Tier 2: 55%
  - Tier 3: 31.5%
- **Custom Tier Names**: Use descriptive names like "Conservative", "Balanced", "Growth"

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

// Calculate recursive PAYE tax
taxAmount, taxSteps := payroll.CalculateTaxRecursive(employee.BasicSalary, country.TaxBrackets)

// Calculate pension with dynamic tier management
pensionCalc := pension.NewCalculator(&employee)

// Option 1: Use external tiers
pensionCalc.Calculate(tiers)

// Option 2: Use internal tiers
pensionCalc.AddTier("Conservative", 0.40)
pensionCalc.AddTier("Balanced", 0.35)
pensionCalc.AddTier("Growth", 0.25)
pensionCalc.CalculateWithInternalTiers()

// Generate report with tax calculation breakdown
payroll.PrintPayrollReport(employee, country, originalBasic, gross, taxAmount, pensionCalc, net, taxSteps)
```

## Progressive Subtraction Tax Calculation

The system implements a progressive subtraction tax calculation where threshold amounts are progressively subtracted from the basic salary. **Minimum wage is automatically taxed at 0%** before any other calculations:

1. **Minimum Wage Protection**: Minimum wage is automatically set as a 0% tax bracket
2. **Progressive Subtraction**: Each threshold amount is subtracted from the remaining salary
3. **Tax Calculation**: Each threshold amount × its tax rate = tax for that amount
4. **Salary Reduction**: Salary decreases progressively with each threshold subtraction
5. **Remaining Salary Handling**: When remaining salary < threshold, apply threshold rate to remaining amount and end calculation
6. **Total Tax**: Sum of all individual threshold taxes
7. **Net Salary**: Original salary minus total tax

### Example Tax Bracket Setup:
```
Minimum Wage: 800 (0% tax) - Automatically applied
Threshold 110: 5% tax = 5.50 tax
Threshold 130: 10% tax = 13.00 tax
Threshold 3166: 17.5% tax = 554.05 tax
Threshold 16000: 25% tax = 4,000.00 tax
Threshold 30520: 30% tax = 9,156.00 tax
```

**Note**: The minimum wage threshold is automatically added as a 0% tax bracket, ensuring that minimum wage is always tax-free.

### Progressive Subtraction Example:
Starting salary: 50,416.67
- Subtract 800 (min wage) → Remaining: 49,616.67, Tax: 0.00
- Subtract 110 → Remaining: 49,506.67, Tax: 5.50
- Subtract 130 → Remaining: 49,376.67, Tax: 13.00
- Subtract 3166 → Remaining: 46,210.67, Tax: 554.05
- Subtract 16000 → Remaining: 30,210.67, Tax: 4,000.00
- Subtract 30520 → Remaining: 0.00, Tax: 9,156.00

**Total Tax: 13,728.55**

**Key Features**:
- Minimum wage (800) is automatically taxed at 0%
- When remaining salary < threshold, the threshold rate is applied to the remaining amount
- Calculation ends automatically when all salary is taxed

## Pension Tier Management

The system provides flexible pension tier management with the following features:

### Dynamic Tier Operations
- **AddTier(name, percentage)**: Add new tiers or update existing ones
- **RemoveTier(name)**: Remove tiers by name
- **SetTiers(tiers)**: Set complete tier list at once
- **GetTiers()**: Retrieve current tier configuration
- **ValidateTiers()**: Ensure total percentage equals 100%

### Tier Calculation Methods
- **Calculate(tiers)**: Use external tier list
- **CalculateWithInternalTiers()**: Use internally stored tiers
- **GetTierBreakdown()**: Get contribution allocation per tier

### Example Usage
```go
// Create calculator
calc := pension.NewCalculator(employee)

// Add tiers dynamically
calc.AddTier("Conservative", 0.40)  // 40%
calc.AddTier("Balanced", 0.35)      // 35%
calc.AddTier("Growth", 0.25)        // 25%

// Validate and calculate
if calc.ValidateTiers() {
    calc.CalculateWithInternalTiers()
    breakdown := calc.GetTierBreakdown()
}
```

### Interactive Tier Setup
During program execution, users can:
- Enter custom tier names
- Specify percentage allocations
- Get real-time validation feedback
- Use default tiers if validation fails

## Enhanced Tax Calculation Logic

The system now implements enhanced tax calculation logic that automatically handles edge cases:

### Automatic Minimum Wage Protection
- Minimum wage is automatically added as a 0% tax bracket
- No need to manually configure minimum wage tax brackets
- Ensures minimum wage is always tax-free

### Smart Remaining Salary Handling
- **Full Threshold**: When `remainingSalary >= threshold`, tax at threshold rate and subtract threshold amount
- **Partial Threshold**: When `remainingSalary < threshold`, apply threshold rate to remaining salary and end calculation
- **Automatic Termination**: Calculation ends when all salary is processed

### Benefits
- **Simplified Setup**: No need to manually add 0% minimum wage brackets
- **Accurate Calculation**: Handles partial threshold amounts correctly
- **Efficient Processing**: Automatically stops when no more salary remains
- **Flexible Configuration**: Works with any tax bracket configuration

## Dependencies

- Go 1.23.5 or higher

## License

[Add your license information here] 