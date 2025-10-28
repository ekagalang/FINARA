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

type InventoryHandler struct {
	inventoryService services.InventoryService
}

func NewInventoryHandler(inventoryService services.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

// Product Handlers
type CreateProductRequest struct {
	Code        string              `json:"code" binding:"required"`
	Name        string              `json:"name" binding:"required"`
	Description string              `json:"description"`
	Category    string              `json:"category"`
	Unit        string              `json:"unit" binding:"required"`
	CostMethod  models.CostMethod   `json:"cost_method"`
	MinStock    float64             `json:"min_stock"`
}

func (h *InventoryHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	companyID, _ := c.Get("company_id")

	if req.CostMethod == "" {
		req.CostMethod = models.CostMethodFIFO
	}

	product := &models.Product{
		CompanyID:   companyID.(uint),
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Unit:        req.Unit,
		CostMethod:  req.CostMethod,
		MinStock:    req.MinStock,
		IsActive:    true,
	}

	if err := h.inventoryService.CreateProduct(product); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create product", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Product created successfully", product)
}

func (h *InventoryHandler) GetAllProducts(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	products, err := h.inventoryService.GetProductsByCompanyID(companyID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve products", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Products retrieved successfully", products)
}

func (h *InventoryHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	product, err := h.inventoryService.GetProductByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Product not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product retrieved successfully", product)
}

func (h *InventoryHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	companyID, _ := c.Get("company_id")

	product := &models.Product{
		CompanyID:   companyID.(uint),
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Unit:        req.Unit,
		CostMethod:  req.CostMethod,
		MinStock:    req.MinStock,
	}

	if err := h.inventoryService.UpdateProduct(uint(id), product); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to update product", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product updated successfully", product)
}

func (h *InventoryHandler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	if err := h.inventoryService.DeleteProduct(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete product", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Product deleted successfully", nil)
}

// Stock Movement Handlers
type CreateStockMovementRequest struct {
	ProductID    uint    `json:"product_id" binding:"required"`
	MovementDate string  `json:"movement_date" binding:"required"`
	Type         string  `json:"type" binding:"required"` // in, out, adjustment
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
	UnitCost     float64 `json:"unit_cost" binding:"required,gte=0"`
	Reference    string  `json:"reference"`
	Notes        string  `json:"notes"`
}

func (h *InventoryHandler) CreateStockMovement(c *gin.Context) {
	var req CreateStockMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	companyID, _ := c.Get("company_id")
	userID, _ := c.Get("user_id")

	movementDate, err := time.Parse("2006-01-02", req.MovementDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid date format", err)
		return
	}

	movement := &models.StockMovement{
		CompanyID:    companyID.(uint),
		ProductID:    req.ProductID,
		MovementDate: movementDate,
		Type:         req.Type,
		Quantity:     req.Quantity,
		UnitCost:     req.UnitCost,
		Reference:    req.Reference,
		Notes:        req.Notes,
		CreatedBy:    userID.(uint),
	}

	var createErr error
	if req.Type == "in" {
		createErr = h.inventoryService.CreateStockIn(movement)
	} else if req.Type == "out" {
		createErr = h.inventoryService.CreateStockOut(movement)
	} else {
		createErr = h.inventoryService.CreateStockMovement(movement)
	}

	if createErr != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create stock movement", createErr)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Stock movement created successfully", movement)
}

func (h *InventoryHandler) GetStockMovements(c *gin.Context) {
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

	movements, err := h.inventoryService.GetStockMovementsByCompany(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve stock movements", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stock movements retrieved successfully", movements)
}

func (h *InventoryHandler) GetStockBalances(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	balances, err := h.inventoryService.GetAllStockBalances(companyID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve stock balances", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stock balances retrieved successfully", balances)
}

func (h *InventoryHandler) GetStockBalance(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid product ID", err)
		return
	}

	balance, err := h.inventoryService.GetStockBalance(uint(productID))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve stock balance", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stock balance retrieved successfully", balance)
}

// Stock Opname Handlers
type CreateStockOpnameRequest struct {
	OpnameDate string                     `json:"opname_date" binding:"required"`
	Notes      string                     `json:"notes"`
	Items      []StockOpnameItemRequest   `json:"items" binding:"required,min=1"`
}

type StockOpnameItemRequest struct {
	ProductID        uint    `json:"product_id" binding:"required"`
	SystemQuantity   float64 `json:"system_quantity" binding:"required"`
	PhysicalQuantity float64 `json:"physical_quantity" binding:"required"`
	Notes            string  `json:"notes"`
}

func (h *InventoryHandler) CreateStockOpname(c *gin.Context) {
	var req CreateStockOpnameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	companyID, _ := c.Get("company_id")
	userID, _ := c.Get("user_id")

	opnameDate, err := time.Parse("2006-01-02", req.OpnameDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid date format", err)
		return
	}

	// Convert items
	items := make([]models.StockOpnameItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = models.StockOpnameItem{
			ProductID:        item.ProductID,
			SystemQuantity:   item.SystemQuantity,
			PhysicalQuantity: item.PhysicalQuantity,
			Notes:            item.Notes,
		}
	}

	opname := &models.StockOpname{
		CompanyID:  companyID.(uint),
		OpnameDate: opnameDate,
		Notes:      req.Notes,
		CreatedBy:  userID.(uint),
		Items:      items,
	}

	if err := h.inventoryService.CreateStockOpname(opname); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create stock opname", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Stock opname created successfully", opname)
}

func (h *InventoryHandler) GetStockOpnames(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	opnames, err := h.inventoryService.GetStockOpnamesByCompany(companyID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve stock opnames", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stock opnames retrieved successfully", opnames)
}

func (h *InventoryHandler) ApproveStockOpname(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid stock opname ID", err)
		return
	}

	userID, _ := c.Get("user_id")

	if err := h.inventoryService.ApproveStockOpname(uint(id), userID.(uint)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to approve stock opname", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stock opname approved successfully", nil)
}

// HPP Calculation
func (h *InventoryHandler) CalculateHPP(c *gin.Context) {
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

	hpp, err := h.inventoryService.CalculateHPP(companyID.(uint), startDate, endDate)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to calculate HPP", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "HPP calculated successfully", hpp)
}