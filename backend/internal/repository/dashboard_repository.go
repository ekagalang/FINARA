package repository

import (
	"finara-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetDashboardSummary(companyID uint, asOfDate time.Time) (*models.DashboardSummary, error)
	GetMonthlyRevenue(companyID uint, year int) ([]models.MonthlyRevenue, error)
	GetMonthlyExpense(companyID uint, year int) ([]models.MonthlyExpense, error)
	GetExpenseByCategory(companyID uint, startDate, endDate time.Time) ([]models.ExpenseByCategory, error)
	GetRevenueByCategory(companyID uint, startDate, endDate time.Time) ([]models.RevenueByCategory, error)
	GetFinancialRatios(companyID uint, asOfDate time.Time, startDate, endDate time.Time) (*models.FinancialRatio, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetDashboardSummary(companyID uint, asOfDate time.Time) (*models.DashboardSummary, error) {
	var summary models.DashboardSummary

	// Total Revenue (current year)
	r.db.Raw(`
		SELECT COALESCE(SUM(l.credit - l.debit), 0) as total_revenue
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ?
			AND a.type = 'revenue'
			AND YEAR(j.transaction_date) = YEAR(?)
			AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&summary.TotalRevenue)

	// Total Expense (current year)
	r.db.Raw(`
		SELECT COALESCE(SUM(l.debit - l.credit), 0) as total_expense
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ?
			AND a.type = 'expense'
			AND YEAR(j.transaction_date) = YEAR(?)
			AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&summary.TotalExpense)

	summary.NetIncome = summary.TotalRevenue - summary.TotalExpense

	// Total Assets
	r.db.Raw(`
		SELECT COALESCE(SUM(l.debit - l.credit), 0) as total_assets
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ?
			AND a.type = 'asset'
			AND j.transaction_date <= ?
			AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&summary.TotalAssets)

	// Total Liabilities
	r.db.Raw(`
		SELECT COALESCE(SUM(l.credit - l.debit), 0) as total_liabilities
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ?
			AND a.type = 'liability'
			AND j.transaction_date <= ?
			AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&summary.TotalLiabilities)

	// Total Equity
	r.db.Raw(`
		SELECT COALESCE(SUM(l.credit - l.debit), 0) as total_equity
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ?
			AND a.type = 'equity'
			AND j.transaction_date <= ?
			AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&summary.TotalEquity)

	// Cash Balance
	r.db.Raw(`
		SELECT COALESCE(SUM(l.debit - l.credit), 0) as cash_balance
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ?
			AND a.code = '1-1100'
			AND j.transaction_date <= ?
			AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&summary.CashBalance)

	// Bank Balance
	r.db.Raw(`
		SELECT COALESCE(SUM(l.debit - l.credit), 0) as bank_balance
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ?
			AND a.code = '1-1200'
			AND j.transaction_date <= ?
			AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&summary.BankBalance)

	return &summary, nil
}

func (r *dashboardRepository) GetMonthlyRevenue(companyID uint, year int) ([]models.MonthlyRevenue, error) {
	var revenues []models.MonthlyRevenue

	err := r.db.Raw(`
		SELECT 
			DATE_FORMAT(j.transaction_date, '%Y-%m') as month,
			COALESCE(SUM(l.credit - l.debit), 0) as revenue
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ?
			AND a.type = 'revenue'
			AND YEAR(j.transaction_date) = ?
			AND j.status = 'posted'
		GROUP BY DATE_FORMAT(j.transaction_date, '%Y-%m')
		ORDER BY month ASC
	`, companyID, year).Scan(&revenues).Error

	return revenues, err
}

func (r *dashboardRepository) GetMonthlyExpense(companyID uint, year int) ([]models.MonthlyExpense, error) {
	var expenses []models.MonthlyExpense

	err := r.db.Raw(`
		SELECT 
			DATE_FORMAT(j.transaction_date, '%Y-%m') as month,
			COALESCE(SUM(l.debit - l.credit), 0) as expense
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ?
			AND a.type = 'expense'
			AND YEAR(j.transaction_date) = ?
			AND j.status = 'posted'
		GROUP BY DATE_FORMAT(j.transaction_date, '%Y-%m')
		ORDER BY month ASC
	`, companyID, year).Scan(&expenses).Error

	return expenses, err
}

func (r *dashboardRepository) GetExpenseByCategory(companyID uint, startDate, endDate time.Time) ([]models.ExpenseByCategory, error) {
	var expenses []models.ExpenseByCategory

	err := r.db.Raw(`
		SELECT 
			a.category as category,
			COALESCE(SUM(l.debit - l.credit), 0) as amount
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ?
			AND a.type = 'expense'
			AND j.transaction_date BETWEEN ? AND ?
			AND j.status = 'posted'
		GROUP BY a.category
		ORDER BY amount DESC
	`, companyID, startDate, endDate).Scan(&expenses).Error

	return expenses, err
}

func (r *dashboardRepository) GetRevenueByCategory(companyID uint, startDate, endDate time.Time) ([]models.RevenueByCategory, error) {
	var revenues []models.RevenueByCategory

	err := r.db.Raw(`
		SELECT 
			a.category as category,
			COALESCE(SUM(l.credit - l.debit), 0) as amount
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ?
			AND a.type = 'revenue'
			AND j.transaction_date BETWEEN ? AND ?
			AND j.status = 'posted'
		GROUP BY a.category
		ORDER BY amount DESC
	`, companyID, startDate, endDate).Scan(&revenues).Error

	return revenues, err
}

func (r *dashboardRepository) GetFinancialRatios(companyID uint, asOfDate time.Time, startDate, endDate time.Time) (*models.FinancialRatio, error) {
	var ratios models.FinancialRatio

	// Get required values
	var currentAssets, currentLiabilities, totalAssets, totalLiabilities, totalEquity, inventory float64
	var revenue, netIncome float64

	// Current Assets
	r.db.Raw(`
		SELECT COALESCE(SUM(l.debit - l.credit), 0)
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ? AND a.type = 'asset' AND a.category = 'current_asset'
			AND j.transaction_date <= ? AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&currentAssets)

	// Current Liabilities
	r.db.Raw(`
		SELECT COALESCE(SUM(l.credit - l.debit), 0)
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ? AND a.type = 'liability' AND a.category = 'current_liability'
			AND j.transaction_date <= ? AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&currentLiabilities)

	// Inventory
	r.db.Raw(`
		SELECT COALESCE(SUM(l.debit - l.credit), 0)
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ? AND a.code = '1-1400'
			AND j.transaction_date <= ? AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&inventory)

	// Total Assets
	r.db.Raw(`
		SELECT COALESCE(SUM(l.debit - l.credit), 0)
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ? AND a.type = 'asset'
			AND j.transaction_date <= ? AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&totalAssets)

	// Total Liabilities
	r.db.Raw(`
		SELECT COALESCE(SUM(l.credit - l.debit), 0)
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ? AND a.type = 'liability'
			AND j.transaction_date <= ? AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&totalLiabilities)

	// Total Equity
	r.db.Raw(`
		SELECT COALESCE(SUM(l.credit - l.debit), 0)
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ? AND a.type = 'equity'
			AND j.transaction_date <= ? AND j.status = 'posted'
	`, companyID, asOfDate).Scan(&totalEquity)

	// Revenue
	r.db.Raw(`
		SELECT COALESCE(SUM(l.credit - l.debit), 0)
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ? AND a.type = 'revenue'
			AND j.transaction_date BETWEEN ? AND ? AND j.status = 'posted'
	`, companyID, startDate, endDate).Scan(&revenue)

	// Net Income
	var expense float64
	r.db.Raw(`
		SELECT COALESCE(SUM(l.debit - l.credit), 0)
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		JOIN accounts a ON l.account_id = a.id
		WHERE a.company_id = ? AND a.type = 'expense'
			AND j.transaction_date BETWEEN ? AND ? AND j.status = 'posted'
	`, companyID, startDate, endDate).Scan(&expense)
	netIncome = revenue - expense

	// Calculate ratios
	if currentLiabilities > 0 {
		ratios.CurrentRatio = currentAssets / currentLiabilities
		ratios.QuickRatio = (currentAssets - inventory) / currentLiabilities
	}

	if totalEquity > 0 {
		ratios.DebtToEquityRatio = totalLiabilities / totalEquity
		ratios.ReturnOnEquity = (netIncome / totalEquity) * 100
	}

	if totalAssets > 0 {
		ratios.DebtToAssetRatio = totalLiabilities / totalAssets
		ratios.ReturnOnAssets = (netIncome / totalAssets) * 100
	}

	if revenue > 0 {
		ratios.ProfitMargin = (netIncome / revenue) * 100
	}

	return &ratios, nil
}