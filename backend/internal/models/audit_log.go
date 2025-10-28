package models

import "time"

type AuditAction string

const (
	ActionCreate AuditAction = "create"
	ActionUpdate AuditAction = "update"
	ActionDelete AuditAction = "delete"
	ActionView   AuditAction = "view"
	ActionLogin  AuditAction = "login"
	ActionLogout AuditAction = "logout"
	ActionExport AuditAction = "export"
	ActionPost   AuditAction = "post"
	ActionVoid   AuditAction = "void"
)

type AuditLog struct {
	BaseModel
	CompanyID    uint        `gorm:"not null;index" json:"company_id"`
	Company      Company     `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	UserID       uint        `gorm:"not null;index" json:"user_id"`
	User         User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action       AuditAction `gorm:"type:varchar(20);not null;index" json:"action"`
	Module       string      `gorm:"size:50;not null;index" json:"module"` // journal, account, tax, etc
	RecordID     *uint       `gorm:"index" json:"record_id"`
	RecordType   string      `gorm:"size:50" json:"record_type"` // journal, account, tax, etc
	OldValue     string      `gorm:"type:text" json:"old_value"` // JSON
	NewValue     string      `gorm:"type:text" json:"new_value"` // JSON
	IPAddress    string      `gorm:"size:45" json:"ip_address"`
	UserAgent    string      `gorm:"type:text" json:"user_agent"`
	Description  string      `gorm:"type:text" json:"description"`
	PerformedAt  time.Time   `gorm:"not null;index" json:"performed_at"`
}

type AuditLogResponse struct {
	ID          uint        `json:"id"`
	CompanyID   uint        `json:"company_id"`
	UserID      uint        `json:"user_id"`
	UserName    string      `json:"user_name"`
	Action      AuditAction `json:"action"`
	Module      string      `json:"module"`
	RecordID    *uint       `json:"record_id"`
	RecordType  string      `json:"record_type"`
	Description string      `json:"description"`
	IPAddress   string      `json:"ip_address"`
	PerformedAt string      `json:"performed_at"`
}