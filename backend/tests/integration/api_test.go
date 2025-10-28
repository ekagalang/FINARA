package integration

import (
	"bytes"
	"encoding/json"
	"finara-backend/internal/config"
	"finara-backend/internal/database"
	"finara-backend/internal/handlers"
	"finara-backend/internal/middleware"
	"finara-backend/internal/repository"
	"finara-backend/internal/services"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Load test environment variables
	if err := godotenv.Load("../../.env.test"); err != nil {
		// If .env.test not found, try .env
		if err := godotenv.Load("../../.env"); err != nil {
			log.Println("No .env file found, using default values")
		}
	}

	// Override config for testing
	cfg := &config.Config{
		AppName:       os.Getenv("APP_NAME"),
		AppEnv:        "test",
		AppPort:       "8081",
		DBDriver:      os.Getenv("DB_DRIVER"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		JWTSecret:     "test_jwt_secret",
		JWTExpiration: 24 * time.Hour, // Set expiration to 24 hours for testing
	}

	// Set default values if not set
	if cfg.DBDriver == "" {
		cfg.DBDriver = "mysql"
	}
	if cfg.DBHost == "" {
		cfg.DBHost = "localhost"
	}
	if cfg.DBPort == "" {
		cfg.DBPort = "3306"
	}
	if cfg.DBUser == "" {
		cfg.DBUser = "root"
	}
	if cfg.DBPassword == "" {
		cfg.DBPassword = "admin"
	}
	if cfg.DBName == "" {
		cfg.DBName = "finara_test"
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

	// Clean database before tests (optional)
	cleanupDatabase(db)

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
	inventoryRepo := repository.NewInventoryRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)
	backupRepo := repository.NewBackupRepository(db)

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
	inventoryService := services.NewInventoryService(inventoryRepo, journalRepo, accountRepo)
	auditLogService := services.NewAuditLogService(auditLogRepo)
	exportService := services.NewExportService()
	backupService := services.NewBackupService(backupRepo, dbConfig)

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
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	auditLogHandler := handlers.NewAuditLogHandler(auditLogService)
	exportHandler := handlers.NewExportHandler(exportService, journalService, ledgerService, reportService)
	backupHandler := handlers.NewBackupHandler(backupService)

	// Setup router
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
		{
			// Profile
			protected.GET("/profile", userHandler.GetProfile)
			protected.PUT("/profile", userHandler.UpdateProfile)

			// User management
			users := protected.Group("/users")
			users.Use(middleware.RoleMiddleware("admin"))
			{
				users.GET("", userHandler.GetAllUsers)
				users.GET("/:id", userHandler.GetUserByID)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
			}

			// Company management
			companies := protected.Group("/companies")
			companies.Use(middleware.RoleMiddleware("admin"))
			{
				companies.POST("", companyHandler.CreateCompany)
				companies.GET("", companyHandler.GetAllCompanies)
				companies.GET("/:id", companyHandler.GetCompanyByID)
				companies.PUT("/:id", companyHandler.UpdateCompany)
				companies.DELETE("/:id", companyHandler.DeleteCompany)
			}

			// Chart of Accounts
			accounts := protected.Group("/accounts")
			{
				accounts.POST("", accountHandler.CreateAccount)
				accounts.GET("", accountHandler.GetAllAccounts)
				accounts.GET("/:id", accountHandler.GetAccountByID)
				accounts.PUT("/:id", accountHandler.UpdateAccount)
				accounts.DELETE("/:id", accountHandler.DeleteAccount)
			}

			// Journals
			journals := protected.Group("/journals")
			{
				journals.POST("", journalHandler.CreateJournal)
				journals.GET("", journalHandler.GetJournalsByPeriod)
				journals.GET("/:id", journalHandler.GetJournalByID)
				journals.PUT("/:id", journalHandler.UpdateJournal)
				journals.DELETE("/:id", journalHandler.DeleteJournal)
			}

			// Ledgers
			ledgers := protected.Group("/ledgers")
			{
				ledgers.GET("/account/:account_id", ledgerHandler.GetLedgerByAccount)
				ledgers.GET("/trial-balance", ledgerHandler.GetTrialBalance)
			}

			// Reports
			reports := protected.Group("/reports")
			{
				reports.GET("/income-statement", reportHandler.GetIncomeStatement)
				reports.GET("/balance-sheet", reportHandler.GetBalanceSheet)
				reports.GET("/cash-flow", reportHandler.GetCashFlow)
			}

			// Cash & Bank
			cashBank := protected.Group("/cash-bank")
			{
				cashBank.POST("", cashBankHandler.CreateTransaction)
				cashBank.GET("", cashBankHandler.GetTransactionsByPeriod)
			}

			// Tax
			taxes := protected.Group("/taxes")
			{
				taxes.POST("", taxHandler.CreateTax)
				taxes.GET("/period", taxHandler.GetTaxesByPeriod)
			}

			// Dashboard
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/summary", dashboardHandler.GetDashboardSummary)
			}

			// Notifications
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", notificationHandler.GetNotifications)
			}

			// Inventory
			inventory := protected.Group("/inventory")
			{
				inventory.POST("/products", inventoryHandler.CreateProduct)
				inventory.GET("/products", inventoryHandler.GetAllProducts)
			}

			// Audit Logs
			auditLogs := protected.Group("/audit-logs")
			auditLogs.Use(middleware.RoleMiddleware("admin"))
			{
				auditLogs.GET("", auditLogHandler.GetAuditLogs)
			}

			// Export
			exports := protected.Group("/export")
			{
				exports.GET("/journals", exportHandler.ExportJournals)
			}

			// Backups
			backups := protected.Group("/backups")
			backups.Use(middleware.RoleMiddleware("admin"))
			{
				backups.POST("", backupHandler.CreateBackup)
				backups.GET("", backupHandler.GetBackups)
			}
		}
	}

	return r
}

