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

type JournalHandler struct {
	journalService services.JournalService
}

func NewJournalHandler(journalService services.JournalService) *JournalHandler {
	return &JournalHandler{journalService: journalService}
}

type CreateJournalRequest struct {
	TransactionDate string                     `json:"transaction_date" binding:"required"`
	Description     string                     `json:"description" binding:"required"`
	Entries         []CreateJournalEntryRequest `json:"entries" binding:"required,min=2"`
}

type CreateJournalEntryRequest struct {
	AccountID   uint    `json:"account_id" binding:"required"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Position    int     `json:"position" binding:"required"`
}

func (h *JournalHandler) CreateJournal(c *gin.Context) {
	var req CreateJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	companyID, _ := c.Get("company_id")
	userID, _ := c.Get("user_id")

	// Parse transaction date
	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid date format, use YYYY-MM-DD", err)
		return
	}

	// Create journal entries
	entries := make([]models.JournalEntry, len(req.Entries))
	for i, entry := range req.Entries {
		entries[i] = models.JournalEntry{
			AccountID:   entry.AccountID,
			Description: entry.Description,
			Debit:       entry.Debit,
			Credit:      entry.Credit,
			Position:    entry.Position,
		}
	}

	journal := &models.Journal{
		CompanyID:       companyID.(uint),
		TransactionDate: transactionDate,
		Description:     req.Description,
		CreatedBy:       userID.(uint),
		Status:          models.JournalStatusDraft,
		Entries:         entries,
	}

	if err := h.journalService.CreateJournal(journal); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create journal", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Journal created successfully", journal)
}

func (h *JournalHandler) GetJournalByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID", err)
		return
	}

	journal, err := h.journalService.GetJournalByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Journal not found", err)
		return
	}

	// Build response
	response := buildJournalResponse(journal)
	utils.SuccessResponse(c, http.StatusOK, "Journal retrieved successfully", response)
}

func (h *JournalHandler) GetJournalsByPeriod(c *gin.Context) {
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

	journals, err := h.journalService.GetJournalsByCompanyID(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve journals", err)
		return
	}

	// Build responses
	responses := make([]models.JournalResponse, len(journals))
	for i, journal := range journals {
		responses[i] = buildJournalResponse(&journal)
	}

	utils.SuccessResponse(c, http.StatusOK, "Journals retrieved successfully", responses)
}

func (h *JournalHandler) GetJournalsByStatus(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	statusStr := c.Param("status")

	status := models.JournalStatus(statusStr)
	if status != models.JournalStatusDraft && status != models.JournalStatusPosted && status != models.JournalStatusVoided {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid status. Use: draft, posted, or voided", nil)
		return
	}

	journals, err := h.journalService.GetJournalsByStatus(companyID.(uint), status)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve journals", err)
		return
	}

	// Build responses
	responses := make([]models.JournalResponse, len(journals))
	for i, journal := range journals {
		responses[i] = buildJournalResponse(&journal)
	}

	utils.SuccessResponse(c, http.StatusOK, "Journals retrieved successfully", responses)
}

func (h *JournalHandler) UpdateJournal(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID", err)
		return
	}

	var req CreateJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	companyID, _ := c.Get("company_id")

	transactionDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid date format", err)
		return
	}

	entries := make([]models.JournalEntry, len(req.Entries))
	for i, entry := range req.Entries {
		entries[i] = models.JournalEntry{
			AccountID:   entry.AccountID,
			Description: entry.Description,
			Debit:       entry.Debit,
			Credit:      entry.Credit,
			Position:    entry.Position,
		}
	}

	journal := &models.Journal{
		CompanyID:       companyID.(uint),
		TransactionDate: transactionDate,
		Description:     req.Description,
		Entries:         entries,
	}

	if err := h.journalService.UpdateJournal(uint(id), journal); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update journal", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journal updated successfully", journal)
}

func (h *JournalHandler) DeleteJournal(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID", err)
		return
	}

	if err := h.journalService.DeleteJournal(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete journal", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journal deleted successfully", nil)
}

func (h *JournalHandler) PostJournal(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID", err)
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.journalService.PostJournal(uint(id), userID.(uint)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to post journal", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journal posted successfully", nil)
}

func (h *JournalHandler) VoidJournal(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid journal ID", err)
		return
	}

	if err := h.journalService.VoidJournal(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to void journal", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Journal voided successfully", nil)
}

// Helper function to build journal response
func buildJournalResponse(journal *models.Journal) models.JournalResponse {
	entries := make([]models.JournalEntryResponse, len(journal.Entries))
	for i, entry := range journal.Entries {
		entries[i] = models.JournalEntryResponse{
			ID:          entry.ID,
			AccountID:   entry.AccountID,
			AccountCode: entry.Account.Code,
			AccountName: entry.Account.Name,
			Description: entry.Description,
			Debit:       entry.Debit,
			Credit:      entry.Credit,
			Position:    entry.Position,
		}
	}

	response := models.JournalResponse{
		ID:              journal.ID,
		CompanyID:       journal.CompanyID,
		JournalNumber:   journal.JournalNumber,
		TransactionDate: journal.TransactionDate.Format("2006-01-02"),
		Description:     journal.Description,
		Status:          journal.Status,
		TotalDebit:      journal.TotalDebit,
		TotalCredit:     journal.TotalCredit,
		CreatedBy:       journal.CreatedBy,
		CreatedByName:   journal.User.FullName,
		Entries:         entries,
		CreatedAt:       journal.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if journal.PostedAt != nil {
		postedAt := journal.PostedAt.Format("2006-01-02 15:04:05")
		response.PostedAt = &postedAt
	}

	return response
}