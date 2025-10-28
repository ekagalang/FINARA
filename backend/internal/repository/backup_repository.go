package repository

import (
	"finara-backend/internal/models"

	"gorm.io/gorm"
)

type BackupRepository interface {
	Create(backup *models.Backup) error
	FindByID(id uint) (*models.Backup, error)
	FindByCompanyID(companyID uint) ([]models.Backup, error)
	Update(backup *models.Backup) error
	Delete(id uint) error
}

type backupRepository struct {
	db *gorm.DB
}

func NewBackupRepository(db *gorm.DB) BackupRepository {
	return &backupRepository{db: db}
}

func (r *backupRepository) Create(backup *models.Backup) error {
	return r.db.Create(backup).Error
}

func (r *backupRepository) FindByID(id uint) (*models.Backup, error) {
	var backup models.Backup
	err := r.db.Preload("User").First(&backup, id).Error
	return &backup, err
}

func (r *backupRepository) FindByCompanyID(companyID uint) ([]models.Backup, error) {
	var backups []models.Backup
	err := r.db.Where("company_id = ?", companyID).
		Order("created_at DESC").
		Preload("User").
		Find(&backups).Error
	return backups, err
}

func (r *backupRepository) Update(backup *models.Backup) error {
	return r.db.Save(backup).Error
}

func (r *backupRepository) Delete(id uint) error {
	return r.db.Delete(&models.Backup{}, id).Error
}