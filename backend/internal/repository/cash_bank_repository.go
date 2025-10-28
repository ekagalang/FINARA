package repository

import (
	"finara-backend/internal/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CashBankRepository interface {
	Create(transaction *models.CashBankTransaction) error
	FindByID(id uint) (*models.CashBankTransaction, error)
	FindByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.CashBankTransaction, error)
	FindByAccountID(accountID uint, startDate, endDate time.Time) ([]models.CashBankTransaction, error)
	Update(transaction *models.CashBankTransaction) error
	Delete(id uint) error
	GenerateTransactionNumber(companyID uint, transactionType models.TransactionType, date time.Time) (string, error)
	GetCashPosition(companyID uint, endDate time.Time) (map[string]float64, error)
}

type cashBankRepository struct {
	db *gorm.DB
}

func NewCashBankRepository(db *gorm.DB) CashBankRepository {
	return &cashBankRepository{db: db}
}

func (r *cashBankRepository) Create(transaction *models.CashBankTransaction) error {
	return r.db.Create(transaction).Error
}

func (r *cashBankRepository) FindByID(id uint) (*models.CashBankTransaction, error) {
	var transaction models.CashBankTransaction
	err := r.db.Preload("Account").
		Preload("Journal").
		Preload("User").
		First(&transaction, id).Error
	return &transaction, err
}

func (r *cashBankRepository) FindByCompanyID(companyID uint, startDate, endDate time.Time) ([]models.CashBankTransaction, error) {
	var transactions []models.CashBankTransaction
	err := r.db.Where("company_id = ? AND transaction_date BETWEEN ? AND ?", companyID, startDate, endDate).
		Order("transaction_date DESC, id DESC").
		Preload("Account").
		Preload("User").
		Find(&transactions).Error
	return transactions, err
}

func (r *cashBankRepository) FindByAccountID(accountID uint, startDate, endDate time.Time) ([]models.CashBankTransaction, error) {
	var transactions []models.CashBankTransaction
	err := r.db.Where("account_id = ? AND transaction_date BETWEEN ? AND ?", accountID, startDate, endDate).
		Order("transaction_date DESC").
		Preload("Account").
		Preload("User").
		Find(&transactions).Error
	return transactions, err
}

func (r *cashBankRepository) Update(transaction *models.CashBankTransaction) error {
	return r.db.Save(transaction).Error
}

func (r *cashBankRepository) Delete(id uint) error {
	return r.db.Delete(&models.CashBankTransaction{}, id).Error
}

func (r *cashBankRepository) GenerateTransactionNumber(companyID uint, transactionType models.TransactionType, date time.Time) (string, error) {
	var count int64
	var prefix string

	switch transactionType {
	case models.TransactionTypeIn:
		prefix = "CI/" // Cash In
	case models.TransactionTypeOut:
		prefix = "CO/" // Cash Out
	case models.TransactionTypeTransfer:
		prefix = "TRF/" // Transfer
	default:
		prefix = "TRX/"
	}

	prefix += date.Format("200601/")

	err := r.db.Model(&models.CashBankTransaction{}).
		Where("company_id = ? AND transaction_number LIKE ?", companyID, prefix+"%").
		Count(&count).Error

	if err != nil {
		return "", err
	}

	return prefix + fmt.Sprintf("%04d", count+1), nil
}

func (r *cashBankRepository) GetCashPosition(companyID uint, endDate time.Time) (map[string]float64, error) {
	type CashPosition struct {
		AccountCode string  `json:"account_code"`
		AccountName string  `json:"account_name"`
		Balance     float64 `json:"balance"`
	}

	var positions []CashPosition
	err := r.db.Raw(`
		SELECT 
			a.code as account_code,
			a.name as account_name,
			COALESCE(SUM(l.debit - l.credit), 0) as balance
		FROM accounts a
		LEFT JOIN ledgers l ON a.id = l.account_id
		LEFT JOIN journals j ON l.journal_id = j.id
		WHERE a.company_id = ?
			AND a.code IN ('1-1100', '1-1200')
			AND a.is_active = true
			AND (j.transaction_date IS NULL OR j.transaction_date <= ?)
			AND (j.status IS NULL OR j.status = 'posted')
		GROUP BY a.id, a.code, a.name
		ORDER BY a.code ASC
	`, companyID, endDate).Scan(&positions).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]float64)
	var totalCash float64

	for _, pos := range positions {
		result[pos.AccountName] = pos.Balance
		totalCash += pos.Balance
	}

	result["Total"] = totalCash

	return result, nil
}