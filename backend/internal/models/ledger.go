package models

type Ledger struct {
	BaseModel
	CompanyID   uint         `gorm:"not null;index" json:"company_id"`
	Company     Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	AccountID   uint         `gorm:"not null;index" json:"account_id"`
	Account     Account      `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	JournalID   uint         `gorm:"not null;index" json:"journal_id"`
	Journal     Journal      `gorm:"foreignKey:JournalID" json:"journal,omitempty"`
	EntryID     uint         `gorm:"not null;index" json:"entry_id"`
	Entry       JournalEntry `gorm:"foreignKey:EntryID" json:"entry,omitempty"`
	Debit       float64      `gorm:"type:decimal(20,2);default:0" json:"debit"`
	Credit      float64      `gorm:"type:decimal(20,2);default:0" json:"credit"`
	Balance     float64      `gorm:"type:decimal(20,2);default:0" json:"balance"`
	Description string       `gorm:"type:text" json:"description"`
}

type LedgerResponse struct {
	ID              uint    `json:"id"`
	AccountID       uint    `json:"account_id"`
	AccountCode     string  `json:"account_code"`
	AccountName     string  `json:"account_name"`
	JournalNumber   string  `json:"journal_number"`
	TransactionDate string  `json:"transaction_date"`
	Description     string  `json:"description"`
	Debit           float64 `json:"debit"`
	Credit          float64 `json:"credit"`
	Balance         float64 `json:"balance"`
	CreatedAt       string  `json:"created_at"`
}

type TrialBalanceResponse struct {
	AccountCode   string  `json:"account_code"`
	AccountName   string  `json:"account_name"`
	AccountType   string  `json:"account_type"`
	Debit         float64 `json:"debit"`
	Credit        float64 `json:"credit"`
	DebitBalance  float64 `json:"debit_balance"`
	CreditBalance float64 `json:"credit_balance"`
}