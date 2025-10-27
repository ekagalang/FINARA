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
		&models.Account{},        // NEW
		&models.Journal{},        // NEW
		&models.JournalEntry{},   // NEW
		&models.Ledger{},         // NEW
	)

	if err != nil {
		return err
	}

	log.Println("Migrations completed successfully")
	return nil
}