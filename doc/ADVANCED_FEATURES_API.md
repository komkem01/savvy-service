# Advanced Features API Documentation

## 📊 Data Visualization APIs

### 1. Expense Pie Chart
แสดงข้อมูลค่าใช้จ่ายในรูปแบบ Pie Chart ตามหมวดหมู่

**Endpoint:** `GET /api/v1/analytics/pie/expenses`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `year` (optional): ปี เช่น 2024
- `month` (optional): เดือน (1-12)

**Response:**
```json
{
  "type": "pie",
  "data": [
    {
      "label": "อาหาร",
      "value": 15000.00,
      "color": "#FF6B6B"
    },
    {
      "label": "เดินทาง",
      "value": 8000.00,
      "color": "#4ECDC4"
    }
  ]
}
```

### 2. Income vs Expense Bar Chart
แสดงข้อมูลเปรียบเทียบรายรับ-รายจ่ายรายเดือน

**Endpoint:** `GET /api/v1/analytics/bar/income-expense`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `year` (optional): ปี เช่น 2024

**Response:**
```json
{
  "type": "bar",
  "data": [
    {
      "label": "January",
      "income": 50000.00,
      "expense": 35000.00
    },
    {
      "label": "February",
      "income": 55000.00,
      "expense": 40000.00
    }
  ]
}
```

### 3. Category Trend Line Chart
แสดงแนวโน้มการใช้จ่ายของหมวดหมู่เฉพาะ

**Endpoint:** `GET /api/v1/analytics/trend/category/{category_id}`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Path Parameters:**
- `category_id`: UUID ของหมวดหมู่

**Query Parameters:**
- `year` (optional): ปี เช่น 2024

**Response:**
```json
{
  "type": "line",
  "data": [
    {
      "month": "January",
      "amount": 5000.00
    },
    {
      "month": "February",
      "amount": 6500.00
    }
  ]
}
```

### 4. Top Categories Chart
แสดงหมวดหมู่ที่ใช้งานมากที่สุด

**Endpoint:** `GET /api/v1/analytics/top/categories`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "type": "horizontal-bar",
  "data": [
    {
      "name": "อาหาร",
      "count": 45,
      "color": "#FF6B6B"
    },
    {
      "name": "เดินทาง",
      "count": 32,
      "color": "#4ECDC4"
    }
  ]
}
```

## 🔍 Advanced Transaction Search

### Enhanced Transaction Filtering
ปรับปรุง API การดึงข้อมูลธุรกรรมให้รองรับการค้นหาขั้นสูง

**Endpoint:** `GET /api/v1/transactions`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `search` (optional): ค้นหาจาก note หรือชื่อหมวดหมู่
- `min_amount` (optional): จำนวนเงินต่ำสุด
- `max_amount` (optional): จำนวนเงินสูงสุด
- `account_id` (optional): UUID ของบัญชี
- `category_id` (optional): UUID ของหมวดหมู่
- `type` (optional): income หรือ expense
- `start_date` (optional): วันที่เริ่มต้น (YYYY-MM-DD)
- `end_date` (optional): วันที่สิ้นสุด (YYYY-MM-DD)
- `limit` (optional): จำนวนรายการต่อหน้า (default: 20, max: 100)
- `offset` (optional): เริ่มจากรายการที่ (default: 0)

**Example Request:**
```
GET /api/v1/transactions?search=กาแฟ&min_amount=100&max_amount=500&type=expense
```

**Response:**
```json
{
  "transactions": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "category_id": "uuid",
      "account_id": "uuid",
      "amount": "350.00",
      "type": "expense",
      "note": "กาแฟ Starbucks",
      "transaction_date": "2024-01-15T00:00:00Z",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  ],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "count": 1
  },
  "filters": {
    "search": "กาแฟ",
    "min_amount": "100.00",
    "max_amount": "500.00",
    "account_id": null,
    "category_id": null,
    "type": "expense",
    "start_date": null,
    "end_date": null
  }
}
```

## 📂 Advanced Category Management

### Get Categories with Filter
ดึงข้อมูลหมวดหมู่พร้อมตัวกรองขั้นสูง

**Endpoint:** `GET /api/v1/categories/filter`

**Headers:**
```
Authorization: Bearer <jwt_token>
```

**Query Parameters:**
- `is_system` (optional): true สำหรับหมวดหมู่ระบบ, false สำหรับหมวดหมู่ผู้ใช้
- `is_archived` (optional): true/false สำหรับสถานะ archive
- `type` (optional): income หรือ expense
- `search` (optional): ค้นหาจากชื่อหมวดหมู่

### Archive/Unarchive Category
จัดการสถานะ archive ของหมวดหมู่

**Archive Category:**
```
PUT /api/v1/categories/{id}/archive
```

**Unarchive Category:**
```
PUT /api/v1/categories/{id}/unarchive
```

### Category Usage Statistics
ดูสถิติการใช้งานหมวดหมู่

**Endpoint:** `GET /api/v1/categories/usage-stats`

**Response:**
```json
{
  "category_usage": {
    "uuid1": 45,
    "uuid2": 32,
    "uuid3": 28
  }
}
```

## 🎯 Use Cases สำหรับฟีเจอร์ใหม่

### 1. Dashboard Analytics
- แสดง Pie Chart ของค่าใช้จ่ายเดือนปัจจุบัน
- แสดง Bar Chart เปรียบเทียบรายรับ-รายจ่าย 12 เดือน
- แสดง Top 5 หมวดหมู่ที่ใช้มากที่สุด

### 2. Advanced Search
- ค้นหาธุรกรรมจาก "กาแฟ" หรือ "Starbucks"
- กรองธุรกรรมในช่วงเงิน 100-500 บาท
- ค้นหาค่าใช้จ่ายประเภทอาหารในเดือนมกราคม

### 3. Category Analysis
- ดูแนวโน้มการใช้จ่ายหมวดหมู่ "อาหาร" ตลอดปี
- วิเคราะห์หมวดหมู่ที่ใช้บ่อยที่สุด
- จัดการหมวดหมู่ที่ไม่ใช้แล้วโดย Archive

### 4. Financial Insights
- เปรียบเทียบรายจ่ายแต่ละเดือน
- ติดตามแนวโน้มการใช้จ่ายในหมวดหมู่สำคัญ
- วิเคราะห์ balance รายเดือนผ่าน Bar Chart

## 🔐 Security Notes

- ทุก API ต้องใช้ JWT Authentication
- ข้อมูลจะถูกกรองตาม user_id อัตโนมัติ
- ไม่สามารถเข้าถึงข้อมูลของผู้ใช้คนอื่นได้

## 📈 Performance Considerations

- การใช้ Index ในฐานข้อมูลสำหรับการค้นหา
- Pagination สำหรับข้อมูลจำนวนมาก
- Cache ข้อมูล Category ที่ใช้บ่อย
- Limit การค้นหาสูงสุด 100 รายการต่อครั้ง
