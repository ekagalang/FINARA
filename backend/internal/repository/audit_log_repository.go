package repository

import (
	"finara-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type AuditLogRepository interface {
	Create(log *models.AuditLog) error
	FindByID(id uint) (*models.AuditLog, error)
	FindByCompanyID(companyID uint, limit int) ([]models.AuditLog, error)
	FindByUserID(userID uint, limit int) ([]models.AuditLog, error)
	FindByModule(companyID uint, module string, limit int) ([]models.AuditLog, error)
	FindByDateRange(companyID uint, startDate, endDate time.Time) ([]models.AuditLog, error)
	FindByAction(companyID uint, action models.AuditAction, limit int) ([]models.AuditLog, error)
	Delete(id uint) error
	DeleteOldLogs(beforeDate time.Time) error
}

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) AuditLogRepository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditLogRepository) FindByID(id uint) (*models.AuditLog, error) {
	var log models.AuditLog
	err := r.db.Preload("User").
		Preload("Company").
		First(&log, id).Error
	return &log, err
}

func (r *auditLogRepository) FindByCompanyID(companyID uint, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := r.db.Where("company_id = ?", companyID).
		Order("performed_at DESC").
		Preload("User")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindByUserID(userID uint, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := r.db.Where("user_id = ?", userID).
		Order("performed_at DESC").
		Preload("User")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindByModule(companyID uint, module string, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := r.db.Where("company_id = ? AND module = ?", companyID, module).
		Order("performed_at DESC").
		Preload("User")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindByDateRange(companyID uint, startDate, endDate time.Time) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	err := r.db.Where("company_id = ? AND performed_at BETWEEN ? AND ?", companyID, startDate, endDate).
		Order("performed_at DESC").
		Preload("User").
		Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) FindByAction(companyID uint, action models.AuditAction, limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	query := r.db.Where("company_id = ? AND action = ?", companyID, action).
		Order("performed_at DESC").
		Preload("User")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&logs).Error
	return logs, err
}

func (r *auditLogRepository) Delete(id uint) error {
	return r.db.Delete(&models.AuditLog{}, id).Error
}

func (r *auditLogRepository) DeleteOldLogs(beforeDate time.Time) error {
	return r.db.Where("performed_at < ?", beforeDate).Delete(&models.AuditLog{}).Error
}