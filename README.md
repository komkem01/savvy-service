# Savvy Backend API

ระบบจัดการการเงินส่วนบุคคล (Personal Finance Management System) ที่พัฒนาด้วย Go และ PostgreSQL

## โครงสร้างโปรเจกต์

```
backend/
├── cmd/
│   └── server/              # Entry point ของแอปพลิเคชัน
│       └── main.go
├── internal/
│   ├── config/              # Configuration management
│   ├── domain/
│   │   ├── entity/          # Domain entities (User, Account, Transaction, etc.)
│   │   └── repository/      # Repository interfaces
│   ├── usecase/             # Business logic layer
│   ├── delivery/
│   │   └── http/            # HTTP handlers และ routes
│   └── infrastructure/
│       └── database/        # Database connections และ implementations
├── pkg/
│   └── utils/               # Shared utilities
├── go.mod
├── .env.example
└── README.md
```

## Features

### ✅ ที่เสร็จแล้ว (Must-Have MVP Features)

#### 🔐 ระบบสมาชิก (User Authentication) 
- ✅ สมัครสมาชิก/เข้าสู่ระบบ
- ✅ JWT Authentication & Token refresh
- ✅ Password hashing ด้วย bcrypt
- ✅ User profile management

#### 💰 การจัดการธุรกรรม (Transaction Management)
- ✅ **เพิ่มรายการ**: บันทึกรายรับ-รายจ่าย พร้อมหมวดหมู่และบันทึก
- ✅ **ดูรายการทั้งหมด**: แสดงประวัติธุรกรรม เรียงตามวันที่ล่าสุด
- ✅ **แก้ไขรายการ**: แก้ไขธุรกรรมที่บันทึกไว้แล้ว
- ✅ **ลบรายการ**: ลบรายการที่บันทึกผิดพลาด
- ✅ **กรองข้อมูล**: กรองตามประเภท, หมวดหมู่, บัญชี, วันที่
- ✅ **Pagination**: แบ่งหน้าสำหรับข้อมูลจำนวนมาก
- ✅ **ค้นหาขั้นสูง**: ค้นหาจาก note หรือชื่อหมวดหมู่
- ✅ **กรองตามจำนวนเงิน**: กำหนดช่วงเงินขั้นต่ำ-สูงสุด

#### 📊 หน้าสรุปผล (Dashboard)
- ✅ **สรุปยอดประจำเดือน**: แสดงรายรับ, รายจ่าย, ยอดคงเหลือของเดือนปัจจุบัน
- ✅ **ธุรกรรมล่าสุด**: แสดงรายการธุรกรรมล่าสุด พร้อมรายละเอียด
- ✅ **รายจ่ายตามหมวดหมู่**: แสดงการใช้จ่ายแยกตามหมวดหมู่
- ✅ **Dashboard รวม**: รวมข้อมูลสำคัญทั้งหมดในหน้าเดียว

#### 📈 Data Visualization (ใหม่!)
- ✅ **Pie Chart**: แสดงสัดส่วนค่าใช้จ่ายตามหมวดหมู่
- ✅ **Bar Chart**: เปรียบเทียบรายรับ-รายจ่ายรายเดือน
- ✅ **Line Chart**: แนวโน้มการใช้จ่ายของหมวดหมู่เฉพาะ
- ✅ **Top Categories**: หมวดหมู่ที่ใช้งานมากที่สุด

#### 📂 ระบบหมวดหมู่ (Category Management)
- ✅ **หมวดหมู่พื้นฐาน**: อาหาร, เดินทาง, ช้อปปิ้ง, เงินเดือน ฯลฯ
- ✅ **หมวดหมู่ส่วนตัว**: สร้างหมวดหมู่เพิ่มตามต้องการ
- ✅ **จัดการหมวดหมู่**: เก็บ (Archive) หมวดหมู่ที่ไม่ใช้
- ✅ **ไอคอนและสี**: กำหนดไอคอนและสีสำหรับหมวดหมู่
- ✅ **ค้นหาหมวดหมู่**: ค้นหาด้วยชื่อหมวดหมู่
- ✅ **สถิติการใช้งาน**: ดูความถี่การใช้หมวดหมู่แต่ละอัน

