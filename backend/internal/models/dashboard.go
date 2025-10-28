package models

type DashboardSummary struct {
	TotalRevenue     float64 `json:"total_revenue"`
	TotalExpense     float64 `json:"total_expense"`
	NetIncome        float64 `json:"net_income"`
	TotalAssets      float64 `json:"total_assets"`
	TotalLiabilities float64 `json:"total_liabilities"`
	TotalEquity      float64 `json:"total_equity"`
	CashBalance      float64 `json:"cash_balance"`
	BankBalance      float64 `json:"bank_balance"`
}

type MonthlyRevenue struct {
	Month   string  `json:"month"`
	Revenue float64 `json:"revenue"`
}

type MonthlyExpense struct {
	Month   string  `json:"month"`
	Expense float64 `json:"expense"`
}

type ExpenseByCategory struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

type RevenueByCategory struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

type FinancialRatio struct {
	CurrentRatio      float64 `json:"current_ratio"`        // Aset Lancar / Liabilitas Lancar
	QuickRatio        float64 `json:"quick_ratio"`          // (Aset Lancar - Persediaan) / Liabilitas Lancar
	DebtToEquityRatio float64 `json:"debt_to_equity_ratio"` // Total Liabilitas / Total Ekuitas
	DebtToAssetRatio  float64 `json:"debt_to_asset_ratio"`  // Total Liabilitas / Total Aset
	ProfitMargin      float64 `json:"profit_margin"`        // (Laba Bersih / Pendapatan) * 100
	ReturnOnAssets    float64 `json:"return_on_assets"`     // (Laba Bersih / Total Aset) * 100
	ReturnOnEquity    float64 `json:"return_on_equity"`     // (Laba Bersih / Total Ekuitas) * 100
}