package services

import (
	"errors"
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"time"
)

type JournalService interface {
	CreateJournal(journal *models.Journal) error
	GetJournalByID(id uint) (*models.Journal, error)
	GetJournalsByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Journal, error)
	GetJournalsByStatus(companyID uint, status models.JournalStatus) ([]models.Journal, error)
	UpdateJournal(id uint, journal *models.Journal) error
	DeleteJournal(id uint) error
	PostJournal(id uint, postedBy uint) error
	VoidJournal(id uint) error
}

type journalService struct {
	journalRepo repository.JournalRepository
	ledgerRepo  repository.LedgerRepository
	accountRepo repository.AccountRepository
}

func NewJournalService(
	journalRepo repository.JournalRepository,
	ledgerRepo repository.LedgerRepository,
	accountRepo repository.AccountRepository,
) JournalService {
	return &journalService{
		journalRepo: journalRepo,
		ledgerRepo:  ledgerRepo,
		accountRepo: accountRepo,
	}
}

func (s *journalService) CreateJournal(journal *models.Journal) error {
	// Validasi: entries harus ada minimal 2
	if len(journal.Entries) < 2 {
		return errors.New("journal must have at least 2 entries")
	}

	// Validasi: total debit harus sama dengan total credit
	var totalDebit, totalCredit float64
	for _, entry := range journal.Entries {
		totalDebit += entry.Debit
		totalCredit += entry.Credit

		// Validasi: entry tidak boleh debit dan credit bersamaan
		if entry.Debit > 0 && entry.Credit > 0 {
			return errors.New("entry cannot have both debit and credit")
		}

		// Validasi: entry harus punya salah satu debit atau credit
		if entry.Debit == 0 && entry.Credit == 0 {
			return errors.New("entry must have either debit or credit")
		}
	}

	// Validasi keseimbangan dengan toleransi floating point
	if totalDebit != totalCredit {
		return errors.New("total debit must equal total credit")
	}

	journal.TotalDebit = totalDebit
	journal.TotalCredit = totalCredit

	// Generate journal number
	journalNumber, err := s.journalRepo.GenerateJournalNumber(journal.CompanyID, journal.TransactionDate)
	if err != nil {
		return err
	}
	journal.JournalNumber = journalNumber

	// Set status default
	if journal.Status == "" {
		journal.Status = models.JournalStatusDraft
	}

	return s.journalRepo.Create(journal)
}

func (s *journalService) GetJournalByID(id uint) (*models.Journal, error) {
	return s.journalRepo.FindByID(id)
}

func (s *journalService) GetJournalsByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Journal, error) {
	return s.journalRepo.FindByCompanyID(companyID, startDate, endDate)
}

func (s *journalService) GetJournalsByStatus(companyID uint, status models.JournalStatus) ([]models.Journal, error) {
	return s.journalRepo.FindByStatus(companyID, status)
}

func (s *journalService) UpdateJournal(id uint, updatedJournal *models.Journal) error {
	journal, err := s.journalRepo.FindByID(id)
	if err != nil {
		return errors.New("journal not found")
	}

	// Validasi: hanya draft yang bisa diupdate
	if journal.Status != models.JournalStatusDraft {
		return errors.New("only draft journals can be updated")
	}

	// Validasi entries
	if len(updatedJournal.Entries) < 2 {
		return errors.New("journal must have at least 2 entries")
	}

	var totalDebit, totalCredit float64
	for _, entry := range updatedJournal.Entries {
		totalDebit += entry.Debit
		totalCredit += entry.Credit
	}

	if totalDebit != totalCredit {
		return errors.New("total debit must equal total credit")
	}

	updatedJournal.TotalDebit = totalDebit
	updatedJournal.TotalCredit = totalCredit
	updatedJournal.ID = id

	return s.journalRepo.Update(updatedJournal)
}

func (s *journalService) DeleteJournal(id uint) error {
	journal, err := s.journalRepo.FindByID(id)
	if err != nil {
		return errors.New("journal not found")
	}

	// Validasi: hanya draft yang bisa dihapus
	if journal.Status != models.JournalStatusDraft {
		return errors.New("only draft journals can be deleted")
	}

	return s.journalRepo.Delete(id)
}

func (s *journalService) PostJournal(id uint, postedBy uint) error {
	journal, err := s.journalRepo.FindByID(id)
	if err != nil {
		return errors.New("journal not found")
	}

	// Validasi: hanya draft yang bisa dipost
	if journal.Status != models.JournalStatusDraft {
		return errors.New("only draft journals can be posted")
	}

	// Update status journal
	now := time.Now()
	journal.Status = models.JournalStatusPosted
	journal.PostedAt = &now
	journal.PostedBy = &postedBy

	if err := s.journalRepo.Update(journal); err != nil {
		return err
	}

	// Posting ke buku besar
	for _, entry := range journal.Entries {
		// Get current account balance
		currentBalance, _ := s.ledgerRepo.GetAccountBalance(entry.AccountID, time.Now())

		// Calculate new balance based on account type
		account, err := s.accountRepo.FindByID(entry.AccountID)
		if err != nil {
			return err
		}

		var newBalance float64
		if account.Type == models.AccountTypeAsset || account.Type == models.AccountTypeExpense {
			// Debit increases asset and expense
			newBalance = currentBalance + entry.Debit - entry.Credit
		} else {
			// Credit increases liability, equity, and revenue
			newBalance = currentBalance + entry.Credit - entry.Debit
		}

		ledger := &models.Ledger{
			CompanyID:   journal.CompanyID,
			AccountID:   entry.AccountID,
			JournalID:   journal.ID,
			EntryID:     entry.ID,
			Debit:       entry.Debit,
			Credit:      entry.Credit,
			Balance:     newBalance,
			Description: entry.Description,
		}

		if err := s.ledgerRepo.Create(ledger); err != nil {
			return err
		}

		// Update account balance
		account.Balance = newBalance
		if err := s.accountRepo.Update(account); err != nil {
			return err
		}
	}

	return nil
}

func (s *journalService) VoidJournal(id uint) error {
	journal, err := s.journalRepo.FindByID(id)
	if err != nil {
		return errors.New("journal not found")
	}

	// Validasi: hanya posted yang bisa divoid
	if journal.Status != models.JournalStatusPosted {
		return errors.New("only posted journals can be voided")
	}

	// Update status
	journal.Status = models.JournalStatusVoided
	return s.journalRepo.Update(journal)
}