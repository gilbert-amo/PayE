package pension

import (
	"github.com/gilbert-amo/PayE/types"
	_ "math"
)

type Calculator struct {
	Employee             *types.Employee
	EmployeeContribution float64
	EmployerContribution float64
	TotalMandatory       float64
	TierContributions    map[string]float64
}

func NewCalculator(employee *types.Employee) *Calculator {
	return &Calculator{
		Employee:          employee,
		TierContributions: make(map[string]float64),
	}
}

func (pc *Calculator) Calculate(tiers []types.Tier) {
	basicSalary := pc.Employee.BasicSalary

	// Calculate contributions
	pc.EmployeeContribution = basicSalary * 0.055 // 5.5%
	pc.EmployerContribution = basicSalary * 0.13  // 13%
	pc.TotalMandatory = pc.EmployeeContribution + pc.EmployerContribution

	// Calculate tier allocations
	for _, tier := range tiers {
		pc.TierContributions[tier.Name] = pc.TotalMandatory * tier.Percentage
	}
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
