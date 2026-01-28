package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/taxsmart/taxsmart-api/internal/middleware"
	"github.com/taxsmart/taxsmart-api/internal/model"
	"github.com/taxsmart/taxsmart-api/internal/service/classifier"
	"github.com/taxsmart/taxsmart-api/internal/service/parser"
	"github.com/taxsmart/taxsmart-api/internal/service/tax"
	"github.com/taxsmart/taxsmart-api/pkg/response"
)

// Handler holds all dependencies for HTTP handlers
type Handler struct {
	csvParser  *parser.CSVParser
	classifier *classifier.Classifier
	taxEngine  *tax.Engine
}

// NewHandler creates a new handler with all dependencies
func NewHandler(aiProvider, aiAPIKey string) *Handler {
	return &Handler{
		csvParser:  parser.NewCSVParser(),
		classifier: classifier.NewClassifier(aiProvider, aiAPIKey),
		taxEngine:  tax.NewEngine(),
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.Success(w, map[string]string{
		"status":  "healthy",
		"service": "taxsmart-api",
	})
}

// ParseFile handles file upload and parsing
func (h *Handler) ParseFile(w http.ResponseWriter, r *http.Request) {
	// Limit file size to 10MB
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("file")
	if err != nil {
		response.BadRequest(w, "Failed to read uploaded file")
		return
	}
	defer file.Close()

	// Determine file type
	filename := header.Filename
	var transactions []model.ParsedTransaction
	var bankFormat string

	if strings.HasSuffix(strings.ToLower(filename), ".csv") {
		txs, format, err := h.csvParser.Parse(file)
		if err != nil {
			response.BadRequest(w, "Failed to parse CSV: "+err.Error())
			return
		}
		transactions = txs
		bankFormat = string(format)
	} else if strings.HasSuffix(strings.ToLower(filename), ".pdf") {
		// PDF parsing would go here
		response.BadRequest(w, "PDF parsing not yet implemented")
		return
	} else {
		response.BadRequest(w, "Unsupported file format. Please upload CSV or PDF")
		return
	}

	response.Success(w, map[string]interface{}{
		"transactions": transactions,
		"count":        len(transactions),
		"bank_format":  bankFormat,
		"filename":     filename,
	})
}

// ClassifyTransactions handles transaction classification
func (h *Handler) ClassifyTransactions(w http.ResponseWriter, r *http.Request) {
	var transactions []model.ParsedTransaction

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.BadRequest(w, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &transactions); err != nil {
		response.BadRequest(w, "Invalid JSON format")
		return
	}

	// Classify transactions
	results := h.classifier.ClassifyBatch(r.Context(), transactions)

	// Build response with transactions and their classifications
	classified := make([]map[string]interface{}, len(transactions))
	for i, tx := range transactions {
		classified[i] = map[string]interface{}{
			"date":        tx.Date,
			"description": tx.Description,
			"amount":      tx.Amount,
			"type":        tx.Type,
			"category":    results[i].Category,
			"confidence":  results[i].Confidence,
			"method":      results[i].Method,
		}
	}

	response.Success(w, map[string]interface{}{
		"transactions": classified,
		"count":        len(classified),
	})
}

// CalculateTax handles tax calculation requests
func (h *Handler) CalculateTax(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "User not authenticated")
		return
	}

	var req struct {
		TaxYear      int                 `json:"tax_year"`
		Transactions []model.Transaction `json:"transactions"`
		Reliefs      model.ReliefInput   `json:"reliefs"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.BadRequest(w, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		response.BadRequest(w, "Invalid JSON format")
		return
	}

	// Build calculation request
	calcReq := model.TaxCalculationRequest{
		TaxYear:      req.TaxYear,
		Transactions: req.Transactions,
		Reliefs:      req.Reliefs,
	}

	// Parse user ID
	_ = userID // Will be used when saving to database

	// Calculate tax
	report, err := h.taxEngine.CalculateTax(calcReq)
	if err != nil {
		response.InternalError(w, "Tax calculation failed: "+err.Error())
		return
	}

	response.Success(w, report)
}

// QuickCalculatePIT handles simple PIT calculation
func (h *Handler) QuickCalculatePIT(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AnnualIncome float64 `json:"annual_income"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response.BadRequest(w, "Failed to read request body")
		return
	}

	if err := json.Unmarshal(body, &req); err != nil {
		response.BadRequest(w, "Invalid JSON format")
		return
	}

	pitAmount := h.taxEngine.QuickCalculatePIT(req.AnnualIncome)

	response.Success(w, map[string]interface{}{
		"annual_income":  req.AnnualIncome,
		"pit_amount":     pitAmount,
		"effective_rate": pitAmount / req.AnnualIncome * 100,
	})
}
