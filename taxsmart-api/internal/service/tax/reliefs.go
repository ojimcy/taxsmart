package tax

import (
	"github.com/taxsmart/taxsmart-api/internal/model"
)

// ReliefCalculator calculates tax reliefs based on Nigeria 2026 rules
type ReliefCalculator struct{}

// NewReliefCalculator creates a new relief calculator
func NewReliefCalculator() *ReliefCalculator {
	return &ReliefCalculator{}
}

// CalculateReliefs computes all applicable reliefs
func (c *ReliefCalculator) CalculateReliefs(input model.ReliefInput) (float64, map[string]float64) {
	reliefs := make(map[string]float64)
	var total float64

	// Rent Relief: 20% of annual rent, capped at ₦500,000
	if input.AnnualRent > 0 {
		rentRelief := input.AnnualRent * RentReliefPercentage
		if rentRelief > RentReliefCap {
			rentRelief = RentReliefCap
		}
		reliefs["rent_relief"] = rentRelief
		total += rentRelief
	}

	// Pension contributions - fully deductible
	if input.PensionContribution > 0 {
		reliefs["pension"] = input.PensionContribution
		total += input.PensionContribution
	}

	// NHIS contributions - fully deductible
	if input.NHISContribution > 0 {
		reliefs["nhis"] = input.NHISContribution
		total += input.NHISContribution
	}

	// NHF contributions - fully deductible
	if input.NHFContribution > 0 {
		reliefs["nhf"] = input.NHFContribution
		total += input.NHFContribution
	}

	return total, reliefs
}

// CalculateRentRelief calculates rent relief only
func (c *ReliefCalculator) CalculateRentRelief(annualRent float64) float64 {
	if annualRent <= 0 {
		return 0
	}
	relief := annualRent * RentReliefPercentage
	if relief > RentReliefCap {
		return RentReliefCap
	}
	return relief
}
