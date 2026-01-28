package classifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/taxsmart/taxsmart-api/internal/model"
)

// AIClassifier classifies transactions using AI APIs
type AIClassifier struct {
	provider   string
	apiKey     string
	httpClient *http.Client
}

// NewAIClassifier creates a new AI classifier
func NewAIClassifier(provider, apiKey string) *AIClassifier {
	return &AIClassifier{
		provider: provider,
		apiKey:   apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsAvailable returns true if AI classification is configured
func (c *AIClassifier) IsAvailable() bool {
	return c.apiKey != ""
}

// Classify uses AI to classify a transaction
func (c *AIClassifier) Classify(ctx context.Context, description string, txType string, amount float64) (model.ClassificationResult, error) {
	if !c.IsAvailable() {
		return model.ClassificationResult{}, fmt.Errorf("AI classifier not configured")
	}

	prompt := fmt.Sprintf(`Classify this Nigerian bank transaction into one of these categories:
- employment_income: Salary, wages, payroll from employer
- freelance_income: Payments from freelance platforms or clients
- rental_income: Rent received from tenants
- investment_income: Dividends, investment returns
- crypto_income: Cryptocurrency trading/sales
- interest_income: Bank interest
- other_income: Other income sources
- expense: General expenses
- rent_expense: Rent payments to landlord
- transfer: Money transfers between accounts
- uncategorized: Cannot determine

Transaction: "%s"
Type: %s
Amount: %.2f NGN

Respond with ONLY a JSON object like: {"category": "category_name", "confidence": 0.85}`, description, txType, amount)

	var result model.ClassificationResult

	switch c.provider {
	case "gemini":
		result, _ = c.classifyWithGemini(ctx, prompt)
	case "openai":
		result, _ = c.classifyWithOpenAI(ctx, prompt)
	default:
		return model.ClassificationResult{}, fmt.Errorf("unsupported AI provider: %s", c.provider)
	}

	result.Method = "ai"
	return result, nil
}

func (c *AIClassifier) classifyWithGemini(ctx context.Context, prompt string) (model.ClassificationResult, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", c.apiKey)

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.1,
			"maxOutputTokens": 100,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return model.ClassificationResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return model.ClassificationResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return model.ClassificationResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.ClassificationResult{}, err
	}

	return c.parseAIResponse(body)
}

func (c *AIClassifier) classifyWithOpenAI(ctx context.Context, prompt string) (model.ClassificationResult, error) {
	url := "https://api.openai.com/v1/chat/completions"

	requestBody := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.1,
		"max_tokens":  100,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return model.ClassificationResult{}, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return model.ClassificationResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return model.ClassificationResult{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.ClassificationResult{}, err
	}

	return c.parseAIResponse(body)
}

func (c *AIClassifier) parseAIResponse(body []byte) (model.ClassificationResult, error) {
	// Try to extract JSON from response
	var result struct {
		Category   string  `json:"category"`
		Confidence float64 `json:"confidence"`
	}

	// Simple JSON extraction - find the JSON object in the response
	start := bytes.Index(body, []byte("{"))
	end := bytes.LastIndex(body, []byte("}"))
	if start != -1 && end != -1 && end > start {
		jsonData := body[start : end+1]
		if err := json.Unmarshal(jsonData, &result); err == nil {
			return model.ClassificationResult{
				Category:   model.Category(result.Category),
				Confidence: result.Confidence,
				Method:     "ai",
			}, nil
		}
	}

	return model.ClassificationResult{
		Category:   model.CategoryUncategorized,
		Confidence: 0,
		Method:     "ai",
	}, fmt.Errorf("failed to parse AI response")
}
