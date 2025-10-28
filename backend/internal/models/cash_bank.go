package models

import "time"

type TransactionType string
type TransactionCategory string

const (
	TransactionTypeIn  TransactionType = "in"
	TransactionTypeOut TransactionType = "out"
	TransactionTypeTransfer TransactionType = "transfer"

	CategoryCashSales     TransactionCategory = "cash_sales"
	CategoryCashPurchase  TransactionCategory = "cash_purchase"
	CategoryExpense       TransactionCategory = "expense"
	CategoryWithdrawal    TransactionCategory = "withdrawal"
	CategoryDeposit       TransactionCategory = "deposit"
	CategoryTransfer      TransactionCategory = "transfer"
	CategoryOther         TransactionCategory = "other"
)

type CashBankTransaction struct {
	BaseModel
	CompanyID        uint                `gorm:"not null;index" json:"company_id"`
	Company          Company             `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	AccountID        uint                `gorm:"not null;index" json:"account_id"`
	Account          Account             `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	TransactionNumber string             `gorm:"uniqueIndex;size:50;not null" json:"transaction_number"`
	TransactionDate  time.Time           `gorm:"not null;index" json:"transaction_date"`
	Type             TransactionType     `gorm:"type:varchar(20);not null" json:"type"`
	Category         TransactionCategory `gorm:"type:varchar(50);not null" json:"category"`
	Amount           float64             `gorm:"type:decimal(20,2);not null" json:"amount"`
	Description      string              `gorm:"type:text;not null" json:"description"`
	Reference        string              `gorm:"size:100" json:"reference"`
	JournalID        *uint               `gorm:"index" json:"journal_id"`
	Journal          *Journal            `gorm:"foreignKey:JournalID" json:"journal,omitempty"`
	CreatedBy        uint                `gorm:"not null" json:"created_by"`
	User             User                `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
}

type BankReconciliation struct {
	BaseModel
	CompanyID          uint      `gorm:"not null;index" json:"company_id"`
	Company            Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	BankAccountID      uint      `gorm:"not null;index" json:"bank_account_id"`
	BankAccount        Account   `gorm:"foreignKey:BankAccountID" json:"bank_account,omitempty"`
	ReconciliationDate time.Time `gorm:"not null" json:"reconciliation_date"`
	StatementBalance   float64   `gorm:"type:decimal(20,2);not null" json:"statement_balance"`
	BookBalance        float64   `gorm:"type:decimal(20,2);not null" json:"book_balance"`
	Difference         float64   `gorm:"type:decimal(20,2)" json:"difference"`
	Notes              string    `gorm:"type:text" json:"notes"`
	IsReconciled       bool      `gorm:"default:false" json:"is_reconciled"`
	ReconciledBy       *uint     `json:"reconciled_by"`
	ReconciledAt       *time.Time `json:"reconciled_at"`
}

type CashBankTransactionResponse struct {
	ID                uint                `json:"id"`
	CompanyID         uint                `json:"company_id"`
	AccountID         uint                `json:"account_id"`
	AccountCode       string              `json:"account_code"`
	AccountName       string              `json:"account_name"`
	TransactionNumber string              `json:"transaction_number"`
	TransactionDate   string              `json:"transaction_date"`
	Type              TransactionType     `json:"type"`
	Category          TransactionCategory `json:"category"`
	Amount            float64             `json:"amount"`
	Description       string              `json:"description"`
	Reference         string              `json:"reference"`
	JournalID         *uint               `json:"journal_id"`
	CreatedBy         uint                `json:"created_by"`
	CreatedByName     string              `json:"created_by_name"`
	CreatedAt         string              `json:"created_at"`
}