# Savvy Backend API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
ใช้ JWT Token ในรูปแบบ Bearer Token ใน Header:
```
Authorization: Bearer <your_jwt_token>
```

---

## 📱 ฟังก์ชัน Must-Have Features (MVP)

### 1. 🔐 ระบบสมาชิก (Authentication)

#### สมัครสมาชิก
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "display_name": "John Doe" // optional
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "display_name": "John Doe",
    "currency_preference": "THB",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### เข้าสู่ระบบ
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### ต่ออายุ Token
```http
POST /auth/refresh
Authorization: Bearer <current_token>
```

---

### 2. 💰 การจัดการธุรกรรม (Transaction Management)

#### เพิ่มรายการใหม่
```http
POST /transactions
Authorization: Bearer <token>
Content-Type: application/json

{
  "category_id": "uuid",
  "account_id": "uuid", 
  "amount": "500.00",
  "type": "expense", // "income" หรือ "expense"
  "note": "ค่าอาหารเย็น", // optional
  "transaction_date": "2024-01-01"
}
```

#### ดูรายการทั้งหมด (มี Filter)
```http
GET /transactions?limit=20&offset=0&type=expense&start_date=2024-01-01&end_date=2024-01-31
Authorization: Bearer <token>
```

**Query Parameters:**
- `limit`: จำนวนรายการต่อหน้า (default: 20, max: 100)
- `offset`: เริ่มต้นจากรายการที่ (default: 0)
- `type`: `income` หรือ `expense`
- `account_id`: UUID ของบัญชี
- `category_id`: UUID ของหมวดหมู่
- `start_date`: วันที่เริ่มต้น (YYYY-MM-DD)
- `end_date`: วันที่สิ้นสุด (YYYY-MM-DD)

**Response:**
```json
{
  "transactions": [
    {
      "id": "uuid",
      "amount": "500.00",
      "type": "expense",
      "note": "ค่าอาหารเย็น",
      "transaction_date": "2024-01-01T00:00:00Z",
      "created_at": "2024-01-01T10:30:00Z"
    }
  ],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "count": 15
  }
}
```

#### แก้ไขรายการ
```http
PUT /transactions/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "category_id": "uuid",
  "account_id": "uuid",
  "amount": "600.00",
  "type": "expense",
  "note": "ค่าอาหารเย็น (แก้ไข)",
  "transaction_date": "2024-01-01"
}
```

#### ลบรายการ
```http
DELETE /transactions/:id
Authorization: Bearer <token>
```

---

### 3. 📊 หน้าสรุปผล (Dashboard)

#### ดู Dashboard แบบเบเสิก
```http
GET /dashboard
Authorization: Bearer <token>
```

**Response:**
```json
{
  "monthly_summary": {
    "year": 2024,
    "month": 1,
    "total_income": "50000.00",
    "total_expense": "25000.00", 
    "balance": "25000.00",
    "start_date": "2024-01-01T00:00:00Z",
    "end_date": "2024-01-31T23:59:59Z"
  },
  "recent_transactions": [
    {
      "transaction": { /* transaction object */ },
      "category_name": "อาหาร",
      "account_name": "บัญชีเงินสด"
    }
  ],
  "spending_by_category": [
    {
      "category_id": "uuid",
      "category_name": "อาหาร",
      "amount": "5000.00",
      "icon_name": "🍽️",
      "color_hex": "#FF6B6B"
    }
  ]
}
```

#### สรุปยอดประจำเดือน
```http
GET /dashboard/summary
Authorization: Bearer <token>

# หรือระบุเดือน
GET /dashboard/summary/monthly?year=2024&month=1
Authorization: Bearer <token>
```

#### ธุรกรรมล่าสุด
```http
GET /dashboard/transactions/recent?limit=5
Authorization: Bearer <token>
```

#### รายจ่ายแยกตามหมวดหมู่
```http
GET /dashboard/spending/category?year=2024&month=1
Authorization: Bearer <token>
```

