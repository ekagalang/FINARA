package handlers

import (
	"finara-backend/internal/services"
	"finara-backend/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BackupHandler struct {
	backupService services.BackupService
}

func NewBackupHandler(backupService services.BackupService) *BackupHandler {
	return &BackupHandler{backupService: backupService}
}

func (h *BackupHandler) CreateBackup(c *gin.Context) {
	companyID, _ := c.Get("company_id")
	userID, _ := c.Get("user_id")

	backup, err := h.backupService.CreateBackup(companyID.(uint), userID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create backup", err)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Backup created successfully", backup)
}

func (h *BackupHandler) GetBackups(c *gin.Context) {
	companyID, _ := c.Get("company_id")

	backups, err := h.backupService.GetBackupsByCompanyID(companyID.(uint))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve backups", err)
		return
	}

	// Format file sizes
	for i := range backups {
		backups[i].BackupName = backups[i].BackupName + " (" + services.FormatFileSize(backups[i].FileSize) + ")"
	}

	utils.SuccessResponse(c, http.StatusOK, "Backups retrieved successfully", backups)
}

func (h *BackupHandler) GetBackupByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid backup ID", err)
		return
	}

	backup, err := h.backupService.GetBackupByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Backup not found", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Backup retrieved successfully", backup)
}

func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid backup ID", err)
		return
	}

	if err := h.backupService.RestoreBackup(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to restore backup", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Backup restored successfully", nil)
}

func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid backup ID", err)
		return
	}

	if err := h.backupService.DeleteBackup(uint(id)); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete backup", err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Backup deleted successfully", nil)
}

func (h *BackupHandler) DownloadBackup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid backup ID", err)
		return
	}

	backup, err := h.backupService.GetBackupByID(uint(id))
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Backup not found", err)
		return
	}

	// Return file for download
	c.FileAttachment(backup.BackupPath, backup.BackupName)
}