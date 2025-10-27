package services

import (
	"errors"
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
)

type CompanyService interface {
	CreateCompany(company *models.Company) error
	GetCompanyByID(id uint) (*models.Company, error)
	GetAllCompanies() ([]models.Company, error)
	UpdateCompany(id uint, updates map[string]interface{}) (*models.Company, error)
	DeleteCompany(id uint) error
}

type companyService struct {
	companyRepo repository.CompanyRepository
}

func NewCompanyService(companyRepo repository.CompanyRepository) CompanyService {
	return &companyService{companyRepo: companyRepo}
}

func (s *companyService) CreateCompany(company *models.Company) error {
	return s.companyRepo.Create(company)
}

func (s *companyService) GetCompanyByID(id uint) (*models.Company, error) {
	return s.companyRepo.FindByID(id)
}

func (s *companyService) GetAllCompanies() ([]models.Company, error) {
	return s.companyRepo.FindAll()
}

func (s *companyService) UpdateCompany(id uint, updates map[string]interface{}) (*models.Company, error) {
	company, err := s.companyRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("company not found")
	}

	if name, ok := updates["name"].(string); ok {
		company.Name = name
	}
	if address, ok := updates["address"].(string); ok {
		company.Address = address
	}
	if phone, ok := updates["phone"].(string); ok {
		company.Phone = phone
	}
	if email, ok := updates["email"].(string); ok {
		company.Email = email
	}
	if description, ok := updates["description"].(string); ok {
		company.Description = description
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		company.IsActive = isActive
	}

	if err := s.companyRepo.Update(company); err != nil {
		return nil, err
	}

	return company, nil
}

func (s *companyService) DeleteCompany(id uint) error {
	return s.companyRepo.Delete(id)
}