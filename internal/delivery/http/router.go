package http

import (
	"savvy-backend/internal/usecase"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	authUsecase usecase.AuthUsecase,
	transactionUsecase usecase.TransactionUsecase,
	accountUsecase usecase.AccountUsecase,
	categoryUsecase usecase.CategoryUsecase,
	dashboardUsecase usecase.DashboardUsecase,
	budgetUsecase usecase.BudgetUsecase,
	recurringUsecase usecase.RecurringTransactionUsecase,
	aiInsightUsecase usecase.AIInsightUsecase,
) *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(CORSMiddleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api/v1")

	// Auth routes
	authHandler := NewAuthHandler(authUsecase)
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes
	protected := api.Group("/")
	protected.Use(AuthMiddleware(authUsecase))
	{
		// Account routes
		accountHandler := NewAccountHandler(accountUsecase)
		accounts := protected.Group("/accounts")
		{
			accounts.POST("/", accountHandler.CreateAccount)
			accounts.GET("/", accountHandler.GetAccounts)
			accounts.GET("/:id", accountHandler.GetAccount)
			accounts.DELETE("/:id", accountHandler.DeleteAccount)
		}

		// Transaction routes
		transactionHandler := NewTransactionHandler(transactionUsecase)
		transactions := protected.Group("/transactions")
		{
			transactions.POST("/", transactionHandler.CreateTransaction)
			transactions.GET("/", transactionHandler.GetTransactions)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.PUT("/:id", transactionHandler.UpdateTransaction)
			transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
		}

		// Category routes
		categoryHandler := NewCategoryHandler(categoryUsecase)
		categories := protected.Group("/categories")
		{
			categories.POST("/", categoryHandler.CreateCategory)
			categories.GET("/", categoryHandler.GetCategories)
			categories.GET("/user", categoryHandler.GetUserCategories)
			categories.PUT("/:id/archive", categoryHandler.ArchiveCategory)
		}

		// Dashboard routes
		dashboardHandler := NewDashboardHandler(dashboardUsecase)
		dashboard := protected.Group("/dashboard")
		{
			dashboard.GET("/", dashboardHandler.GetDashboard)
			dashboard.GET("/summary", dashboardHandler.GetCurrentMonthlySummary)
			dashboard.GET("/summary/monthly", dashboardHandler.GetMonthlySummary)
			dashboard.GET("/transactions/recent", dashboardHandler.GetRecentTransactions)
			dashboard.GET("/spending/category", dashboardHandler.GetSpendingByCategory)
		}

		// Analytics routes for Data Visualization
		analyticsHandler := NewAnalyticsHandler(dashboardUsecase, transactionUsecase, categoryUsecase)
		analytics := protected.Group("/analytics")
		{
			analytics.GET("/pie/expenses", analyticsHandler.GetExpensePieChart)
			analytics.GET("/bar/income-expense", analyticsHandler.GetIncomeExpenseBarChart)
			analytics.GET("/trend/category/:category_id", analyticsHandler.GetCategoryTrendChart)
			analytics.GET("/top/categories", analyticsHandler.GetTopCategoriesChart)
		}

		// Budget routes
		budgetHandler := NewBudgetHandler(budgetUsecase)
		budgets := protected.Group("/budgets")
		{
			budgets.POST("/", budgetHandler.CreateBudget)
			budgets.GET("/", budgetHandler.GetBudgets)
			budgets.GET("/:id", budgetHandler.GetBudget)
			budgets.PUT("/:id", budgetHandler.UpdateBudget)
			budgets.DELETE("/:id", budgetHandler.DeleteBudget)
			budgets.GET("/progress", budgetHandler.GetBudgetProgress)
			budgets.GET("/progress/current", budgetHandler.GetCurrentMonthProgress)
			budgets.POST("/alerts/check", budgetHandler.CheckBudgetAlerts)
		}

		// Recurring Transaction routes
		recurringHandler := NewRecurringTransactionHandler(recurringUsecase)
		recurring := protected.Group("/recurring-transactions")
		{
			recurring.POST("/", recurringHandler.CreateRecurringTransaction)
			recurring.GET("/", recurringHandler.GetRecurringTransactions)
			recurring.GET("/:id", recurringHandler.GetRecurringTransaction)
			recurring.PUT("/:id", recurringHandler.UpdateRecurringTransaction)
			recurring.DELETE("/:id", recurringHandler.DeleteRecurringTransaction)
			recurring.POST("/:id/execute", recurringHandler.ExecuteRecurringTransaction)
			recurring.GET("/due", recurringHandler.GetDueTransactions)
		}

		// AI Insights routes
		aiInsightHandler := NewAIInsightHandler(aiInsightUsecase)
		aiInsights := protected.Group("/ai-insights")
		{
			// Generate insights
			aiInsights.POST("/spending-anomalies/generate", aiInsightHandler.GenerateSpendingAnomalyInsights)
			aiInsights.POST("/spending-patterns/generate", aiInsightHandler.GenerateSpendingPatternInsights)
			aiInsights.POST("/category-recommendations/generate", aiInsightHandler.GenerateCategoryRecommendations)
			aiInsights.POST("/savings-recommendations/generate", aiInsightHandler.GenerateSavingsRecommendations)

			// Process insights
			aiInsights.POST("/weekly/process", aiInsightHandler.ProcessWeeklyInsights)

			// Get processed insights (placeholder endpoints)
			aiInsights.GET("/spending-anomalies", aiInsightHandler.GetSpendingAnomalies)
			aiInsights.GET("/spending-patterns", aiInsightHandler.GetSpendingPatterns)
		}
	}

	// Admin/Setup routes (should be protected in production)
	setup := api.Group("/setup")
	{
		categoryHandler := NewCategoryHandler(categoryUsecase)
		setup.POST("/categories/default", categoryHandler.InitializeDefaultCategories)
	}

	// System routes (for background processing)
	system := api.Group("/system")
	{
		// These should be protected with admin authentication in production
		recurringHandler := NewRecurringTransactionHandler(recurringUsecase)
		system.POST("/recurring-transactions/process-all", recurringHandler.ProcessAllDueTransactions)

		aiInsightHandler := NewAIInsightHandler(aiInsightUsecase)
		system.POST("/ai-insights/process-all", aiInsightHandler.ProcessAllUsersInsights)
	}

	return r
}
