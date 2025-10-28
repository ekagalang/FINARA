package models

import "time"

type CostMethod string

const (
	CostMethodFIFO    CostMethod = "fifo"
	CostMethodLIFO    CostMethod = "lifo"
	CostMethodAverage CostMethod = "average"
)

// Product/Item Master
type Product struct {
	BaseModel
	CompanyID   uint       `gorm:"not null;index" json:"company_id"`
	Company     Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Code        string     `gorm:"uniqueIndex:idx_company_product_code;size:50;not null" json:"code"`
	Name        string     `gorm:"size:255;not null" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Category    string     `gorm:"size:100" json:"category"`
	Unit        string     `gorm:"size:50" json:"unit"` // pcs, kg, liter, etc
	CostMethod  CostMethod `gorm:"type:varchar(20);not null;default:'fifo'" json:"cost_method"`
	MinStock    float64    `gorm:"type:decimal(20,2);default:0" json:"min_stock"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
}

// Stock Movement
type StockMovement struct {
	BaseModel
	CompanyID       uint      `gorm:"not null;index" json:"company_id"`
	Company         Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	ProductID       uint      `gorm:"not null;index" json:"product_id"`
	Product         Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	MovementNumber  string    `gorm:"uniqueIndex;size:50;not null" json:"movement_number"`
	MovementDate    time.Time `gorm:"not null;index" json:"movement_date"`
	Type            string    `gorm:"type:varchar(20);not null" json:"type"` // in, out, adjustment
	Quantity        float64   `gorm:"type:decimal(20,2);not null" json:"quantity"`
	UnitCost        float64   `gorm:"type:decimal(20,2);not null" json:"unit_cost"`
	TotalCost       float64   `gorm:"type:decimal(20,2);not null" json:"total_cost"`
	Reference       string    `gorm:"size:100" json:"reference"`
	Notes           string    `gorm:"type:text" json:"notes"`
	JournalID       *uint     `gorm:"index" json:"journal_id"`
	Journal         *Journal  `gorm:"foreignKey:JournalID" json:"journal,omitempty"`
	CreatedBy       uint      `gorm:"not null" json:"created_by"`
	User            User      `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
}

// Stock Balance
type StockBalance struct {
	BaseModel
	CompanyID     uint    `gorm:"not null;index" json:"company_id"`
	Company       Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	ProductID     uint    `gorm:"not null;index" json:"product_id"`
	Product       Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity      float64 `gorm:"type:decimal(20,2);not null;default:0" json:"quantity"`
	AverageCost   float64 `gorm:"type:decimal(20,2);not null;default:0" json:"average_cost"`
	TotalValue    float64 `gorm:"type:decimal(20,2);not null;default:0" json:"total_value"`
}

// Stock Opname (Physical Count)
type StockOpname struct {
	BaseModel
	CompanyID      uint      `gorm:"not null;index" json:"company_id"`
	Company        Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	OpnameNumber   string    `gorm:"uniqueIndex;size:50;not null" json:"opname_number"`
	OpnameDate     time.Time `gorm:"not null" json:"opname_date"`
	Status         string    `gorm:"type:varchar(20);not null;default:'draft'" json:"status"` // draft, approved
	Notes          string    `gorm:"type:text" json:"notes"`
	CreatedBy      uint      `gorm:"not null" json:"created_by"`
	User           User      `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
	ApprovedBy     *uint     `json:"approved_by"`
	ApprovedAt     *time.Time `json:"approved_at"`
	Items          []StockOpnameItem `gorm:"foreignKey:OpnameID" json:"items,omitempty"`
}

type StockOpnameItem struct {
	BaseModel
	OpnameID        uint    `gorm:"not null;index" json:"opname_id"`
	ProductID       uint    `gorm:"not null;index" json:"product_id"`
	Product         Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	SystemQuantity  float64 `gorm:"type:decimal(20,2);not null" json:"system_quantity"`
	PhysicalQuantity float64 `gorm:"type:decimal(20,2);not null" json:"physical_quantity"`
	Difference      float64 `gorm:"type:decimal(20,2);not null" json:"difference"`
	Notes           string  `gorm:"type:text" json:"notes"`
}

// HPP Calculation (Cost of Goods Sold)
type HPPCalculation struct {
	ProductID       uint    `json:"product_id"`
	ProductCode     string  `json:"product_code"`
	ProductName     string  `json:"product_name"`
	BeginningStock  float64 `json:"beginning_stock"`
	BeginningValue  float64 `json:"beginning_value"`
	Purchases       float64 `json:"purchases"`
	PurchaseValue   float64 `json:"purchase_value"`
	Sales           float64 `json:"sales"`
	COGS            float64 `json:"cogs"`
	EndingStock     float64 `json:"ending_stock"`
	EndingValue     float64 `json:"ending_value"`
}