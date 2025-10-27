package repository

import (
	"finara-backend/internal/models"

	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *models.Account) error
	FindByID(id uint) (*models.Account, error)
	FindByCode(companyID uint, code string) (*models.Account, error)
	FindByCompanyID(companyID uint) ([]models.Account, error)
	FindByType(companyID uint, accountType models.AccountType) ([]models.Account, error)
	Update(account *models.Account) error
	Delete(id uint) error
	GetActiveAccounts(companyID uint) ([]models.Account, error)
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *models.Account) error {
	return r.db.Create(account).Error
}

func (r *accountRepository) FindByID(id uint) (*models.Account, error) {
	var account models.Account
	err := r.db.Preload("Parent").Preload("Company").First(&account, id).Error
	return &account, err
}

func (r *accountRepository) FindByCode(companyID uint, code string) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("company_id = ? AND code = ?", companyID, code).First(&account).Error
	return &account, err
}

func (r *accountRepository) FindByCompanyID(companyID uint) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("company_id = ?", companyID).
		Order("code ASC").
		Preload("Parent").
		Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) FindByType(companyID uint, accountType models.AccountType) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("company_id = ? AND type = ?", companyID, accountType).
		Order("code ASC").
		Find(&accounts).Error
	return accounts, err
}

func (r *accountRepository) Update(account *models.Account) error {
	return r.db.Save(account).Error
}

func (r *accountRepository) Delete(id uint) error {
	return r.db.Delete(&models.Account{}, id).Error
}

func (r *accountRepository) GetActiveAccounts(companyID uint) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.Where("company_id = ? AND is_active = ? AND is_header = ?", companyID, true, false).
		Order("code ASC").
		Find(&accounts).Error
	return accounts, err
}