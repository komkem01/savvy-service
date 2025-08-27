package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"savvy-backend/internal/domain/entity"
	"savvy-backend/internal/domain/repository"

	"github.com/google/uuid"
)

type AIInsightUsecase interface {
	GenerateSpendingAnomalyInsights(ctx context.Context, userID uuid.UUID) ([]*entity.Insight, error)
	GenerateSpendingPatternInsights(ctx context.Context, userID uuid.UUID) ([]*entity.Insight, error)
	GenerateCategoryRecommendations(ctx context.Context, userID uuid.UUID, transactionNote string) ([]*entity.Insight, error)
	GenerateSavingsRecommendations(ctx context.Context, userID uuid.UUID) ([]*entity.Insight, error)
	ProcessWeeklyInsights(ctx context.Context, userID uuid.UUID) error
	ProcessAllUsersInsights(ctx context.Context) error
}

type aiInsightUsecase struct {
	insightRepo     repository.InsightRepository
	transactionRepo repository.TransactionRepository
	categoryRepo    repository.CategoryRepository
	budgetRepo      repository.BudgetRepository
}

func NewAIInsightUsecase(
	insightRepo repository.InsightRepository,
	transactionRepo repository.TransactionRepository,
	categoryRepo repository.CategoryRepository,
	budgetRepo repository.BudgetRepository,
) AIInsightUsecase {
	return &aiInsightUsecase{
		insightRepo:     insightRepo,
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
		budgetRepo:      budgetRepo,
	}
}

func (a *aiInsightUsecase) GenerateSpendingAnomalyInsights(ctx context.Context, userID uuid.UUID) ([]*entity.Insight, error) {
	// Look for spending anomalies in the last 6 months
	anomalies, err := a.insightRepo.GetSpendingAnomalies(ctx, userID, 6)
	if err != nil {
		return nil, err
	}

	var insights []*entity.Insight

	for _, anomaly := range anomalies {
		var priority entity.InsightPriority
		var title, message string

		switch anomaly.Severity {
		case "high":
			priority = entity.InsightPriorityHigh
			title = fmt.Sprintf("🚨 การใช้จ่าย %s เพิ่มขึ้นอย่างมาก!", anomaly.CategoryName)
			currentAmount, _ := anomaly.CurrentAmount.Float64()
			avgAmount, _ := anomaly.AverageAmount.Float64()
			message = fmt.Sprintf("เดือนนี้คุณใช้จ่ายในหมวด '%s' ถึง %.2f บาท เพิ่มขึ้น %.1f%% จากค่าเฉลี่ย %.2f บาท ลองตรวจสอบว่ามีรายจ่ายผิดปกติหรือไม่",
				anomaly.CategoryName, currentAmount, anomaly.PercentageIncrease, avgAmount)
		case "medium":
			priority = entity.InsightPriorityMedium
			title = fmt.Sprintf("⚠️ การใช้จ่าย %s เพิ่มขึ้น", anomaly.CategoryName)
			currentAmount, _ := anomaly.CurrentAmount.Float64()
			avgAmount, _ := anomaly.AverageAmount.Float64()
			message = fmt.Sprintf("เดือนนี้คุณใช้จ่ายในหมวด '%s' %.2f บาท เพิ่มขึ้น %.1f%% จากค่าเฉลี่ย %.2f บาท",
				anomaly.CategoryName, currentAmount, anomaly.PercentageIncrease, avgAmount)
		default:
			priority = entity.InsightPriorityLow
			title = fmt.Sprintf("💡 การใช้จ่าย %s เปลี่ยนแปลง", anomaly.CategoryName)
			currentAmount, _ := anomaly.CurrentAmount.Float64()
			avgAmount, _ := anomaly.AverageAmount.Float64()
			message = fmt.Sprintf("เดือนนี้คุณใช้จ่ายในหมวด '%s' %.2f บาท เพิ่มขึ้น %.1f%% จากค่าเฉลี่ย %.2f บาท",
				anomaly.CategoryName, currentAmount, anomaly.PercentageIncrease, avgAmount)
		}

		insight := entity.NewAdvancedInsight(userID, entity.InsightTypeAnomalyDetection, priority, title, message)
		insight.RelatedEntityID = &anomaly.CategoryID
		insight.RelatedEntityType = &[]string{"category"}[0]

		// Convert anomaly to JSON for related data
		anomalyJSON, _ := json.Marshal(anomaly)
		insight.RelatedData = anomalyJSON

		validUntil := time.Now().AddDate(0, 0, 7) // Valid for 7 days
		insight.ValidUntil = &validUntil

		insights = append(insights, insight)
	}

	// Save insights to database
	for _, insight := range insights {
		err := a.insightRepo.Create(ctx, insight)
		if err != nil {
			continue // Log error but don't fail
		}
	}

	return insights, nil
}

