package handlers

import (
	"finara-backend/internal/models"
	"finara-backend/internal/services"
	"finara-backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService services.AccountService
}

func NewAccountHandler(accountService services.AccountService) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

type CreateAccountRequest struct {
	CompanyID   *uint                   `json:"company_id"` // Optional, can be provided if not in token
	Code        string                  `json:"code" binding:"required"`
	Name        string                  `json:"name" binding:"required"`
	Type        models.AccountType      `json:"type" binding:"required"`
	Category    models.AccountCategory  `json:"category" binding:"required"`
	ParentID    *uint                   `json:"parent_id"`
	Level       int                     `json:"level" binding:"required,min=1"`
	IsHeader    bool                    `json:"is_header"`
	Description string                  `json:"description"`
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Get company_id from token
	companyIDFromToken, _ := c.Get("company_id")

	// Use company_id from request if provided and token has no company_id (0)
	var companyID uint
	if req.CompanyID != nil && *req.CompanyID > 0 {
		companyID = *req.CompanyID
	} else {
		companyID = companyIDFromToken.(uint)
	}

	// Validate that we have a valid company_id
	if companyID == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Company ID is required", nil)
		return
	}

	account := &models.Account{
		CompanyID:   companyID,
		Code:        req.Code,
		Name:        req.Name,
		Type:        req.Type,
		Category:    req.Category,
		ParentID:    req.ParentID,
		Level:       req.Level,
		IsHeader:    req.IsHeader,
		Description: req.Description,
		IsActive:    true,
		Balance:     0,
	}

	if err := h.accountService.CreateAccount(account); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create account", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Account created successfully", account)
}

func (h *AccountHandler) GetAllAccounts(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	accounts, err := h.accountService.GetAccountsByCompanyID(companyID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve accounts", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Accounts retrieved successfully", accounts)
}

func (h *AccountHandler) GetActiveAccounts(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	accounts, err := h.accountService.GetActiveAccounts(companyID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve active accounts", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Active accounts retrieved successfully", accounts)
}

func (h *AccountHandler) GetAccountByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID", err)
		return
	}

	account, err := h.accountService.GetAccountByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Account not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Account retrieved successfully", account)
}

func (h *AccountHandler) GetAccountByCode(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	code := c.Param("code")

	if code == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Account code is required", nil)
		return
	}

	account, err := h.accountService.GetAccountByCode(companyID.(uint), code)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Account not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Account retrieved successfully", account)
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID", err)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	account, err := h.accountService.UpdateAccount(uint(id), updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update account", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Account updated successfully", account)
}

func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid account ID", err)
		return
	}

	if err := h.accountService.DeleteAccount(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete account", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Account deleted successfully", nil)
}

func (h *AccountHandler) InitializeDefaultAccounts(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	if err := h.accountService.InitializeDefaultAccounts(companyID.(uint)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to initialize default accounts", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Default accounts initialized successfully", nil)
}