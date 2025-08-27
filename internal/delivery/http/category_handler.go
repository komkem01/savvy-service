package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/usecase"
)

type CategoryHandler struct {
	categoryUsecase usecase.CategoryUsecase
}

type CreateCategoryRequest struct {
	Name     string  `json:"name" binding:"required"`
	Type     string  `json:"type" binding:"required"`
	IconName *string `json:"icon_name,omitempty"`
	ColorHex *string `json:"color_hex,omitempty"`
}

func NewCategoryHandler(categoryUsecase usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{
		categoryUsecase: categoryUsecase,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var categoryType entity.CategoryType
	switch req.Type {
	case "income":
		categoryType = entity.CategoryTypeIncome
	case "expense":
		categoryType = entity.CategoryTypeExpense
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category type"})
		return
	}

	category, err := h.categoryUsecase.CreateCategory(
		c.Request.Context(),
		userID.(uuid.UUID),
		req.Name,
		categoryType,
		req.IconName,
		req.ColorHex,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get all available categories (system + user categories)
	categories, err := h.categoryUsecase.GetAllAvailableCategories(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Group by type for easier frontend consumption
	response := gin.H{
		"income":  []*entity.Category{},
		"expense": []*entity.Category{},
	}

	for _, category := range categories {
		switch category.Type {
		case entity.CategoryTypeIncome:
			response["income"] = append(response["income"].([]*entity.Category), category)
		case entity.CategoryTypeExpense:
			response["expense"] = append(response["expense"].([]*entity.Category), category)
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *CategoryHandler) GetUserCategories(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	categories, err := h.categoryUsecase.GetCategoriesByUserID(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

func (h *CategoryHandler) ArchiveCategory(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	err = h.categoryUsecase.ArchiveCategory(c.Request.Context(), userID.(uuid.UUID), categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category archived successfully"})
}

func (h *CategoryHandler) InitializeDefaultCategories(c *gin.Context) {
	err := h.categoryUsecase.InitializeDefaultCategories(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default categories initialized successfully"})
}
