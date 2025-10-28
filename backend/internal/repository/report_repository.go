package repository

import (
	"finara-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

type ReportRepository interface {
	GetIncomeStatement(companyID uint, startDate, endDate time.Time) (*models.IncomeStatementResponse, error)
	GetBalanceSheet(companyID uint, asOfDate time.Time) (*models.BalanceSheetResponse, error)
	GetCashFlow(companyID uint, startDate, endDate time.Time) (*models.CashFlowResponse, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) GetIncomeStatement(companyID uint, startDate, endDate time.Time) (*models.IncomeStatementResponse, error) {
	var revenues []models.IncomeStatementItem
	var expenses []models.IncomeStatementItem

	// Get Revenue accounts
	err := r.db.Raw(`
		SELECT 
			a.code as account_code,
			a.name as account_name,
			COALESCE(SUM(l.credit - l.debit), 0) as amount
		FROM accounts a
		LEFT JOIN ledgers l ON a.id = l.account_id
		LEFT JOIN journals j ON l.journal_id = j.id
		WHERE a.company_id = ? 
			AND a.type = 'revenue'
			AND a.is_active = true
			AND a.is_header = false
			AND j.transaction_date BETWEEN ? AND ?
			AND j.status = 'posted'
		GROUP BY a.id, a.code, a.name
		HAVING amount != 0
		ORDER BY a.code ASC
	`, companyID, startDate, endDate).Scan(&revenues).Error

	if err != nil {
		return nil, err
	}

	// Get Expense accounts
	err = r.db.Raw(`
		SELECT 
			a.code as account_code,
			a.name as account_name,
			COALESCE(SUM(l.debit - l.credit), 0) as amount
		FROM accounts a
		LEFT JOIN ledgers l ON a.id = l.account_id
		LEFT JOIN journals j ON l.journal_id = j.id
		WHERE a.company_id = ? 
			AND a.type = 'expense'
			AND a.is_active = true
			AND a.is_header = false
			AND j.transaction_date BETWEEN ? AND ?
			AND j.status = 'posted'
		GROUP BY a.id, a.code, a.name
		HAVING amount != 0
		ORDER BY a.code ASC
	`, companyID, startDate, endDate).Scan(&expenses).Error

	if err != nil {
		return nil, err
	}

	// Calculate totals
	var totalRevenue, totalExpense float64
	for _, rev := range revenues {
		totalRevenue += rev.Amount
	}
	for _, exp := range expenses {
		totalExpense += exp.Amount
	}

	netIncome := totalRevenue - totalExpense

	return &models.IncomeStatementResponse{
		Period:        startDate.Format("2006-01") + " to " + endDate.Format("2006-01"),
		StartDate:     startDate.Format("2006-01-02"),
		EndDate:       endDate.Format("2006-01-02"),
		Revenues:      revenues,
		TotalRevenue:  totalRevenue,
		Expenses:      expenses,
		TotalExpense:  totalExpense,
		NetIncome:     netIncome,
		IsProfit:      netIncome > 0,
	}, nil
}

func (r *reportRepository) GetBalanceSheet(companyID uint, asOfDate time.Time) (*models.BalanceSheetResponse, error) {
	var currentAssets []models.BalanceSheetItem
	var fixedAssets []models.BalanceSheetItem
	var currentLiabilities []models.BalanceSheetItem
	var longTermLiabilities []models.BalanceSheetItem
	var equity []models.BalanceSheetItem

	// Current Assets
	err := r.db.Raw(`
		SELECT 
			a.code as account_code,
			a.name as account_name,
			COALESCE(SUM(l.debit - l.credit), 0) as amount
		FROM accounts a
		LEFT JOIN ledgers l ON a.id = l.account_id
		LEFT JOIN journals j ON l.journal_id = j.id
		WHERE a.company_id = ? 
			AND a.type = 'asset'
			AND a.category = 'current_asset'
			AND a.is_active = true
			AND a.is_header = false
			AND (j.transaction_date IS NULL OR j.transaction_date <= ?)
			AND (j.status IS NULL OR j.status = 'posted')
		GROUP BY a.id, a.code, a.name
		ORDER BY a.code ASC
	`, companyID, asOfDate).Scan(&currentAssets).Error

	if err != nil {
		return nil, err
	}

	// Fixed Assets
	err = r.db.Raw(`
		SELECT 
			a.code as account_code,
			a.name as account_name,
			COALESCE(SUM(l.debit - l.credit), 0) as amount
		FROM accounts a
		LEFT JOIN ledgers l ON a.id = l.account_id
		LEFT JOIN journals j ON l.journal_id = j.id
		WHERE a.company_id = ? 
			AND a.type = 'asset'
			AND a.category = 'fixed_asset'
			AND a.is_active = true
			AND a.is_header = false
			AND (j.transaction_date IS NULL OR j.transaction_date <= ?)
			AND (j.status IS NULL OR j.status = 'posted')
		GROUP BY a.id, a.code, a.name
		ORDER BY a.code ASC
	`, companyID, asOfDate).Scan(&fixedAssets).Error

	if err != nil {
		return nil, err
	}

	// Current Liabilities
	err = r.db.Raw(`
		SELECT 
			a.code as account_code,
			a.name as account_name,
			COALESCE(SUM(l.credit - l.debit), 0) as amount
		FROM accounts a
		LEFT JOIN ledgers l ON a.id = l.account_id
		LEFT JOIN journals j ON l.journal_id = j.id
		WHERE a.company_id = ? 
			AND a.type = 'liability'
			AND a.category = 'current_liability'
			AND a.is_active = true
			AND a.is_header = false
			AND (j.transaction_date IS NULL OR j.transaction_date <= ?)
			AND (j.status IS NULL OR j.status = 'posted')
		GROUP BY a.id, a.code, a.name
		ORDER BY a.code ASC
	`, companyID, asOfDate).Scan(&currentLiabilities).Error

	if err != nil {
		return nil, err
	}

	// Long-term Liabilities
	err = r.db.Raw(`
		SELECT 
			a.code as account_code,
			a.name as account_name,
			COALESCE(SUM(l.credit - l.debit), 0) as amount
		FROM accounts a
		LEFT JOIN ledgers l ON a.id = l.account_id
		LEFT JOIN journals j ON l.journal_id = j.id
		WHERE a.company_id = ? 
			AND a.type = 'liability'
			AND a.category = 'long_term_liability'
			AND a.is_active = true
			AND a.is_header = false
			AND (j.transaction_date IS NULL OR j.transaction_date <= ?)
			AND (j.status IS NULL OR j.status = 'posted')
		GROUP BY a.id, a.code, a.name
		ORDER BY a.code ASC
	`, companyID, asOfDate).Scan(&longTermLiabilities).Error

	if err != nil {
		return nil, err
	}

	// Equity
	err = r.db.Raw(`
		SELECT 
			a.code as account_code,
			a.name as account_name,
			COALESCE(SUM(l.credit - l.debit), 0) as amount
		FROM accounts a
		LEFT JOIN ledgers l ON a.id = l.account_id
		LEFT JOIN journals j ON l.journal_id = j.id
		WHERE a.company_id = ? 
			AND a.type = 'equity'
			AND a.is_active = true
			AND a.is_header = false
			AND (j.transaction_date IS NULL OR j.transaction_date <= ?)
			AND (j.status IS NULL OR j.status = 'posted')
		GROUP BY a.id, a.code, a.name
		ORDER BY a.code ASC
	`, companyID, asOfDate).Scan(&equity).Error

	if err != nil {
		return nil, err
	}

	// Calculate totals
	var totalCurrentAssets, totalFixedAssets float64
	var totalCurrentLiabilities, totalLongTermLiabilities float64
	var totalEquity float64

	for _, item := range currentAssets {
		totalCurrentAssets += item.Amount
	}
	for _, item := range fixedAssets {
		totalFixedAssets += item.Amount
	}
	for _, item := range currentLiabilities {
		totalCurrentLiabilities += item.Amount
	}
	for _, item := range longTermLiabilities {
		totalLongTermLiabilities += item.Amount
	}
	for _, item := range equity {
		totalEquity += item.Amount
	}

	totalAssets := totalCurrentAssets + totalFixedAssets
	totalLiabilities := totalCurrentLiabilities + totalLongTermLiabilities
	totalLiabilitiesAndEquity := totalLiabilities + totalEquity

	return &models.BalanceSheetResponse{
		AsOfDate:                 asOfDate.Format("2006-01-02"),
		CurrentAssets:            currentAssets,
		TotalCurrentAssets:       totalCurrentAssets,
		FixedAssets:              fixedAssets,
		TotalFixedAssets:         totalFixedAssets,
		TotalAssets:              totalAssets,
		CurrentLiabilities:       currentLiabilities,
		TotalCurrentLiabilities:  totalCurrentLiabilities,
		LongTermLiabilities:      longTermLiabilities,
		TotalLongTermLiabilities: totalLongTermLiabilities,
		TotalLiabilities:         totalLiabilities,
		Equity:                   equity,
		TotalEquity:              totalEquity,
		TotalLiabilitiesAndEquity: totalLiabilitiesAndEquity,
		IsBalanced:               totalAssets == totalLiabilitiesAndEquity,
	}, nil
}

func (r *reportRepository) GetCashFlow(companyID uint, startDate, endDate time.Time) (*models.CashFlowResponse, error) {
	// Simplified cash flow - tracking cash and bank accounts only
	var cashAccounts []uint
	r.db.Raw(`
		SELECT id FROM accounts 
		WHERE company_id = ? 
		AND code IN ('1-1100', '1-1200')
		AND is_active = true
	`, companyID).Scan(&cashAccounts)

	if len(cashAccounts) == 0 {
		return nil, nil
	}

	// Get beginning balance
	var beginningBalance float64
	r.db.Raw(`
		SELECT COALESCE(SUM(debit - credit), 0) as balance
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		WHERE l.company_id = ?
		AND l.account_id IN ?
		AND j.transaction_date < ?
		AND j.status = 'posted'
	`, companyID, cashAccounts, startDate).Scan(&beginningBalance)

	// Get ending balance
	var endingBalance float64
	r.db.Raw(`
		SELECT COALESCE(SUM(debit - credit), 0) as balance
		FROM ledgers l
		JOIN journals j ON l.journal_id = j.id
		WHERE l.company_id = ?
		AND l.account_id IN ?
		AND j.transaction_date <= ?
		AND j.status = 'posted'
	`, companyID, cashAccounts, endDate).Scan(&endingBalance)

	netIncrease := endingBalance - beginningBalance

	// Simplified: Operating activities = net increase (for now)
	operatingActivities := []models.CashFlowItem{
		{
			Description: "Net cash from operations",
			Amount:      netIncrease,
		},
	}

	return &models.CashFlowResponse{
		Period:                startDate.Format("2006-01") + " to " + endDate.Format("2006-01"),
		StartDate:             startDate.Format("2006-01-02"),
		EndDate:               endDate.Format("2006-01-02"),
		OperatingActivities:   operatingActivities,
		NetCashFromOperating:  netIncrease,
		InvestingActivities:   []models.CashFlowItem{},
		NetCashFromInvesting:  0,
		FinancingActivities:   []models.CashFlowItem{},
		NetCashFromFinancing:  0,
		NetIncreaseInCash:     netIncrease,
		CashAtBeginning:       beginningBalance,
		CashAtEnd:             endingBalance,
	}, nil
}