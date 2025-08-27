package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/domain/repository"
	"savvy-backend/internal/usecase"
)

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
}

type CreateTransactionRequest struct {
	CategoryID      string  `json:"category_id" binding:"required"`
	AccountID       string  `json:"account_id" binding:"required"`
	Amount          string  `json:"amount" binding:"required"`
	Type            string  `json:"type" binding:"required"`
	Note            *string `json:"note,omitempty"`
	TransactionDate string  `json:"transaction_date" binding:"required"`
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	var transactionType entity.TransactionType
	switch req.Type {
	case "income":
		transactionType = entity.TransactionTypeIncome
	case "expense":
		transactionType = entity.TransactionTypeExpense
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
		return
	}

	transaction, err := h.transactionUsecase.CreateTransaction(
		c.Request.Context(),
		userID.(uuid.UUID),
		categoryID,
		accountID,
		amount,
		transactionType,
		req.Note,
		req.TransactionDate,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	filter := repository.TransactionFilter{
		UserID: userID.(uuid.UUID),
		Limit:  20, // Default limit
		Offset: 0,
	}

	// Parse query parameters with improved validation
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filter.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}

	if accountIDStr := c.Query("account_id"); accountIDStr != "" {
		if accountID, err := uuid.Parse(accountIDStr); err == nil {
			filter.AccountID = &accountID
		}
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := uuid.Parse(categoryIDStr); err == nil {
			filter.CategoryID = &categoryID
		}
	}

	if typeStr := c.Query("type"); typeStr != "" {
		switch typeStr {
		case "income":
			transactionType := entity.TransactionTypeIncome
			filter.Type = &transactionType
		case "expense":
			transactionType := entity.TransactionTypeExpense
			filter.Type = &transactionType
		}
	}

	// Date range filtering
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		}
	}

	// Advanced search functionality
	if searchQuery := c.Query("search"); searchQuery != "" {
		filter.SearchQuery = &searchQuery
	}

	// Amount range filtering
	if minAmountStr := c.Query("min_amount"); minAmountStr != "" {
		if minAmount, err := decimal.NewFromString(minAmountStr); err == nil {
			filter.MinAmount = &minAmount
		}
	}

	if maxAmountStr := c.Query("max_amount"); maxAmountStr != "" {
		if maxAmount, err := decimal.NewFromString(maxAmountStr); err == nil {
			filter.MaxAmount = &maxAmount
		}
	}

	transactions, err := h.transactionUsecase.GetTransactionsByFilter(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Enhanced response with pagination info and filter details
	response := gin.H{
		"transactions": transactions,
		"pagination": gin.H{
			"limit":  filter.Limit,
			"offset": filter.Offset,
			"count":  len(transactions),
		},
		"filters": gin.H{
			"search":      filter.SearchQuery,
			"min_amount":  filter.MinAmount,
			"max_amount":  filter.MaxAmount,
			"account_id":  filter.AccountID,
			"category_id": filter.CategoryID,
			"type":        filter.Type,
			"start_date":  filter.StartDate,
			"end_date":    filter.EndDate,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	transaction, err := h.transactionUsecase.GetTransactionByID(c.Request.Context(), userID.(uuid.UUID), transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing transaction
	transaction, err := h.transactionUsecase.GetTransactionByID(c.Request.Context(), userID.(uuid.UUID), transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Parse and validate new values
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	var transactionType entity.TransactionType
	switch req.Type {
	case "income":
		transactionType = entity.TransactionTypeIncome
	case "expense":
		transactionType = entity.TransactionTypeExpense
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type"})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction date format"})
		return
	}

	// Update transaction fields
	transaction.CategoryID = categoryID
	transaction.AccountID = accountID
	transaction.Amount = amount
	transaction.Type = transactionType
	transaction.Note = req.Note
	transaction.TransactionDate = parsedDate

	err = h.transactionUsecase.UpdateTransaction(c.Request.Context(), userID.(uuid.UUID), transaction)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	err = h.transactionUsecase.DeleteTransaction(c.Request.Context(), userID.(uuid.UUID), transactionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
