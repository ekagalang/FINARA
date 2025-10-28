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

type CashBankHandler struct {
	cashBankService services.CashBankService
}

func NewCashBankHandler(cashBankService services.CashBankService) *CashBankHandler {
	return &CashBankHandler{cashBankService: cashBankService}
}

type CreateCashBankTransactionRequest struct {
	AccountID       uint                        `json:"account_id" binding:"required"`
	TransactionDate string                      `json:"transaction_date" binding:"required"`
	Type            models.TransactionType      `json:"type" binding:"required"`
	Category        models.TransactionCategory  `json:"category" binding:"required"`
	Amount          float64                     `json:"amount" binding:"required,gt=0"`
	Description     string                      `json:"description" binding:"required"`
	Reference       string                      `json:"reference"`
	ContraAccountID *uint                       `json:"contra_account_id"` // For auto journal creation
}

func (h *CashBankHandler) CreateTransaction(c *gin.Context) {
	var req CreateCashBankTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	companyID, _ := c.Get("company_id")
	userID, _ := c.Get("user_id")

	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid date format", err)
		return
	}

	transaction := &models.CashBankTransaction{
		CompanyID:       companyID.(uint),
		AccountID:       req.AccountID,
		TransactionDate: transactionDate,
		Type:            req.Type,
		Category:        req.Category,
		Amount:          req.Amount,
		Description:     req.Description,
		Reference:       req.Reference,
		CreatedBy:       userID.(uint),
	}

	// If contra account provided, create with journal
	if req.ContraAccountID != nil {
		if req.Type == models.TransactionTypeIn {
			err = h.cashBankService.CreateCashInWithJournal(transaction, *req.ContraAccountID)
		} else if req.Type == models.TransactionTypeOut {
			err = h.cashBankService.CreateCashOutWithJournal(transaction, *req.ContraAccountID)
		} else {
			err = h.cashBankService.CreateTransaction(transaction)
		}
	} else {
		err = h.cashBankService.CreateTransaction(transaction)
	}

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create transaction", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Transaction created successfully", transaction)
}

func (h *CashBankHandler) GetTransactionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err)
		return
	}

	transaction, err := h.cashBankService.GetTransactionByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Transaction not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction retrieved successfully", transaction)
}

func (h *CashBankHandler) GetTransactionsByPeriod(c *gin.Context) {
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

	transactions, err := h.cashBankService.GetTransactionsByCompanyID(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve transactions", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", transactions)
}

func (h *CashBankHandler) GetTransactionsByAccount(c *gin.Context) {
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

	transactions, err := h.cashBankService.GetTransactionsByAccountID(uint(accountID), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve transactions", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transactions retrieved successfully", transactions)
}

func (h *CashBankHandler) UpdateTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err)
		return
	}

	var req CreateCashBankTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	companyID, _ := c.Get("company_id")
	userID, _ := c.Get("user_id")

	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid date format", err)
		return
	}

	transaction := &models.CashBankTransaction{
		CompanyID:       companyID.(uint),
		AccountID:       req.AccountID,
		TransactionDate: transactionDate,
		Type:            req.Type,
		Category:        req.Category,
		Amount:          req.Amount,
		Description:     req.Description,
		Reference:       req.Reference,
		CreatedBy:       userID.(uint),
	}

	if err := h.cashBankService.UpdateTransaction(uint(id), transaction); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update transaction", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction updated successfully", transaction)
}

func (h *CashBankHandler) DeleteTransaction(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err)
		return
	}

	if err := h.cashBankService.DeleteTransaction(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete transaction", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction deleted successfully", nil)
}

func (h *CashBankHandler) GetCashPosition(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	endDateStr := c.Query("end_date")
	if endDateStr == "" {
		endDateStr = time.Now().Format("2006-01-02")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid end_date format", err)
		return
	}

	position, err := h.cashBankService.GetCashPosition(companyID.(uint), endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve cash position", err)
		return
	}

	response := map[string]interface{}{
		"as_of_date": endDateStr,
		"position":   position,
	}

	utils.SuccessResponse(c, http.StatusOK, "Cash position retrieved successfully", response)
}