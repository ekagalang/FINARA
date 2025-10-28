package services

import (
	"errors"
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"time"
)

type InventoryService interface {
	// Product
	CreateProduct(product *models.Product) error
	GetProductByID(id uint) (*models.Product, error)
	GetProductsByCompanyID(companyID uint) ([]models.Product, error)
	UpdateProduct(id uint, product *models.Product) error
	DeleteProduct(id uint) error

	// Stock Movement
	CreateStockMovement(movement *models.StockMovement) error
	CreateStockIn(movement *models.StockMovement) error
	CreateStockOut(movement *models.StockMovement) error
	GetStockMovementByID(id uint) (*models.StockMovement, error)
	GetStockMovementsByProduct(productID uint, startDate, endDate time.Time) ([]models.StockMovement, error)
	GetStockMovementsByCompany(companyID uint, startDate, endDate time.Time) ([]models.StockMovement, error)

	// Stock Balance
	GetStockBalance(productID uint) (*models.StockBalance, error)
	GetAllStockBalances(companyID uint) ([]models.StockBalance, error)

	// Stock Opname
	CreateStockOpname(opname *models.StockOpname) error
	GetStockOpnameByID(id uint) (*models.StockOpname, error)
	GetStockOpnamesByCompany(companyID uint) ([]models.StockOpname, error)
	ApproveStockOpname(id uint, approvedBy uint) error

	// HPP
	CalculateHPP(companyID uint, startDate, endDate time.Time) ([]models.HPPCalculation, error)
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
	journalRepo   repository.JournalRepository
	accountRepo   repository.AccountRepository
}

func NewInventoryService(
	inventoryRepo repository.InventoryRepository,
	journalRepo repository.JournalRepository,
	accountRepo repository.AccountRepository,
) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
		journalRepo:   journalRepo,
		accountRepo:   accountRepo,
	}
}

// Product methods
func (s *inventoryService) CreateProduct(product *models.Product) error {
	return s.inventoryRepo.CreateProduct(product)
}

func (s *inventoryService) GetProductByID(id uint) (*models.Product, error) {
	return s.inventoryRepo.FindProductByID(id)
}

func (s *inventoryService) GetProductsByCompanyID(companyID uint) ([]models.Product, error) {
	return s.inventoryRepo.FindProductsByCompanyID(companyID)
}

func (s *inventoryService) UpdateProduct(id uint, updatedProduct *models.Product) error {
	product, err := s.inventoryRepo.FindProductByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	updatedProduct.ID = product.ID
	return s.inventoryRepo.UpdateProduct(updatedProduct)
}

func (s *inventoryService) DeleteProduct(id uint) error {
	// Check if product has stock balance
	balance, err := s.inventoryRepo.GetStockBalance(id)
	if err != nil && err.Error() != "record not found" {
		return err
	}

	if balance != nil && balance.Quantity > 0 {
		return errors.New("cannot delete product with existing stock")
	}

	return s.inventoryRepo.DeleteProduct(id)
}

// Stock Movement methods
func (s *inventoryService) CreateStockMovement(movement *models.StockMovement) error {
	// Generate movement number
	movementNumber, err := s.inventoryRepo.GenerateMovementNumber(
		movement.CompanyID,
		movement.Type,
		movement.MovementDate,
	)
	if err != nil {
		return err
	}
	movement.MovementNumber = movementNumber

	// Calculate total cost
	movement.TotalCost = movement.Quantity * movement.UnitCost

	// Create movement
	if err := s.inventoryRepo.CreateStockMovement(movement); err != nil {
		return err
	}

	// Update stock balance
	return s.updateStockBalance(movement)
}

func (s *inventoryService) CreateStockIn(movement *models.StockMovement) error {
	movement.Type = "in"
	return s.CreateStockMovement(movement)
}

func (s *inventoryService) CreateStockOut(movement *models.StockMovement) error {
	// Check if stock is sufficient
	balance, err := s.inventoryRepo.GetStockBalance(movement.ProductID)
	if err != nil {
		return err
	}

	if balance.Quantity < movement.Quantity {
		return errors.New("insufficient stock")
	}

	movement.Type = "out"
	// For stock out, use average cost from balance
	movement.UnitCost = balance.AverageCost
	movement.TotalCost = movement.Quantity * movement.UnitCost

	return s.CreateStockMovement(movement)
}

