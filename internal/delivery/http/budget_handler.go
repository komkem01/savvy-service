package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/usecase"
)

type BudgetHandler struct {
	budgetUsecase usecase.BudgetUsecase
}

func NewBudgetHandler(budgetUsecase usecase.BudgetUsecase) *BudgetHandler {
	return &BudgetHandler{
		budgetUsecase: budgetUsecase,
	}
}

type CreateBudgetRequest struct {
	CategoryID string `json:"category_id" binding:"required,uuid"`
	Amount     string `json:"amount" binding:"required"`
	Period     string `json:"period" binding:"required,oneof=monthly yearly"`
	StartDate  string `json:"start_date" binding:"required"`
	EndDate    string `json:"end_date,omitempty"`
}

type UpdateBudgetRequest struct {
	Amount    *string `json:"amount,omitempty"`
	Period    *string `json:"period,omitempty"`
	StartDate *string `json:"start_date,omitempty"`
	EndDate   *string `json:"end_date,omitempty"`
	IsActive  *bool   `json:"is_active,omitempty"`
}

type BudgetResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CategoryID string    `json:"category_id"`
	Amount     string    `json:"amount"`
	Period     string    `json:"period"`
	StartDate  string    `json:"start_date"`
	EndDate    *string   `json:"end_date,omitempty"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type BudgetProgressResponse struct {
	BudgetID        string  `json:"budget_id"`
	SpentAmount     string  `json:"spent_amount"`
	RemainingAmount string  `json:"remaining_amount"`
	PercentageUsed  float64 `json:"percentage_used"`
	IsOverBudget    bool    `json:"is_over_budget"`
	DaysRemaining   int     `json:"days_remaining"`
	AverageDaily    string  `json:"average_daily_spent"`
	ProjectedTotal  string  `json:"projected_total"`
}

func (h *BudgetHandler) CreateBudget(c *gin.Context) {
	var req CreateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	categoryUUID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
		return
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Expected YYYY-MM-DD"})
		return
	}

	// Parse end date if provided
	var endDate *time.Time
	if req.EndDate != "" {
		parsed, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Expected YYYY-MM-DD"})
			return
		}
		endDate = &parsed
	}

	budget, err := h.budgetUsecase.CreateBudget(c.Request.Context(), userUUID, categoryUUID, amount, entity.BudgetPeriod(req.Period), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.budgetToResponse(budget)
	c.JSON(http.StatusCreated, response)
}

func (h *BudgetHandler) GetBudgets(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	budgets, err := h.budgetUsecase.GetUserBudgets(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]*BudgetResponse, len(budgets))
	for i, budget := range budgets {
		responses[i] = h.budgetToResponse(budget)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *BudgetHandler) GetBudget(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	budgetID := c.Param("id")

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	budgetUUID, err := uuid.Parse(budgetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}

	budget, err := h.budgetUsecase.GetBudgetByID(c.Request.Context(), userUUID, budgetUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := h.budgetToResponse(budget)
	c.JSON(http.StatusOK, response)
}

func (h *BudgetHandler) UpdateBudget(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	budgetID := c.Param("id")

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	budgetUUID, err := uuid.Parse(budgetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}

	var req UpdateBudgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing budget
	budget, err := h.budgetUsecase.GetBudgetByID(c.Request.Context(), userUUID, budgetUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.Amount != nil {
		amount, err := decimal.NewFromString(*req.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
			return
		}
		budget.Amount = amount
	}
	if req.Period != nil {
		budget.Period = entity.BudgetPeriod(*req.Period)
	}
	if req.StartDate != nil {
		startDate, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Expected YYYY-MM-DD"})
			return
		}
		budget.StartDate = startDate
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			budget.EndDate = nil
		} else {
			endDate, err := time.Parse("2006-01-02", *req.EndDate)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Expected YYYY-MM-DD"})
				return
			}
			budget.EndDate = &endDate
		}
	}
	if req.IsActive != nil {
		budget.IsActive = *req.IsActive
	}

	err = h.budgetUsecase.UpdateBudget(c.Request.Context(), userUUID, budget)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.budgetToResponse(budget)
	c.JSON(http.StatusOK, response)
}

func (h *BudgetHandler) DeleteBudget(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	budgetID := c.Param("id")

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	budgetUUID, err := uuid.Parse(budgetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid budget ID"})
		return
	}

	err = h.budgetUsecase.DeleteBudget(c.Request.Context(), userUUID, budgetUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *BudgetHandler) GetBudgetProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Get year and month from query params, default to current month
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

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

	progressList, err := h.budgetUsecase.GetBudgetProgress(c.Request.Context(), userUUID, year, month)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]*BudgetProgressResponse, len(progressList))
	for i, progress := range progressList {
		responses[i] = &BudgetProgressResponse{
			BudgetID:        progress.BudgetID.String(),
			SpentAmount:     progress.SpentAmount.String(),
			RemainingAmount: progress.RemainingAmount.String(),
			PercentageUsed:  progress.ProgressPercentage,
			IsOverBudget:    progress.IsOverBudget,
			DaysRemaining:   0,   // Will be calculated based on period
			AverageDaily:    "0", // Will be calculated based on period
			ProjectedTotal:  progress.BudgetAmount.String(),
		}
	}

	c.JSON(http.StatusOK, responses)
}

func (h *BudgetHandler) GetCurrentMonthProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	progressList, err := h.budgetUsecase.GetCurrentMonthBudgetProgress(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]*BudgetProgressResponse, len(progressList))
	for i, progress := range progressList {
		responses[i] = &BudgetProgressResponse{
			BudgetID:        progress.BudgetID.String(),
			SpentAmount:     progress.SpentAmount.String(),
			RemainingAmount: progress.RemainingAmount.String(),
			PercentageUsed:  progress.ProgressPercentage,
			IsOverBudget:    progress.IsOverBudget,
			DaysRemaining:   0,   // Will be calculated based on period
			AverageDaily:    "0", // Will be calculated based on period
			ProjectedTotal:  progress.BudgetAmount.String(),
		}
	}

	c.JSON(http.StatusOK, responses)
}

func (h *BudgetHandler) CheckBudgetAlerts(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	insights, err := h.budgetUsecase.CheckBudgetAlerts(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Budget alerts checked successfully",
		"alerts":  len(insights),
	})
}

func (h *BudgetHandler) budgetToResponse(budget *entity.Budget) *BudgetResponse {
	var endDate *string
	if budget.EndDate != nil {
		endDateStr := budget.EndDate.Format("2006-01-02")
		endDate = &endDateStr
	}

	return &BudgetResponse{
		ID:         budget.ID.String(),
		UserID:     budget.UserID.String(),
		CategoryID: budget.CategoryID.String(),
		Amount:     budget.Amount.String(),
		Period:     string(budget.Period),
		StartDate:  budget.StartDate.Format("2006-01-02"),
		EndDate:    endDate,
		IsActive:   budget.IsActive,
		CreatedAt:  budget.CreatedAt,
		UpdatedAt:  budget.UpdatedAt,
	}
}
