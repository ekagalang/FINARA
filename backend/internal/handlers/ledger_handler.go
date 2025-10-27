package handlers

import (
	"finara-backend/internal/models"
	"finara-backend/internal/services"
	"finara-backend/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type LedgerHandler struct {
	ledgerService services.LedgerService
}

func NewLedgerHandler(ledgerService services.LedgerService) *LedgerHandler {
	return &LedgerHandler{ledgerService: ledgerService}
}

func (h *LedgerHandler) GetLedgerByAccount(c *gin.Context) {
	accountID, err := strconv.ParseUint(c.Param("account_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID", err)
		return
	}

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

	ledgers, err := h.ledgerService.GetLedgerByAccountID(uint(accountID), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve ledger", err)
		return
	}

	// Build response
	responses := make([]models.LedgerResponse, len(ledgers))
	for i, ledger := range ledgers {
		responses[i] = models.LedgerResponse{
			ID:              ledger.ID,
			AccountID:       ledger.AccountID,
			AccountCode:     ledger.Account.Code,
			AccountName:     ledger.Account.Name,
			JournalNumber:   ledger.Journal.JournalNumber,
			TransactionDate: ledger.Journal.TransactionDate.Format("2006-01-02"),
			Description:     ledger.Description,
			Debit:           ledger.Debit,
			Credit:          ledger.Credit,
			Balance:         ledger.Balance,
			CreatedAt:       ledger.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Ledger retrieved successfully", responses)
}

func (h *LedgerHandler) GetLedgerByCompany(c *gin.Context) {
	companyID, _ := c.Get("company_id")

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

	ledgers, err := h.ledgerService.GetLedgerByCompanyID(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve ledger", err)
		return
	}

	// Build response
	responses := make([]models.LedgerResponse, len(ledgers))
	for i, ledger := range ledgers {
		responses[i] = models.LedgerResponse{
			ID:              ledger.ID,
			AccountID:       ledger.AccountID,
			AccountCode:     ledger.Account.Code,
			AccountName:     ledger.Account.Name,
			JournalNumber:   ledger.Journal.JournalNumber,
			TransactionDate: ledger.Journal.TransactionDate.Format("2006-01-02"),
			Description:     ledger.Description,
			Debit:           ledger.Debit,
			Credit:          ledger.Credit,
			Balance:         ledger.Balance,
			CreatedAt:       ledger.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	utils.SuccessResponse(c, http.StatusOK, "Ledger retrieved successfully", responses)
}

func (h *LedgerHandler) GetTrialBalance(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	endDateStr := c.Query("end_date")
	if endDateStr == "" {
		// Default to today
		endDateStr = time.Now().Format("2006-01-02")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end_date format", err)
		return
	}

	trialBalance, err := h.ledgerService.GetTrialBalance(companyID.(uint), endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve trial balance", err)
		return
	}

	// Calculate totals
	var totalDebit, totalCredit, totalDebitBalance, totalCreditBalance float64
	for _, item := range trialBalance {
		totalDebit += item.Debit
		totalCredit += item.Credit
		totalDebitBalance += item.DebitBalance
		totalCreditBalance += item.CreditBalance
	}

	response := map[string]interface{}{
		"end_date":             endDateStr,
		"accounts":             trialBalance,
		"total_debit":          totalDebit,
		"total_credit":         totalCredit,
		"total_debit_balance":  totalDebitBalance,
		"total_credit_balance": totalCreditBalance,
		"is_balanced":          totalDebitBalance == totalCreditBalance,
	}

	utils.SuccessResponse(c, http.StatusOK, "Trial balance retrieved successfully", response)
}

func (h *LedgerHandler) GetAccountBalance(c *gin.Context) {
	accountID, err := strconv.ParseUint(c.Param("account_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID", err)
		return
	}

	endDateStr := c.Query("end_date")
	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01-02")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end_date format", err)
		return
	}

	balance, err := h.ledgerService.GetAccountBalance(uint(accountID), endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve account balance", err)
		return
	}

	response := map[string]interface{}{
		"account_id": accountID,
		"end_date":   endDateStr,
		"balance":    balance,
	}

	utils.SuccessResponse(c, http.StatusOK, "Account balance retrieved successfully", response)
}