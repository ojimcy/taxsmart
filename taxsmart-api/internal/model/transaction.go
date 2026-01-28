package model

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a single financial transaction
type Transaction struct {
	ID              uuid.UUID `json:"id"`
	UploadID        uuid.UUID `json:"upload_id"`
	UserID          uuid.UUID `json:"user_id"`
	TransactionDate time.Time `json:"transaction_date"`
	Description     string    `json:"description"`
	Amount          float64   `json:"amount"`
	TransactionType string    `json:"transaction_type"` // "credit" or "debit"
	Category        Category  `json:"category"`
	Confidence      float64   `json:"confidence"`
	IsManual        bool      `json:"is_manual"`
	RawData         string    `json:"raw_data,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// Category represents transaction classification
type Category string

const (
	CategoryEmployment    Category = "employment_income"
	CategoryFreelance     Category = "freelance_income"
	CategoryRental        Category = "rental_income"
	CategoryInvestment    Category = "investment_income"
	CategoryCrypto        Category = "crypto_income"
	CategoryInterest      Category = "interest_income"
	CategoryOtherIncome   Category = "other_income"
	CategoryExpense       Category = "expense"
	CategoryRentExpense   Category = "rent_expense"
	CategoryTransfer      Category = "transfer"
	CategoryUncategorized Category = "uncategorized"
)

// IsIncome returns true if the category represents income
func (c Category) IsIncome() bool {
	incomeCategories := map[Category]bool{
		CategoryEmployment:  true,
		CategoryFreelance:   true,
		CategoryRental:      true,
		CategoryInvestment:  true,
		CategoryCrypto:      true,
		CategoryInterest:    true,
		CategoryOtherIncome: true,
	}
	return incomeCategories[c]
}

// IsTaxable returns true if the category is taxable
func (c Category) IsTaxable() bool {
	return c.IsIncome()
}

// ParsedTransaction represents a transaction parsed from a file before classification
type ParsedTransaction struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"` // "credit" or "debit"
	Balance     float64   `json:"balance,omitempty"`
	Reference   string    `json:"reference,omitempty"`
}

// ClassificationResult represents the result of classifying a transaction
type ClassificationResult struct {
	Category   Category `json:"category"`
	Confidence float64  `json:"confidence"`
	Method     string   `json:"method"` // "ai" or "rules"
}
