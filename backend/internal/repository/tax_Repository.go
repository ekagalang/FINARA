package repository

import (
	"finara-backend/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type TaxRepository interface {
	Create(tax *models.Tax) error
	FindByID(id uint) (*models.Tax, error)
	FindByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Tax, error)
	FindByPeriod(companyID uint, period string) ([]models.Tax, error)
	FindByType(companyID uint, taxType models.TaxType) ([]models.Tax, error)
	FindDueTaxes(companyID uint, dueDate time.Time) ([]models.Tax, error)
	Update(tax *models.Tax) error
	Delete(id uint) error
	GenerateTaxNumber(companyID uint, taxType models.TaxType, period string) (string, error)
	GetTaxSummary(companyID uint, period string) ([]models.TaxSummary, error)
}

type taxRepository struct {
	db *gorm.DB
}

func NewTaxRepository(db *gorm.DB) TaxRepository {
	return &taxRepository{db: db}
}

func (r *taxRepository) Create(tax *models.Tax) error {
	return r.db.Create(tax).Error
}

func (r *taxRepository) FindByID(id uint) (*models.Tax, error) {
	var tax models.Tax
	err := r.db.Preload("Company").
		Preload("Journal").
		Preload("User").
		First(&tax, id).Error
	return &tax, err
}

func (r *taxRepository) FindByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Tax, error) {
	var taxes []models.Tax
	err := r.db.Where("company_id = ? AND due_date BETWEEN ? AND ?", companyID, startDate, endDate).
		Order("due_date ASC").
		Preload("User").
		Find(&taxes).Error
	return taxes, err
}

func (r *taxRepository) FindByPeriod(companyID uint, period string) ([]models.Tax, error) {
	var taxes []models.Tax
	err := r.db.Where("company_id = ? AND tax_period = ?", companyID, period).
		Order("tax_type ASC").
		Find(&taxes).Error
	return taxes, err
}

func (r *taxRepository) FindByType(companyID uint, taxType models.TaxType) ([]models.Tax, error) {
	var taxes []models.Tax
	err := r.db.Where("company_id = ? AND tax_type = ?", companyID, taxType).
		Order("tax_period DESC").
		Find(&taxes).Error
	return taxes, err
}

func (r *taxRepository) FindDueTaxes(companyID uint, dueDate time.Time) ([]models.Tax, error) {
	var taxes []models.Tax
	err := r.db.Where("company_id = ? AND due_date <= ? AND status != ?", companyID, dueDate, models.TaxStatusPaid).
		Order("due_date ASC").
		Find(&taxes).Error
	return taxes, err
}

func (r *taxRepository) Update(tax *models.Tax) error {
	return r.db.Save(tax).Error
}

func (r *taxRepository) Delete(id uint) error {
	return r.db.Delete(&models.Tax{}, id).Error
}

func (r *taxRepository) GenerateTaxNumber(companyID uint, taxType models.TaxType, period string) (string, error) {
	var count int64
	var prefix string

	switch taxType {
	case models.TaxTypePPNIn:
		prefix = "PPNM/"
	case models.TaxTypePPNOut:
		prefix = "PPNK/"
	case models.TaxTypePPh21:
		prefix = "PPH21/"
	case models.TaxTypePPh23:
		prefix = "PPH23/"
	case models.TaxTypePPh25:
		prefix = "PPH25/"
	case models.TaxTypePPh4Ayat2:
		prefix = "PPH4A2/"
	default:
		prefix = "TAX/"
	}

	prefix += period + "/"

	err := r.db.Model(&models.Tax{}).
		Where("company_id = ? AND tax_number LIKE ?", companyID, prefix+"%").
		Count(&count).Error

	if err != nil {
		return "", err
	}

	return prefix + fmt.Sprintf("%04d", count+1), nil
}

func (r *taxRepository) GetTaxSummary(companyID uint, period string) ([]models.TaxSummary, error) {
	var summaries []models.TaxSummary

	err := r.db.Raw(`
		SELECT 
			tax_type,
			tax_period as period,
			COALESCE(SUM(taxable_amount), 0) as total_taxable,
			COALESCE(SUM(tax_amount), 0) as total_tax,
			COALESCE(SUM(CASE WHEN status IN ('reported', 'paid') THEN tax_amount ELSE 0 END), 0) as total_reported,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN tax_amount ELSE 0 END), 0) as total_paid,
			COALESCE(SUM(CASE WHEN status != 'paid' THEN tax_amount ELSE 0 END), 0) as outstanding
		FROM taxes
		WHERE company_id = ? AND tax_period = ?
		GROUP BY tax_type, tax_period
		ORDER BY tax_type ASC
	`, companyID, period).Scan(&summaries).Error

	return summaries, err
}