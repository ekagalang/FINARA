package models

type IncomeStatementItem struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Amount      float64 `json:"amount"`
}

type IncomeStatementResponse struct {
	Period       string                `json:"period"`
	StartDate    string                `json:"start_date"`
	EndDate      string                `json:"end_date"`
	Revenues     []IncomeStatementItem `json:"revenues"`
	TotalRevenue float64               `json:"total_revenue"`
	Expenses     []IncomeStatementItem `json:"expenses"`
	TotalExpense float64               `json:"total_expense"`
	NetIncome    float64               `json:"net_income"`
	IsProfit     bool                  `json:"is_profit"`
}

type BalanceSheetItem struct {
	AccountCode string  `json:"account_code"`
	AccountName string  `json:"account_name"`
	Amount      float64 `json:"amount"`
}

type BalanceSheetResponse struct {
	AsOfDate                  string             `json:"as_of_date"`
	CurrentAssets             []BalanceSheetItem `json:"current_assets"`
	TotalCurrentAssets        float64            `json:"total_current_assets"`
	FixedAssets               []BalanceSheetItem `json:"fixed_assets"`
	TotalFixedAssets          float64            `json:"total_fixed_assets"`
	TotalAssets               float64            `json:"total_assets"`
	CurrentLiabilities        []BalanceSheetItem `json:"current_liabilities"`
	TotalCurrentLiabilities   float64            `json:"total_current_liabilities"`
	LongTermLiabilities       []BalanceSheetItem `json:"long_term_liabilities"`
	TotalLongTermLiabilities  float64            `json:"total_long_term_liabilities"`
	TotalLiabilities          float64            `json:"total_liabilities"`
	Equity                    []BalanceSheetItem `json:"equity"`
	TotalEquity               float64            `json:"total_equity"`
	TotalLiabilitiesAndEquity float64            `json:"total_liabilities_and_equity"`
	IsBalanced                bool               `json:"is_balanced"`
}

type CashFlowItem struct {
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
}

type CashFlowResponse struct {
	Period               string         `json:"period"`
	StartDate            string         `json:"start_date"`
	EndDate              string         `json:"end_date"`
	OperatingActivities  []CashFlowItem `json:"operating_activities"`
	NetCashFromOperating float64        `json:"net_cash_from_operating"`
	InvestingActivities  []CashFlowItem `json:"investing_activities"`
	NetCashFromInvesting float64        `json:"net_cash_from_investing"`
	FinancingActivities  []CashFlowItem `json:"financing_activities"`
	NetCashFromFinancing float64        `json:"net_cash_from_financing"`
	NetIncreaseInCash    float64        `json:"net_increase_in_cash"`
	CashAtBeginning      float64        `json:"cash_at_beginning"`
	CashAtEnd            float64        `json:"cash_at_end"`
}