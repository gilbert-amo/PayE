package pension

import (
	_ "math"
	_ "strconv"

	"github.com/gilbert-amo/PayE/types"
)

type Calculator struct {
	Employee             *types.Employee
	EmployeeContribution float64
	EmployerContribution float64
	TotalMandatory       float64
	TierContributions    map[string]float64
	Tiers                []types.Tier
}

func NewCalculator(employee *types.Employee) *Calculator {
	return &Calculator{
		Employee:          employee,
		TierContributions: make(map[string]float64),
		Tiers:             make([]types.Tier, 0),
	}
}

// AddTier adds a new tier to the pension calculator
func (pc *Calculator) AddTier(name string, percentage float64) {
	// Validate percentage is between 0 and 1
	if percentage < 0 || percentage > 1 {
		return
	}

	// Check if tier name already exists
	for i, tier := range pc.Tiers {
		if tier.Name == name {
			// Update existing tier
			pc.Tiers[i].Percentage = percentage
			return
		}
	}

	// Add new tier
	pc.Tiers = append(pc.Tiers, types.Tier{
		Name:       name,
		Percentage: percentage,
	})
}

// RemoveTier removes a tier by name
func (pc *Calculator) RemoveTier(name string) {
	for i, tier := range pc.Tiers {
		if tier.Name == name {
			// Remove tier by appending slices
			pc.Tiers = append(pc.Tiers[:i], pc.Tiers[i+1:]...)
			break
		}
	}
}

// SetTiers sets the complete list of tiers
func (pc *Calculator) SetTiers(tiers []types.Tier) {
	pc.Tiers = make([]types.Tier, len(tiers))
	copy(pc.Tiers, tiers)
}

// GetTiers returns the current list of tiers
func (pc *Calculator) GetTiers() []types.Tier {
	return pc.Tiers
}

// ValidateTiers checks if the total percentage equals 1.0 (100%)
func (pc *Calculator) ValidateTiers() bool {
	total := 0.0
	for _, tier := range pc.Tiers {
		total += tier.Percentage
	}
	// Allow for small floating point precision differences
	return total >= 0.999 && total <= 1.001
}

func (pc *Calculator) Calculate(tiers []types.Tier) {
	// Use provided tiers if available, otherwise use internal tiers
	tiersToUse := tiers
	if len(tiers) == 0 {
		tiersToUse = pc.Tiers
	}

	basicSalary := pc.Employee.BasicSalary

	// Calculate contributions
	pc.EmployeeContribution = basicSalary * 0.055 // 5.5%
	pc.EmployerContribution = basicSalary * 0.13  // 13%
	pc.TotalMandatory = pc.EmployeeContribution + pc.EmployerContribution

	// Clear previous tier contributions
	pc.TierContributions = make(map[string]float64)

	// Calculate tier allocations
	for _, tier := range tiersToUse {
		pc.TierContributions[tier.Name] = pc.TotalMandatory * tier.Percentage
	}
}

// CalculateWithInternalTiers calculates pension using the internal tiers
func (pc *Calculator) CalculateWithInternalTiers() {
	pc.Calculate(pc.Tiers)
}

func (pc *Calculator) GetContributionBreakdown() map[string]float64 {
	return map[string]float64{
		"Basic Salary":          pc.Employee.BasicSalary,
		"Employee Contribution": pc.EmployeeContribution,
		"Employer Contribution": pc.EmployerContribution,
		"Total Mandatory":       pc.TotalMandatory,
	}
}

func (pc *Calculator) GetTierBreakdown() map[string]float64 {
	return pc.TierContributions
}
