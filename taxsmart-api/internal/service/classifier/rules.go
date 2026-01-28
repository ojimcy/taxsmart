package classifier

import (
	"regexp"
	"strings"

	"github.com/taxsmart/taxsmart-api/internal/model"
)

// RuleEngine classifies transactions based on pattern matching
type RuleEngine struct {
	incomePatterns  map[string]model.Category
	expensePatterns map[string]model.Category
}

// NewRuleEngine creates a new rule-based classifier
func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		incomePatterns: map[string]model.Category{
			// Employment
			"SALARY":       model.CategoryEmployment,
			"PAY":          model.CategoryEmployment,
			"WAGES":        model.CategoryEmployment,
			"PAYROLL":      model.CategoryEmployment,
			"REMUNERATION": model.CategoryEmployment,

			// Freelance
			"UPWORK":     model.CategoryFreelance,
			"FIVERR":     model.CategoryFreelance,
			"PAYONEER":   model.CategoryFreelance,
			"WISE":       model.CategoryFreelance,
			"TOPTAL":     model.CategoryFreelance,
			"FREELANCER": model.CategoryFreelance,
			"CONTRA":     model.CategoryFreelance,

			// Crypto
			"BINANCE":     model.CategoryCrypto,
			"LUNO":        model.CategoryCrypto,
			"QUIDAX":      model.CategoryCrypto,
			"PAXFUL":      model.CategoryCrypto,
			"COINBASE":    model.CategoryCrypto,
			"KRAKEN":      model.CategoryCrypto,
			"BYBIT":       model.CategoryCrypto,
			"KUCOIN":      model.CategoryCrypto,
			"ROQQU":       model.CategoryCrypto,
			"PATRICIA":    model.CategoryCrypto,
			"BUSHA":       model.CategoryCrypto,
			"YELLOW CARD": model.CategoryCrypto,
			"NOONES":      model.CategoryCrypto,

			// Investment
			"DIVIDEND":          model.CategoryInvestment,
			"INVESTMENT RETURN": model.CategoryInvestment,
			"BAMBOO":            model.CategoryInvestment,
			"RISEVEST":          model.CategoryInvestment,
			"TROVE":             model.CategoryInvestment,
			"CHAKA":             model.CategoryInvestment,

			// Interest
			"INTEREST":   model.CategoryInterest,
			"INT CREDIT": model.CategoryInterest,

			// Rental
			"RENT RECEIVED": model.CategoryRental,
			"TENANT":        model.CategoryRental,
			"RENTAL INCOME": model.CategoryRental,
		},
		expensePatterns: map[string]model.Category{
			// Rent expense
			"RENT PAYMENT":  model.CategoryRentExpense,
			"LANDLORD":      model.CategoryRentExpense,
			"HOUSE RENT":    model.CategoryRentExpense,
			"ACCOMMODATION": model.CategoryRentExpense,
		},
	}
}

// Classify attempts to classify a transaction based on its description
func (r *RuleEngine) Classify(description string, txType string) model.ClassificationResult {
	upperDesc := strings.ToUpper(description)

	// Check income patterns for credits
	if txType == "credit" {
		for pattern, category := range r.incomePatterns {
			if strings.Contains(upperDesc, pattern) {
				return model.ClassificationResult{
					Category:   category,
					Confidence: 0.85,
					Method:     "rules",
				}
			}
		}
	}

	// Check expense patterns for debits
	if txType == "debit" {
		for pattern, category := range r.expensePatterns {
			if strings.Contains(upperDesc, pattern) {
				return model.ClassificationResult{
					Category:   category,
					Confidence: 0.85,
					Method:     "rules",
				}
			}
		}

		// General expenses
		expenseKeywords := []string{"POS", "ATM", "WITHDRAWAL", "TRANSFER", "PAYMENT", "PURCHASE"}
		for _, keyword := range expenseKeywords {
			if strings.Contains(upperDesc, keyword) {
				return model.ClassificationResult{
					Category:   model.CategoryExpense,
					Confidence: 0.70,
					Method:     "rules",
				}
			}
		}
	}

	// Check for transfers (could be income or expense)
	transferPatterns := regexp.MustCompile(`(?i)(NIP|TRANSFER|TRF)`)
	if transferPatterns.MatchString(description) {
		if txType == "credit" {
			return model.ClassificationResult{
				Category:   model.CategoryUncategorized,
				Confidence: 0.50,
				Method:     "rules",
			}
		}
		return model.ClassificationResult{
			Category:   model.CategoryTransfer,
			Confidence: 0.60,
			Method:     "rules",
		}
	}

	return model.ClassificationResult{
		Category:   model.CategoryUncategorized,
		Confidence: 0.0,
		Method:     "rules",
	}
}
