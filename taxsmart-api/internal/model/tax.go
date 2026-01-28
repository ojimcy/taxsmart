package model

import (
	"time"

	"github.com/google/uuid"
)

// TaxReport represents a calculated tax report for a user
type TaxReport struct {
	ID      uuid.UUID `json:"id"`
	UserID  uuid.UUID `json:"user_id"`
	TaxYear int       `json:"tax_year"`

	// Income totals
	TotalIncome      float64 `json:"total_income"`
	EmploymentIncome float64 `json:"employment_income"`
	FreelanceIncome  float64 `json:"freelance_income"`
	RentalIncome     float64 `json:"rental_income"`
	InvestmentIncome float64 `json:"investment_income"`
	CryptoIncome     float64 `json:"crypto_income"`
	OtherIncome      float64 `json:"other_income"`

	// Reliefs
	RentRelief       float64 `json:"rent_relief"`
	PensionDeduction float64 `json:"pension_deduction"`
	NHISDeduction    float64 `json:"nhis_deduction"`
	NHFDeduction     float64 `json:"nhf_deduction"`
	TotalReliefs     float64 `json:"total_reliefs"`

	// Tax calculations
	TaxableIncome float64 `json:"taxable_income"`
	PITAmount     float64 `json:"pit_amount"`
	CGTAmount     float64 `json:"cgt_amount"`
	TotalTax      float64 `json:"total_tax"`

	// Breakdown details
	Breakdown *TaxBreakdown `json:"breakdown,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaxBreakdown provides detailed breakdown of tax calculations
type TaxBreakdown struct {
	PITBreakdown     []BracketDetail    `json:"pit_breakdown"`
	IncomeByCategory map[string]float64 `json:"income_by_category"`
	ReliefsApplied   map[string]float64 `json:"reliefs_applied"`
}

// BracketDetail shows tax applied in each bracket
type BracketDetail struct {
	BracketMin       float64 `json:"bracket_min"`
	BracketMax       float64 `json:"bracket_max"`
	Rate             float64 `json:"rate"`
	TaxableInBracket float64 `json:"taxable_in_bracket"`
	TaxAmount        float64 `json:"tax_amount"`
}

// TaxBracket represents a single tax bracket
type TaxBracket struct {
	Min  float64
	Max  float64
	Rate float64
}

// ReliefInput represents user-provided relief information
type ReliefInput struct {
	AnnualRent          float64 `json:"annual_rent"`
	PensionContribution float64 `json:"pension_contribution"`
	NHISContribution    float64 `json:"nhis_contribution"`
	NHFContribution     float64 `json:"nhf_contribution"`
}

// TaxCalculationRequest represents a request to calculate tax
type TaxCalculationRequest struct {
	UserID       uuid.UUID     `json:"user_id"`
	TaxYear      int           `json:"tax_year"`
	Transactions []Transaction `json:"transactions"`
	Reliefs      ReliefInput   `json:"reliefs"`
}
