# Advanced Personal Finance Features API Documentation

This document describes the advanced features added to the Savvy personal finance backend, including budget management, recurring transactions, and AI-powered insights.

## New Features Overview

### 1. การจัดการงบประมาณ (Budget Management)
- Create and manage budgets by category
- Track budget progress and spending
- Automated budget alerts when approaching or exceeding limits

### 2. รายการที่เกิดขึ้นประจำ (Recurring Transactions)
- Set up automated recurring income/expense transactions
- Flexible scheduling (daily, weekly, monthly, yearly)
- Manual or automatic execution

### 3. ระบบเป้าหมายการออมอัจฉริยะ & ผู้ช่วย AI วิเคราะห์พฤติกรรมการใช้จ่าย
- AI-powered spending anomaly detection
- Spending pattern analysis
- Category recommendations for transactions
- Savings recommendations based on spending behavior

## API Endpoints

### Budget Management APIs

#### Create Budget
```http
POST /api/v1/budgets
Authorization: Bearer <token>
Content-Type: application/json

{
  "category_id": "uuid",
  "amount": "1000.00",
  "period": "monthly"
}
```

#### Get All Budgets
```http
GET /api/v1/budgets
Authorization: Bearer <token>

Query Parameters:
- category_id: Filter by category (optional)
- period: Filter by period (monthly/yearly) (optional)
- active_only: true/false (optional)
```

#### Get Budget by ID
```http
GET /api/v1/budgets/{id}
Authorization: Bearer <token>
```

#### Update Budget
```http
PUT /api/v1/budgets/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "amount": "1200.00",
  "period": "monthly",
  "is_active": true
}
```

#### Delete Budget
```http
DELETE /api/v1/budgets/{id}
Authorization: Bearer <token>
```

#### Get Budget Progress
```http
GET /api/v1/budgets/progress?year=2024&month=8
Authorization: Bearer <token>
```

#### Get Current Month Budget Progress
```http
GET /api/v1/budgets/progress/current
Authorization: Bearer <token>
```

#### Check Budget Alerts
```http
POST /api/v1/budgets/alerts/check
Authorization: Bearer <token>
```

### Recurring Transaction APIs

#### Create Recurring Transaction
```http
POST /api/v1/recurring-transactions
Authorization: Bearer <token>
Content-Type: application/json

{
  "category_id": "uuid",
  "account_id": "uuid",
  "amount": "500.00",
  "type": "expense",
  "note": "Monthly rent payment",
  "frequency": "monthly",
  "start_date": "2024-01-01",
  "end_date": "2024-12-31",
  "auto_execute": false,
  "remaining_executions": 12
}
```

#### Get All Recurring Transactions
```http
GET /api/v1/recurring-transactions
Authorization: Bearer <token>

Query Parameters:
- frequency: Filter by frequency (daily/weekly/monthly/yearly) (optional)
- active_only: true/false (optional)
```

#### Get Recurring Transaction by ID
```http
GET /api/v1/recurring-transactions/{id}
Authorization: Bearer <token>
```

#### Update Recurring Transaction
```http
PUT /api/v1/recurring-transactions/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "amount": "550.00",
  "frequency": "monthly",
  "auto_execute": true,
  "is_active": true
}
```

#### Delete Recurring Transaction
```http
DELETE /api/v1/recurring-transactions/{id}
Authorization: Bearer <token>
```

#### Execute Recurring Transaction Manually
```http
POST /api/v1/recurring-transactions/{id}/execute
Authorization: Bearer <token>
```

#### Get Due Recurring Transactions
```http
GET /api/v1/recurring-transactions/due?limit=50
Authorization: Bearer <token>
```

### AI Insights APIs

#### Generate Spending Anomaly Insights
```http
POST /api/v1/ai-insights/spending-anomalies/generate
Authorization: Bearer <token>
```

#### Generate Spending Pattern Insights
```http
POST /api/v1/ai-insights/spending-patterns/generate
Authorization: Bearer <token>
```

#### Generate Category Recommendations
```http
POST /api/v1/ai-insights/category-recommendations/generate
Authorization: Bearer <token>
Content-Type: application/json

{
  "transaction_note": "Coffee at Starbucks"
}
```

