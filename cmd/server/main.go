package main

import (
	"log"

	"github.com/joho/godotenv"

	"savvy-backend/internal/config"
	"savvy-backend/internal/delivery/http"
	"savvy-backend/internal/infrastructure/database"
	"savvy-backend/internal/usecase"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Load configuration
	cfg := config.Load()

	// Database connection
	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := database.NewUserRepository(db)
	accountRepo := database.NewAccountRepository(db)
	categoryRepo := database.NewCategoryRepository(db)
	transactionRepo := database.NewTransactionRepository(db)
	budgetRepo := database.NewBudgetRepository(db)
	recurringRepo := database.NewRecurringTransactionRepository(db)
	insightRepo := database.NewInsightRepository(db)

	// Initialize use cases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWT.Secret, cfg.JWT.Expiry)
	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo, accountRepo, categoryRepo)
	accountUsecase := usecase.NewAccountUsecase(accountRepo)
	categoryUsecase := usecase.NewCategoryUsecase(categoryRepo)
	dashboardUsecase := usecase.NewDashboardUsecase(transactionRepo, categoryRepo, accountRepo)
	budgetUsecase := usecase.NewBudgetUsecase(budgetRepo, categoryRepo, insightRepo)
	recurringUsecase := usecase.NewRecurringTransactionUsecase(recurringRepo, transactionRepo, categoryRepo, accountRepo)
	aiInsightUsecase := usecase.NewAIInsightUsecase(insightRepo, transactionRepo, categoryRepo, budgetRepo)

	// Setup routes
	router := http.SetupRoutes(authUsecase, transactionUsecase, accountUsecase, categoryUsecase, dashboardUsecase, budgetUsecase, recurringUsecase, aiInsightUsecase)

	// Start server
	serverAddr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
