package services

import (
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"time"
)

type ReportService interface {
	GetIncomeStatement(companyID uint, startDate, endDate time.Time) (*models.IncomeStatementResponse, error)
	GetBalanceSheet(companyID uint, asOfDate time.Time) (*models.BalanceSheetResponse, error)
	GetCashFlow(companyID uint, startDate, endDate time.Time) (*models.CashFlowResponse, error)
}

type reportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(reportRepo repository.ReportRepository) ReportService {
	return &reportService{reportRepo: reportRepo}
}

func (s *reportService) GetIncomeStatement(companyID uint, startDate, endDate time.Time) (*models.IncomeStatementResponse, error) {
	return s.reportRepo.GetIncomeStatement(companyID, startDate, endDate)
}

func (s *reportService) GetBalanceSheet(companyID uint, asOfDate time.Time) (*models.BalanceSheetResponse, error) {
	return s.reportRepo.GetBalanceSheet(companyID, asOfDate)
}

func (s *reportService) GetCashFlow(companyID uint, startDate, endDate time.Time) (*models.CashFlowResponse, error) {
	return s.reportRepo.GetCashFlow(companyID, startDate, endDate)
}