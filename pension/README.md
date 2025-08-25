# Pension Calculator with Tier Rate Entry

This package provides a flexible pension calculator that allows for dynamic tier rate entry and management.

## Features

- **Dynamic Tier Management**: Add, remove, and update pension tiers at runtime
- **Validation**: Ensures tier percentages total to 100%
- **Flexible Calculation**: Use external tiers or internal tiers stored in the calculator
- **Real-time Updates**: Modify tier rates without recreating the calculator

## Usage

### Basic Usage

```go
import "github.com/gilbert-amo/PayE/pension"

// Create a pension calculator
calc := pension.NewCalculator(employee)

// Add tiers
calc.AddTier("Tier 1", 0.135)  // 13.5%
calc.AddTier("Tier 2", 0.55)   // 55%
calc.AddTier("Tier 3", 0.315)  // 31.5%

// Calculate pension
calc.CalculateWithInternalTiers()
```

### Tier Management Methods

#### AddTier(name string, percentage float64)
Adds a new tier or updates an existing one with the same name.

```go
calc.AddTier("Conservative", 0.40)  // 40%
calc.AddTier("Balanced", 0.35)      // 35%
calc.AddTier("Growth", 0.25)        // 25%
```

#### RemoveTier(name string)
Removes a tier by name.

```go
calc.RemoveTier("Growth")
```

#### SetTiers(tiers []types.Tier)
Sets the complete list of tiers at once.

```go
tiers := []types.Tier{
    {Name: "Tier 1", Percentage: 0.135},
    {Name: "Tier 2", Percentage: 0.55},
    {Name: "Tier 3", Percentage: 0.315},
}
calc.SetTiers(tiers)
```

#### GetTiers() []types.Tier
Returns the current list of tiers.

```go
currentTiers := calc.GetTiers()
for _, tier := range currentTiers {
    fmt.Printf("%s: %.1f%%\n", tier.Name, tier.Percentage*100)
}
```

#### ValidateTiers() bool
Checks if the total percentage equals 100%.

```go
if calc.ValidateTiers() {
    fmt.Println("Tiers are valid")
} else {
    fmt.Println("Tiers do not total 100%")
}
```

### Calculation Methods

#### Calculate(tiers []types.Tier)
Calculates pension using provided tiers.

```go
externalTiers := []types.Tier{
    {Name: "Custom Tier", Percentage: 1.0},
}
calc.Calculate(externalTiers)
```

#### CalculateWithInternalTiers()
Calculates pension using the internal tiers stored in the calculator.

```go
calc.CalculateWithInternalTiers()
```

### Example Programs

- **main.go**: Interactive tier setup during program execution
- **example.go**: Demonstrates all tier management methods

## Validation Rules

- Tier percentages must be between 0 and 1 (0% to 100%)
- Total percentage must equal 1.0 (100%)
- Small floating-point precision differences are allowed (±0.001)

## Error Handling

- Invalid percentages (negative or > 100%) are ignored
- Tiers that would exceed 100% total are rejected
- The calculator gracefully handles missing or invalid tiers

## Integration

The pension calculator integrates seamlessly with the existing PayE system:

- Works with the existing `types.Employee` structure
- Compatible with the current payroll calculation flow
- Maintains backward compatibility with hardcoded tiers
