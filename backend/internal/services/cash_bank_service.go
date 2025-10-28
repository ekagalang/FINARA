package services

import (
	"errors"
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"time"
)

type CashBankService interface {
	CreateTransaction(transaction *models.CashBankTransaction) error
	GetTransactionByID(id uint) (*models.CashBankTransaction, error)
	GetTransactionsByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.CashBankTransaction, error)
	GetTransactionsByAccountID(accountID uint, startDate, endDate time.Time) ([]models.CashBankTransaction, error)
	UpdateTransaction(id uint, transaction *models.CashBankTransaction) error
	DeleteTransaction(id uint) error
	GetCashPosition(companyID uint, endDate time.Time) (map[string]float64, error)
	CreateCashInWithJournal(transaction *models.CashBankTransaction, contraAccountID uint) error
	CreateCashOutWithJournal(transaction *models.CashBankTransaction, contraAccountID uint) error
}

type cashBankService struct {
	cashBankRepo repository.CashBankRepository
	journalRepo  repository.JournalRepository
	ledgerRepo   repository.LedgerRepository
	accountRepo  repository.AccountRepository
}

func NewCashBankService(
	cashBankRepo repository.CashBankRepository,
	journalRepo repository.JournalRepository,
	ledgerRepo repository.LedgerRepository,
	accountRepo repository.AccountRepository,
) CashBankService {
	return &cashBankService{
		cashBankRepo: cashBankRepo,
		journalRepo:  journalRepo,
		ledgerRepo:   ledgerRepo,
		accountRepo:  accountRepo,
	}
}

func (s *cashBankService) CreateTransaction(transaction *models.CashBankTransaction) error {
	// Generate transaction number
	transactionNumber, err := s.cashBankRepo.GenerateTransactionNumber(
		transaction.CompanyID,
		transaction.Type,
		transaction.TransactionDate,
	)
	if err != nil {
		return err
	}
	transaction.TransactionNumber = transactionNumber

	return s.cashBankRepo.Create(transaction)
}

func (s *cashBankService) GetTransactionByID(id uint) (*models.CashBankTransaction, error) {
	return s.cashBankRepo.FindByID(id)
}

func (s *cashBankService) GetTransactionsByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.CashBankTransaction, error) {
	return s.cashBankRepo.FindByCompanyID(companyID, startDate, endDate)
}

func (s *cashBankService) GetTransactionsByAccountID(accountID uint, startDate, endDate time.Time) ([]models.CashBankTransaction, error) {
	return s.cashBankRepo.FindByAccountID(accountID, startDate, endDate)
}

func (s *cashBankService) UpdateTransaction(id uint, updatedTransaction *models.CashBankTransaction) error {
	transaction, err := s.cashBankRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Validasi: jika sudah terhubung dengan journal, tidak bisa diupdate
	if transaction.JournalID != nil {
		return errors.New("cannot update transaction that is already linked to a journal")
	}

	updatedTransaction.ID = id
	return s.cashBankRepo.Update(updatedTransaction)
}

func (s *cashBankService) DeleteTransaction(id uint) error {
	transaction, err := s.cashBankRepo.FindByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Validasi: jika sudah terhubung dengan journal, tidak bisa dihapus
	if transaction.JournalID != nil {
		return errors.New("cannot delete transaction that is already linked to a journal")
	}

	return s.cashBankRepo.Delete(id)
}

func (s *cashBankService) GetCashPosition(companyID uint, endDate time.Time) (map[string]float64, error) {
	return s.cashBankRepo.GetCashPosition(companyID, endDate)
}

func (s *cashBankService) CreateCashInWithJournal(transaction *models.CashBankTransaction, contraAccountID uint) error {
	// Generate transaction number
	transactionNumber, err := s.cashBankRepo.GenerateTransactionNumber(
		transaction.CompanyID,
		models.TransactionTypeIn,
		transaction.TransactionDate,
	)
	if err != nil {
		return err
	}
	transaction.TransactionNumber = transactionNumber
	transaction.Type = models.TransactionTypeIn

	// Create journal entry
	journal := &models.Journal{
		CompanyID:       transaction.CompanyID,
		TransactionDate: transaction.TransactionDate,
		Description:     transaction.Description,
		CreatedBy:       transaction.CreatedBy,
		Status:          models.JournalStatusDraft,
		Entries: []models.JournalEntry{
			{
				AccountID:   transaction.AccountID, // Cash/Bank account (Debit)
				Description: transaction.Description,
				Debit:       transaction.Amount,
				Credit:      0,
				Position:    1,
			},
			{
				AccountID:   contraAccountID, // Contra account (Credit)
				Description: transaction.Description,
				Debit:       0,
				Credit:      transaction.Amount,
				Position:    2,
			},
		},
	}

	// Generate journal number
	journalNumber, err := s.journalRepo.GenerateJournalNumber(journal.CompanyID, journal.TransactionDate)
	if err != nil {
		return err
	}
	journal.JournalNumber = journalNumber
	journal.TotalDebit = transaction.Amount
	journal.TotalCredit = transaction.Amount

	// Create journal
	if err := s.journalRepo.Create(journal); err != nil {
		return err
	}

	// Link transaction to journal
	transaction.JournalID = &journal.ID

	// Create transaction
	return s.cashBankRepo.Create(transaction)
}

func (s *cashBankService) CreateCashOutWithJournal(transaction *models.CashBankTransaction, contraAccountID uint) error {
	// Generate transaction number
	transactionNumber, err := s.cashBankRepo.GenerateTransactionNumber(
		transaction.CompanyID,
		models.TransactionTypeOut,
		transaction.TransactionDate,
	)
	if err != nil {
		return err
	}
	transaction.TransactionNumber = transactionNumber
	transaction.Type = models.TransactionTypeOut

	// Create journal entry
	journal := &models.Journal{
		CompanyID:       transaction.CompanyID,
		TransactionDate: transaction.TransactionDate,
		Description:     transaction.Description,
		CreatedBy:       transaction.CreatedBy,
		Status:          models.JournalStatusDraft,
		Entries: []models.JournalEntry{
			{
				AccountID:   contraAccountID, // Expense/Contra account (Debit)
				Description: transaction.Description,
				Debit:       transaction.Amount,
				Credit:      0,
				Position:    1,
			},
			{
				AccountID:   transaction.AccountID, // Cash/Bank account (Credit)
				Description: transaction.Description,
				Debit:       0,
				Credit:      transaction.Amount,
				Position:    2,
			},
		},
	}

	// Generate journal number
	journalNumber, err := s.journalRepo.GenerateJournalNumber(journal.CompanyID, journal.TransactionDate)
	if err != nil {
		return err
	}
	journal.JournalNumber = journalNumber
	journal.TotalDebit = transaction.Amount
	journal.TotalCredit = transaction.Amount

	// Create journal
	if err := s.journalRepo.Create(journal); err != nil {
		return err
	}

	// Link transaction to journal
	transaction.JournalID = &journal.ID

	// Create transaction
	return s.cashBankRepo.Create(transaction)
}