package models

import "time"

type TaxType string
type TaxStatus string

const (
	TaxTypePPNIn      TaxType = "ppn_in"      // PPN Masukan
	TaxTypePPNOut     TaxType = "ppn_out"     // PPN Keluaran
	TaxTypePPh21      TaxType = "pph21"       // PPh Pasal 21
	TaxTypePPh23      TaxType = "pph23"       // PPh Pasal 23
	TaxTypePPh25      TaxType = "pph25"       // PPh Pasal 25
	TaxTypePPh4Ayat2  TaxType = "pph4ayat2"   // PPh Pasal 4 Ayat 2

	TaxStatusDraft    TaxStatus = "draft"
	TaxStatusReported TaxStatus = "reported"
	TaxStatusPaid     TaxStatus = "paid"
)

type Tax struct {
	BaseModel
	CompanyID     uint      `gorm:"not null;index" json:"company_id"`
	Company       Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	TaxNumber     string    `gorm:"uniqueIndex;size:50;not null" json:"tax_number"`
	TaxType       TaxType   `gorm:"type:varchar(20);not null;index" json:"tax_type"`
	TaxPeriod     string    `gorm:"size:7;not null;index" json:"tax_period"` // Format: YYYY-MM
	TaxableAmount float64   `gorm:"type:decimal(20,2);not null" json:"taxable_amount"`
	TaxRate       float64   `gorm:"type:decimal(5,2);not null" json:"tax_rate"` // Percentage
	TaxAmount     float64   `gorm:"type:decimal(20,2);not null" json:"tax_amount"`
	Status        TaxStatus `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`
	DueDate       time.Time `gorm:"not null" json:"due_date"`
	ReportedDate  *time.Time `json:"reported_date"`
	PaidDate      *time.Time `json:"paid_date"`
	Description   string    `gorm:"type:text" json:"description"`
	JournalID     *uint     `gorm:"index" json:"journal_id"`
	Journal       *Journal  `gorm:"foreignKey:JournalID" json:"journal,omitempty"`
	CreatedBy     uint      `gorm:"not null" json:"created_by"`
	User          User      `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
}

type TaxResponse struct {
	ID            uint      `json:"id"`
	CompanyID     uint      `json:"company_id"`
	TaxNumber     string    `json:"tax_number"`
	TaxType       TaxType   `json:"tax_type"`
	TaxPeriod     string    `json:"tax_period"`
	TaxableAmount float64   `json:"taxable_amount"`
	TaxRate       float64   `json:"tax_rate"`
	TaxAmount     float64   `json:"tax_amount"`
	Status        TaxStatus `json:"status"`
	DueDate       string    `json:"due_date"`
	ReportedDate  *string   `json:"reported_date"`
	PaidDate      *string   `json:"paid_date"`
	Description   string    `json:"description"`
	CreatedBy     uint      `json:"created_by"`
	CreatedAt     string    `json:"created_at"`
}

type TaxSummary struct {
	TaxType       TaxType `json:"tax_type"`
	Period        string  `json:"period"`
	TotalTaxable  float64 `json:"total_taxable"`
	TotalTax      float64 `json:"total_tax"`
	TotalReported float64 `json:"total_reported"`
	TotalPaid     float64 `json:"total_paid"`
	Outstanding   float64 `json:"outstanding"`
}