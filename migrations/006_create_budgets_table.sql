-- Migration: Add budgets table
-- Description: Create budgets table for budget management

CREATE TABLE IF NOT EXISTS budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    period VARCHAR(20) NOT NULL CHECK (period IN ('monthly', 'yearly')),
    start_date DATE NOT NULL,
    end_date DATE NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Ensure only one active budget per user-category combination
    UNIQUE(user_id, category_id, is_active) DEFERRABLE INITIALLY DEFERRED
);

-- Add indexes for better performance
CREATE INDEX idx_budgets_user_id ON budgets(user_id);
CREATE INDEX idx_budgets_category_id ON budgets(category_id);
CREATE INDEX idx_budgets_active ON budgets(is_active);
CREATE INDEX idx_budgets_period ON budgets(period);
CREATE INDEX idx_budgets_start_date ON budgets(start_date);

-- Add comments
COMMENT ON TABLE budgets IS 'User budgets for expense categories';
COMMENT ON COLUMN budgets.amount IS 'Budget amount in the specified period';
COMMENT ON COLUMN budgets.period IS 'Budget period: monthly or yearly';
COMMENT ON COLUMN budgets.start_date IS 'When the budget starts';
COMMENT ON COLUMN budgets.end_date IS 'When the budget ends (NULL for indefinite)';
COMMENT ON COLUMN budgets.is_active IS 'Whether the budget is currently active';
