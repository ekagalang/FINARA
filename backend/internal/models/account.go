package models

type AccountType string
type AccountCategory string

const (
	// Account Types
	AccountTypeAsset     AccountType = "asset"
	AccountTypeLiability AccountType = "liability"
	AccountTypeEquity    AccountType = "equity"
	AccountTypeRevenue   AccountType = "revenue"
	AccountTypeExpense   AccountType = "expense"

	// Account Categories
	CategoryCurrentAsset      AccountCategory = "current_asset"
	CategoryFixedAsset        AccountCategory = "fixed_asset"
	CategoryCurrentLiability  AccountCategory = "current_liability"
	CategoryLongTermLiability AccountCategory = "long_term_liability"
	CategoryEquity            AccountCategory = "equity"
	CategoryOperatingRevenue  AccountCategory = "operating_revenue"
	CategoryOtherRevenue      AccountCategory = "other_revenue"
	CategoryOperatingExpense  AccountCategory = "operating_expense"
	CategoryOtherExpense      AccountCategory = "other_expense"
)

type Account struct {
	BaseModel
	CompanyID   uint            `gorm:"not null;index" json:"company_id"`
	Company     Company         `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Code        string          `gorm:"uniqueIndex:idx_company_account_code;size:50;not null" json:"code"`
	Name        string          `gorm:"not null;size:255" json:"name"`
	Type        AccountType     `gorm:"type:varchar(20);not null" json:"type"`
	Category    AccountCategory `gorm:"type:varchar(50);not null" json:"category"`
	ParentID    *uint           `gorm:"index" json:"parent_id"`
	Parent      *Account        `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Level       int             `gorm:"not null;default:1" json:"level"`
	IsHeader    bool            `gorm:"default:false" json:"is_header"`
	IsActive    bool            `gorm:"default:true" json:"is_active"`
	Balance     float64         `gorm:"type:decimal(20,2);default:0" json:"balance"`
	Description string          `gorm:"type:text" json:"description"`
}

type AccountResponse struct {
	ID          uint            `json:"id"`
	CompanyID   uint            `json:"company_id"`
	Code        string          `json:"code"`
	Name        string          `json:"name"`
	Type        AccountType     `json:"type"`
	Category    AccountCategory `json:"category"`
	ParentID    *uint           `json:"parent_id"`
	Level       int             `json:"level"`
	IsHeader    bool            `json:"is_header"`
	IsActive    bool            `json:"is_active"`
	Balance     float64         `json:"balance"`
	Description string          `json:"description"`
	CreatedAt   string          `json:"created_at"`
}