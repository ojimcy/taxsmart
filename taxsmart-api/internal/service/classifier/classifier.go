package classifier

import (
	"context"

	"github.com/taxsmart/taxsmart-api/internal/model"
)

// Classifier combines AI and rule-based classification
type Classifier struct {
	ai    *AIClassifier
	rules *RuleEngine
}

// NewClassifier creates a new hybrid classifier
func NewClassifier(aiProvider, aiAPIKey string) *Classifier {
	var ai *AIClassifier
	if aiAPIKey != "" {
		ai = NewAIClassifier(aiProvider, aiAPIKey)
	}

	return &Classifier{
		ai:    ai,
		rules: NewRuleEngine(),
	}
}

// Classify classifies a transaction using AI with rule-based fallback
func (c *Classifier) Classify(ctx context.Context, description string, txType string, amount float64) model.ClassificationResult {
	// Try AI first if available
	if c.ai != nil && c.ai.IsAvailable() {
		result, err := c.ai.Classify(ctx, description, txType, amount)
		if err == nil && result.Confidence > 0.7 {
			return result
		}
	}

	// Fallback to rules
	return c.rules.Classify(description, txType)
}

// ClassifyBatch classifies multiple transactions
func (c *Classifier) ClassifyBatch(ctx context.Context, transactions []model.ParsedTransaction) []model.ClassificationResult {
	results := make([]model.ClassificationResult, len(transactions))

	for i, tx := range transactions {
		results[i] = c.Classify(ctx, tx.Description, tx.Type, tx.Amount)
	}

	return results
}