Or via query parameter:
```http
POST /api/v1/ai-insights/category-recommendations/generate?transaction_note=Coffee+at+Starbucks
Authorization: Bearer <token>
```

#### Generate Savings Recommendations
```http
POST /api/v1/ai-insights/savings-recommendations/generate
Authorization: Bearer <token>
```

#### Process Weekly Insights for User
```http
POST /api/v1/ai-insights/weekly/process
Authorization: Bearer <token>
```

#### Get Spending Anomalies (Information Endpoint)
```http
GET /api/v1/ai-insights/spending-anomalies
Authorization: Bearer <token>
```

#### Get Spending Patterns (Information Endpoint)
```http
GET /api/v1/ai-insights/spending-patterns
Authorization: Bearer <token>
```

### System Administration APIs

#### Process All Due Recurring Transactions (System-wide)
```http
POST /api/v1/system/recurring-transactions/process-all
```
*Note: Should be protected with admin authentication in production*

#### Process AI Insights for All Users
```http
POST /api/v1/system/ai-insights/process-all
```
*Note: Should be protected with admin authentication in production*

## Response Examples

### Budget Response
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "category_id": "uuid",
  "amount": "1000.00",
  "period": "monthly",
  "is_active": true,
  "created_at": "2024-08-28T10:00:00Z",
  "updated_at": "2024-08-28T10:00:00Z"
}
```

### Budget Progress Response
```json
[
  {
    "budget_id": "uuid",
    "spent_amount": "750.00",
    "remaining_amount": "250.00",
    "percentage_used": 75.0,
    "is_over_budget": false,
    "days_remaining": 10,
    "average_daily": "25.00",
    "projected_total": "1000.00"
  }
]
```

### Recurring Transaction Response
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "category_id": "uuid",
  "account_id": "uuid",
  "amount": "500.00",
  "type": "expense",
  "note": "Monthly rent payment",
  "frequency": "monthly",
  "start_date": "2024-01-01",
  "end_date": "2024-12-31",
  "next_execution_date": "2024-09-01T00:00:00Z",
  "last_execution_date": "2024-08-01T00:00:00Z",
  "is_active": true,
  "auto_execute": false,
  "remaining_executions": 4,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-08-01T10:00:00Z"
}
```

### AI Insight Response
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "type": "spending_anomaly",
  "priority": "high",
  "title": "Unusual Spending Detected",
  "content": "Your grocery spending this month is 150% higher than usual",
  "action_text": "Review your recent grocery purchases",
  "is_read": false,
  "related_entity_id": "category_uuid",
  "related_entity_type": "category",
  "related_data": {
    "category_name": "Groceries",
    "current_amount": "450.00",
    "average_amount": "300.00",
    "percentage_increase": 50.0
  },
  "valid_until": "2024-09-28T10:00:00Z",
  "created_at": "2024-08-28T10:00:00Z",
  "updated_at": "2024-08-28T10:00:00Z"
}
```

## Database Migrations

Three new migration files have been created:

1. `006_create_budgets_table.sql` - Creates budgets table with progress tracking
2. `007_create_recurring_transactions_table.sql` - Creates recurring transactions table
3. `008_enhance_insights_for_ai.sql` - Enhances insights table and adds anomaly/pattern tables

Run migrations:
```bash
# Apply all pending migrations
migrate -path migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up
```

## Error Handling

All endpoints return standard HTTP status codes:

- `200 OK` - Success
- `201 Created` - Resource created successfully
- `204 No Content` - Success with no content (e.g., delete operations)
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

Error responses include descriptive messages:
```json
{
  "error": "Invalid amount format"
}
```

## Authentication

All protected endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

The token is obtained through the existing authentication endpoints:
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/register`

## Rate Limiting & Performance

- AI insight generation endpoints may take longer to process due to analysis complexity
- Consider implementing rate limiting for compute-intensive AI endpoints
- Budget progress calculations are optimized with database indexes
- Recurring transaction processing is designed for batch execution

## Production Considerations

1. **Security**: System administration endpoints should be protected with admin authentication
2. **Scheduling**: Set up cron jobs or background workers for:
   - Processing recurring transactions daily
   - Generating weekly AI insights
   - Checking budget alerts
3. **Monitoring**: Monitor AI insight generation performance and accuracy
4. **Data Privacy**: Ensure AI analysis respects user privacy and data retention policies
