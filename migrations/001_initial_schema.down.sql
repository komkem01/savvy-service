-- Drop triggers
DROP TRIGGER IF EXISTS update_savings_goals_updated_at ON savings_goals;
DROP TRIGGER IF EXISTS update_transactions_updated_at ON transactions;
DROP TRIGGER IF EXISTS update_categories_updated_at ON categories;
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_insights_is_read;
DROP INDEX IF EXISTS idx_insights_user_id;
DROP INDEX IF EXISTS idx_goal_deposits_user_id;
DROP INDEX IF EXISTS idx_goal_deposits_goal_id;
DROP INDEX IF EXISTS idx_savings_goals_user_id;
DROP INDEX IF EXISTS idx_transactions_user_date_type;
DROP INDEX IF EXISTS idx_transactions_type;
DROP INDEX IF EXISTS idx_transactions_account_id;
DROP INDEX IF EXISTS idx_transactions_category_id;
DROP INDEX IF EXISTS idx_transactions_date;
DROP INDEX IF EXISTS idx_transactions_user_id;
DROP INDEX IF EXISTS idx_categories_type;
DROP INDEX IF EXISTS idx_categories_user_id;
DROP INDEX IF EXISTS idx_accounts_user_id;

-- Drop tables in reverse order of dependencies
DROP TABLE IF EXISTS insights;
DROP TABLE IF EXISTS goal_deposits;
DROP TABLE IF EXISTS savings_goals;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS users;

-- Drop enums
DROP TYPE IF EXISTS insight_type;
DROP TYPE IF EXISTS goal_status;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS category_type;
DROP TYPE IF EXISTS account_type;
