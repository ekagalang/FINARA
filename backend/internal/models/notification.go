package models

import "time"

type NotificationType string
type NotificationStatus string

const (
	NotificationTypeTaxDue       NotificationType = "tax_due"
	NotificationTypeLowCash      NotificationType = "low_cash"
	NotificationTypeHighExpense  NotificationType = "high_expense"
	NotificationTypeJournalDraft NotificationType = "journal_draft"

	NotificationStatusUnread NotificationStatus = "unread"
	NotificationStatusRead   NotificationStatus = "read"
)

type Notification struct {
	BaseModel
	CompanyID   uint               `gorm:"not null;index" json:"company_id"`
	Company     Company            `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	UserID      uint               `gorm:"not null;index" json:"user_id"`
	User        User               `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type        NotificationType   `gorm:"type:varchar(30);not null" json:"type"`
	Title       string             `gorm:"size:255;not null" json:"title"`
	Message     string             `gorm:"type:text;not null" json:"message"`
	Status      NotificationStatus `gorm:"type:varchar(20);not null;default:'unread'" json:"status"`
	RelatedID   *uint              `json:"related_id"` // ID terkait (tax_id, journal_id, etc)
	RelatedType string             `gorm:"size:50" json:"related_type"`
	ReadAt      *time.Time         `json:"read_at"`
}

type NotificationResponse struct {
	ID          uint               `json:"id"`
	CompanyID   uint               `json:"company_id"`
	UserID      uint               `json:"user_id"`
	Type        NotificationType   `json:"type"`
	Title       string             `json:"title"`
	Message     string             `json:"message"`
	Status      NotificationStatus `json:"status"`
	RelatedID   *uint              `json:"related_id"`
	RelatedType string             `json:"related_type"`
	ReadAt      *string            `json:"read_at"`
	CreatedAt   string             `json:"created_at"`
}