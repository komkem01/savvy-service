package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/usecase"
)

type AIInsightHandler struct {
	aiInsightUsecase usecase.AIInsightUsecase
}

func NewAIInsightHandler(aiInsightUsecase usecase.AIInsightUsecase) *AIInsightHandler {
	return &AIInsightHandler{
		aiInsightUsecase: aiInsightUsecase,
	}
}

type InsightResponse struct {
	ID                string                 `json:"id"`
	UserID            string                 `json:"user_id"`
	Type              string                 `json:"type"`
	Priority          string                 `json:"priority"`
	Title             string                 `json:"title"`
	Content           string                 `json:"content"`
	ActionText        *string                `json:"action_text,omitempty"`
	IsRead            bool                   `json:"is_read"`
	RelatedEntityID   *string                `json:"related_entity_id,omitempty"`
	RelatedEntityType *string                `json:"related_entity_type,omitempty"`
	RelatedData       map[string]interface{} `json:"related_data,omitempty"`
	ValidUntil        *string                `json:"valid_until,omitempty"`
	CreatedAt         string                 `json:"created_at"`
	UpdatedAt         string                 `json:"updated_at"`
}

type SpendingAnomalyResponse struct {
	CategoryID         string  `json:"category_id"`
	CategoryName       string  `json:"category_name"`
	CurrentAmount      string  `json:"current_amount"`
	AverageAmount      string  `json:"average_amount"`
	PercentageIncrease float64 `json:"percentage_increase"`
	Period             string  `json:"period"`
	Severity           string  `json:"severity"`
}

type SpendingPatternResponse struct {
	CategoryID     string `json:"category_id"`
	CategoryName   string `json:"category_name"`
	DayOfWeek      string `json:"day_of_week"`
	TimeOfDay      string `json:"time_of_day"`
	FrequencyCount int    `json:"frequency_count"`
	AverageAmount  string `json:"average_amount"`
	Period         string `json:"period"`
}

func (h *AIInsightHandler) GenerateSpendingAnomalyInsights(c *gin.Context) {
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

	insights, err := h.aiInsightUsecase.GenerateSpendingAnomalyInsights(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]*InsightResponse, len(insights))
	for i, insight := range insights {
		responses[i] = h.insightToResponse(insight)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *AIInsightHandler) GenerateSpendingPatternInsights(c *gin.Context) {
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

	insights, err := h.aiInsightUsecase.GenerateSpendingPatternInsights(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]*InsightResponse, len(insights))
	for i, insight := range insights {
		responses[i] = h.insightToResponse(insight)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *AIInsightHandler) GenerateCategoryRecommendations(c *gin.Context) {
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

	// Get transaction note from request body or query param
	transactionNote := c.Query("transaction_note")
	if transactionNote == "" {
		var requestBody struct {
			TransactionNote string `json:"transaction_note"`
		}
		if err := c.ShouldBindJSON(&requestBody); err == nil {
			transactionNote = requestBody.TransactionNote
		}
	}

	insights, err := h.aiInsightUsecase.GenerateCategoryRecommendations(c.Request.Context(), userUUID, transactionNote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]*InsightResponse, len(insights))
	for i, insight := range insights {
		responses[i] = h.insightToResponse(insight)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *AIInsightHandler) GenerateSavingsRecommendations(c *gin.Context) {
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

	insights, err := h.aiInsightUsecase.GenerateSavingsRecommendations(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]*InsightResponse, len(insights))
	for i, insight := range insights {
		responses[i] = h.insightToResponse(insight)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *AIInsightHandler) ProcessWeeklyInsights(c *gin.Context) {
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

	err = h.aiInsightUsecase.ProcessWeeklyInsights(c.Request.Context(), userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Weekly insights processed successfully",
	})
}

func (h *AIInsightHandler) ProcessAllUsersInsights(c *gin.Context) {
	// This endpoint might be restricted to admin users in production
	err := h.aiInsightUsecase.ProcessAllUsersInsights(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "All users insights processed successfully",
	})
}

func (h *AIInsightHandler) GetSpendingAnomalies(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	_, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Note: This uses the repository method directly since the usecase doesn't expose it
	// In a real implementation, you might want to add this to the usecase interface
	c.JSON(http.StatusOK, gin.H{
		"message":    "Use the insight generation endpoints to get anomaly insights",
		"suggestion": "Call /ai-insights/spending-anomalies/generate instead",
	})
}

func (h *AIInsightHandler) GetSpendingPatterns(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	_, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Note: This uses the repository method directly since the usecase doesn't expose it
	// In a real implementation, you might want to add this to the usecase interface
	c.JSON(http.StatusOK, gin.H{
		"message":    "Use the insight generation endpoints to get pattern insights",
		"suggestion": "Call /ai-insights/spending-patterns/generate instead",
	})
}

func (h *AIInsightHandler) insightToResponse(insight *entity.Insight) *InsightResponse {
	response := &InsightResponse{
		ID:        insight.ID.String(),
		UserID:    insight.UserID.String(),
		Type:      string(insight.Type),
		Priority:  string(insight.Priority),
		Title:     insight.Title,
		Content:   insight.Content,
		IsRead:    insight.IsRead,
		CreatedAt: insight.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: insight.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if insight.ActionText != nil {
		response.ActionText = insight.ActionText
	}

	if insight.RelatedEntityID != nil {
		entityID := insight.RelatedEntityID.String()
		response.RelatedEntityID = &entityID
	}

	if insight.RelatedEntityType != nil {
		response.RelatedEntityType = insight.RelatedEntityType
	}

	if insight.RelatedData != nil {
		var relatedData map[string]interface{}
		if err := json.Unmarshal(insight.RelatedData, &relatedData); err == nil {
			response.RelatedData = relatedData
		}
	}

	if insight.ValidUntil != nil {
		validUntil := insight.ValidUntil.Format("2006-01-02T15:04:05Z")
		response.ValidUntil = &validUntil
	}

	return response
}

func (h *AIInsightHandler) anomalyToResponse(anomaly *entity.SpendingAnomaly) *SpendingAnomalyResponse {
	return &SpendingAnomalyResponse{
		CategoryID:         anomaly.CategoryID.String(),
		CategoryName:       anomaly.CategoryName,
		CurrentAmount:      anomaly.CurrentAmount.String(),
		AverageAmount:      anomaly.AverageAmount.String(),
		PercentageIncrease: anomaly.PercentageIncrease,
		Period:             anomaly.Period,
		Severity:           anomaly.Severity,
	}
}

func (h *AIInsightHandler) patternToResponse(pattern *entity.SpendingPattern) *SpendingPatternResponse {
	return &SpendingPatternResponse{
		CategoryID:     pattern.CategoryID.String(),
		CategoryName:   pattern.CategoryName,
		DayOfWeek:      pattern.DayOfWeek,
		TimeOfDay:      pattern.TimeOfDay,
		FrequencyCount: pattern.FrequencyCount,
		AverageAmount:  pattern.AverageAmount.String(),
		Period:         pattern.Period,
	}
}
