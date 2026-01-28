package tax

import (
	"time"

	"github.com/google/uuid"
	"github.com/taxsmart/taxsmart-api/internal/model"
)

// Engine orchestrates all tax calculations
type Engine struct {
	pitCalculator    *PITCalculator
	reliefCalculator *ReliefCalculator
}

// NewEngine creates a new tax calculation engine
func NewEngine() *Engine {
	return &Engine{
		pitCalculator:    NewPITCalculator(),
		reliefCalculator: NewReliefCalculator(),
	}
}

// CalculateTax computes the full tax report for given transactions and reliefs
func (e *Engine) CalculateTax(req model.TaxCalculationRequest) (*model.TaxReport, error) {
	report := &model.TaxReport{
		ID:        uuid.New(),
		UserID:    req.UserID,
		TaxYear:   req.TaxYear,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Aggregate income by category
	incomeByCategory := make(map[string]float64)
	for _, tx := range req.Transactions {
		if tx.Category.IsIncome() && tx.TransactionType == "credit" {
			incomeByCategory[string(tx.Category)] += tx.Amount
		}
	}

	// Calculate totals by category
	report.EmploymentIncome = incomeByCategory[string(model.CategoryEmployment)]
	report.FreelanceIncome = incomeByCategory[string(model.CategoryFreelance)]
	report.RentalIncome = incomeByCategory[string(model.CategoryRental)]
	report.InvestmentIncome = incomeByCategory[string(model.CategoryInvestment)]
	report.CryptoIncome = incomeByCategory[string(model.CategoryCrypto)]
	report.OtherIncome = incomeByCategory[string(model.CategoryOtherIncome)] + incomeByCategory[string(model.CategoryInterest)]

	report.TotalIncome = report.EmploymentIncome + report.FreelanceIncome +
		report.RentalIncome + report.InvestmentIncome +
		report.CryptoIncome + report.OtherIncome

	// Calculate reliefs
	totalReliefs, reliefsApplied := e.reliefCalculator.CalculateReliefs(req.Reliefs)
	report.RentRelief = reliefsApplied["rent_relief"]
	report.PensionDeduction = reliefsApplied["pension"]
	report.NHISDeduction = reliefsApplied["nhis"]
	report.NHFDeduction = reliefsApplied["nhf"]
	report.TotalReliefs = totalReliefs

	// Calculate taxable income (income - reliefs, but not below 0)
	report.TaxableIncome = report.TotalIncome - report.TotalReliefs
	if report.TaxableIncome < 0 {
		report.TaxableIncome = 0
	}

	// Calculate PIT
	pitAmount, pitBreakdown := e.pitCalculator.Calculate(report.TaxableIncome)
	report.PITAmount = pitAmount

	// For now, CGT is handled separately (crypto gains would need cost basis tracking)
	// Simplified: we don't calculate CGT here as it requires more complex tracking
	report.CGTAmount = 0

	// Total tax
	report.TotalTax = report.PITAmount + report.CGTAmount

	// Build breakdown
	report.Breakdown = &model.TaxBreakdown{
		PITBreakdown:     pitBreakdown,
		IncomeByCategory: incomeByCategory,
		ReliefsApplied:   reliefsApplied,
	}

	return report, nil
}

// QuickCalculatePIT is a convenience method for quick PIT calculation
func (e *Engine) QuickCalculatePIT(annualIncome float64) float64 {
	return e.pitCalculator.CalculateSimple(annualIncome)
}

// CalculateRentRelief is a convenience method for rent relief
func (e *Engine) CalculateRentRelief(annualRent float64) float64 {
	return e.reliefCalculator.CalculateRentRelief(annualRent)
}
