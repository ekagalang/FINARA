package repository

import (
	"finara-backend/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type InventoryRepository interface {
	// Product
	CreateProduct(product *models.Product) error
	FindProductByID(id uint) (*models.Product, error)
	FindProductsByCompanyID(companyID uint) ([]models.Product, error)
	UpdateProduct(product *models.Product) error
	DeleteProduct(id uint) error

	// Stock Movement
	CreateStockMovement(movement *models.StockMovement) error
	FindStockMovementByID(id uint) (*models.StockMovement, error)
	FindStockMovementsByProduct(productID uint, startDate, endDate time.Time) ([]models.StockMovement, error)
	FindStockMovementsByCompany(companyID uint, startDate, endDate time.Time) ([]models.StockMovement, error)
	GenerateMovementNumber(companyID uint, movementType string, date time.Time) (string, error)

	// Stock Balance
	GetStockBalance(productID uint) (*models.StockBalance, error)
	UpdateStockBalance(balance *models.StockBalance) error
	GetAllStockBalances(companyID uint) ([]models.StockBalance, error)

	// Stock Opname
	CreateStockOpname(opname *models.StockOpname) error
	FindStockOpnameByID(id uint) (*models.StockOpname, error)
	FindStockOpnamesByCompany(companyID uint) ([]models.StockOpname, error)
	UpdateStockOpname(opname *models.StockOpname) error
	ApproveStockOpname(id uint, approvedBy uint) error
	GenerateOpnameNumber(companyID uint, date time.Time) (string, error)

	// HPP Calculation
	CalculateHPP(companyID uint, startDate, endDate time.Time) ([]models.HPPCalculation, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

// Product methods
func (r *inventoryRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *inventoryRepository) FindProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Company").First(&product, id).Error
	return &product, err
}

func (r *inventoryRepository) FindProductsByCompanyID(companyID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("company_id = ?", companyID).
		Order("code ASC").
		Find(&products).Error
	return products, err
}

func (r *inventoryRepository) UpdateProduct(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *inventoryRepository) DeleteProduct(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// Stock Movement methods
func (r *inventoryRepository) CreateStockMovement(movement *models.StockMovement) error {
	return r.db.Create(movement).Error
}

func (r *inventoryRepository) FindStockMovementByID(id uint) (*models.StockMovement, error) {
	var movement models.StockMovement
	err := r.db.Preload("Product").
		Preload("Journal").
		Preload("User").
		First(&movement, id).Error
	return &movement, err
}

func (r *inventoryRepository) FindStockMovementsByProduct(productID uint, startDate, endDate time.Time) ([]models.StockMovement, error) {
	var movements []models.StockMovement
	err := r.db.Where("product_id = ? AND movement_date BETWEEN ? AND ?", productID, startDate, endDate).
		Order("movement_date ASC, id ASC").
		Preload("Product").
		Find(&movements).Error
	return movements, err
}

func (r *inventoryRepository) FindStockMovementsByCompany(companyID uint, startDate, endDate time.Time) ([]models.StockMovement, error) {
	var movements []models.StockMovement
	err := r.db.Where("company_id = ? AND movement_date BETWEEN ? AND ?", companyID, startDate, endDate).
		Order("movement_date DESC, id DESC").
		Preload("Product").
		Find(&movements).Error
	return movements, err
}

func (r *inventoryRepository) GenerateMovementNumber(companyID uint, movementType string, date time.Time) (string, error) {
	var count int64
	var prefix string

	switch movementType {
	case "in":
		prefix = "IN/"
	case "out":
		prefix = "OUT/"
	case "adjustment":
		prefix = "ADJ/"
	default:
		prefix = "MOV/"
	}

	prefix += date.Format("200601/")

	err := r.db.Model(&models.StockMovement{}).
		Where("company_id = ? AND movement_number LIKE ?", companyID, prefix+"%").
		Count(&count).Error

	if err != nil {
		return "", err
	}

	return prefix + fmt.Sprintf("%04d", count+1), nil
}

// Stock Balance methods
func (r *inventoryRepository) GetStockBalance(productID uint) (*models.StockBalance, error) {
	var balance models.StockBalance
	err := r.db.Where("product_id = ?", productID).
		Preload("Product").
		First(&balance).Error
	
	if err == gorm.ErrRecordNotFound {
		return &models.StockBalance{
			ProductID:   productID,
			Quantity:    0,
			AverageCost: 0,
			TotalValue:  0,
		}, nil
	}
	
	return &balance, err
}

func (r *inventoryRepository) UpdateStockBalance(balance *models.StockBalance) error {
	return r.db.Save(balance).Error
}

func (r *inventoryRepository) GetAllStockBalances(companyID uint) ([]models.StockBalance, error) {
	var balances []models.StockBalance
	err := r.db.Where("company_id = ?", companyID).
		Preload("Product").
		Order("product_id ASC").
		Find(&balances).Error
	return balances, err
}

// Stock Opname methods
func (r *inventoryRepository) CreateStockOpname(opname *models.StockOpname) error {
	return r.db.Create(opname).Error
}

func (r *inventoryRepository) FindStockOpnameByID(id uint) (*models.StockOpname, error) {
	var opname models.StockOpname
	err := r.db.Preload("Items.Product").
		Preload("User").
		First(&opname, id).Error
	return &opname, err
}

func (r *inventoryRepository) FindStockOpnamesByCompany(companyID uint) ([]models.StockOpname, error) {
	var opnames []models.StockOpname
	err := r.db.Where("company_id = ?", companyID).
		Order("opname_date DESC").
		Preload("User").
		Find(&opnames).Error
	return opnames, err
}

func (r *inventoryRepository) UpdateStockOpname(opname *models.StockOpname) error {
	return r.db.Save(opname).Error
}

func (r *inventoryRepository) ApproveStockOpname(id uint, approvedBy uint) error {
	now := time.Now()
	return r.db.Model(&models.StockOpname{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      "approved",
			"approved_by": approvedBy,
			"approved_at": now,
		}).Error
}

func (r *inventoryRepository) GenerateOpnameNumber(companyID uint, date time.Time) (string, error) {
	var count int64
	prefix := "OPN/" + date.Format("200601/")

	err := r.db.Model(&models.StockOpname{}).
		Where("company_id = ? AND opname_number LIKE ?", companyID, prefix+"%").
		Count(&count).Error

	if err != nil {
		return "", err
	}

	return prefix + fmt.Sprintf("%04d", count+1), nil
}

// HPP Calculation
func (r *inventoryRepository) CalculateHPP(companyID uint, startDate, endDate time.Time) ([]models.HPPCalculation, error) {
	var calculations []models.HPPCalculation

	// This is a simplified calculation - you may need to adjust based on cost method
	err := r.db.Raw(`
		SELECT 
			p.id as product_id,
			p.code as product_code,
			p.name as product_name,
			COALESCE(beginning.qty, 0) as beginning_stock,
			COALESCE(beginning.value, 0) as beginning_value,
			COALESCE(purchases.qty, 0) as purchases,
			COALESCE(purchases.value, 0) as purchase_value,
			COALESCE(sales.qty, 0) as sales,
			COALESCE(sales.value, 0) as cogs,
			COALESCE(ending.qty, 0) as ending_stock,
			COALESCE(ending.value, 0) as ending_value
		FROM products p
		LEFT JOIN (
			SELECT product_id, 
				SUM(CASE WHEN type = 'in' THEN quantity ELSE -quantity END) as qty,
				SUM(CASE WHEN type = 'in' THEN total_cost ELSE -total_cost END) as value
			FROM stock_movements
			WHERE company_id = ? AND movement_date < ?
			GROUP BY product_id
		) beginning ON p.id = beginning.product_id
		LEFT JOIN (
			SELECT product_id,
				SUM(quantity) as qty,
				SUM(total_cost) as value
			FROM stock_movements
			WHERE company_id = ? AND type = 'in' AND movement_date BETWEEN ? AND ?
			GROUP BY product_id
		) purchases ON p.id = purchases.product_id
		LEFT JOIN (
			SELECT product_id,
				SUM(quantity) as qty,
				SUM(total_cost) as value
			FROM stock_movements
			WHERE company_id = ? AND type = 'out' AND movement_date BETWEEN ? AND ?
			GROUP BY product_id
		) sales ON p.id = sales.product_id
		LEFT JOIN (
			SELECT product_id,
				SUM(CASE WHEN type = 'in' THEN quantity ELSE -quantity END) as qty,
				SUM(CASE WHEN type = 'in' THEN total_cost ELSE -total_cost END) as value
			FROM stock_movements
			WHERE company_id = ? AND movement_date <= ?
			GROUP BY product_id
		) ending ON p.id = ending.product_id
		WHERE p.company_id = ?
		ORDER BY p.code ASC
	`, companyID, startDate, companyID, startDate, endDate, companyID, startDate, endDate, companyID, endDate, companyID).Scan(&calculations).Error

	return calculations, err
}