#### 🏦 การจัดการบัญชี (Account Management)
- ✅ **บัญชีหลากหลาย**: เงินสด, ธนาคาร, เครดิต, ออม
- ✅ **CRUD Operations**: สร้าง, ดู, แก้ไข, ลบบัญชี
- ✅ **ความปลอดภัย**: ตรวจสอบความเป็นเจ้าของก่อนดำเนินการ

#### 🛠 Technical Features
- ✅ **Clean Architecture**: แยก layers ชัดเจน
- ✅ **PostgreSQL**: ฐานข้อมูลพร้อม indexes เพื่อประสิทธิภาพ
- ✅ **Repository Pattern**: ง่ายต่อการ test และ maintain
- ✅ **Decimal Precision**: ความแม่นยำสำหรับการคำนวณเงิน
- ✅ **Database Migrations**: จัดการ schema versions
- ✅ **Docker Support**: Ready สำหรับ deployment
- ✅ **Comprehensive API**: RESTful APIs ครบครัน

### 🚧 ที่ต้องทำต่อ (Nice-to-Have Features)
- **Savings Goals Management**: ระบบเป้าหมายการออม
- **Budget Planning**: วางแผนงบประมาณ
- **Insights & Analytics**: ระบบวิเคราะห์และให้คำแนะนำอัจฉริยะ
- **Export Data**: ส่งออกข้อมูลเป็น CSV/Excel
- **Monthly/Yearly Reports**: รายงานสรุปรายเดือน/ปี
- **Data Backup & Sync**: สำรองและซิงค์ข้อมูล
- **Multi-currency Support**: รองรับหลายสกุลเงิน

### 🔧 Technical Improvements Needed
- **Unit Tests**: เขียน test cases ครอบคลุม
- **Integration Tests**: ทดสอบการทำงานร่วมกัน  
- **API Documentation**: Swagger/OpenAPI spec
- **Error Handling**: ปรับปรุงการจัดการข้อผิดพลาด
- **Logging**: ระบบ logging ที่สมบูรณ์
- **Rate Limiting**: จำกัดอัตราการเรียก API
- **Caching**: ระบบ cache เพื่อประสิทธิภาพ
- **Performance Optimization**: ปรับปรุงประสิทธิภาพ

## Database Schema

ระบบใช้ PostgreSQL และมี tables หลักดังนี้:

- `users` - ข้อมูลผู้ใช้
- `accounts` - บัญชีการเงิน (เงินสด, ธนาคาร, เครดิต)
- `categories` - หมวดหมู่รายรับ-รายจ่าย
- `transactions` - รายการรับ-จ่าย
- `savings_goals` - เป้าหมายการออม
- `goal_deposits` - การฝากเงินเข้าเป้าหมาย
- `insights` - ข้อมูลวิเคราะห์และคำแนะนำ

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - สมัครสมาชิก
- `POST /api/v1/auth/login` - เข้าสู่ระบบ  
- `POST /api/v1/auth/refresh` - ต่ออายุ token

### Accounts (Protected)
- `POST /api/v1/accounts` - สร้างบัญชี
- `GET /api/v1/accounts` - ดูบัญชีทั้งหมด
- `GET /api/v1/accounts/:id` - ดูบัญชีเฉพาะ
- `DELETE /api/v1/accounts/:id` - ลบบัญชี

### Transactions (Protected)
- `POST /api/v1/transactions` - สร้างรายการใหม่
- `GET /api/v1/transactions` - ดูรายการทั้งหมด (มี filter)
- `GET /api/v1/transactions/:id` - ดูรายการเฉพาะ
- `PUT /api/v1/transactions/:id` - แก้ไขรายการ
- `DELETE /api/v1/transactions/:id` - ลบรายการ

