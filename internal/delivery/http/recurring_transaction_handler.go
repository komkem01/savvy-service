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

type RecurringTransactionHandler struct {
	recurringUsecase usecase.RecurringTransactionUsecase
}

func NewRecurringTransactionHandler(recurringUsecase usecase.RecurringTransactionUsecase) *RecurringTransactionHandler {
	return &RecurringTransactionHandler{
		recurringUsecase: recurringUsecase,
	}
}

type CreateRecurringTransactionRequest struct {
	CategoryID          string  `json:"category_id" binding:"required,uuid"`
	AccountID           string  `json:"account_id" binding:"required,uuid"`
	Amount              string  `json:"amount" binding:"required"`
	Type                string  `json:"type" binding:"required,oneof=income expense"`
	Note                *string `json:"note,omitempty"`
	Frequency           string  `json:"frequency" binding:"required,oneof=daily weekly monthly yearly"`
	StartDate           string  `json:"start_date" binding:"required"`
	EndDate             *string `json:"end_date,omitempty"`
	AutoExecute         bool    `json:"auto_execute"`
	RemainingExecutions *int    `json:"remaining_executions,omitempty"`
}

type UpdateRecurringTransactionRequest struct {
	Amount              *string `json:"amount,omitempty"`
	Note                *string `json:"note,omitempty"`
	Frequency           *string `json:"frequency,omitempty"`
	EndDate             *string `json:"end_date,omitempty"`
	AutoExecute         *bool   `json:"auto_execute,omitempty"`
	RemainingExecutions *int    `json:"remaining_executions,omitempty"`
	IsActive            *bool   `json:"is_active,omitempty"`
}

