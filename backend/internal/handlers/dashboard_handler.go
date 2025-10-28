package handlers

import (
	"finara-backend/internal/services"
	"finara-backend/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService services.DashboardService
}

func NewDashboardHandler(dashboardService services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
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

	summary, err := h.dashboardService.GetDashboardSummary(companyID.(uint), asOfDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve dashboard summary", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Dashboard summary retrieved successfully", summary)
}

func (h *DashboardHandler) GetMonthlyRevenue(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	yearStr := c.Query("year")
	if yearStr == "" {
		yearStr = strconv.Itoa(time.Now().Year())
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid year format", err)
		return
	}

	revenues, err := h.dashboardService.GetMonthlyRevenue(companyID.(uint), year)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve monthly revenue", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Monthly revenue retrieved successfully", revenues)
}

func (h *DashboardHandler) GetMonthlyExpense(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	yearStr := c.Query("year")
	if yearStr == "" {
		yearStr = strconv.Itoa(time.Now().Year())
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid year format", err)
		return
	}

	expenses, err := h.dashboardService.GetMonthlyExpense(companyID.(uint), year)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve monthly expense", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Monthly expense retrieved successfully", expenses)
}

func (h *DashboardHandler) GetExpenseByCategory(c *gin.Context) {
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

	expenses, err := h.dashboardService.GetExpenseByCategory(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve expense by category", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Expense by category retrieved successfully", expenses)
}

func (h *DashboardHandler) GetRevenueByCategory(c *gin.Context) {
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

	revenues, err := h.dashboardService.GetRevenueByCategory(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve revenue by category", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Revenue by category retrieved successfully", revenues)
}

func (h *DashboardHandler) GetFinancialRatios(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	asOfDateStr := c.Query("as_of_date")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if asOfDateStr == "" {
		asOfDateStr = time.Now().Format("2006-01-02")
	}
	if startDateStr == "" || endDateStr == "" {
		// Default to current year
		now := time.Now()
		startDateStr = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		endDateStr = asOfDateStr
	}

	asOfDate, err := time.Parse("2006-01-02", asOfDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid as_of_date format", err)
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

	ratios, err := h.dashboardService.GetFinancialRatios(companyID.(uint), asOfDate, startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve financial ratios", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Financial ratios retrieved successfully", ratios)
}