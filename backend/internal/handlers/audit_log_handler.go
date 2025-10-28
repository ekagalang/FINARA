package handlers

import (
	"finara-backend/internal/services"
	"finara-backend/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AuditLogHandler struct {
	auditLogService services.AuditLogService
}

func NewAuditLogHandler(auditLogService services.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{auditLogService: auditLogService}
}

func (h *AuditLogHandler) GetAuditLogs(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := h.auditLogService.GetAuditLogsByCompanyID(companyID.(uint), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve audit logs", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Audit logs retrieved successfully", logs)
}

func (h *AuditLogHandler) GetAuditLogsByUser(c *gin.Context) {
	userIDParam := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := h.auditLogService.GetAuditLogsByUserID(uint(userID), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve audit logs", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Audit logs retrieved successfully", logs)
}

func (h *AuditLogHandler) GetAuditLogsByModule(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	module := c.Param("module")

	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	logs, err := h.auditLogService.GetAuditLogsByModule(companyID.(uint), module, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve audit logs", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Audit logs retrieved successfully", logs)
}

func (h *AuditLogHandler) GetAuditLogsByDateRange(c *gin.Context) {
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

	logs, err := h.auditLogService.GetAuditLogsByDateRange(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve audit logs", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Audit logs retrieved successfully", logs)
}