---

### 4. 📂 การจัดการหมวดหมู่ (Categories)

#### ดูหมวดหมู่ทั้งหมด (รวมระบบ + ส่วนตัว)
```http
GET /categories
Authorization: Bearer <token>
```

**Response:**
```json
{
  "income": [
    {
      "id": "uuid",
      "name": "เงินเดือน",
      "type": "income",
      "icon_name": "💰",
      "color_hex": "#00D2D3",
      "user_id": null // null = หมวดหมู่ระบบ
    }
  ],
  "expense": [
    {
      "id": "uuid", 
      "name": "อาหาร",
      "type": "expense",
      "icon_name": "🍽️",
      "color_hex": "#FF6B6B",
      "user_id": null
    }
  ]
}
```

#### สร้างหมวดหมู่ส่วนตัว
```http
POST /categories
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "ค่าเดินทางทำงาน",
  "type": "expense", // "income" หรือ "expense"
  "icon_name": "🚌", // optional
  "color_hex": "#123456" // optional
}
```

#### เก็บหมวดหมู่ (Archive)
```http
PUT /categories/:id/archive
Authorization: Bearer <token>
```

---

### 5. 🏦 การจัดการบัญชี (Accounts)

#### สร้างบัญชีใหม่
```http
POST /accounts
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "กระเป๋าเงิน",
  "type": "cash" // "cash", "bank", "credit", "savings"
}
```

#### ดูบัญชีทั้งหมด
```http
GET /accounts
Authorization: Bearer <token>
```

#### ลบบัญชี
```http
DELETE /accounts/:id
Authorization: Bearer <token>
```

---

## 🔧 Setup & Admin Endpoints

#### สร้างหมวดหมู่เริ่มต้น (ครั้งแรกเท่านั้น)
```http
POST /setup/categories/default
```

---

## 🚨 Error Responses

ทุก API จะ return error ในรูปแบบ:
```json
{
  "error": "รายละเอียดข้อผิดพลาด"
}
```

**HTTP Status Codes:**
- `200`: สำเร็จ
- `201`: สร้างสำเร็จ  
- `400`: ข้อมูลไม่ถูกต้อง
- `401`: ไม่ได้ระบุตัวตน
- `403`: ไม่มีสิทธิ์เข้าถึง
- `404`: ไม่พบข้อมูล
- `500`: ข้อผิดพลาดของเซิร์ฟเวอร์

---

## 📅 รูปแบบวันที่

- **Date**: `YYYY-MM-DD` (เช่น `2024-01-01`)
- **DateTime**: RFC3339 format (เช่น `2024-01-01T10:30:00Z`)

---

## 💡 Tips สำหรับการใช้งาน

1. **สร้างบัญชีก่อน**: ก่อนบันทึกธุรกรรม ต้องมีบัญชีอย่างน้อย 1 บัญชี
2. **หมวดหมู่เริ่มต้น**: เรียก `/setup/categories/default` ครั้งแรกเพื่อสร้างหมวดหมู่พื้นฐาน  
3. **Pagination**: ใช้ `limit` และ `offset` สำหรับ pagination
4. **Date Filter**: ใช้ `start_date` และ `end_date` เพื่อกรองธุรกรรมตามช่วงวันที่
5. **จำนวนเงิน**: ใช้ string เพื่อความแม่นยำ (เช่น `"123.45"`)

---

## 🔄 Workflow ที่แนะนำ

1. **สมัครสมาชิก/เข้าสู่ระบบ** → ได้ JWT Token
2. **สร้างหมวดหมู่เริ่มต้น** → `/setup/categories/default`
3. **สร้างบัญชี** → อย่างน้อย 1 บัญชี
4. **บันทึกธุรกรรม** → เริ่มใช้งานได้
5. **ดู Dashboard** → ติดตามสรุปผล