### Categories (Protected)
- `POST /api/v1/categories` - สร้างหมวดหมู่ส่วนตัว
- `GET /api/v1/categories` - ดูหมวดหมู่ทั้งหมด (ระบบ + ส่วนตัว)
- `GET /api/v1/categories/user` - ดูหมวดหมู่ส่วนตัวเท่านั้น
- `PUT /api/v1/categories/:id/archive` - เก็บหมวดหมู่

### Dashboard (Protected)
- `GET /api/v1/dashboard` - ดู Dashboard แบบเบเสิก
- `GET /api/v1/dashboard/summary` - สรุปยอดเดือนปัจจุบัน
- `GET /api/v1/dashboard/summary/monthly` - สรุปยอดประจำเดือน (ระบุเดือน)
- `GET /api/v1/dashboard/transactions/recent` - ธุรกรรมล่าสุด
- `GET /api/v1/dashboard/spending/category` - รายจ่ายแยกตามหมวดหมู่

### Setup/Admin
- `POST /api/v1/setup/categories/default` - สร้างหมวดหมู่เริ่มต้น

> 📖 ดู [API Documentation](API_DOCUMENTATION.md) สำหรับรายละเอียดและตอวอย่างการใช้งาน

## การติดตั้งและรัน

### 🐳 วิธีที่ 1: ใช้ Docker Compose (แนะนำ)

1. Clone repository:
```bash
git clone <repository-url>
cd backend
```

2. รันด้วย Docker Compose:
```bash
docker-compose up -d
```

3. สร้างหมวดหมู่เริ่มต้น:
```bash
curl -X POST http://localhost:8080/api/v1/setup/categories/default
```

4. เข้าใช้งานที่: http://localhost:8080

### 🖥️ วิธีที่ 2: รันแบบ Local Development

1. ติดตั้ง dependencies:
```bash
go mod tidy
```

2. ตั้งค่า PostgreSQL และสร้าง database:
```sql
CREATE DATABASE savvy;
```

3. รัน migrations:
```bash
# ใช้ migrate tool หรือรันไฟล์ SQL โดยตรง
psql -U postgres -d savvy -f migrations/001_initial_schema.up.sql
psql -U postgres -d savvy -f migrations/002_default_categories.up.sql
```

4. ตั้งค่า environment variables:
```bash
cp .env.example .env
# แก้ไขค่าใน .env ตามต้องการ
```

5. รันแอปพลิเคชัน:
```bash
# Development mode (with hot reload)
make dev

# หรือรันปกติ
make run
```

### 🛠️ Development Tools

```bash
# ติดตั้ง development tools
make install-tools

# Format code
make fmt

# Run tests
make test

# Build
make build

# Show all available commands
make help
```

## Architecture Principles

โปรเจกต์นี้ใช้ **Clean Architecture** ซึ่งแบ่งเป็น layers:

1. **Domain Layer** (`internal/domain/`)
   - Entities: โครงสร้างข้อมูลหลัก
   - Repository Interfaces: นิยาม contracts สำหรับ data access

2. **Use Case Layer** (`internal/usecase/`)
   - Business logic และ application rules
   - ใช้ repository interfaces เพื่อ access data

3. **Interface Layer** (`internal/delivery/`)
   - HTTP handlers, middleware, routing
   - รับ input จาก client และส่งต่อไปยัง use cases

4. **Infrastructure Layer** (`internal/infrastructure/`)
   - Database implementations, external services
   - Concrete implementations ของ repository interfaces

## Dependencies

- **Gin** - HTTP web framework
- **PostgreSQL Driver** (lib/pq) - Database driver
- **JWT** - Authentication
- **UUID** - Unique identifiers
- **Decimal** - ความแม่นยำในการคำนวณเงิน
- **Bcrypt** - Password hashing
- **Godotenv** - Environment variable loading

## Contributing

1. สร้าง feature branch
2. เขียน tests สำหรับ code ใหม่
3. ให้แน่ใจว่า tests ผ่านทั้งหมด
4. สร้าง pull request

## License

MIT License
# savvy-service
