package services

import (
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type BackupService interface {
	CreateBackup(companyID uint, createdBy uint) (*models.Backup, error)
	GetBackupByID(id uint) (*models.Backup, error)
	GetBackupsByCompanyID(companyID uint) ([]models.Backup, error)
	RestoreBackup(id uint) error
	DeleteBackup(id uint) error
}

type backupService struct {
	backupRepo repository.BackupRepository
	dbConfig   *DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewBackupService(backupRepo repository.BackupRepository, dbConfig *DatabaseConfig) BackupService {
	return &backupService{
		backupRepo: backupRepo,
		dbConfig:   dbConfig,
	}
}

func (s *backupService) CreateBackup(companyID uint, createdBy uint) (*models.Backup, error) {
	// Create backup directory if not exists
	backupDir := "backups"
	os.MkdirAll(backupDir, 0755)

	// Generate backup filename
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("backup_company_%d_%s.sql", companyID, timestamp)
	backupPath := fmt.Sprintf("%s/%s", backupDir, backupName)

	// Create backup record
	backup := &models.Backup{
		CompanyID:  companyID,
		BackupName: backupName,
		BackupPath: backupPath,
		Status:     models.BackupStatusInProgress,
		StartedAt:  time.Now(),
		CreatedBy:  createdBy,
	}

	if err := s.backupRepo.Create(backup); err != nil {
		return nil, err
	}

	// Perform backup in background
	go s.performBackup(backup)

	return backup, nil
}

func (s *backupService) performBackup(backup *models.Backup) {
	// mysqldump command
	cmd := exec.Command("mysqldump",
		"-h", s.dbConfig.Host,
		"-P", s.dbConfig.Port,
		"-u", s.dbConfig.User,
		fmt.Sprintf("-p%s", s.dbConfig.Password),
		s.dbConfig.DBName,
	)

	// Create output file
	outfile, err := os.Create(backup.BackupPath)
	if err != nil {
		backup.Status = models.BackupStatusFailed
		backup.ErrorMessage = err.Error()
		s.backupRepo.Update(backup)
		return
	}
	defer outfile.Close()

	cmd.Stdout = outfile

	// Execute backup
	if err := cmd.Run(); err != nil {
		backup.Status = models.BackupStatusFailed
		backup.ErrorMessage = err.Error()
		s.backupRepo.Update(backup)
		return
	}

	// Get file size
	fileInfo, err := os.Stat(backup.BackupPath)
	if err != nil {
		backup.Status = models.BackupStatusFailed
		backup.ErrorMessage = err.Error()
		s.backupRepo.Update(backup)
		return
	}

	// Update backup record
	now := time.Now()
	backup.Status = models.BackupStatusCompleted
	backup.FileSize = fileInfo.Size()
	backup.CompletedAt = &now
	s.backupRepo.Update(backup)
}

func (s *backupService) GetBackupByID(id uint) (*models.Backup, error) {
	return s.backupRepo.FindByID(id)
}

func (s *backupService) GetBackupsByCompanyID(companyID uint) ([]models.Backup, error) {
	return s.backupRepo.FindByCompanyID(companyID)
}

func (s *backupService) RestoreBackup(id uint) error {
	backup, err := s.backupRepo.FindByID(id)
	if err != nil {
		return err
	}

	if backup.Status != models.BackupStatusCompleted {
		return fmt.Errorf("cannot restore incomplete backup")
	}

	// Check if backup file exists
	if _, err := os.Stat(backup.BackupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found")
	}

	// mysql restore command
	cmd := exec.Command("mysql",
		"-h", s.dbConfig.Host,
		"-P", s.dbConfig.Port,
		"-u", s.dbConfig.User,
		fmt.Sprintf("-p%s", s.dbConfig.Password),
		s.dbConfig.DBName,
	)

	// Open backup file
	infile, err := os.Open(backup.BackupPath)
	if err != nil {
		return err
	}
	defer infile.Close()

	cmd.Stdin = infile

	// Execute restore
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (s *backupService) DeleteBackup(id uint) error {
	backup, err := s.backupRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Delete backup file
	if err := os.Remove(backup.BackupPath); err != nil {
		// Continue even if file deletion fails
	}

	// Delete backup record
	return s.backupRepo.Delete(id)
}

// Helper function to format file size
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}