package repository

import (
	"finara-backend/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type JournalRepository interface {
	Create(journal *models.Journal) error
	FindByID(id uint) (*models.Journal, error)
	FindByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Journal, error)
	FindByStatus(companyID uint, status models.JournalStatus) ([]models.Journal, error)
	Update(journal *models.Journal) error
	Delete(id uint) error
	GenerateJournalNumber(companyID uint, date time.Time) (string, error)
}

type journalRepository struct {
	db *gorm.DB
}

func NewJournalRepository(db *gorm.DB) JournalRepository {
	return &journalRepository{db: db}
}

func (r *journalRepository) Create(journal *models.Journal) error {
	return r.db.Create(journal).Error
}

func (r *journalRepository) FindByID(id uint) (*models.Journal, error) {
	var journal models.Journal
	err := r.db.Preload("Entries.Account").
		Preload("User").
		First(&journal, id).Error
	return &journal, err
}

func (r *journalRepository) FindByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Journal, error) {
	var journals []models.Journal
	err := r.db.Where("company_id = ? AND transaction_date BETWEEN ? AND ?", companyID, startDate, endDate).
		Order("transaction_date DESC, journal_number DESC").
		Preload("Entries.Account").
		Preload("User").
		Find(&journals).Error
	return journals, err
}

func (r *journalRepository) FindByStatus(companyID uint, status models.JournalStatus) ([]models.Journal, error) {
	var journals []models.Journal
	err := r.db.Where("company_id = ? AND status = ?", companyID, status).
		Order("transaction_date DESC").
		Preload("Entries.Account").
		Find(&journals).Error
	return journals, err
}

func (r *journalRepository) Update(journal *models.Journal) error {
	return r.db.Save(journal).Error
}

func (r *journalRepository) Delete(id uint) error {
	return r.db.Select("Entries").Delete(&models.Journal{BaseModel: models.BaseModel{ID: id}}).Error
}

func (r *journalRepository) GenerateJournalNumber(companyID uint, date time.Time) (string, error) {
	var count int64
	prefix := date.Format("JRN/200601/")
	
	err := r.db.Model(&models.Journal{}).
		Where("company_id = ? AND journal_number LIKE ?", companyID, prefix+"%").
		Count(&count).Error
	
	if err != nil {
		return "", err
	}
	
	return prefix + fmt.Sprintf("%04d", count+1), nil
}