func cleanupDatabase(db *gorm.DB) {
	log.Println("🧹 Cleaning up test database...")

	// Get database driver name
	dbDriver := db.Dialector.Name()

	// List of tables to truncate (in order for databases without CASCADE support)
	tables := []string{
		"audit_logs",
		"backups",
		"stock_opname_items",
		"stock_opnames",
		"stock_movements",
		"stock_balances",
		"products",
		"notifications",
		"taxes",
		"bank_reconciliations",
		"cash_bank_transactions",
		"ledgers",
		"journal_entries",
		"journals",
		"accounts",
		"companies",
		"users",
	}

	switch dbDriver {
	case "mysql":
		// MySQL: Disable foreign key checks
		db.Exec("SET FOREIGN_KEY_CHECKS = 0;")

		for _, table := range tables {
			db.Exec("TRUNCATE TABLE " + table)
		}

		// Re-enable foreign key checks
		db.Exec("SET FOREIGN_KEY_CHECKS = 1;")

	case "postgres":
		// PostgreSQL: Use CASCADE with TRUNCATE
		for _, table := range tables {
			db.Exec("TRUNCATE TABLE " + table + " CASCADE")
		}

	default:
		// Fallback: Delete all records (slower but works for any database)
		log.Printf("⚠️  Unknown database driver '%s', using DELETE fallback", dbDriver)
		for _, table := range tables {
			db.Exec("DELETE FROM " + table)
		}
	}

	log.Println("✅ Database cleanup completed")
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "ok" {
		t.Errorf("Expected status ok, got %v", response["status"])
	}
}

func TestRegisterEndpoint(t *testing.T) {
	router := setupTestRouter()

	registerData := map[string]string{
		"full_name": "Test User", // Changed from "name" to "full_name"
		"email":     "testintegration@example.com",
		"password":  "password123",
		"role":      "admin",
	}

	jsonData, _ := json.Marshal(registerData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestLoginEndpoint(t *testing.T) {
	router := setupTestRouter()

	// First register a user
	registerData := map[string]string{
		"full_name": "Test User Login", // Changed from "name" to "full_name"
		"email":     "testlogin@example.com",
		"password":  "password123",
		"role":      "admin",
	}
	jsonData, _ := json.Marshal(registerData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req)

	// Then login
	loginData := map[string]string{
		"email":    "testlogin@example.com",
		"password": "password123",
	}
	jsonData, _ = json.Marshal(loginData)
	w := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Check if response has success field
	if success, ok := response["success"].(bool); !ok || !success {
		t.Errorf("Expected success: true, got %v", response)
	}

	// Check if token exists in data
	if data, ok := response["data"].(map[string]interface{}); ok {
		if token, exists := data["token"]; !exists || token == nil {
			t.Error("Expected token in response data")
		}
	} else {
		t.Error("Expected data field in response")
	}
}

func TestProtectedEndpoint_WithoutToken(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/profile", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestProtectedEndpoint_WithToken(t *testing.T) {
	router := setupTestRouter()

	// First register and login to get token
	registerData := map[string]string{
		"full_name": "Test Protected",
		"email":     "testprotected@example.com",
		"password":  "password123",
		"role":      "admin",
	}
	jsonData, _ := json.Marshal(registerData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req)

	// Login
	loginData := map[string]string{
		"email":    "testprotected@example.com",
		"password": "password123",
	}
	jsonData, _ = json.Marshal(loginData)
	w := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)

	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	// Now access protected endpoint with token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}
}

func TestCreateAccount(t *testing.T) {
	router := setupTestRouter()

	// Register and login first
	registerData := map[string]string{
		"full_name": "Test Account Creator",
		"email":     "testaccount@example.com",
		"password":  "password123",
		"role":      "admin",
	}
	jsonData, _ := json.Marshal(registerData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(httptest.NewRecorder(), req)

	// Login to get token
	loginData := map[string]string{
		"email":    "testaccount@example.com",
		"password": "password123",
	}
	jsonData, _ = json.Marshal(loginData)
	w := httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	token := loginResponse["data"].(map[string]interface{})["token"].(string)

	// Create company first (required for account)
	companyData := map[string]interface{}{
		"name":           "Test Company",
		"address":        "Test Address",
		"phone":          "081234567890",
		"email":          "company@example.com",
		"tax_id":         "12.345.678.9-012.345",
		"currency":       "IDR",
		"fiscal_year_end": "12-31",
	}
	jsonData, _ = json.Marshal(companyData)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/companies", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	var companyResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &companyResponse)
	companyID := uint(companyResponse["data"].(map[string]interface{})["id"].(float64))

	// Create account with company_id
	accountData := map[string]interface{}{
		"company_id": companyID,
		"code":       "1-1100",
		"name":       "Kas",
		"type":       "asset",
		"category":   "current_asset",
		"level":      3,
		"is_header":  false,
		"is_active":  true,
	}
	jsonData, _ = json.Marshal(accountData)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/accounts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}
