package models

import "time"

type BackupStatus string

const (
	BackupStatusPending   BackupStatus = "pending"
	BackupStatusInProgress BackupStatus = "in_progress"
	BackupStatusCompleted BackupStatus = "completed"
	BackupStatusFailed    BackupStatus = "failed"
)

type Backup struct {
	BaseModel
	CompanyID    uint         `gorm:"not null;index" json:"company_id"`
	Company      Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	BackupName   string       `gorm:"size:255;not null" json:"backup_name"`
	BackupPath   string       `gorm:"type:text;not null" json:"backup_path"`
	FileSize     int64        `json:"file_size"` // in bytes
	Status       BackupStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	StartedAt    time.Time    `json:"started_at"`
	CompletedAt  *time.Time   `json:"completed_at"`
	ErrorMessage string       `gorm:"type:text" json:"error_message"`
	CreatedBy    uint         `gorm:"not null" json:"created_by"`
	User         User         `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
}

type BackupResponse struct {
	ID           uint         `json:"id"`
	CompanyID    uint         `json:"company_id"`
	BackupName   string       `json:"backup_name"`
	FileSize     int64        `json:"file_size"`
	FileSizeHuman string      `json:"file_size_human"`
	Status       BackupStatus `json:"status"`
	StartedAt    string       `json:"started_at"`
	CompletedAt  *string      `json:"completed_at"`
	CreatedBy    uint         `json:"created_by"`
	CreatedByName string      `json:"created_by_name"`
}