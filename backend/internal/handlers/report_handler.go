package handlers

import (
	"finara-backend/internal/services"
	"finara-backend/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService services.ReportService
}

func NewReportHandler(reportService services.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) GetIncomeStatement(c *gin.Context) {
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

	report, err := h.reportService.GetIncomeStatement(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate income statement", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Income statement generated successfully", report)
}

func (h *ReportHandler) GetBalanceSheet(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	asOfDateStr := c.Query("as_of_date")
	if asOfDateStr == "" {
		asOfDateStr = time.Now().Format("2006-01-02")
	}

	asOfDate, err := time.Parse("2006-01-02", asOfDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid as_of_date format", err)
		return
	}

	report, err := h.reportService.GetBalanceSheet(companyID.(uint), asOfDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate balance sheet", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Balance sheet generated successfully", report)
}

func (h *ReportHandler) GetCashFlow(c *gin.Context) {
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

	report, err := h.reportService.GetCashFlow(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate cash flow statement", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Cash flow statement generated successfully", report)
}