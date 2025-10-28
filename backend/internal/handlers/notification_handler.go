package handlers

import (
	"finara-backend/internal/services"
	"finara-backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler(notificationService services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")

	limitStr := c.Query("limit")
	limit := 20
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	notifications, err := h.notificationService.GetNotificationsByUserID(userID.(uint), limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve notifications", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notifications retrieved successfully", notifications)
}

func (h *NotificationHandler) GetUnreadNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")

	notifications, err := h.notificationService.GetUnreadNotifications(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve unread notifications", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Unread notifications retrieved successfully", notifications)
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	count, err := h.notificationService.GetUnreadCount(userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve unread count", err)
		return
	}

	response := map[string]interface{}{
		"unread_count": count,
	}

	utils.SuccessResponse(c, http.StatusOK, "Unread count retrieved successfully", response)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID", err)
		return
	}

	if err := h.notificationService.MarkAsRead(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to mark notification as read", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification marked as read successfully", nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")

	if err := h.notificationService.MarkAllAsRead(userID.(uint)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to mark all notifications as read", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "All notifications marked as read successfully", nil)
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid notification ID", err)
		return
	}

	if err := h.notificationService.DeleteNotification(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to delete notification", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notification deleted successfully", nil)
}