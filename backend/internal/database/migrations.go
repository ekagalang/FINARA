package database

import (
	"finara-backend/internal/models"
	"log"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Company{},
		&models.Account{},
		&models.Journal{},
		&models.JournalEntry{},
		&models.Ledger{},
		&models.CashBankTransaction{},
		&models.BankReconciliation{},
		&models.Tax{},
		&models.Notification{},
		&models.Product{},
		&models.StockMovement{},
		&models.StockBalance{},
		&models.StockOpname{},
		&models.StockOpnameItem{},
		&models.AuditLog{},
		&models.Backup{},
	)

	if err != nil {
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}