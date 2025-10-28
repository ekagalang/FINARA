package services

import (
	"errors"
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"time"
)

type TaxService interface {
	CreateTax(tax *models.Tax) error
	GetTaxByID(id uint) (*models.Tax, error)
	GetTaxesByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Tax, error)
	GetTaxesByPeriod(companyID uint, period string) ([]models.Tax, error)
	GetTaxesByType(companyID uint, taxType models.TaxType) ([]models.Tax, error)
	GetDueTaxes(companyID uint, dueDate time.Time) ([]models.Tax, error)
	UpdateTax(id uint, tax *models.Tax) error
	DeleteTax(id uint) error
	MarkAsReported(id uint) error
	MarkAsPaid(id uint) error
	GetTaxSummary(companyID uint, period string) ([]models.TaxSummary, error)
	CalculatePPN(companyID uint, period string) error
	CalculatePPh21(companyID uint, period string) error
}

type taxService struct {
	taxRepo repository.TaxRepository
}

func NewTaxService(taxRepo repository.TaxRepository) TaxService {
	return &taxService{taxRepo: taxRepo}
}

func (s *taxService) CreateTax(tax *models.Tax) error {
	// Generate tax number
	taxNumber, err := s.taxRepo.GenerateTaxNumber(tax.CompanyID, tax.TaxType, tax.TaxPeriod)
	if err != nil {
		return err
	}
	tax.TaxNumber = taxNumber

	// Calculate tax amount
	tax.TaxAmount = tax.TaxableAmount * (tax.TaxRate / 100)

	// Set default status
	if tax.Status == "" {
		tax.Status = models.TaxStatusDraft
	}

	return s.taxRepo.Create(tax)
}

func (s *taxService) GetTaxByID(id uint) (*models.Tax, error) {
	return s.taxRepo.FindByID(id)
}

func (s *taxService) GetTaxesByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Tax, error) {
	return s.taxRepo.FindByCompanyID(companyID, startDate, endDate)
}

func (s *taxService) GetTaxesByPeriod(companyID uint, period string) ([]models.Tax, error) {
	return s.taxRepo.FindByPeriod(companyID, period)
}

func (s *taxService) GetTaxesByType(companyID uint, taxType models.TaxType) ([]models.Tax, error) {
	return s.taxRepo.FindByType(companyID, taxType)
}

func (s *taxService) GetDueTaxes(companyID uint, dueDate time.Time) ([]models.Tax, error) {
	return s.taxRepo.FindDueTaxes(companyID, dueDate)
}

func (s *taxService) UpdateTax(id uint, updatedTax *models.Tax) error {
	tax, err := s.taxRepo.FindByID(id)
	if err != nil {
		return errors.New("tax not found")
	}

	// Validasi: hanya draft yang bisa diupdate
	if tax.Status != models.TaxStatusDraft {
		return errors.New("only draft taxes can be updated")
	}

	// Recalculate tax amount
	updatedTax.TaxAmount = updatedTax.TaxableAmount * (updatedTax.TaxRate / 100)
	updatedTax.ID = id

	return s.taxRepo.Update(updatedTax)
}

func (s *taxService) DeleteTax(id uint) error {
	tax, err := s.taxRepo.FindByID(id)
	if err != nil {
		return errors.New("tax not found")
	}

	// Validasi: hanya draft yang bisa dihapus
	if tax.Status != models.TaxStatusDraft {
		return errors.New("only draft taxes can be deleted")
	}

	return s.taxRepo.Delete(id)
}

func (s *taxService) MarkAsReported(id uint) error {
	tax, err := s.taxRepo.FindByID(id)
	if err != nil {
		return errors.New("tax not found")
	}

	now := time.Now()
	tax.Status = models.TaxStatusReported
	tax.ReportedDate = &now

	return s.taxRepo.Update(tax)
}

func (s *taxService) MarkAsPaid(id uint) error {
	tax, err := s.taxRepo.FindByID(id)
	if err != nil {
		return errors.New("tax not found")
	}

	now := time.Now()
	tax.Status = models.TaxStatusPaid
	tax.PaidDate = &now

	return s.taxRepo.Update(tax)
}

func (s *taxService) GetTaxSummary(companyID uint, period string) ([]models.TaxSummary, error) {
	return s.taxRepo.GetTaxSummary(companyID, period)
}

func (s *taxService) CalculatePPN(companyID uint, period string) error {
	// This is a simplified calculation
	// In real implementation, you would calculate from invoices/transactions
	
	// For now, return not implemented
	return errors.New("PPN calculation not yet implemented")
}

func (s *taxService) CalculatePPh21(companyID uint, period string) error {
	// This is a simplified calculation
	// In real implementation, you would calculate from payroll data
	
	// For now, return not implemented
	return errors.New("PPh21 calculation not yet implemented")
}