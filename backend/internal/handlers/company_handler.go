package handlers

import (
	"finara-backend/internal/models"
	"finara-backend/internal/services"
	"finara-backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CompanyHandler struct {
	companyService services.CompanyService
	userService    services.UserService
}

func NewCompanyHandler(companyService services.CompanyService, userService services.UserService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
		userService:    userService,
	}
}

type CreateCompanyRequest struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	TaxID       string `json:"tax_id"`
	Description string `json:"description"`
}

func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	company := &models.Company{
		Name:        req.Name,
		Address:     req.Address,
		Phone:       req.Phone,
		Email:       req.Email,
		TaxID:       req.TaxID,
		Description: req.Description,
		IsActive:    true,
	}

	if err := h.companyService.CreateCompany(company); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create company", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Company created successfully", company)
}

func (h *CompanyHandler) GetAllCompanies(c *gin.Context) {
	companies, err := h.companyService.GetAllCompanies()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve companies", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Companies retrieved successfully", companies)
}

func (h *CompanyHandler) GetCompanyByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid company ID", err)
		return
	}

	company, err := h.companyService.GetCompanyByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Company not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Company retrieved successfully", company)
}

func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid company ID", err)
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	company, err := h.companyService.UpdateCompany(uint(id), updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update company", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Company updated successfully", company)
}

func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid company ID", err)
		return
	}

	if err := h.companyService.DeleteCompany(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete company", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Company deleted successfully", nil)
}

func (h *CompanyHandler) AssignUserToCompany(c *gin.Context) {
	var req struct {
		UserID    uint `json:"user_id" binding:"required"`
		CompanyID uint `json:"company_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	// Update user's company_id
	updates := map[string]interface{}{
		"company_id": req.CompanyID,
	}

	user, err := h.userService.UpdateUser(req.UserID, updates)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to assign user to company", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User assigned to company successfully", user)
}