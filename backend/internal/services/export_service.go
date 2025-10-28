package services

import (
	"encoding/csv"
	"errors"
	"finara-backend/internal/models"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

type ExportService interface {
	ExportJournalsToCSV(journals []models.Journal, filename string) (string, error)
	ExportJournalsToExcel(journals []models.Journal, filename string) (string, error)
	ExportTrialBalanceToCSV(trialBalance *models.TrialBalanceReport, filename string) (string, error)
	ExportTrialBalanceToExcel(trialBalance *models.TrialBalanceReport, filename string) (string, error)
	ExportIncomeStatementToCSV(incomeStatement *models.IncomeStatementResponse, filename string) (string, error)
	ExportIncomeStatementToExcel(incomeStatement *models.IncomeStatementResponse, filename string) (string, error)
	ExportBalanceSheetToCSV(balanceSheet *models.BalanceSheetResponse, filename string) (string, error)
	ExportBalanceSheetToExcel(balanceSheet *models.BalanceSheetResponse, filename string) (string, error)
}

type exportService struct{}

func NewExportService() ExportService {
	return &exportService{}
}

// Journals Export
func (s *exportService) ExportJournalsToCSV(journals []models.Journal, filename string) (string, error) {
	filepath := fmt.Sprintf("exports/%s", filename)
	
	// Create exports directory if not exists
	os.MkdirAll("exports", 0755)

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Journal Number", "Date", "Description", "Status", "Total Debit", "Total Credit"}
	writer.Write(header)

	// Write data
	for _, journal := range journals {
		row := []string{
			journal.JournalNumber,
			journal.TransactionDate.Format("2006-01-02"),
			journal.Description,
			string(journal.Status),
			fmt.Sprintf("%.2f", journal.TotalDebit),
			fmt.Sprintf("%.2f", journal.TotalCredit),
		}
		writer.Write(row)
	}

	return filepath, nil
}

func (s *exportService) ExportJournalsToExcel(journals []models.Journal, filename string) (string, error) {
	filepath := fmt.Sprintf("exports/%s", filename)
	
	// Create exports directory if not exists
	os.MkdirAll("exports", 0755)

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Journals"
	index, _ := f.NewSheet(sheetName)

	// Set headers
	headers := []string{"Journal Number", "Date", "Description", "Status", "Total Debit", "Total Credit"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, journal := range journals {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), journal.JournalNumber)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), journal.TransactionDate.Format("2006-01-02"))
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), journal.Description)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), journal.Status)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), journal.TotalDebit)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), journal.TotalCredit)
	}

	f.SetActiveSheet(index)

	if err := f.SaveAs(filepath); err != nil {
		return "", err
	}

	return filepath, nil
}

// Trial Balance Export
func (s *exportService) ExportTrialBalanceToCSV(trialBalance *models.TrialBalanceReport, filename string) (string, error) {
	filepath := fmt.Sprintf("exports/%s", filename)
	
	os.MkdirAll("exports", 0755)

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Account Code", "Account Name", "Debit", "Credit", "Debit Balance", "Credit Balance"}
	writer.Write(header)

	// Write data
	for _, account := range trialBalance.Accounts {
		row := []string{
			account.AccountCode,
			account.AccountName,
			fmt.Sprintf("%.2f", account.Debit),
			fmt.Sprintf("%.2f", account.Credit),
			fmt.Sprintf("%.2f", account.DebitBalance),
			fmt.Sprintf("%.2f", account.CreditBalance),
		}
		writer.Write(row)
	}

	// Write totals
	writer.Write([]string{})
	writer.Write([]string{"TOTAL", "", 
		fmt.Sprintf("%.2f", trialBalance.TotalDebit),
		fmt.Sprintf("%.2f", trialBalance.TotalCredit),
		fmt.Sprintf("%.2f", trialBalance.TotalDebitBalance),
		fmt.Sprintf("%.2f", trialBalance.TotalCreditBalance),
	})

	return filepath, nil
}

