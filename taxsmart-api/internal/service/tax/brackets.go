package tax

import (
	"github.com/taxsmart/taxsmart-api/internal/model"
)

// Nigeria 2026 Personal Income Tax Brackets
var PITBrackets2026 = []model.TaxBracket{
	{Min: 0, Max: 800_000, Rate: 0.00},
	{Min: 800_001, Max: 3_000_000, Rate: 0.15},
	{Min: 3_000_001, Max: 12_000_000, Rate: 0.18},
	{Min: 12_000_001, Max: 25_000_000, Rate: 0.21},
	{Min: 25_000_001, Max: 50_000_000, Rate: 0.23},
	{Min: 50_000_001, Max: 1e18, Rate: 0.25}, // Effectively infinite
}

// Constants for reliefs
const (
	// Rent relief is 20% of annual rent, capped at 500,000
	RentReliefPercentage = 0.20
	RentReliefCap        = 500_000.0

	// Tax-free threshold
	TaxFreeThreshold = 800_000.0

	// CGT exemption thresholds
	CGTExemptProceedsLimit = 150_000_000.0
	CGTExemptGainLimit     = 10_000_000.0
)
