package repository

import (
	"finara-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	FindByID(id uint) (*models.Notification, error)
	FindByUserID(userID uint, limit int) ([]models.Notification, error)
	FindUnreadByUserID(userID uint) ([]models.Notification, error)
	MarkAsRead(id uint) error
	MarkAllAsRead(userID uint) error
	Delete(id uint) error
	CountUnread(userID uint) (int64, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindByID(id uint) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.First(&notification, id).Error
	return &notification, err
}

func (r *notificationRepository) FindByUserID(userID uint, limit int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) FindUnreadByUserID(userID uint) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ? AND status = ?", userID, models.NotificationStatusUnread).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) MarkAsRead(id uint) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":  models.NotificationStatusRead,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) MarkAllAsRead(userID uint) error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).
		Where("user_id = ? AND status = ?", userID, models.NotificationStatusUnread).
		Updates(map[string]interface{}{
			"status":  models.NotificationStatusRead,
			"read_at": now,
		}).Error
}

func (r *notificationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Notification{}, id).Error
}

func (r *notificationRepository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).
		Where("user_id = ? AND status = ?", userID, models.NotificationStatusUnread).
		Count(&count).Error
	return count, err
}