func (a *aiInsightUsecase) GenerateSpendingPatternInsights(ctx context.Context, userID uuid.UUID) ([]*entity.Insight, error) {
	// Analyze spending patterns in the last 30 days
	patterns, err := a.insightRepo.GetSpendingPatterns(ctx, userID, 30)
	if err != nil {
		return nil, err
	}

	var insights []*entity.Insight

	for _, pattern := range patterns {
		if pattern.FrequencyCount >= 5 { // Only show patterns with significant frequency
			title := fmt.Sprintf("📊 พฤติกรรมการใช้จ่าย %s", pattern.CategoryName)
			avgAmount, _ := pattern.AverageAmount.Float64()
			dayOfWeek := strings.TrimSpace(pattern.DayOfWeek)

			message := fmt.Sprintf("คุณมักใช้จ่ายหมวด '%s' ในช่วง%s ของวัน%s โดยเฉลี่ย %.2f บาท/ครั้ง (%d ครั้งในเดือนที่ผ่านมา)",
				pattern.CategoryName, pattern.TimeOfDay, dayOfWeek, avgAmount, pattern.FrequencyCount)

			insight := entity.NewAdvancedInsight(userID, entity.InsightTypeSpendingPattern, entity.InsightPriorityLow, title, message)
			insight.RelatedEntityID = &pattern.CategoryID
			insight.RelatedEntityType = &[]string{"category"}[0]

			// Convert pattern to JSON for related data
			patternJSON, _ := json.Marshal(pattern)
			insight.RelatedData = patternJSON

			validUntil := time.Now().AddDate(0, 0, 14) // Valid for 14 days
			insight.ValidUntil = &validUntil

			insights = append(insights, insight)
		}
	}

	// Save insights to database
	for _, insight := range insights {
		err := a.insightRepo.Create(ctx, insight)
		if err != nil {
			continue
		}
	}

	return insights, nil
}

func (a *aiInsightUsecase) GenerateCategoryRecommendations(ctx context.Context, userID uuid.UUID, transactionNote string) ([]*entity.Insight, error) {
	if transactionNote == "" {
		return nil, nil
	}

	// Simple keyword-based category recommendation
	// In a real AI system, this would use NLP/ML models
	categoryKeywords := map[string][]string{
		"อาหาร":    {"กาแฟ", "ข้าว", "อาหาร", "ร้านอาหาร", "เซเว่น", "แมค", "kfc", "starbucks", "coffee", "food"},
		"เดินทาง":  {"grab", "taxi", "รถไฟ", "bts", "mrt", "น้ำมัน", "ปตท", "shell", "uber", "transport"},
		"ช้อปปิ้ง": {"เสื้อผ้า", "รองเท้า", "กระเป๋า", "shopping", "mall", "เซ็นทรัล", "robinson", "shopee", "lazada"},
		"บันเทิง":  {"หนัง", "เกม", "คอนเสิร์ต", "netflix", "spotify", "cinema", "game", "entertainment"},
		"สุขภาพ":   {"โรงพยาบาล", "คลินิก", "ยา", "วิตามิน", "ออกกำลังกาย", "fitness", "hospital", "pharmacy"},
		"การศึกษา": {"หนังสือ", "เรียน", "course", "คอร์ส", "udemy", "school", "university", "education"},
		"บิล":      {"ไฟฟ้า", "น้ำ", "เน็ต", "true", "ais", "dtac", "electric", "water", "internet", "bill"},
	}

	note := strings.ToLower(transactionNote)
	var suggestedCategories []string

	for categoryName, keywords := range categoryKeywords {
		for _, keyword := range keywords {
			if strings.Contains(note, strings.ToLower(keyword)) {
				suggestedCategories = append(suggestedCategories, categoryName)
				break
			}
		}
	}

	var insights []*entity.Insight

	if len(suggestedCategories) > 0 {
		title := "💡 คำแนะนำหมวดหมู่"
		message := fmt.Sprintf("จากรายการ '%s' ระบบแนะนำหมวดหมู่: %s",
			transactionNote, strings.Join(suggestedCategories, ", "))

		insight := entity.NewAdvancedInsight(userID, entity.InsightTypeCategoryRecommendation, entity.InsightPriorityLow, title, message)

		// Add suggested categories to related data
		suggestionData := map[string]interface{}{
			"original_note":        transactionNote,
			"suggested_categories": suggestedCategories,
		}
		suggestionJSON, _ := json.Marshal(suggestionData)
		insight.RelatedData = suggestionJSON

		validUntil := time.Now().AddDate(0, 0, 1) // Valid for 1 day
		insight.ValidUntil = &validUntil

		insights = append(insights, insight)

		// Save insight to database
		err := a.insightRepo.Create(ctx, insight)
		if err != nil {
			return nil, err
		}
	}

	return insights, nil
}

