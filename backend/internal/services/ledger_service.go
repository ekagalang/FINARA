package services

import (
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"time"
)

type LedgerService interface {
	GetLedgerByAccountID(accountID uint, startDate, endDate time.Time) ([]models.Ledger, error)
	GetLedgerByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Ledger, error)
	GetTrialBalance(companyID uint, endDate time.Time) ([]models.TrialBalanceResponse, error)
	GetAccountBalance(accountID uint, endDate time.Time) (float64, error)
}

type ledgerService struct {
	ledgerRepo repository.LedgerRepository
}

func NewLedgerService(ledgerRepo repository.LedgerRepository) LedgerService {
	return &ledgerService{ledgerRepo: ledgerRepo}
}

func (s *ledgerService) GetLedgerByAccountID(accountID uint, startDate, endDate time.Time) ([]models.Ledger, error) {
	return s.ledgerRepo.FindByAccountID(accountID, startDate, endDate)
}

func (s *ledgerService) GetLedgerByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Ledger, error) {
	return s.ledgerRepo.FindByCompanyID(companyID, startDate, endDate)
}

func (s *ledgerService) GetTrialBalance(companyID uint, endDate time.Time) ([]models.TrialBalanceResponse, error) {
	return s.ledgerRepo.GetTrialBalance(companyID, endDate)
}

func (s *ledgerService) GetAccountBalance(accountID uint, endDate time.Time) (float64, error) {
	return s.ledgerRepo.GetAccountBalance(accountID, endDate)
}