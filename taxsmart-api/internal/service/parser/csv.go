package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/taxsmart/taxsmart-api/internal/model"
)

// CSVParser parses bank statement CSV files
type CSVParser struct{}

// NewCSVParser creates a new CSV parser
func NewCSVParser() *CSVParser {
	return &CSVParser{}
}

// BankFormat represents detected bank format
type BankFormat string

const (
	FormatGTBank    BankFormat = "gtbank"
	FormatAccess    BankFormat = "access"
	FormatFirstBank BankFormat = "firstbank"
	FormatUBA       BankFormat = "uba"
	FormatZenith    BankFormat = "zenith"
	FormatGeneric   BankFormat = "generic"
)

// Parse parses a CSV file and returns transactions
func (p *CSVParser) Parse(reader io.Reader) ([]model.ParsedTransaction, BankFormat, error) {
	csvReader := csv.NewReader(reader)
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true

	// Read all records
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, "", fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, "", fmt.Errorf("CSV file too short")
	}

	// Detect bank format from headers
	headers := records[0]
	format := p.detectFormat(headers)

	// Parse based on format
	var transactions []model.ParsedTransaction
	for i, record := range records[1:] {
		tx, err := p.parseRecord(record, headers, format)
		if err != nil {
			// Skip unparseable rows
			continue
		}
		if tx.Amount != 0 {
			transactions = append(transactions, tx)
		}
		_ = i // Avoid unused variable warning
	}

	return transactions, format, nil
}

func (p *CSVParser) detectFormat(headers []string) BankFormat {
	headerStr := strings.ToUpper(strings.Join(headers, " "))

	if strings.Contains(headerStr, "NARRATION") || strings.Contains(headerStr, "TRANS DATE") {
		return FormatGTBank
	}
	if strings.Contains(headerStr, "WITHDRAWALS") && strings.Contains(headerStr, "LODGEMENTS") {
		return FormatAccess
	}
	if strings.Contains(headerStr, "MONEY OUT") && strings.Contains(headerStr, "MONEY IN") {
		return FormatUBA
	}
	if strings.Contains(headerStr, "DEBIT AMOUNT") && strings.Contains(headerStr, "CREDIT AMOUNT") {
		return FormatFirstBank
	}

	return FormatGeneric
}

func (p *CSVParser) parseRecord(record, headers []string, format BankFormat) (model.ParsedTransaction, error) {
	// Create header index map
	headerIndex := make(map[string]int)
	for i, h := range headers {
		headerIndex[strings.ToUpper(strings.TrimSpace(h))] = i
	}

	var tx model.ParsedTransaction
	var err error

	switch format {
	case FormatGTBank:
		tx, err = p.parseGTBank(record, headerIndex)
	case FormatAccess:
		tx, err = p.parseAccess(record, headerIndex)
	case FormatUBA:
		tx, err = p.parseUBA(record, headerIndex)
	case FormatFirstBank:
		tx, err = p.parseFirstBank(record, headerIndex)
	default:
		tx, err = p.parseGeneric(record, headerIndex)
	}

	return tx, err
}

func (p *CSVParser) parseGTBank(record []string, idx map[string]int) (model.ParsedTransaction, error) {
	tx := model.ParsedTransaction{}

	// Date
	if i, ok := idx["TRANS DATE"]; ok && i < len(record) {
		tx.Date = p.parseDate(record[i])
	} else if i, ok := idx["DATE"]; ok && i < len(record) {
		tx.Date = p.parseDate(record[i])
	}

	// Description
	if i, ok := idx["NARRATION"]; ok && i < len(record) {
		tx.Description = record[i]
	} else if i, ok := idx["DESCRIPTION"]; ok && i < len(record) {
		tx.Description = record[i]
	}

	// Amount - check debit and credit
	if i, ok := idx["DEBIT"]; ok && i < len(record) {
		if amount := p.parseAmount(record[i]); amount > 0 {
			tx.Amount = amount
			tx.Type = "debit"
		}
	}
	if i, ok := idx["CREDIT"]; ok && i < len(record) {
		if amount := p.parseAmount(record[i]); amount > 0 {
			tx.Amount = amount
			tx.Type = "credit"
		}
	}

	// Balance
	if i, ok := idx["BALANCE"]; ok && i < len(record) {
		tx.Balance = p.parseAmount(record[i])
	}

	return tx, nil
}

func (p *CSVParser) parseAccess(record []string, idx map[string]int) (model.ParsedTransaction, error) {
	tx := model.ParsedTransaction{}

	if i, ok := idx["TRANSACTION DATE"]; ok && i < len(record) {
		tx.Date = p.parseDate(record[i])
	}
	if i, ok := idx["NARRATION"]; ok && i < len(record) {
		tx.Description = record[i]
	}
	if i, ok := idx["WITHDRAWALS"]; ok && i < len(record) {
		if amount := p.parseAmount(record[i]); amount > 0 {
			tx.Amount = amount
			tx.Type = "debit"
		}
	}
	if i, ok := idx["LODGEMENTS"]; ok && i < len(record) {
		if amount := p.parseAmount(record[i]); amount > 0 {
			tx.Amount = amount
			tx.Type = "credit"
		}
	}

	return tx, nil
}

