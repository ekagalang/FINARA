package handlers

import (
	"finara-backend/internal/services"
	"finara-backend/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ExportHandler struct {
	exportService  services.ExportService
	journalService services.JournalService
	ledgerService  services.LedgerService
	reportService  services.ReportService
}

func NewExportHandler(
	exportService services.ExportService,
	journalService services.JournalService,
	ledgerService services.LedgerService,
	reportService services.ReportService,
) *ExportHandler {
	return &ExportHandler{
		exportService:  exportService,
		journalService: journalService,
		ledgerService:  ledgerService,
		reportService:  reportService,
	}
}

func (h *ExportHandler) ExportJournals(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	format := c.Query("format") // csv or excel

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "start_date and end_date are required", nil)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid start_date format", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end_date format", err)
		return
	}

	// Get journals
	journals, err := h.journalService.GetJournalsByCompanyID(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve journals", err)
		return
	}

	// Generate filename
	timestamp := time.Now().Format("20060102_150405")
	var filepath string
	var exportErr error

	if format == "excel" || format == "xlsx" {
		filename := fmt.Sprintf("journals_%s.xlsx", timestamp)
		filepath, exportErr = h.exportService.ExportJournalsToExcel(journals, filename)
	} else {
		filename := fmt.Sprintf("journals_%s.csv", timestamp)
		filepath, exportErr = h.exportService.ExportJournalsToCSV(journals, filename)
	}

	if exportErr != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to export journals", exportErr)
		return
	}

	// Return file
	c.FileAttachment(filepath, filepath)
}

func (h *ExportHandler) ExportTrialBalance(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	format := c.Query("format") // csv or excel

	endDateStr := c.Query("end_date")
	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01-02")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end_date format", err)
		return
	}

	// Get trial balance
	trialBalance, err := h.ledgerService.GetTrialBalance(companyID.(uint), endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve trial balance", err)
		return
	}

	// Generate filename
	timestamp := time.Now().Format("20060102_150405")
	var filepath string
	var exportErr error

	if format == "excel" || format == "xlsx" {
		filename := fmt.Sprintf("trial_balance_%s.xlsx", timestamp)
		filepath, exportErr = h.exportService.ExportTrialBalanceToExcel(trialBalance, filename)
	} else {
		filename := fmt.Sprintf("trial_balance_%s.csv", timestamp)
		filepath, exportErr = h.exportService.ExportTrialBalanceToCSV(trialBalance, filename)
	}

	if exportErr != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to export trial balance", exportErr)
		return
	}

	// Return file
	c.FileAttachment(filepath, filepath)
}

func (h *ExportHandler) ExportIncomeStatement(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	format := c.Query("format") // csv or excel

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "start_date and end_date are required", nil)
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid start_date format", err)
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end_date format", err)
		return
	}

	// Get income statement
	incomeStatement, err := h.reportService.GetIncomeStatement(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve income statement", err)
		return
	}

	// Generate filename
	timestamp := time.Now().Format("20060102_150405")
	var filepath string
	var exportErr error

	if format == "excel" || format == "xlsx" {
		filename := fmt.Sprintf("income_statement_%s.xlsx", timestamp)
		filepath, exportErr = h.exportService.ExportIncomeStatementToExcel(incomeStatement, filename)
	} else {
		filename := fmt.Sprintf("income_statement_%s.csv", timestamp)
		filepath, exportErr = h.exportService.ExportIncomeStatementToCSV(incomeStatement, filename)
	}

	if exportErr != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to export income statement", exportErr)
		return
	}

	// Return file
	c.FileAttachment(filepath, filepath)
}

func (h *ExportHandler) ExportBalanceSheet(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	format := c.Query("format") // csv or excel

	asOfDateStr := c.Query("as_of_date")
	if asOfDateStr == "" {
		asOfDateStr = time.Now().Format("2006-01-02")
	}

	asOfDate, err := time.Parse("2006-01-02", asOfDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid as_of_date format", err)
		return
	}

	// Get balance sheet
	balanceSheet, err := h.reportService.GetBalanceSheet(companyID.(uint), asOfDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve balance sheet", err)
		return
	}

	// Generate filename
	timestamp := time.Now().Format("20060102_150405")
	var filepath string
	var exportErr error

	if format == "excel" || format == "xlsx" {
		filename := fmt.Sprintf("balance_sheet_%s.xlsx", timestamp)
		filepath, exportErr = h.exportService.ExportBalanceSheetToExcel(balanceSheet, filename)
	} else {
		filename := fmt.Sprintf("balance_sheet_%s.csv", timestamp)
		filepath, exportErr = h.exportService.ExportBalanceSheetToCSV(balanceSheet, filename)
	}

	if exportErr != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to export balance sheet", exportErr)
		return
	}

	// Return file
	c.FileAttachment(filepath, filepath)
}