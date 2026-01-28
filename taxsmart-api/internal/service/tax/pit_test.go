package tax

import (
	"testing"
)

func TestPITCalculator_Calculate(t *testing.T) {
	calc := NewPITCalculator()

	tests := []struct {
		name        string
		income      float64
		expectedTax float64
	}{
		{
			name:        "Zero income",
			income:      0,
			expectedTax: 0,
		},
		{
			name:        "Below threshold - no tax",
			income:      800_000,
			expectedTax: 0,
		},
		{
			name:   "Just above threshold",
			income: 1_000_000,
			// (1,000,000 - 800,001) * 0.15 = 199,999 * 0.15 = 29,999.85
			// Or simpler: taxable = 200,000, but bracket starts at 800,001
			// Amount in 15% bracket = 1,000,000 - 800,001 = 199,999
			expectedTax: 29999.85,
		},
		{
			name:   "At 3 million",
			income: 3_000_000,
			// 15% bracket: 3,000,000 - 800,001 = 2,199,999 * 0.15 = 329,999.85
			expectedTax: 329999.85,
		},
		{
			name:   "At 5 million",
			income: 5_000_000,
			// 15% bracket: 3,000,000 - 800,001 = 2,199,999 * 0.15 = 329,999.85
			// 18% bracket: 5,000,000 - 3,000,001 = 1,999,999 * 0.18 = 359,999.82
			// Total: 689,999.67
			expectedTax: 689999.67,
		},
		{
			name:   "At 20 million",
			income: 20_000_000,
			// 15%: 2,199,999 * 0.15 = 329,999.85
			// 18%: 8,999,999 * 0.18 = 1,619,999.82
			// 21%: 7,999,999 * 0.21 = 1,679,999.79
			// Total: 3,629,999.46
			expectedTax: 3629999.46,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tax, _ := calc.Calculate(tt.income)

			// Allow small floating point difference
			diff := tax - tt.expectedTax
			if diff < 0 {
				diff = -diff
			}
			if diff > 0.01 {
				t.Errorf("Expected tax %.2f, got %.2f (diff: %.2f)", tt.expectedTax, tax, diff)
			}
		})
	}
}

func TestPITCalculator_ZeroTaxThreshold(t *testing.T) {
	calc := NewPITCalculator()

	// Test various amounts at and below threshold
	testCases := []float64{0, 100_000, 500_000, 800_000}

	for _, income := range testCases {
		tax := calc.CalculateSimple(income)
		if tax != 0 {
			t.Errorf("Income %.0f should have 0 tax, got %.2f", income, tax)
		}
	}
}

func TestPITCalculator_BreakdownCount(t *testing.T) {
	calc := NewPITCalculator()

	// 5 million should span 2 taxable brackets (15% and 18%)
	_, breakdown := calc.Calculate(5_000_000)
	if len(breakdown) != 2 {
		t.Errorf("Expected 2 brackets for 5M income, got %d", len(breakdown))
	}

	// 50 million should span 4 taxable brackets
	_, breakdown = calc.Calculate(50_000_000)
	if len(breakdown) != 4 {
		t.Errorf("Expected 4 brackets for 50M income, got %d", len(breakdown))
	}

	// 100 million should span 5 taxable brackets
	_, breakdown = calc.Calculate(100_000_000)
	if len(breakdown) != 5 {
		t.Errorf("Expected 5 brackets for 100M income, got %d", len(breakdown))
	}
}
