package repository

import (
	"finara-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type LedgerRepository interface {
	Create(ledger *models.Ledger) error
	FindByAccountID(accountID uint, startDate, endDate time.Time) ([]models.Ledger, error)
	FindByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Ledger, error)
	GetTrialBalance(companyID uint, endDate time.Time) ([]models.TrialBalanceResponse, error)
	GetAccountBalance(accountID uint, endDate time.Time) (float64, error)
}

type ledgerRepository struct {
	db *gorm.DB
}

func NewLedgerRepository(db *gorm.DB) LedgerRepository {
	return &ledgerRepository{db: db}
}

func (r *ledgerRepository) Create(ledger *models.Ledger) error {
	return r.db.Create(ledger).Error
}

func (r *ledgerRepository) FindByAccountID(accountID uint, startDate, endDate time.Time) ([]models.Ledger, error) {
	var ledgers []models.Ledger
	err := r.db.Joins("JOIN journals ON journals.id = ledgers.journal_id").
		Where("ledgers.account_id = ? AND journals.transaction_date BETWEEN ? AND ?", accountID, startDate, endDate).
		Order("journals.transaction_date ASC, ledgers.id ASC").
		Preload("Account").
		Preload("Journal").
		Find(&ledgers).Error
	return ledgers, err
}

func (r *ledgerRepository) FindByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.Ledger, error) {
	var ledgers []models.Ledger
	err := r.db.Joins("JOIN journals ON journals.id = ledgers.journal_id").
		Where("ledgers.company_id = ? AND journals.transaction_date BETWEEN ? AND ?", companyID, startDate, endDate).
		Order("journals.transaction_date ASC").
		Preload("Account").
		Preload("Journal").
		Find(&ledgers).Error
	return ledgers, err
}

func (r *ledgerRepository) GetTrialBalance(companyID uint, endDate time.Time) ([]models.TrialBalanceResponse, error) {
	var results []models.TrialBalanceResponse
	
	err := r.db.Raw(`
		SELECT 
			a.code as account_code,
			a.name as account_name,
			a.type as account_type,
			COALESCE(SUM(l.debit), 0) as debit,
			COALESCE(SUM(l.credit), 0) as credit,
			CASE 
				WHEN a.type IN ('asset', 'expense') THEN COALESCE(SUM(l.debit - l.credit), 0)
				ELSE 0
			END as debit_balance,
			CASE 
				WHEN a.type IN ('liability', 'equity', 'revenue') THEN COALESCE(SUM(l.credit - l.debit), 0)
				ELSE 0
			END as credit_balance
		FROM accounts a
		LEFT JOIN ledgers l ON a.id = l.account_id
		LEFT JOIN journals j ON l.journal_id = j.id
		WHERE a.company_id = ? 
			AND a.is_active = true 
			AND a.is_header = false
			AND (j.transaction_date IS NULL OR j.transaction_date <= ?)
			AND (j.status IS NULL OR j.status = 'posted')
		GROUP BY a.id, a.code, a.name, a.type
		ORDER BY a.code ASC
	`, companyID, endDate).Scan(&results).Error
	
	return results, err
}

func (r *ledgerRepository) GetAccountBalance(accountID uint, endDate time.Time) (float64, error) {
	var result struct {
		Balance float64
	}
	
	err := r.db.Raw(`
		SELECT COALESCE(SUM(debit - credit), 0) as balance
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		WHERE l.account_id = ? 
			AND j.transaction_date <= ?
			AND j.status = 'posted'
	`, accountID, endDate).Scan(&result).Error
	
	return result.Balance, err
}