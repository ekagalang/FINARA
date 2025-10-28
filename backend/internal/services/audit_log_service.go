package services

import (
	"encoding/json"
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"time"
)

type AuditLogService interface {
	LogAction(log *models.AuditLog) error
	GetAuditLogByID(id uint) (*models.AuditLog, error)
	GetAuditLogsByCompanyID(companyID uint, limit int) ([]models.AuditLog, error)
	GetAuditLogsByUserID(userID uint, limit int) ([]models.AuditLog, error)
	GetAuditLogsByModule(companyID uint, module string, limit int) ([]models.AuditLog, error)
	GetAuditLogsByDateRange(companyID uint, startDate, endDate time.Time) ([]models.AuditLog, error)
	GetAuditLogsByAction(companyID uint, action models.AuditAction, limit int) ([]models.AuditLog, error)
	DeleteAuditLog(id uint) error
	CleanupOldLogs(days int) error
}

type auditLogService struct {
	auditLogRepo repository.AuditLogRepository
}

func NewAuditLogService(auditLogRepo repository.AuditLogRepository) AuditLogService {
	return &auditLogService{auditLogRepo: auditLogRepo}
}

func (s *auditLogService) LogAction(log *models.AuditLog) error {
	log.PerformedAt = time.Now()
	return s.auditLogRepo.Create(log)
}

func (s *auditLogService) GetAuditLogByID(id uint) (*models.AuditLog, error) {
	return s.auditLogRepo.FindByID(id)
}

func (s *auditLogService) GetAuditLogsByCompanyID(companyID uint, limit int) ([]models.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.auditLogRepo.FindByCompanyID(companyID, limit)
}

func (s *auditLogService) GetAuditLogsByUserID(userID uint, limit int) ([]models.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.auditLogRepo.FindByUserID(userID, limit)
}

func (s *auditLogService) GetAuditLogsByModule(companyID uint, module string, limit int) ([]models.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.auditLogRepo.FindByModule(companyID, module, limit)
}

func (s *auditLogService) GetAuditLogsByDateRange(companyID uint, startDate, endDate time.Time) ([]models.AuditLog, error) {
	return s.auditLogRepo.FindByDateRange(companyID, startDate, endDate)
}

func (s *auditLogService) GetAuditLogsByAction(companyID uint, action models.AuditAction, limit int) ([]models.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	return s.auditLogRepo.FindByAction(companyID, action, limit)
}

func (s *auditLogService) DeleteAuditLog(id uint) error {
	return s.auditLogRepo.Delete(id)
}

func (s *auditLogService) CleanupOldLogs(days int) error {
	beforeDate := time.Now().AddDate(0, 0, -days)
	return s.auditLogRepo.DeleteOldLogs(beforeDate)
}

// Helper function to create audit log from any object
func CreateAuditLog(companyID, userID uint, action models.AuditAction, module string, recordID *uint, recordType string, oldValue, newValue interface{}, description, ipAddress, userAgent string) *models.AuditLog {
	var oldJSON, newJSON string

	if oldValue != nil {
		if bytes, err := json.Marshal(oldValue); err == nil {
			oldJSON = string(bytes)
		}
	}

	if newValue != nil {
		if bytes, err := json.Marshal(newValue); err == nil {
			newJSON = string(bytes)
		}
	}

	return &models.AuditLog{
		CompanyID:   companyID,
		UserID:      userID,
		Action:      action,
		Module:      module,
		RecordID:    recordID,
		RecordType:  recordType,
		OldValue:    oldJSON,
		NewValue:    newJSON,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Description: description,
		PerformedAt: time.Now(),
	}
}