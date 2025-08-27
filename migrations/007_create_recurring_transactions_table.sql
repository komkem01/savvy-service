-- Migration: Add recurring_transactions table
-- Description: Create recurring_transactions table for automated recurring transactions

CREATE TABLE IF NOT EXISTS recurring_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    account_id UUID NOT NULL REFERENCES accounts(id) ON DELETE RESTRICT,
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    type VARCHAR(10) NOT NULL CHECK (type IN ('income', 'expense')),
    note TEXT,
    frequency VARCHAR(20) NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly', 'yearly')),
    start_date DATE NOT NULL,
    end_date DATE NULL,
    next_execution_date TIMESTAMP WITH TIME ZONE NOT NULL,
    last_execution_date TIMESTAMP WITH TIME ZONE NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    auto_execute BOOLEAN NOT NULL DEFAULT false,
    remaining_executions INTEGER NULL CHECK (remaining_executions IS NULL OR remaining_executions >= 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    -- Ensure logical constraints
    CHECK (end_date IS NULL OR end_date >= start_date),
    CHECK (last_execution_date IS NULL OR last_execution_date >= start_date)
);

-- Add indexes for better performance
CREATE INDEX idx_recurring_transactions_user_id ON recurring_transactions(user_id);
CREATE INDEX idx_recurring_transactions_category_id ON recurring_transactions(category_id);
CREATE INDEX idx_recurring_transactions_account_id ON recurring_transactions(account_id);
CREATE INDEX idx_recurring_transactions_active ON recurring_transactions(is_active);
CREATE INDEX idx_recurring_transactions_frequency ON recurring_transactions(frequency);
CREATE INDEX idx_recurring_transactions_next_execution ON recurring_transactions(next_execution_date);
CREATE INDEX idx_recurring_transactions_due ON recurring_transactions(next_execution_date, is_active) 
    WHERE is_active = true;

-- Add comments
COMMENT ON TABLE recurring_transactions IS 'Recurring transactions that can be automatically executed';
COMMENT ON COLUMN recurring_transactions.frequency IS 'How often the transaction repeats: daily, weekly, monthly, yearly';
COMMENT ON COLUMN recurring_transactions.next_execution_date IS 'When the transaction should next be executed';
COMMENT ON COLUMN recurring_transactions.last_execution_date IS 'When the transaction was last executed';
COMMENT ON COLUMN recurring_transactions.auto_execute IS 'Whether to automatically create transactions or require manual approval';
COMMENT ON COLUMN recurring_transactions.remaining_executions IS 'Number of executions remaining (NULL for unlimited)';
COMMENT ON COLUMN recurring_transactions.end_date IS 'When the recurring transaction should stop (NULL for indefinite)';
