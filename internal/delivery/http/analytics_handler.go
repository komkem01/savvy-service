package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/usecase"
)

type AnalyticsHandler struct {
	dashboardUsecase   usecase.DashboardUsecase
	transactionUsecase usecase.TransactionUsecase
	categoryUsecase    usecase.CategoryUsecase
}

type PieChartData struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
	Color string  `json:"color"`
}

type BarChartData struct {
	Label   string  `json:"label"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
}

type AnalyticsResponse struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func NewAnalyticsHandler(
	dashboardUsecase usecase.DashboardUsecase,
	transactionUsecase usecase.TransactionUsecase,
	categoryUsecase usecase.CategoryUsecase,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		dashboardUsecase:   dashboardUsecase,
		transactionUsecase: transactionUsecase,
		categoryUsecase:    categoryUsecase,
	}
}

// GetExpensePieChart - แสดง pie chart ของค่าใช้จ่ายตามหมวดหมู่
func (h *AnalyticsHandler) GetExpensePieChart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse year and month from query params
	year := time.Now().Year()
	month := int(time.Now().Month())

	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	if monthStr := c.Query("month"); monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			month = m
		}
	}

	// Get spending by category
	categorySpending, err := h.dashboardUsecase.GetSpendingByCategory(c.Request.Context(), userID.(uuid.UUID), year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build pie chart data
	var pieData []PieChartData
	for _, spending := range categorySpending {
		amount, _ := spending.Amount.Float64()
		color := "#007BFF" // default color
		if spending.ColorHex != nil {
			color = *spending.ColorHex
		}

		pieData = append(pieData, PieChartData{
			Label: spending.CategoryName,
			Value: amount,
			Color: color,
		})
	}

	response := AnalyticsResponse{
		Type: "pie",
		Data: pieData,
	}

	c.JSON(http.StatusOK, response)
}

// GetIncomeExpenseBarChart - แสดง bar chart เปรียบเทียบรายรับ-รายจ่ายรายเดือน
func (h *AnalyticsHandler) GetIncomeExpenseBarChart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse year from query param
	year := time.Now().Year()
	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	var barData []BarChartData

	// Get data for each month
	for month := 1; month <= 12; month++ {
		summary, err := h.dashboardUsecase.GetMonthlySummary(c.Request.Context(), userID.(uuid.UUID), year, month)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		monthName := time.Month(month).String()

		// Convert decimal to float64
		income, _ := summary.TotalIncome.Float64()
		expense, _ := summary.TotalExpense.Float64()

		barData = append(barData, BarChartData{
			Label:   monthName,
			Income:  income,
			Expense: expense,
		})
	}

	response := AnalyticsResponse{
		Type: "bar",
		Data: barData,
	}

	c.JSON(http.StatusOK, response)
}

// GetCategoryTrendChart - แสดงแนวโน้มการใช้จ่ายของหมวดหมู่เฉพาะ
func (h *AnalyticsHandler) GetCategoryTrendChart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	categoryIDStr := c.Param("category_id")
	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Parse year from query param
	year := time.Now().Year()
	if yearStr := c.Query("year"); yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}

	type TrendData struct {
		Month  string  `json:"month"`
		Amount float64 `json:"amount"`
	}

	var trendData []TrendData

	// Get data for each month for specific category
	for month := 1; month <= 12; month++ {
		categorySpending, err := h.dashboardUsecase.GetSpendingByCategory(c.Request.Context(), userID.(uuid.UUID), year, month)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		amount := 0.0
		for _, spending := range categorySpending {
			if spending.CategoryID == categoryID {
				amount, _ = spending.Amount.Float64()
				break
			}
		}

		monthName := time.Month(month).String()
		trendData = append(trendData, TrendData{
			Month:  monthName,
			Amount: amount,
		})
	}

	response := AnalyticsResponse{
		Type: "line",
		Data: trendData,
	}

	c.JSON(http.StatusOK, response)
}

// GetTopCategoriesChart - แสดง top spending categories
func (h *AnalyticsHandler) GetTopCategoriesChart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get category usage stats
	usageStats, err := h.categoryUsecase.GetCategoryUsageStats(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get categories to map names
	categories, err := h.categoryUsecase.GetUserCategories(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	categoryMap := make(map[uuid.UUID]*entity.Category)
	for _, cat := range categories {
		categoryMap[cat.ID] = cat
	}

	type TopCategoryData struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
		Color string `json:"color"`
	}

	var topData []TopCategoryData
	for categoryID, count := range usageStats {
		if category, exists := categoryMap[categoryID]; exists {
			color := "#007BFF" // default color
			if category.ColorHex != nil {
				color = *category.ColorHex
			}

			topData = append(topData, TopCategoryData{
				Name:  category.Name,
				Count: count,
				Color: color,
			})
		}
	}

	// Sort by count (top 10)
	if len(topData) > 10 {
		topData = topData[:10]
	}

	response := AnalyticsResponse{
		Type: "horizontal-bar",
		Data: topData,
	}

	c.JSON(http.StatusOK, response)
}
