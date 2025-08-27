package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"savvy-backend/internal/usecase"
)

type DashboardHandler struct {
	dashboardUsecase usecase.DashboardUsecase
}

func NewDashboardHandler(dashboardUsecase usecase.DashboardUsecase) *DashboardHandler {
	return &DashboardHandler{
		dashboardUsecase: dashboardUsecase,
	}
}

func (h *DashboardHandler) GetCurrentMonthlySummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	summary, err := h.dashboardUsecase.GetCurrentMonthlySummary(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *DashboardHandler) GetMonthlySummary(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse year and month from query parameters
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	var year, month int
	var err error

	if yearStr == "" || monthStr == "" {
		// Use current year and month if not provided
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	} else {
		year, err = strconv.Atoi(yearStr)
		if err != nil || year < 1900 || year > 2100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
			return
		}

		month, err = strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
			return
		}
	}

	summary, err := h.dashboardUsecase.GetMonthlySummary(c.Request.Context(), userID.(uuid.UUID), year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (h *DashboardHandler) GetRecentTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse limit from query parameter
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10 // Default to 10 if invalid
	}

	transactions, err := h.dashboardUsecase.GetRecentTransactions(c.Request.Context(), userID.(uuid.UUID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

func (h *DashboardHandler) GetSpendingByCategory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse year and month from query parameters
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	var year, month int
	var err error

	if yearStr == "" || monthStr == "" {
		// Use current year and month if not provided
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	} else {
		year, err = strconv.Atoi(yearStr)
		if err != nil || year < 1900 || year > 2100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year"})
			return
		}

		month, err = strconv.Atoi(monthStr)
		if err != nil || month < 1 || month > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid month"})
			return
		}
	}

	spending, err := h.dashboardUsecase.GetSpendingByCategory(c.Request.Context(), userID.(uuid.UUID), year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"spending_by_category": spending})
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get current month summary
	summary, err := h.dashboardUsecase.GetCurrentMonthlySummary(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get monthly summary"})
		return
	}

	// Get recent transactions
	recentTransactions, err := h.dashboardUsecase.GetRecentTransactions(c.Request.Context(), userID.(uuid.UUID), 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recent transactions"})
		return
	}

	// Get spending by category for current month
	now := time.Now()
	spendingByCategory, err := h.dashboardUsecase.GetSpendingByCategory(c.Request.Context(), userID.(uuid.UUID), now.Year(), int(now.Month()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get spending by category"})
		return
	}

	response := gin.H{
		"monthly_summary":      summary,
		"recent_transactions":  recentTransactions,
		"spending_by_category": spendingByCategory,
	}

	c.JSON(http.StatusOK, response)
}
