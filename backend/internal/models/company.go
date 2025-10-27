package models

type Company struct {
	BaseModel
	Name        string `gorm:"not null;size:255" json:"name"`
	Address     string `gorm:"type:text" json:"address"`
	Phone       string `gorm:"size:50" json:"phone"`
	Email       string `gorm:"size:255" json:"email"`
	TaxID       string `gorm:"uniqueIndex;size:100" json:"tax_id"` // NPWP
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}