func (s *exportService) ExportTrialBalanceToExcel(trialBalance *models.TrialBalanceReport, filename string) (string, error) {
	filepath := fmt.Sprintf("exports/%s", filename)
	
	os.MkdirAll("exports", 0755)

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Trial Balance"
	index, _ := f.NewSheet(sheetName)

	// Set headers
	headers := []string{"Account Code", "Account Name", "Debit", "Credit", "Debit Balance", "Credit Balance"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Write data
	for i, account := range trialBalance.Accounts {
		row := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), account.AccountCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), account.AccountName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), account.Debit)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), account.Credit)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), account.DebitBalance)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), account.CreditBalance)
	}

	// Write totals
	totalRow := len(trialBalance.Accounts) + 3
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", totalRow), "TOTAL")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", totalRow), trialBalance.TotalDebit)
	f.SetCellValue(sheetName, fmt.Sprintf("D%d", totalRow), trialBalance.TotalCredit)
	f.SetCellValue(sheetName, fmt.Sprintf("E%d", totalRow), trialBalance.TotalDebitBalance)
	f.SetCellValue(sheetName, fmt.Sprintf("F%d", totalRow), trialBalance.TotalCreditBalance)

	f.SetActiveSheet(index)

	if err := f.SaveAs(filepath); err != nil {
		return "", err
	}

	return filepath, nil
}

// Income Statement Export
func (s *exportService) ExportIncomeStatementToCSV(incomeStatement *models.IncomeStatementResponse, filename string) (string, error) {
	filepath := fmt.Sprintf("exports/%s", filename)
	
	os.MkdirAll("exports", 0755)

	file, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Title
	writer.Write([]string{"INCOME STATEMENT"})
	writer.Write([]string{fmt.Sprintf("Period: %s to %s", incomeStatement.StartDate, incomeStatement.EndDate)})
	writer.Write([]string{})

	// Revenues
	writer.Write([]string{"REVENUES"})
	for _, rev := range incomeStatement.Revenues {
		writer.Write([]string{rev.AccountCode, rev.AccountName, fmt.Sprintf("%.2f", rev.Amount)})
	}
	writer.Write([]string{"", "Total Revenue", fmt.Sprintf("%.2f", incomeStatement.TotalRevenue)})
	writer.Write([]string{})

	// Expenses
	writer.Write([]string{"EXPENSES"})
	for _, exp := range incomeStatement.Expenses {
		writer.Write([]string{exp.AccountCode, exp.AccountName, fmt.Sprintf("%.2f", exp.Amount)})
	}
	writer.Write([]string{"", "Total Expense", fmt.Sprintf("%.2f", incomeStatement.TotalExpense)})
	writer.Write([]string{})

	// Net Income
	writer.Write([]string{"", "NET INCOME", fmt.Sprintf("%.2f", incomeStatement.NetIncome)})

	return filepath, nil
}

func (s *exportService) ExportIncomeStatementToExcel(incomeStatement *models.IncomeStatementResponse, filename string) (string, error) {
	filepath := fmt.Sprintf("exports/%s", filename)
	
	os.MkdirAll("exports", 0755)

	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Income Statement"
	index, _ := f.NewSheet(sheetName)

	row := 1

	// Title
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "INCOME STATEMENT")
	row++
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("Period: %s to %s", incomeStatement.StartDate, incomeStatement.EndDate))
	row += 2

	// Revenues
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "REVENUES")
	row++
	for _, rev := range incomeStatement.Revenues {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), rev.AccountCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), rev.AccountName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), rev.Amount)
		row++
	}
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Total Revenue")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), incomeStatement.TotalRevenue)
	row += 2

	// Expenses
	f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), "EXPENSES")
	row++
	for _, exp := range incomeStatement.Expenses {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), exp.AccountCode)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), exp.AccountName)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), exp.Amount)
		row++
	}
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "Total Expense")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), incomeStatement.TotalExpense)
	row += 2

	// Net Income
	f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), "NET INCOME")
	f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), incomeStatement.NetIncome)

	f.SetActiveSheet(index)

	if err := f.SaveAs(filepath); err != nil {
		return "", err
	}

	return filepath, nil
}

// Balance Sheet Export
func (s *exportService) ExportBalanceSheetToCSV(balanceSheet *models.BalanceSheetResponse, filename string) (string, error) {
	return "", errors.New("not implemented yet")
}

func (s *exportService) ExportBalanceSheetToExcel(balanceSheet *models.BalanceSheetResponse, filename string) (string, error) {
	return "", errors.New("not implemented yet")
}