type RecurringTransactionResponse struct {
	ID                  string     `json:"id"`
	UserID              string     `json:"user_id"`
	CategoryID          string     `json:"category_id"`
	AccountID           string     `json:"account_id"`
	Amount              string     `json:"amount"`
	Type                string     `json:"type"`
	Note                *string    `json:"note,omitempty"`
	Frequency           string     `json:"frequency"`
	StartDate           string     `json:"start_date"`
	EndDate             *string    `json:"end_date,omitempty"`
	NextExecutionDate   time.Time  `json:"next_execution_date"`
	LastExecutionDate   *time.Time `json:"last_execution_date,omitempty"`
	IsActive            bool       `json:"is_active"`
	AutoExecute         bool       `json:"auto_execute"`
	RemainingExecutions *int       `json:"remaining_executions,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func (h *RecurringTransactionHandler) CreateRecurringTransaction(c *gin.Context) {
	var req CreateRecurringTransactionRequest
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

	accountUUID, err := uuid.Parse(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount format"})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format, expected YYYY-MM-DD"})
		return
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, expected YYYY-MM-DD"})
			return
		}
		endDate = &parsed
	}

	// Create recurring transaction entity
	recurringTx := &entity.RecurringTransaction{
		ID:                  uuid.New(),
		UserID:              userUUID,
		CategoryID:          categoryUUID,
		AccountID:           accountUUID,
		Amount:              amount,
		Type:                entity.TransactionType(req.Type),
		Note:                req.Note,
		Frequency:           entity.RecurringFrequency(req.Frequency),
		StartDate:           startDate,
		EndDate:             endDate,
		AutoExecute:         req.AutoExecute,
		RemainingExecutions: req.RemainingExecutions,
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Calculate next execution date
	recurringTx.NextExecutionDate = recurringTx.CalculateNextExecutionDate()

	recurringTransaction, err := h.recurringUsecase.CreateRecurringTransaction(c.Request.Context(), userUUID, recurringTx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.recurringToResponse(recurringTransaction)
	c.JSON(http.StatusCreated, response)
}

func (h *RecurringTransactionHandler) GetRecurringTransactions(c *gin.Context) {
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

	// Parse query parameters - we'll filter manually since the interface doesn't support these params
	frequency := c.Query("frequency")
	activeOnly := c.Query("active_only") == "true"

	recurringTransactions, err := h.recurringUsecase.GetUserRecurringTransactions(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Filter results manually
	filteredTransactions := make([]*entity.RecurringTransaction, 0)
	for _, tx := range recurringTransactions {
		// Filter by frequency
		if frequency != "" && string(tx.Frequency) != frequency {
			continue
		}
		// Filter by active status
		if activeOnly && !tx.IsActive {
			continue
		}
		filteredTransactions = append(filteredTransactions, tx)
	}

	responses := make([]*RecurringTransactionResponse, len(filteredTransactions))
	for i, tx := range filteredTransactions {
		responses[i] = h.recurringToResponse(tx)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *RecurringTransactionHandler) GetRecurringTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	txID := c.Param("id")

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	txUUID, err := uuid.Parse(txID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	recurringTransaction, err := h.recurringUsecase.GetRecurringTransactionByID(c.Request.Context(), userUUID, txUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := h.recurringToResponse(recurringTransaction)
	c.JSON(http.StatusOK, response)
}

func (h *RecurringTransactionHandler) UpdateRecurringTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	txID := c.Param("id")

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	txUUID, err := uuid.Parse(txID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var req UpdateRecurringTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing recurring transaction
	recurringTransaction, err := h.recurringUsecase.GetRecurringTransactionByID(c.Request.Context(), userUUID, txUUID)
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
		recurringTransaction.Amount = amount
	}
	if req.Note != nil {
		recurringTransaction.Note = req.Note
	}
	if req.Frequency != nil {
		recurringTransaction.Frequency = entity.RecurringFrequency(*req.Frequency)
		// Recalculate next execution date if frequency changed
		recurringTransaction.NextExecutionDate = recurringTransaction.CalculateNextExecutionDate()
	}
	if req.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format, expected YYYY-MM-DD"})
			return
		}
		recurringTransaction.EndDate = &endDate
	}
	if req.AutoExecute != nil {
		recurringTransaction.AutoExecute = *req.AutoExecute
	}
	if req.RemainingExecutions != nil {
		recurringTransaction.RemainingExecutions = req.RemainingExecutions
	}
	if req.IsActive != nil {
		recurringTransaction.IsActive = *req.IsActive
	}

	err = h.recurringUsecase.UpdateRecurringTransaction(c.Request.Context(), userUUID, recurringTransaction)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := h.recurringToResponse(recurringTransaction)
	c.JSON(http.StatusOK, response)
}

func (h *RecurringTransactionHandler) DeleteRecurringTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	txID := c.Param("id")

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	txUUID, err := uuid.Parse(txID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	err = h.recurringUsecase.DeleteRecurringTransaction(c.Request.Context(), userUUID, txUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *RecurringTransactionHandler) ExecuteRecurringTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	txID := c.Param("id")

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	txUUID, err := uuid.Parse(txID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.recurringUsecase.ExecuteRecurringTransaction(c.Request.Context(), userUUID, txUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Recurring transaction executed successfully",
		"transaction": transaction,
	})
}

func (h *RecurringTransactionHandler) GetDueTransactions(c *gin.Context) {
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

	// Parse limit parameter
	limit := 50 // default
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	dueTransactions, err := h.recurringUsecase.GetDueTransactions(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Apply limit manually
	if len(dueTransactions) > limit {
		dueTransactions = dueTransactions[:limit]
	}

	responses := make([]*RecurringTransactionResponse, len(dueTransactions))
	for i, tx := range dueTransactions {
		responses[i] = h.recurringToResponse(tx)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *RecurringTransactionHandler) ProcessAllDueTransactions(c *gin.Context) {
	// This endpoint processes all due transactions system-wide, not user-specific
	// In a production system, this might be restricted to admin users or scheduled jobs

	err := h.recurringUsecase.ProcessAllDueTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All due transactions processed successfully",
	})
}

func (h *RecurringTransactionHandler) recurringToResponse(tx *entity.RecurringTransaction) *RecurringTransactionResponse {
	response := &RecurringTransactionResponse{
		ID:                  tx.ID.String(),
		UserID:              tx.UserID.String(),
		CategoryID:          tx.CategoryID.String(),
		AccountID:           tx.AccountID.String(),
		Amount:              tx.Amount.String(),
		Type:                string(tx.Type),
		Note:                tx.Note,
		Frequency:           string(tx.Frequency),
		StartDate:           tx.StartDate.Format("2006-01-02"),
		NextExecutionDate:   tx.NextExecutionDate,
		LastExecutionDate:   tx.LastExecutionDate,
		IsActive:            tx.IsActive,
		AutoExecute:         tx.AutoExecute,
		RemainingExecutions: tx.RemainingExecutions,
		CreatedAt:           tx.CreatedAt,
		UpdatedAt:           tx.UpdatedAt,
	}

	if tx.EndDate != nil {
		endDate := tx.EndDate.Format("2006-01-02")
		response.EndDate = &endDate
	}

	return response
}
