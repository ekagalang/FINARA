package services

import (
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
	"time"
)

type DashboardService interface {
	GetDashboardSummary(companyID uint, asOfDate time.Time) (*models.DashboardSummary, error)
	GetMonthlyRevenue(companyID uint, year int) ([]models.MonthlyRevenue, error)
	GetMonthlyExpense(companyID uint, year int) ([]models.MonthlyExpense, error)
	GetExpenseByCategory(companyID uint, startDate, endDate time.Time) ([]models.ExpenseByCategory, error)
	GetRevenueByCategory(companyID uint, startDate, endDate time.Time) ([]models.RevenueByCategory, error)
	GetFinancialRatios(companyID uint, asOfDate time.Time, startDate, endDate time.Time) (*models.FinancialRatio, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardService(dashboardRepo repository.DashboardRepository) DashboardService {
	return &dashboardService{dashboardRepo: dashboardRepo}
}

func (s *dashboardService) GetDashboardSummary(companyID uint, asOfDate time.Time) (*models.DashboardSummary, error) {
	return s.dashboardRepo.GetDashboardSummary(companyID, asOfDate)
}

func (s *dashboardService) GetMonthlyRevenue(companyID uint, year int) ([]models.MonthlyRevenue, error) {
	return s.dashboardRepo.GetMonthlyRevenue(companyID, year)
}

func (s *dashboardService) GetMonthlyExpense(companyID uint, year int) ([]models.MonthlyExpense, error) {
	return s.dashboardRepo.GetMonthlyExpense(companyID, year)
}

func (s *dashboardService) GetExpenseByCategory(companyID uint, startDate, endDate time.Time) ([]models.ExpenseByCategory, error) {
	return s.dashboardRepo.GetExpenseByCategory(companyID, startDate, endDate)
}

func (s *dashboardService) GetRevenueByCategory(companyID uint, startDate, endDate time.Time) ([]models.RevenueByCategory, error) {
	return s.dashboardRepo.GetRevenueByCategory(companyID, startDate, endDate)
}

func (s *dashboardService) GetFinancialRatios(companyID uint, asOfDate time.Time, startDate, endDate time.Time) (*models.FinancialRatio, error) {
	return s.dashboardRepo.GetFinancialRatios(companyID, asOfDate, startDate, endDate)
}