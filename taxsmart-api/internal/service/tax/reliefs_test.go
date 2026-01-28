package tax

import (
	"testing"

	"github.com/taxsmart/taxsmart-api/internal/model"
)

func TestReliefCalculator_CalculateRentRelief(t *testing.T) {
	calc := NewReliefCalculator()

	tests := []struct {
		name     string
		rent     float64
		expected float64
	}{
		{
			name:     "No rent",
			rent:     0,
			expected: 0,
		},
		{
			name:     "Low rent - 20% applies",
			rent:     1_000_000,
			expected: 200_000, // 20% of 1M
		},
		{
			name:     "High rent - capped at 500k",
			rent:     5_000_000,
			expected: 500_000, // Cap applies
		},
		{
			name:     "Exactly at cap threshold",
			rent:     2_500_000,
			expected: 500_000, // 20% of 2.5M = 500k (exactly at cap)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			relief := calc.CalculateRentRelief(tt.rent)
			if relief != tt.expected {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, relief)
			}
		})
	}
}

func TestReliefCalculator_CalculateReliefs(t *testing.T) {
	calc := NewReliefCalculator()

	input := model.ReliefInput{
		AnnualRent:          2_000_000,
		PensionContribution: 100_000,
		NHISContribution:    50_000,
		NHFContribution:     25_000,
	}

	total, reliefs := calc.CalculateReliefs(input)

	// Rent relief: 20% of 2M = 400,000
	expectedRent := 400_000.0
	if reliefs["rent_relief"] != expectedRent {
		t.Errorf("Expected rent relief %.2f, got %.2f", expectedRent, reliefs["rent_relief"])
	}

	// Pension: 100,000
	if reliefs["pension"] != 100_000 {
		t.Errorf("Expected pension 100,000, got %.2f", reliefs["pension"])
	}

	// Total: 400k + 100k + 50k + 25k = 575,000
	expectedTotal := 575_000.0
	if total != expectedTotal {
		t.Errorf("Expected total %.2f, got %.2f", expectedTotal, total)
	}
}
