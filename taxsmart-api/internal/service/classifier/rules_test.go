package classifier

import (
	"testing"

	"github.com/taxsmart/taxsmart-api/internal/model"
)

func TestRuleEngine_Classify(t *testing.T) {
	engine := NewRuleEngine()

	tests := []struct {
		name        string
		description string
		txType      string
		expected    model.Category
	}{
		// Employment patterns
		{
			name:        "Salary payment",
			description: "SALARY FOR JANUARY 2026",
			txType:      "credit",
			expected:    model.CategoryEmployment,
		},
		{
			name:        "Payroll",
			description: "PAYROLL/ABC COMPANY LTD",
			txType:      "credit",
			expected:    model.CategoryEmployment,
		},

		// Freelance patterns
		{
			name:        "Upwork payment",
			description: "UPWORK ESCROW INC",
			txType:      "credit",
			expected:    model.CategoryFreelance,
		},
		{
			name:        "Payoneer transfer",
			description: "TRF FROM PAYONEER",
			txType:      "credit",
			expected:    model.CategoryFreelance,
		},

		// Crypto patterns
		{
			name:        "Binance withdrawal",
			description: "BINANCE FIAT WITHDRAWAL",
			txType:      "credit",
			expected:    model.CategoryCrypto,
		},
		{
			name:        "Luno trade",
			description: "LUNO NIGERIA LTD",
			txType:      "credit",
			expected:    model.CategoryCrypto,
		},
		{
			name:        "Quidax",
			description: "QUIDAX TECHNOLOGIES",
			txType:      "credit",
			expected:    model.CategoryCrypto,
		},

		// Rent patterns
		{
			name:        "Rent expense",
			description: "RENT PAYMENT TO LANDLORD",
			txType:      "debit",
			expected:    model.CategoryRentExpense,
		},

		// General expense
		{
			name:        "POS purchase",
			description: "POS/SHOPRITE LAGOS",
			txType:      "debit",
			expected:    model.CategoryExpense,
		},
		{
			name:        "ATM withdrawal",
			description: "ATM WITHDRAWAL",
			txType:      "debit",
			expected:    model.CategoryExpense,
		},

		// Uncategorized
		{
			name:        "Unknown transfer",
			description: "NIP TRF TO JOHN DOE",
			txType:      "debit",
			expected:    model.CategoryTransfer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.Classify(tt.description, tt.txType)
			if result.Category != tt.expected {
				t.Errorf("Expected category %s, got %s", tt.expected, result.Category)
			}
			if result.Method != "rules" {
				t.Errorf("Expected method 'rules', got '%s'", result.Method)
			}
		})
	}
}
