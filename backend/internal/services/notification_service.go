package services

import (
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"time"
)

type NotificationService interface {
	CreateNotification(notification *models.Notification) error
	GetNotificationByID(id uint) (*models.Notification, error)
	GetNotificationsByUserID(userID uint, limit int) ([]models.Notification, error)
	GetUnreadNotifications(userID uint) ([]models.Notification, error)
	MarkAsRead(id uint) error
	MarkAllAsRead(userID uint) error
	DeleteNotification(id uint) error
	GetUnreadCount(userID uint) (int64, error)
	CheckAndCreateTaxDueNotifications(companyID uint) error
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
	taxRepo          repository.TaxRepository
	userRepo         repository.UserRepository
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	taxRepo repository.TaxRepository,
	userRepo repository.UserRepository,
) NotificationService {
	return &notificationService{
		notificationRepo: notificationRepo,
		taxRepo:          taxRepo,
		userRepo:         userRepo,
	}
}

func (s *notificationService) CreateNotification(notification *models.Notification) error {
	if notification.Status == "" {
		notification.Status = models.NotificationStatusUnread
	}
	return s.notificationRepo.Create(notification)
}

func (s *notificationService) GetNotificationByID(id uint) (*models.Notification, error) {
	return s.notificationRepo.FindByID(id)
}

func (s *notificationService) GetNotificationsByUserID(userID uint, limit int) ([]models.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.notificationRepo.FindByUserID(userID, limit)
}

func (s *notificationService) GetUnreadNotifications(userID uint) ([]models.Notification, error) {
	return s.notificationRepo.FindUnreadByUserID(userID)
}

func (s *notificationService) MarkAsRead(id uint) error {
	return s.notificationRepo.MarkAsRead(id)
}

func (s *notificationService) MarkAllAsRead(userID uint) error {
	return s.notificationRepo.MarkAllAsRead(userID)
}

func (s *notificationService) DeleteNotification(id uint) error {
	return s.notificationRepo.Delete(id)
}

func (s *notificationService) GetUnreadCount(userID uint) (int64, error) {
	return s.notificationRepo.CountUnread(userID)
}

func (s *notificationService) CheckAndCreateTaxDueNotifications(companyID uint) error {
	// Check for taxes due in next 7 days
	dueDate := time.Now().AddDate(0, 0, 7)
	
	taxes, err := s.taxRepo.FindDueTaxes(companyID, dueDate)
	if err != nil {
		return err
	}

	// Get admin/accountant users for this company
	users, err := s.userRepo.FindAll()
	if err != nil {
		return err
	}

	for _, tax := range taxes {
		for _, user := range users {
			if user.CompanyID == companyID && (user.Role == models.RoleAdmin || user.Role == models.RoleAccountant) {
				// Check if notification already exists
				// (in real implementation, you'd want to prevent duplicate notifications)
				
				notification := &models.Notification{
					CompanyID:   companyID,
					UserID:      user.ID,
					Type:        models.NotificationTypeTaxDue,
					Title:       "Pajak Jatuh Tempo",
					Message:     "Pajak " + string(tax.TaxType) + " periode " + tax.TaxPeriod + " akan jatuh tempo pada " + tax.DueDate.Format("2006-01-02"),
					Status:      models.NotificationStatusUnread,
					RelatedID:   &tax.ID,
					RelatedType: "tax",
				}

				s.notificationRepo.Create(notification)
			}
		}
	}

	return nil
}