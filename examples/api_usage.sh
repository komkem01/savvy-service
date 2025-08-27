#!/bin/bash

# Savvy API Usage Examples
# ตัวอย่างการใช้งาน API ของระบบ Savvy

BASE_URL="http://localhost:8080/api/v1"

echo "🚀 Savvy API Usage Examples"
echo "================================"

# 1. Setup default categories (ครั้งแรกเท่านั้น)
echo "📂 Setting up default categories..."
curl -X POST "${BASE_URL}/setup/categories/default" \
  -H "Content-Type: application/json"

echo -e "\n"

# 2. Register a new user
echo "👤 Registering new user..."
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "demo@savvy.com",
    "password": "password123",
    "display_name": "Demo User"
  }')

echo "Register Response: $REGISTER_RESPONSE"

# Extract token from response
TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "Token: $TOKEN"

echo -e "\n"

# 3. Create an account
echo "🏦 Creating cash account..."
ACCOUNT_RESPONSE=$(curl -s -X POST "${BASE_URL}/accounts" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "กระเป๋าเงิน",
    "type": "cash"
  }')

echo "Account Response: $ACCOUNT_RESPONSE"

# Extract account ID
ACCOUNT_ID=$(echo $ACCOUNT_RESPONSE | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Account ID: $ACCOUNT_ID"

echo -e "\n"

# 4. Get all categories
echo "📂 Getting all categories..."
CATEGORIES_RESPONSE=$(curl -s -X GET "${BASE_URL}/categories" \
  -H "Authorization: Bearer $TOKEN")

echo "Categories Response: $CATEGORIES_RESPONSE"

# Extract a category ID for expense (assuming it's the first expense category)
EXPENSE_CATEGORY_ID=$(echo $CATEGORIES_RESPONSE | grep -o '"expense":\[[^]]*' | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
echo "Expense Category ID: $EXPENSE_CATEGORY_ID"

echo -e "\n"

# 5. Create some transactions
echo "💰 Creating income transaction..."
curl -s -X POST "${BASE_URL}/transactions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"category_id\": \"$(echo $CATEGORIES_RESPONSE | grep -o '"income":\[[^]]*' | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)\",
    \"account_id\": \"$ACCOUNT_ID\",
    \"amount\": \"50000.00\",
    \"type\": \"income\",
    \"note\": \"เงินเดือนประจำเดือน\",
    \"transaction_date\": \"$(date +%Y-%m-01)\"
  }"

echo -e "\n"

echo "🛍️ Creating expense transaction..."
curl -s -X POST "${BASE_URL}/transactions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"category_id\": \"$EXPENSE_CATEGORY_ID\",
    \"account_id\": \"$ACCOUNT_ID\",
    \"amount\": \"300.00\",
    \"type\": \"expense\",
    \"note\": \"ค่าอาหารเย็น\",
    \"transaction_date\": \"$(date +%Y-%m-%d)\"
  }"

echo -e "\n"

# 6. Get dashboard
echo "📊 Getting dashboard..."
DASHBOARD_RESPONSE=$(curl -s -X GET "${BASE_URL}/dashboard" \
  -H "Authorization: Bearer $TOKEN")

echo "Dashboard Response: $DASHBOARD_RESPONSE" | jq '.' 2>/dev/null || echo "$DASHBOARD_RESPONSE"

echo -e "\n"

# 7. Get recent transactions
echo "📋 Getting recent transactions..."
TRANSACTIONS_RESPONSE=$(curl -s -X GET "${BASE_URL}/transactions?limit=5" \
  -H "Authorization: Bearer $TOKEN")

echo "Recent Transactions: $TRANSACTIONS_RESPONSE" | jq '.' 2>/dev/null || echo "$TRANSACTIONS_RESPONSE"

echo -e "\n"

echo "✅ Demo completed! You can now:"
echo "   • View dashboard: GET ${BASE_URL}/dashboard"
echo "   • Add more transactions: POST ${BASE_URL}/transactions"
echo "   • View all transactions: GET ${BASE_URL}/transactions"
echo "   • Check monthly summary: GET ${BASE_URL}/dashboard/summary"
echo ""
echo "📖 See API_DOCUMENTATION.md for more details"
