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

type TaxHandler struct {
	taxService services.TaxService
}

func NewTaxHandler(taxService services.TaxService) *TaxHandler {
	return &TaxHandler{taxService: taxService}
}

type CreateTaxRequest struct {
	TaxType       models.TaxType `json:"tax_type" binding:"required"`
	TaxPeriod     string         `json:"tax_period" binding:"required"` // Format: YYYY-MM
	TaxableAmount float64        `json:"taxable_amount" binding:"required,gt=0"`
	TaxRate       float64        `json:"tax_rate" binding:"required,gt=0"`
	DueDate       string         `json:"due_date" binding:"required"`
	Description   string         `json:"description"`
}

func (h *TaxHandler) CreateTax(c *gin.Context) {
	var req CreateTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	companyID, _ := c.Get("company_id")
	userID, _ := c.Get("user_id")

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid due_date format", err)
		return
	}

	tax := &models.Tax{
		CompanyID:     companyID.(uint),
		TaxType:       req.TaxType,
		TaxPeriod:     req.TaxPeriod,
		TaxableAmount: req.TaxableAmount,
		TaxRate:       req.TaxRate,
		DueDate:       dueDate,
		Description:   req.Description,
		CreatedBy:     userID.(uint),
	}

	if err := h.taxService.CreateTax(tax); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create tax", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Tax created successfully", tax)
}

func (h *TaxHandler) GetTaxByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid tax ID", err)
		return
	}

	tax, err := h.taxService.GetTaxByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Tax not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tax retrieved successfully", tax)
}

func (h *TaxHandler) GetTaxesByPeriod(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	period := c.Query("period")

	if period == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "period is required (format: YYYY-MM)", nil)
		return
	}

	taxes, err := h.taxService.GetTaxesByPeriod(companyID.(uint), period)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve taxes", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Taxes retrieved successfully", taxes)
}

func (h *TaxHandler) GetTaxesByType(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	taxType := models.TaxType(c.Param("type"))

	taxes, err := h.taxService.GetTaxesByType(companyID.(uint), taxType)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve taxes", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Taxes retrieved successfully", taxes)
}

func (h *TaxHandler) GetDueTaxes(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	
	dueDateStr := c.Query("due_date")
	if dueDateStr == "" {
		dueDateStr = time.Now().Format("2006-01-02")
	}

	dueDate, err := time.Parse("2006-01-02", dueDateStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid due_date format", err)
		return
	}

	taxes, err := h.taxService.GetDueTaxes(companyID.(uint), dueDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve due taxes", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Due taxes retrieved successfully", taxes)
}

func (h *TaxHandler) UpdateTax(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid tax ID", err)
		return
	}

	var req CreateTaxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	companyID, _ := c.Get("company_id")
	userID, _ := c.Get("user_id")

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid due_date format", err)
		return
	}

	tax := &models.Tax{
		CompanyID:     companyID.(uint),
		TaxType:       req.TaxType,
		TaxPeriod:     req.TaxPeriod,
		TaxableAmount: req.TaxableAmount,
		TaxRate:       req.TaxRate,
		DueDate:       dueDate,
		Description:   req.Description,
		CreatedBy:     userID.(uint),
	}

	if err := h.taxService.UpdateTax(uint(id), tax); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update tax", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tax updated successfully", tax)
}

func (h *TaxHandler) DeleteTax(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid tax ID", err)
		return
	}

	if err := h.taxService.DeleteTax(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete tax", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tax deleted successfully", nil)
}

func (h *TaxHandler) MarkAsReported(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid tax ID", err)
		return
	}

	if err := h.taxService.MarkAsReported(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to mark as reported", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tax marked as reported successfully", nil)
}

func (h *TaxHandler) MarkAsPaid(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid tax ID", err)
		return
	}

	if err := h.taxService.MarkAsPaid(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to mark as paid", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tax marked as paid successfully", nil)
}

func (h *TaxHandler) GetTaxSummary(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	period := c.Query("period")

	if period == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "period is required (format: YYYY-MM)", nil)
		return
	}

	summary, err := h.taxService.GetTaxSummary(companyID.(uint), period)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve tax summary", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Tax summary retrieved successfully", summary)
}