func (s *inventoryService) updateStockBalance(movement *models.StockMovement) error {
	balance, err := s.inventoryRepo.GetStockBalance(movement.ProductID)
	if err != nil {
		return err
	}

	// If balance doesn't exist, create new
	if balance.ID == 0 {
		balance.CompanyID = movement.CompanyID
		balance.ProductID = movement.ProductID
		balance.Quantity = 0
		balance.AverageCost = 0
		balance.TotalValue = 0
	}

	// Update balance based on movement type
	if movement.Type == "in" {
		// Weighted average cost method
		totalValue := balance.TotalValue + movement.TotalCost
		totalQty := balance.Quantity + movement.Quantity
		balance.Quantity = totalQty
		balance.TotalValue = totalValue
		if totalQty > 0 {
			balance.AverageCost = totalValue / totalQty
		}
	} else if movement.Type == "out" {
		balance.Quantity -= movement.Quantity
		balance.TotalValue -= movement.TotalCost
		if balance.Quantity < 0 {
			return errors.New("negative stock not allowed")
		}
	} else if movement.Type == "adjustment" {
		// For adjustment, recalculate everything
		balance.Quantity = movement.Quantity
		balance.TotalValue = movement.TotalCost
		if balance.Quantity > 0 {
			balance.AverageCost = balance.TotalValue / balance.Quantity
		}
	}

	return s.inventoryRepo.UpdateStockBalance(balance)
}

func (s *inventoryService) GetStockMovementByID(id uint) (*models.StockMovement, error) {
	return s.inventoryRepo.FindStockMovementByID(id)
}

func (s *inventoryService) GetStockMovementsByProduct(productID uint, startDate, endDate time.Time) ([]models.StockMovement, error) {
	return s.inventoryRepo.FindStockMovementsByProduct(productID, startDate, endDate)
}

func (s *inventoryService) GetStockMovementsByCompany(companyID uint, startDate, endDate time.Time) ([]models.StockMovement, error) {
	return s.inventoryRepo.FindStockMovementsByCompany(companyID, startDate, endDate)
}

// Stock Balance methods
func (s *inventoryService) GetStockBalance(productID uint) (*models.StockBalance, error) {
	return s.inventoryRepo.GetStockBalance(productID)
}

func (s *inventoryService) GetAllStockBalances(companyID uint) ([]models.StockBalance, error) {
	return s.inventoryRepo.GetAllStockBalances(companyID)
}

// Stock Opname methods
func (s *inventoryService) CreateStockOpname(opname *models.StockOpname) error {
	// Generate opname number
	opnameNumber, err := s.inventoryRepo.GenerateOpnameNumber(opname.CompanyID, opname.OpnameDate)
	if err != nil {
		return err
	}
	opname.OpnameNumber = opnameNumber

	// Calculate differences for each item
	for i := range opname.Items {
		opname.Items[i].Difference = opname.Items[i].PhysicalQuantity - opname.Items[i].SystemQuantity
	}

	return s.inventoryRepo.CreateStockOpname(opname)
}

func (s *inventoryService) GetStockOpnameByID(id uint) (*models.StockOpname, error) {
	return s.inventoryRepo.FindStockOpnameByID(id)
}

func (s *inventoryService) GetStockOpnamesByCompany(companyID uint) ([]models.StockOpname, error) {
	return s.inventoryRepo.FindStockOpnamesByCompany(companyID)
}

func (s *inventoryService) ApproveStockOpname(id uint, approvedBy uint) error {
	opname, err := s.inventoryRepo.FindStockOpnameByID(id)
	if err != nil {
		return errors.New("stock opname not found")
	}

	if opname.Status != "draft" {
		return errors.New("only draft stock opname can be approved")
	}

	// Approve opname
	if err := s.inventoryRepo.ApproveStockOpname(id, approvedBy); err != nil {
		return err
	}

	// Create adjustment movements for differences
	for _, item := range opname.Items {
		if item.Difference != 0 {
			movement := &models.StockMovement{
				CompanyID:      opname.CompanyID,
				ProductID:      item.ProductID,
				MovementDate:   opname.OpnameDate,
				Type:           "adjustment",
				Quantity:       item.PhysicalQuantity,
				Reference:      opname.OpnameNumber,
				Notes:          "Stock opname adjustment: " + item.Notes,
				CreatedBy:      approvedBy,
			}

			// Get current balance for unit cost
			balance, _ := s.inventoryRepo.GetStockBalance(item.ProductID)
			movement.UnitCost = balance.AverageCost
			movement.TotalCost = item.PhysicalQuantity * balance.AverageCost

			s.CreateStockMovement(movement)
		}
	}

	return nil
}

// HPP Calculation
func (s *inventoryService) CalculateHPP(companyID uint, startDate, endDate time.Time) ([]models.HPPCalculation, error) {
	return s.inventoryRepo.CalculateHPP(companyID, startDate, endDate)
}