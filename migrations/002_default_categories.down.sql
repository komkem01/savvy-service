-- Delete default categories (system categories have user_id = NULL)
DELETE FROM categories WHERE user_id IS NULL;
