# Trading Office AI Dashboard — Backend

REST API สำหรับ Trading Office AI Dashboard พัฒนาด้วย **Go + Fiber** เชื่อมต่อ **PostgreSQL** และรองรับการ integrate กับ **Claude AI** และข้อมูลตลาดจาก **Binance**

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.24 |
| Web Framework | [Fiber v2](https://github.com/gofiber/fiber) |
| Database | PostgreSQL 16 |
| DB Driver | `lib/pq` |
| Migration | `golang-migrate/migrate v4` |
| Config | `spf13/viper` + `joho/godotenv` |
| API Docs | Swagger (`swaggo/swag` + `swaggo/fiber-swagger`) |
| AI | Anthropic Claude API |
| Market Data | Binance REST API |
| Containerization | Docker + Docker Compose |

---

## Project Structure

```
trading_office_backend/
├── main.go                   # Entry point — parse -env flag, start server
├── config/
│   ├── config.go             # Config struct + LoadConfig() — โหลด .env.{env} และ config.yaml
│   └── config.yaml           # Default config values
├── route/
│   ├── boostrap.go           # Bootstrap() — init DB, migration, server, graceful shutdown
│   └── route_apiserver.go    # NewAPIServer() — register middleware + routes
├── handler/
│   └── dashboard_handler.go  # DashboardHandler — /live health check
├── middleware/
│   ├── auth.go               # Bearer token authentication
│   ├── cors.go               # CORS configuration
│   └── logger.go             # Request logger (timezone Asia/Bangkok)
├── db/
│   └── db.go                 # NewDatabase() + RunMigrations()
├── migrations/
│   ├── 0001_init_project.up.sql    # สร้าง schema + tables
│   └── 0001_init_project.down.sql  # Rollback
├── model/
│   └── response.go           # BaseResponse, ErrorDetail, PaginatedResponse
├── utils/
│   └── response_error.go     # ErrorResponse(), SuccessResponse(), TooManyRequests()
├── locales/
│   └── global_error.json     # Error messages
├── docs/                     # Auto-generated Swagger docs (swag init)
├── Dockerfile                # Multi-stage build
├── docker-compose.yml        # postgres + backend services
├── .env.example              # Template environment variables
└── go.mod / go.sum
```

---

## Libraries

### Direct Dependencies

| Package | Version | ใช้ทำอะไร |
|---|---|---|
| `github.com/gofiber/fiber/v2` | v2.52.13 | HTTP web framework |
| `github.com/golang-migrate/migrate/v4` | v4.19.1 | Database migration |
| `github.com/lib/pq` | v1.10.9 | PostgreSQL driver |
| `github.com/spf13/viper` | v1.19.0 | Config management (YAML + env vars) |
| `github.com/joho/godotenv` | v1.5.1 | โหลด `.env` files |
| `github.com/swaggo/fiber-swagger` | v1.3.0 | Swagger UI handler สำหรับ Fiber |
| `github.com/swaggo/swag` | v1.16.6 | Generate Swagger docs จาก annotations |
| `github.com/google/uuid` | v1.6.0 | UUID generation |

### Indirect / Transitive

| Package | ใช้ทำอะไร |
|---|---|
| `github.com/valyala/fasthttp` | HTTP engine ข้างใต้ Fiber |
| `github.com/klauspost/compress` | Compression สำหรับ HTTP responses |
| `github.com/go-openapi/*` | OpenAPI spec generation (ใช้โดย swag) |
| `golang.org/x/net` | Extended networking utilities |

---

## Setup

### Prerequisites

- Go 1.24+
- PostgreSQL 16+ (หรือใช้ Docker)
- `swag` CLI (สำหรับ regenerate docs)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 1. Clone & Install Dependencies

```bash
git clone <repo-url>
cd trading_office_backend
go mod download
```

### 2. Environment Variables

คัดลอก `.env.example` แล้วแก้ไขค่าตาม environment:

```bash
cp .env.example .env.dev
```

| Variable | Default | Required | Description |
|---|---|---|---|
| `APP_PORT` | `8080` | ❌ | Port ที่ server จะ listen |
| `APP_NAME` | `trading-office-api` | ❌ | ชื่อ application |
| `DB_HOST` | `localhost` | ✅ | PostgreSQL host |
| `DB_PORT` | `5432` | ❌ | PostgreSQL port |
| `DB_NAME` | `postgres` | ✅ | Database name |
| `DB_USER` | — | ✅ | Database username |
| `DB_PASS` | — | ❌ | Database password |
| `DB_SSLMODE` | `disable` | ❌ | SSL mode (`disable` / `require`) |
| `CLAUDE_API_KEY` | — | ❌ | Anthropic Claude API key |
| `CLAUDE_MODEL` | `claude-sonnet-4-20250514` | ❌ | Claude model string |
| `EXCHANGE_BASE_URL` | `https://api.binance.com` | ❌ | Binance API base URL |
| `EXCHANGE_FETCH_INTERVAL` | `5` | ❌ | ดึงข้อมูลตลาดทุกกี่วินาที |

Config โหลดแบบ layered: `.env.{APP_ENV}` → `config.yaml` → OS environment variables (ค่าหลังทับค่าก่อน)

### 3. Run Locally

```bash
# dev environment (โหลด .env.dev)
go run main.go -env dev

# production
go run main.go -env prod
```

Flags:

| Flag | Default | Description |
|---|---|---|
| `-env` | `dev` | Environment: `dev`, `qa`, `uat`, `prod` |

### 4. Run with Docker Compose

```bash
# สร้าง .env.prod หรือ set CLAUDE_API_KEY ก่อน
CLAUDE_API_KEY=sk-ant-xxx docker-compose up --build
```

Services ที่ถูก start:
- `postgres` — PostgreSQL 16 บน port `5433` (map จาก container port 5432)
- `backend` — API server บน port `8080`

### 5. Database Migration

Migration รันอัตโนมัติเมื่อ server start ผ่าน `db.RunMigrations()`

หากต้องการรัน manual:

```bash
migrate -path ./migrations -database "postgres://user:pass@localhost:5432/trading_office?sslmode=disable" up
migrate -path ./migrations -database "..." down
```

### 6. Generate Swagger Docs

```bash
swag init -g main.go -o docs
```

---

## API Endpoints

Base URL: `http://localhost:8080`

### Public Endpoints

#### `GET /live` — Health Check

ตรวจสอบว่า server ทำงานปกติ

**Response `200 OK`**
```json
{
  "status": "ok",
  "version": "0.1.0"
}
```

---

#### `POST /auth/token` — Get Access Token

ออก JWT token สำหรับ authentication

**Rate Limit:** 5 requests/minute ต่อ IP

> ⚠️ Endpoint นี้ยังไม่ได้ implement logic — เป็น placeholder สำหรับ phase ถัดไป

---

#### `GET /swagger/*` — Swagger UI

เปิด API documentation แบบ interactive

URL: `http://localhost:8080/swagger/index.html`

---

### Protected Endpoints (`/api/v1/*`)

ทุก route ใต้ `/api/v1` ต้องส่ง `Authorization` header:

```
Authorization: Bearer <token>
```

ถ้าไม่มี header จะได้รับ:

**Response `401 Unauthorized`**
```json
{
  "status": 401,
  "message": "Unauthorized",
  "data": {
    "timestamp": "2026-06-13T10:00:00Z",
    "status": 401,
    "errorCode": "GLB004",
    "error": "UNAUTHORIZED",
    "messageEN": "Unauthorized",
    "messageTH": "ไม่มีสิทธิ์เข้าถึง",
    "path": "/api/v1/..."
  }
}
```

> 📋 Protected routes จะถูกเพิ่มใน upcoming phases

---

## Response Format

### Success

```json
{
  "status": 200,
  "message": "success",
  "data": { ... }
}
```

### Error

```json
{
  "status": 400,
  "message": "Bad Request",
  "data": {
    "timestamp": "2026-06-13T10:00:00Z",
    "status": 400,
    "errorCode": "GLB001",
    "error": "BAD_REQUEST",
    "messageEN": "Invalid request",
    "messageTH": "คำขอไม่ถูกต้อง",
    "path": "/api/v1/resource"
  }
}
```

### Paginated

```json
{
  "status": 200,
  "message": "success",
  "data": [ ... ],
  "total": 100,
  "page": 1,
  "limit": 20
}
```

---

## Database Schema

### Schema: `trading_office`

Migration tracking ใช้ schema `trading_office` แทน `public`

#### `trading_office.market_prices`

เก็บราคาตลาดที่ดึงจาก Binance API

| Column | Type | Description |
|---|---|---|
| `id` | BIGSERIAL PK | Auto-increment primary key |
| `symbol` | VARCHAR(20) | Trading pair เช่น `BTCUSDT` |
| `price` | NUMERIC(20,8) | ราคา (ทศนิยม 8 หลัก) |
| `timestamp` | TIMESTAMPTZ | เวลาที่บันทึกราคา |
| `created_at` | TIMESTAMPTZ | วันที่สร้าง record |
| `created_by` | VARCHAR(100) | ผู้สร้าง (default: `tdo-system`) |
| `updated_at` | TIMESTAMPTZ | วันที่อัปเดตล่าสุด |
| `updated_by` | VARCHAR(100) | ผู้อัปเดต |
| `is_deleted` | BOOLEAN | Soft delete flag |
| `deleted_at` | TIMESTAMPTZ | วันที่ลบ |
| `deleted_by` | VARCHAR(100) | ผู้ลบ |

**Indexes:**
- `idx_market_prices_symbol` — query by symbol
- `idx_market_prices_timestamp DESC` — query ล่าสุดก่อน
- `idx_market_prices_is_deleted` — filter soft-deleted records

---

## Middleware

| Middleware | Config |
|---|---|
| **CORS** | Allow all origins, headers: `Origin, Content-Type, Accept, Authorization`, methods: `GET POST PUT DELETE OPTIONS` |
| **Logger** | Format: `[time] status method path latency`, timezone: Asia/Bangkok |
| **Auth** | ตรวจ `Authorization` header — ถ้าว่างคืน 401 |
| **Rate Limiter** | เฉพาะ `POST /auth/token` — max 5 req/min ต่อ IP |

---

## Error Codes

| Code | HTTP Status | Description |
|---|---|---|
| `GLB004` | 401 | Unauthorized — ไม่มี Authorization header |
| `GLB005` | 429 | Too Many Requests — เกิน rate limit |

---

## Development Notes

- **Config layering:** `.env.{env}` โหลดด้วย `godotenv` → ค่าใน `config.yaml` โหลดด้วย `viper` → OS env vars override ผ่าน `overrideFromEnv()`
- **Migration:** รันอัตโนมัติทุกครั้งที่ server start, ใช้ `golang-migrate` track version ใน `schema_migrations` table ภายใน schema `trading_office`
- **Graceful shutdown:** รับ `SIGINT` (Ctrl+C) → shutdown Fiber server → close DB connection
- **DB pool:** MaxOpenConns=25, MaxIdleConns=10
