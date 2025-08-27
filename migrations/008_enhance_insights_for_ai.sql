-- Migration: Update insights table for AI capabilities
-- Description: Enhance insights table to support AI-driven insights and recommendations

-- First, add new columns to existing insights table
ALTER TABLE insights 
ADD COLUMN IF NOT EXISTS insight_type VARCHAR(50) NOT NULL DEFAULT 'general',
ADD COLUMN IF NOT EXISTS priority VARCHAR(10) NOT NULL DEFAULT 'medium',
ADD COLUMN IF NOT EXISTS action_recommendation TEXT,
ADD COLUMN IF NOT EXISTS related_data JSONB,
ADD COLUMN IF NOT EXISTS valid_until TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS is_read BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN IF NOT EXISTS severity VARCHAR(10);

-- Add constraints for new columns
ALTER TABLE insights 
ADD CONSTRAINT check_insight_type 
    CHECK (insight_type IN ('general', 'spending_anomaly', 'spending_pattern', 'category_recommendation', 'savings_suggestion', 'budget_alert')),
ADD CONSTRAINT check_priority 
    CHECK (priority IN ('low', 'medium', 'high')),
ADD CONSTRAINT check_severity 
    CHECK (severity IS NULL OR severity IN ('low', 'medium', 'high'));

-- Create indexes for AI insights
CREATE INDEX IF NOT EXISTS idx_insights_type ON insights(insight_type);
CREATE INDEX IF NOT EXISTS idx_insights_priority ON insights(priority);
CREATE INDEX IF NOT EXISTS idx_insights_unread ON insights(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_insights_valid ON insights(valid_until) WHERE valid_until IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_insights_severity ON insights(severity) WHERE severity IS NOT NULL;

-- Create spending_anomalies table for detailed anomaly tracking
CREATE TABLE IF NOT EXISTS spending_anomalies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    anomaly_type VARCHAR(50) NOT NULL,
    severity VARCHAR(10) NOT NULL CHECK (severity IN ('low', 'medium', 'high')),
    expected_range_min DECIMAL(15,2),
    expected_range_max DECIMAL(15,2),
    actual_amount DECIMAL(15,2) NOT NULL,
    deviation_percentage DECIMAL(5,2),
    detection_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    description TEXT,
    
    UNIQUE(user_id, transaction_id, anomaly_type)
);

-- Create spending_patterns table for pattern analysis
CREATE TABLE IF NOT EXISTS spending_patterns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    pattern_type VARCHAR(50) NOT NULL,
    frequency VARCHAR(20) NOT NULL,
    average_amount DECIMAL(15,2) NOT NULL,
    trend VARCHAR(20) NOT NULL CHECK (trend IN ('increasing', 'decreasing', 'stable')),
    confidence_score DECIMAL(3,2) NOT NULL CHECK (confidence_score BETWEEN 0 AND 1),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    CHECK (end_date >= start_date)
);

-- Add indexes for performance
CREATE INDEX idx_spending_anomalies_user_id ON spending_anomalies(user_id);
CREATE INDEX idx_spending_anomalies_transaction_id ON spending_anomalies(transaction_id);
CREATE INDEX idx_spending_anomalies_severity ON spending_anomalies(severity);
CREATE INDEX idx_spending_anomalies_detection_date ON spending_anomalies(detection_date);

CREATE INDEX idx_spending_patterns_user_id ON spending_patterns(user_id);
CREATE INDEX idx_spending_patterns_category_id ON spending_patterns(category_id);
CREATE INDEX idx_spending_patterns_type ON spending_patterns(pattern_type);
CREATE INDEX idx_spending_patterns_confidence ON spending_patterns(confidence_score);

-- Add comments
COMMENT ON TABLE spending_anomalies IS 'Detected spending anomalies with severity levels';
COMMENT ON COLUMN spending_anomalies.anomaly_type IS 'Type of anomaly detected (e.g., unusual_amount, frequency_spike)';
COMMENT ON COLUMN spending_anomalies.deviation_percentage IS 'Percentage deviation from expected pattern';

COMMENT ON TABLE spending_patterns IS 'Identified spending patterns and trends';
COMMENT ON COLUMN spending_patterns.pattern_type IS 'Type of pattern (e.g., weekly_grocery, monthly_utilities)';
COMMENT ON COLUMN spending_patterns.confidence_score IS 'AI confidence in pattern identification (0-1)';
COMMENT ON COLUMN spending_patterns.trend IS 'Whether spending in this pattern is increasing, decreasing, or stable';
