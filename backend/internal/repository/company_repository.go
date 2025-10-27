package repository

import (
	"finara-backend/internal/models"

	"gorm.io/gorm"
)

type CompanyRepository interface {
	Create(company *models.Company) error
	FindByID(id uint) (*models.Company, error)
	FindAll() ([]models.Company, error)
	Update(company *models.Company) error
	Delete(id uint) error
}

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) Create(company *models.Company) error {
	return r.db.Create(company).Error
}

func (r *companyRepository) FindByID(id uint) (*models.Company, error) {
	var company models.Company
	err := r.db.First(&company, id).Error
	return &company, err
}

func (r *companyRepository) FindAll() ([]models.Company, error) {
	var companies []models.Company
	err := r.db.Find(&companies).Error
	return companies, err
}

func (r *companyRepository) Update(company *models.Company) error {
	return r.db.Save(company).Error
}

func (r *companyRepository) Delete(id uint) error {
	return r.db.Delete(&models.Company{}, id).Error
}