func (p *CSVParser) parseUBA(record []string, idx map[string]int) (model.ParsedTransaction, error) {
	tx := model.ParsedTransaction{}

	if i, ok := idx["DATE"]; ok && i < len(record) {
		tx.Date = p.parseDate(record[i])
	}
	if i, ok := idx["DESCRIPTION"]; ok && i < len(record) {
		tx.Description = record[i]
	}
	if i, ok := idx["MONEY OUT"]; ok && i < len(record) {
		if amount := p.parseAmount(record[i]); amount > 0 {
			tx.Amount = amount
			tx.Type = "debit"
		}
	}
	if i, ok := idx["MONEY IN"]; ok && i < len(record) {
		if amount := p.parseAmount(record[i]); amount > 0 {
			tx.Amount = amount
			tx.Type = "credit"
		}
	}

	return tx, nil
}

func (p *CSVParser) parseFirstBank(record []string, idx map[string]int) (model.ParsedTransaction, error) {
	tx := model.ParsedTransaction{}

	if i, ok := idx["VALUE DATE"]; ok && i < len(record) {
		tx.Date = p.parseDate(record[i])
	}
	if i, ok := idx["REFERENCE"]; ok && i < len(record) {
		tx.Description = record[i]
	} else if i, ok := idx["DESCRIPTION"]; ok && i < len(record) {
		tx.Description = record[i]
	}
	if i, ok := idx["DEBIT AMOUNT"]; ok && i < len(record) {
		if amount := p.parseAmount(record[i]); amount > 0 {
			tx.Amount = amount
			tx.Type = "debit"
		}
	}
	if i, ok := idx["CREDIT AMOUNT"]; ok && i < len(record) {
		if amount := p.parseAmount(record[i]); amount > 0 {
			tx.Amount = amount
			tx.Type = "credit"
		}
	}

	return tx, nil
}

func (p *CSVParser) parseGeneric(record []string, idx map[string]int) (model.ParsedTransaction, error) {
	tx := model.ParsedTransaction{}

	// Try common date column names
	dateColumns := []string{"DATE", "TRANS DATE", "TRANSACTION DATE", "VALUE DATE", "POST DATE"}
	for _, col := range dateColumns {
		if i, ok := idx[col]; ok && i < len(record) {
			tx.Date = p.parseDate(record[i])
			if !tx.Date.IsZero() {
				break
			}
		}
	}

	// Try common description column names
	descColumns := []string{"DESCRIPTION", "NARRATION", "REMARKS", "REFERENCE", "DETAILS"}
	for _, col := range descColumns {
		if i, ok := idx[col]; ok && i < len(record) {
			tx.Description = record[i]
			if tx.Description != "" {
				break
			}
		}
	}

	// Try common amount column names
	debitColumns := []string{"DEBIT", "WITHDRAWALS", "MONEY OUT", "DEBIT AMOUNT", "DR"}
	for _, col := range debitColumns {
		if i, ok := idx[col]; ok && i < len(record) {
			if amount := p.parseAmount(record[i]); amount > 0 {
				tx.Amount = amount
				tx.Type = "debit"
				break
			}
		}
	}

	creditColumns := []string{"CREDIT", "LODGEMENTS", "MONEY IN", "CREDIT AMOUNT", "CR"}
	for _, col := range creditColumns {
		if i, ok := idx[col]; ok && i < len(record) {
			if amount := p.parseAmount(record[i]); amount > 0 {
				tx.Amount = amount
				tx.Type = "credit"
				break
			}
		}
	}

	// Handle single amount column
	if tx.Amount == 0 {
		if i, ok := idx["AMOUNT"]; ok && i < len(record) {
			tx.Amount = p.parseAmount(record[i])
			// Negative usually means debit
			if tx.Amount < 0 {
				tx.Amount = -tx.Amount
				tx.Type = "debit"
			} else {
				tx.Type = "credit"
			}
		}
	}

	return tx, nil
}

func (p *CSVParser) parseDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}

	// Common date formats
	formats := []string{
		"02-Jan-2006",
		"02/01/2006",
		"2006-01-02",
		"01/02/2006",
		"02-01-2006",
		"2/1/2006",
		"02 Jan 2006",
		"Jan 02, 2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t
		}
	}

	return time.Time{}
}

func (p *CSVParser) parseAmount(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" || s == "-" || s == "0" || s == "0.00" {
		return 0
	}

	// Remove currency symbols and commas
	re := regexp.MustCompile(`[₦NGN,\s]`)
	s = re.ReplaceAllString(s, "")

	// Handle parentheses as negative
	isNegative := false
	if strings.HasPrefix(s, "(") && strings.HasSuffix(s, ")") {
		s = strings.Trim(s, "()")
		isNegative = true
	}

	amount, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	if isNegative {
		amount = -amount
	}

	return amount
}
