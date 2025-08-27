-- Insert default income categories
INSERT INTO categories (id, user_id, name, type, icon_name, color_hex, created_at, updated_at) VALUES
    (gen_random_uuid(), NULL, 'เงินเดือน', 'income', '💰', '#00D2D3', NOW(), NOW()),
    (gen_random_uuid(), NULL, 'โบนัส', 'income', '🎁', '#FF9F43', NOW(), NOW()),
    (gen_random_uuid(), NULL, 'ลงทุน', 'income', '📈', '#5f27cd', NOW(), NOW()),
    (gen_random_uuid(), NULL, 'ธุรกิจ', 'income', '🏢', '#00d2d3', NOW(), NOW()),
    (gen_random_uuid(), NULL, 'อื่นๆ', 'income', '💵', '#2ed573', NOW(), NOW());

-- Insert default expense categories  
INSERT INTO categories (id, user_id, name, type, icon_name, color_hex, created_at, updated_at) VALUES
    (gen_random_uuid(), NULL, 'อาหาร', 'expense', '🍽️', '#FF6B6B', NOW(), NOW()),
    (gen_random_uuid(), NULL, 'เดินทาง', 'expense', '🚗', '#4ECDC4', NOW(), NOW()),
    (gen_random_uuid(), NULL, 'ช้อปปิ้ง', 'expense', '🛍️', '#45B7D1', NOW(), NOW()),
    (gen_random_uuid(), NULL, 'บันเทิง', 'expense', '🎬', '#96CEB4', NOW(), NOW()),
    (gen_random_uuid(), NULL, 'สุขภาพ', 'expense', '🏥', '#FECA57', NOW(), NOW()),
    (gen_random_uuid(), NULL, 'การศึกษา', 'expense', '📚', '#FF9FF3', NOW(), NOW()),
    (gen_random_uuid(), NULL, 'บิล/ค่าใช้จ่าย', 'expense', '💡', '#54A0FF', NOW(), NOW()),
    (gen_random_uuid(), NULL, 'อื่นๆ', 'expense', '📝', '#5F27CD', NOW(), NOW());
