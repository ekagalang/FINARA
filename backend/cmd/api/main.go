package main

import (
	"finara-backend/internal/config"
	"finara-backend/internal/database"
	"finara-backend/internal/handlers"
	"finara-backend/internal/middleware"
	"finara-backend/internal/repository"
	"finara-backend/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode based on environment
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	accountRepo := repository.NewAccountRepository(db)
	journalRepo := repository.NewJournalRepository(db)
	ledgerRepo := repository.NewLedgerRepository(db)
	reportRepo := repository.NewReportRepository(db)
	cashBankRepo := repository.NewCashBankRepository(db)
	taxRepo := repository.NewTaxRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)          // NEW - FASE 5
	auditLogRepo := repository.NewAuditLogRepository(db)            // NEW - FASE 5
	backupRepo := repository.NewBackupRepository(db)                // NEW - FASE 5

	// Database config for backup service
	dbConfig := &services.DatabaseConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
	}

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)
	companyService := services.NewCompanyService(companyRepo)
	accountService := services.NewAccountService(accountRepo)
	journalService := services.NewJournalService(journalRepo, ledgerRepo, accountRepo)
	ledgerService := services.NewLedgerService(ledgerRepo)
	reportService := services.NewReportService(reportRepo)
	cashBankService := services.NewCashBankService(cashBankRepo, journalRepo, ledgerRepo, accountRepo)
	taxService := services.NewTaxService(taxRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)
	notificationService := services.NewNotificationService(notificationRepo, taxRepo, userRepo)
	inventoryService := services.NewInventoryService(inventoryRepo, journalRepo, accountRepo)  // NEW - FASE 5
	auditLogService := services.NewAuditLogService(auditLogRepo)                               // NEW - FASE 5
	exportService := services.NewExportService()                                               // NEW - FASE 5
	backupService := services.NewBackupService(backupRepo, dbConfig)                           // NEW - FASE 5

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	companyHandler := handlers.NewCompanyHandler(companyService, userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	journalHandler := handlers.NewJournalHandler(journalService)
	ledgerHandler := handlers.NewLedgerHandler(ledgerService)
	reportHandler := handlers.NewReportHandler(reportService)
	cashBankHandler := handlers.NewCashBankHandler(cashBankService)
	taxHandler := handlers.NewTaxHandler(taxService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)                                    // NEW - FASE 5
	auditLogHandler := handlers.NewAuditLogHandler(auditLogService)                                       // NEW - FASE 5
	exportHandler := handlers.NewExportHandler(exportService, journalService, ledgerService, reportService) // NEW - FASE 5
	backupHandler := handlers.NewBackupHandler(backupService)                                             // NEW - FASE 5

	// Setup Gin router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"app":    cfg.AppName,
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		protected.Use(middleware.AuditMiddleware(auditLogService)) // NEW - FASE 5
		{
			// Profile
			protected.GET("/profile", userHandler.GetProfile)
			protected.PUT("/profile", userHandler.UpdateProfile)

			// User management (admin only)
			users := protected.Group("/users")
			users.Use(middleware.RoleMiddleware("admin"))
			{
				users.GET("", userHandler.GetAllUsers)
				users.GET("/:id", userHandler.GetUserByID)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
			}

			// Company management (admin only)
			companies := protected.Group("/companies")
			companies.Use(middleware.RoleMiddleware("admin"))
			{
				companies.POST("", companyHandler.CreateCompany)
				companies.GET("", companyHandler.GetAllCompanies)
				companies.GET("/:id", companyHandler.GetCompanyByID)
				companies.PUT("/:id", companyHandler.UpdateCompany)
				companies.DELETE("/:id", companyHandler.DeleteCompany)
				companies.POST("/assign-user", companyHandler.AssignUserToCompany)
			}

			// Chart of Accounts
			accounts := protected.Group("/accounts")
			{
				accounts.POST("", accountHandler.CreateAccount)
				accounts.GET("", accountHandler.GetAllAccounts)
				accounts.GET("/active", accountHandler.GetActiveAccounts)
				accounts.GET("/code/:code", accountHandler.GetAccountByCode)
				accounts.GET("/:id", accountHandler.GetAccountByID)
				accounts.PUT("/:id", accountHandler.UpdateAccount)
				accounts.DELETE("/:id", accountHandler.DeleteAccount)
				accounts.POST("/initialize", accountHandler.InitializeDefaultAccounts)
			}

			// Journal (Jurnal Umum)
			journals := protected.Group("/journals")
			{
				journals.POST("", journalHandler.CreateJournal)
				journals.GET("", journalHandler.GetJournalsByPeriod)
				journals.GET("/status/:status", journalHandler.GetJournalsByStatus)
				journals.GET("/:id", journalHandler.GetJournalByID)
				journals.PUT("/:id", journalHandler.UpdateJournal)
				journals.DELETE("/:id", journalHandler.DeleteJournal)
				journals.POST("/:id/post", journalHandler.PostJournal)
				journals.POST("/:id/void", journalHandler.VoidJournal)
			}

			// Ledger (Buku Besar)
			ledgers := protected.Group("/ledgers")
			{
				ledgers.GET("/account/:account_id", ledgerHandler.GetLedgerByAccount)
				ledgers.GET("/company", ledgerHandler.GetLedgerByCompany)
				ledgers.GET("/trial-balance", ledgerHandler.GetTrialBalance)
				ledgers.GET("/account/:account_id/balance", ledgerHandler.GetAccountBalance)
			}

			// Financial Reports
			reports := protected.Group("/reports")
			{
				reports.GET("/income-statement", reportHandler.GetIncomeStatement)
				reports.GET("/balance-sheet", reportHandler.GetBalanceSheet)
				reports.GET("/cash-flow", reportHandler.GetCashFlow)
			}

			// Cash & Bank Management
			cashBank := protected.Group("/cash-bank")
			{
				cashBank.POST("", cashBankHandler.CreateTransaction)
				cashBank.GET("", cashBankHandler.GetTransactionsByPeriod)
				cashBank.GET("/account/:account_id", cashBankHandler.GetTransactionsByAccount)
				cashBank.GET("/position", cashBankHandler.GetCashPosition)
				cashBank.GET("/:id", cashBankHandler.GetTransactionByID)
				cashBank.PUT("/:id", cashBankHandler.UpdateTransaction)
				cashBank.DELETE("/:id", cashBankHandler.DeleteTransaction)
			}

			// Tax Management
			taxes := protected.Group("/taxes")
			{
				taxes.POST("", taxHandler.CreateTax)
				taxes.GET("/period", taxHandler.GetTaxesByPeriod)
				taxes.GET("/type/:type", taxHandler.GetTaxesByType)
				taxes.GET("/due", taxHandler.GetDueTaxes)
				taxes.GET("/summary", taxHandler.GetTaxSummary)
				taxes.GET("/:id", taxHandler.GetTaxByID)
				taxes.PUT("/:id", taxHandler.UpdateTax)
				taxes.DELETE("/:id", taxHandler.DeleteTax)
				taxes.POST("/:id/report", taxHandler.MarkAsReported)
				taxes.POST("/:id/pay", taxHandler.MarkAsPaid)
			}

			// Dashboard & Analytics
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/summary", dashboardHandler.GetDashboardSummary)
				dashboard.GET("/revenue/monthly", dashboardHandler.GetMonthlyRevenue)
				dashboard.GET("/expense/monthly", dashboardHandler.GetMonthlyExpense)
				dashboard.GET("/expense/category", dashboardHandler.GetExpenseByCategory)
				dashboard.GET("/revenue/category", dashboardHandler.GetRevenueByCategory)
				dashboard.GET("/ratios", dashboardHandler.GetFinancialRatios)
			}

			// Notifications
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", notificationHandler.GetNotifications)
				notifications.GET("/unread", notificationHandler.GetUnreadNotifications)
				notifications.GET("/count", notificationHandler.GetUnreadCount)
				notifications.POST("/:id/read", notificationHandler.MarkAsRead)
				notifications.POST("/read-all", notificationHandler.MarkAllAsRead)
				notifications.DELETE("/:id", notificationHandler.DeleteNotification)
			}

			// Inventory Management (NEW - FASE 5)
			inventory := protected.Group("/inventory")
			{
				// Products
				inventory.POST("/products", inventoryHandler.CreateProduct)
				inventory.GET("/products", inventoryHandler.GetAllProducts)
				inventory.GET("/products/:id", inventoryHandler.GetProductByID)
				inventory.PUT("/products/:id", inventoryHandler.UpdateProduct)
				inventory.DELETE("/products/:id", inventoryHandler.DeleteProduct)

				// Stock Movements
				inventory.POST("/movements", inventoryHandler.CreateStockMovement)
				inventory.GET("/movements", inventoryHandler.GetStockMovements)

				// Stock Balances
				inventory.GET("/balances", inventoryHandler.GetStockBalances)
				inventory.GET("/balances/:product_id", inventoryHandler.GetStockBalance)

				// Stock Opname
				inventory.POST("/opname", inventoryHandler.CreateStockOpname)
				inventory.GET("/opname", inventoryHandler.GetStockOpnames)
				inventory.POST("/opname/:id/approve", inventoryHandler.ApproveStockOpname)

				// HPP Calculation
				inventory.GET("/hpp", inventoryHandler.CalculateHPP)
			}

			// Audit Logs (NEW - FASE 5)
			auditLogs := protected.Group("/audit-logs")
			auditLogs.Use(middleware.RoleMiddleware("admin")) // Only admin can view audit logs
			{
				auditLogs.GET("", auditLogHandler.GetAuditLogs)
				auditLogs.GET("/user/:user_id", auditLogHandler.GetAuditLogsByUser)
				auditLogs.GET("/module/:module", auditLogHandler.GetAuditLogsByModule)
				auditLogs.GET("/date-range", auditLogHandler.GetAuditLogsByDateRange)
			}

			// Export (NEW - FASE 5)
			exports := protected.Group("/export")
			{
				exports.GET("/journals", exportHandler.ExportJournals)                   // ?format=csv/excel&start_date=&end_date=
				exports.GET("/trial-balance", exportHandler.ExportTrialBalance)          // ?format=csv/excel&end_date=
				exports.GET("/income-statement", exportHandler.ExportIncomeStatement)    // ?format=csv/excel&start_date=&end_date=
				exports.GET("/balance-sheet", exportHandler.ExportBalanceSheet)          // ?format=csv/excel&as_of_date=
			}

			// Backup & Restore (NEW - FASE 5)
			backups := protected.Group("/backups")
			backups.Use(middleware.RoleMiddleware("admin")) // Only admin can manage backups
			{
				backups.POST("", backupHandler.CreateBackup)
				backups.GET("", backupHandler.GetBackups)
				backups.GET("/:id", backupHandler.GetBackupByID)
				backups.POST("/:id/restore", backupHandler.RestoreBackup)
				backups.GET("/:id/download", backupHandler.DownloadBackup)
				backups.DELETE("/:id", backupHandler.DeleteBackup)
			}
		}
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}