func (a *aiInsightUsecase) GenerateSavingsRecommendations(ctx context.Context, userID uuid.UUID) ([]*entity.Insight, error) {
	// Get spending patterns to generate savings recommendations
	patterns, err := a.insightRepo.GetSpendingPatterns(ctx, userID, 30)
	if err != nil {
		return nil, err
	}

	var insights []*entity.Insight

	// Find categories with high frequency and suggest savings
	for _, pattern := range patterns {
		if pattern.FrequencyCount >= 10 { // High frequency spending
			avgAmount, _ := pattern.AverageAmount.Float64()
			totalSpent := avgAmount * float64(pattern.FrequencyCount)

			if avgAmount > 100 { // Only for significant amounts
				potentialSavings := avgAmount * 0.2 // Assume 20% savings potential
				monthlySavings := potentialSavings * float64(pattern.FrequencyCount)

				title := fmt.Sprintf("💰 โอกาสประหยัด %s", pattern.CategoryName)
				message := fmt.Sprintf("คุณใช้จ่าย '%s' บ่อย (%d ครั้ง/เดือน) เฉลี่ย %.2f บาท/ครั้ง หากลดลง 20%% จะประหยัดได้ %.2f บาท/เดือน",
					pattern.CategoryName, pattern.FrequencyCount, avgAmount, monthlySavings)

				insight := entity.NewAdvancedInsight(userID, entity.InsightTypeSavingsRecommendation, entity.InsightPriorityMedium, title, message)
				insight.RelatedEntityID = &pattern.CategoryID
				insight.RelatedEntityType = &[]string{"category"}[0]

				recommendationData := map[string]interface{}{
					"category_name":     pattern.CategoryName,
					"current_avg":       avgAmount,
					"frequency":         pattern.FrequencyCount,
					"total_spent":       totalSpent,
					"potential_savings": monthlySavings,
					"reduction_percent": 20,
				}
				recommendationJSON, _ := json.Marshal(recommendationData)
				insight.RelatedData = recommendationJSON

				validUntil := time.Now().AddDate(0, 0, 14) // Valid for 14 days
				insight.ValidUntil = &validUntil

				insights = append(insights, insight)
			}
		}
	}

	// Save insights to database
	for _, insight := range insights {
		err := a.insightRepo.Create(ctx, insight)
		if err != nil {
			continue
		}
	}

	return insights, nil
}

func (a *aiInsightUsecase) ProcessWeeklyInsights(ctx context.Context, userID uuid.UUID) error {
	// Generate comprehensive weekly insights for a user
	_, err := a.GenerateSpendingAnomalyInsights(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to generate anomaly insights: %w", err)
	}

	_, err = a.GenerateSpendingPatternInsights(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to generate pattern insights: %w", err)
	}

	_, err = a.GenerateSavingsRecommendations(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to generate savings recommendations: %w", err)
	}

	return nil
}

func (a *aiInsightUsecase) ProcessAllUsersInsights(ctx context.Context) error {
	// This method would typically be called by a scheduled job
	// For now, we'll just return nil as we don't have a user repository
	// In a real implementation, you'd fetch all active users and process insights for each
	return nil
}
