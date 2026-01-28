package tax

import (
	"github.com/taxsmart/taxsmart-api/internal/model"
)

// PITCalculator calculates Personal Income Tax based on Nigeria 2026 brackets
type PITCalculator struct {
	brackets []model.TaxBracket
}

// NewPITCalculator creates a new PIT calculator with 2026 brackets
func NewPITCalculator() *PITCalculator {
	return &PITCalculator{
		brackets: PITBrackets2026,
	}
}

// Calculate computes PIT for the given annual income
// Returns total tax and breakdown by bracket
func (c *PITCalculator) Calculate(annualIncome float64) (float64, []model.BracketDetail) {
	if annualIncome <= 0 {
		return 0, nil
	}

	// First ₦800,000 is tax-free
	if annualIncome <= TaxFreeThreshold {
		return 0, nil
	}

	var totalTax float64
	var breakdown []model.BracketDetail

	// Tax is calculated on income above the tax-free threshold
	// Using cumulative bracket approach
	taxableRemaining := annualIncome

	for _, bracket := range c.brackets {
		// Skip if income doesn't reach this bracket
		if taxableRemaining <= bracket.Min {
			break
		}

		// Skip 0% bracket
		if bracket.Rate == 0 {
			continue
		}

		// Calculate how much falls in this bracket
		var amountInBracket float64
		if taxableRemaining > bracket.Max {
			amountInBracket = bracket.Max - bracket.Min
		} else {
			amountInBracket = taxableRemaining - bracket.Min
		}

		if amountInBracket <= 0 {
			continue
		}

		taxAmount := amountInBracket * bracket.Rate
		totalTax += taxAmount

		breakdown = append(breakdown, model.BracketDetail{
			BracketMin:       bracket.Min,
			BracketMax:       bracket.Max,
			Rate:             bracket.Rate,
			TaxableInBracket: amountInBracket,
			TaxAmount:        taxAmount,
		})
	}

	return totalTax, breakdown
}

// CalculateSimple returns just the total tax amount
func (c *PITCalculator) CalculateSimple(annualIncome float64) float64 {
	tax, _ := c.Calculate(annualIncome)
	return tax
}
