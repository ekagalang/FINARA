package services

import (
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"time"
)

type LedgerService interface {
	GetLedgerByAccountID(accountID uint, startDate, endDate time.Time) ([]models.Ledger, error)
	GetLedgerByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Ledger, error)
	GetTrialBalance(companyID uint, endDate time.Time) (*models.TrialBalanceReport, error)
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

func (s *ledgerService) GetTrialBalance(companyID uint, endDate time.Time) (*models.TrialBalanceReport, error) {
	accounts, err := s.ledgerRepo.GetTrialBalance(companyID, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate totals
	var totalDebit, totalCredit, totalDebitBalance, totalCreditBalance float64
	for _, account := range accounts {
		totalDebit += account.Debit
		totalCredit += account.Credit
		totalDebitBalance += account.DebitBalance
		totalCreditBalance += account.CreditBalance
	}

	// Build the report
	report := &models.TrialBalanceReport{
		AsOfDate:           endDate.Format("2006-01-02"),
		Accounts:           accounts,
		TotalDebit:         totalDebit,
		TotalCredit:        totalCredit,
		TotalDebitBalance:  totalDebitBalance,
		TotalCreditBalance: totalCreditBalance,
		IsBalanced:         totalDebitBalance == totalCreditBalance,
	}

	return report, nil
}

func (s *ledgerService) GetAccountBalance(accountID uint, endDate time.Time) (float64, error) {
	return s.ledgerRepo.GetAccountBalance(accountID, endDate)
}