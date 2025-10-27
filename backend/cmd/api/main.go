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

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg)
	userService := services.NewUserService(userRepo)
	companyService := services.NewCompanyService(companyRepo)
	accountService := services.NewAccountService(accountRepo)
	journalService := services.NewJournalService(journalRepo, ledgerRepo, accountRepo)
	ledgerService := services.NewLedgerService(ledgerRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	companyHandler := handlers.NewCompanyHandler(companyService, userService)
	accountHandler := handlers.NewAccountHandler(accountService)
	journalHandler := handlers.NewJournalHandler(journalService)
	ledgerHandler := handlers.NewLedgerHandler(ledgerService)

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
				accounts.GET("/:id", accountHandler.GetAccountByID)
				accounts.PUT("/:id", accountHandler.UpdateAccount)
				accounts.DELETE("/:id", accountHandler.DeleteAccount)
				accounts.POST("/initialize", accountHandler.InitializeDefaultAccounts)
			}

			// Journal (Jurnal Umum)
			journals := protected.Group("/journals")
			{
				journals.POST("", journalHandler.CreateJournal)
				journals.GET("", journalHandler.GetJournalsByPeriod) // ?start_date=2025-01-01&end_date=2025-12-31
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
				ledgers.GET("/account/:account_id", ledgerHandler.GetLedgerByAccount) // ?start_date=&end_date=
				ledgers.GET("/company", ledgerHandler.GetLedgerByCompany)             // ?start_date=&end_date=
				ledgers.GET("/trial-balance", ledgerHandler.GetTrialBalance)          // ?end_date=
				ledgers.GET("/account/:account_id/balance", ledgerHandler.GetAccountBalance)
			}
		}
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}