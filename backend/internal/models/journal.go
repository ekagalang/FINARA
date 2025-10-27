package models

import "time"

type JournalStatus string

const (
	JournalStatusDraft    JournalStatus = "draft"
	JournalStatusPosted   JournalStatus = "posted"
	JournalStatusVoided   JournalStatus = "voided"
)

type Journal struct {
	BaseModel
	CompanyID     uint          `gorm:"not null;index" json:"company_id"`
	Company       Company       `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	JournalNumber string        `gorm:"uniqueIndex;size:50;not null" json:"journal_number"`
	TransactionDate time.Time   `gorm:"not null;index" json:"transaction_date"`
	Description   string        `gorm:"type:text;not null" json:"description"`
	Status        JournalStatus `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	TotalDebit    float64       `gorm:"type:decimal(20,2);default:0" json:"total_debit"`
	TotalCredit   float64       `gorm:"type:decimal(20,2);default:0" json:"total_credit"`
	CreatedBy     uint          `gorm:"not null" json:"created_by"`
	User          User          `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
	PostedAt      *time.Time    `json:"posted_at"`
	PostedBy      *uint         `json:"posted_by"`
	Entries       []JournalEntry `gorm:"foreignKey:JournalID" json:"entries,omitempty"`
}

type JournalEntry struct {
	BaseModel
	JournalID   uint    `gorm:"not null;index" json:"journal_id"`
	Journal     Journal `gorm:"foreignKey:JournalID" json:"journal,omitempty"`
	AccountID   uint    `gorm:"not null;index" json:"account_id"`
	Account     Account `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	Description string  `gorm:"type:text" json:"description"`
	Debit       float64 `gorm:"type:decimal(20,2);default:0" json:"debit"`
	Credit      float64 `gorm:"type:decimal(20,2);default:0" json:"credit"`
	Position    int     `gorm:"not null" json:"position"` // urutan entry
}

type JournalResponse struct {
	ID              uint                  `json:"id"`
	CompanyID       uint                  `json:"company_id"`
	JournalNumber   string                `json:"journal_number"`
	TransactionDate string                `json:"transaction_date"`
	Description     string                `json:"description"`
	Status          JournalStatus         `json:"status"`
	TotalDebit      float64               `json:"total_debit"`
	TotalCredit     float64               `json:"total_credit"`
	CreatedBy       uint                  `json:"created_by"`
	CreatedByName   string                `json:"created_by_name"`
	PostedAt        *string               `json:"posted_at"`
	Entries         []JournalEntryResponse `json:"entries"`
	CreatedAt       string                `json:"created_at"`
}

type JournalEntryResponse struct {
	ID          uint    `json:"id"`
	AccountID   uint    `json:"account_id"`
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Description string  `json:"description"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Position    int     `json:"position"`
}