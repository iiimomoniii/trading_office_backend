# Trading Office Backend

Go + Fiber + PostgreSQL + Claude AI

## Project Structure

```
trading_office_backend/
├── config/
│   └── config.go
├── db/
│   └── db.go
├── handler/
│   └── dashboard_handler.go
├── locales/
│   └── global_error.json
├── middleware/
│   ├── auth.go
│   ├── cors.go
│   └── logger.go
├── migrations/
│   ├── 0001_init_project.up.sql
│   └── 0001_init_project.down.sql
├── model/
│   └── response.go
├── repository/
├── route/
│   ├── bootstrap.go
│   └── route_apiserver.go
├── service/
├── utils/
│   └── response_error.go
├── config.yaml
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── main.go
```

## Install Go Packages

```bash
# HTTP Framework
go get github.com/gofiber/fiber/v2

# Database
go get github.com/lib/pq
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file

# Config
go get github.com/spf13/viper
go get github.com/joho/godotenv

# JWT
go get github.com/golang-jwt/jwt/v4

# Swagger
go install github.com/swaggo/swag/cmd/swag@latest
go get github.com/swaggo/fiber-swagger
go get github.com/swaggo/files

# Generate Swagger docs
swag init -g main.go

# ตรวจสอบ go.mod และ go.sum
go mod tidy
```

## Run

```bash
# Development
go run main.go -env dev

# Production
go run main.go -env prod
```

## Docker

```bash
# Run with Docker Compose
docker-compose up --build

# Stop
docker-compose down
```

## Environment Variables

| ตัวแปร | Default | คำอธิบาย |
|---|---|---|
| APP_PORT | :8080 | Port ที่ server รัน |
| APP_NAME | trading-office-api | ชื่อ app |
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_NAME | trading_office | Database name |
| DB_USER | postgres | Database user |
| DB_PASS | — | Database password |
| DB_SSLMODE | disable | SSL mode |
| CLAUDE_API_KEY | — | Claude API key |
| CLAUDE_MODEL | claude-sonnet-4-20250514 | Claude model |
| EXCHANGE_BASE_URL | https://api.binance.com | Exchange API URL |
| EXCHANGE_FETCH_INTERVAL | 5 | Fetch interval (seconds) |

## API Endpoints

| Method | Path | Description |
|---|---|---|
| GET | /live | Health check |
| GET | /events | SSE event stream (realtime) |
| GET | /api/v1/markets/latest | Get latest market prices |
| GET | /api/v1/markets/:symbol/history | Get price history by symbol |
| GET | /api/v1/signals/latest | Get latest signals |
| GET | /api/v1/ai/summary/:symbol | Get AI analysis summary |
| GET | /api/v1/rules | Get all rules |
| POST | /api/v1/rules | Create rule |
| PUT | /api/v1/rules/:id | Update rule |
| DELETE | /api/v1/rules/:id | Delete rule |

## Error Codes

### Global (GLB)

| Code | Name | EN | TH |
|---|---|---|---|
| GLB001 | NOT_FOUND | Record not found | ไม่พบข้อมูล |
| GLB002 | INTERNAL_ERROR | Internal server error | เกิดข้อผิดพลาดภายในระบบ |
| GLB003 | INVALID_BODY | Invalid request body | รูปแบบข้อมูลไม่ถูกต้อง |
| GLB004 | UNAUTHORIZED | Unauthorized | ไม่มีสิทธิ์เข้าถึง |
| GLB005 | TOO_MANY_REQUESTS | Too many requests | คำขอมากเกินไป |

### Market (MKT)

| Code | Name | EN | TH |
|---|---|---|---|
| MKT001 | SYMBOL_INVALID | Invalid symbol | สัญลักษณ์ไม่ถูกต้อง |
| MKT002 | FETCH_FAILED | Failed to fetch market data | ดึงข้อมูลตลาดไม่สำเร็จ |
| MKT003 | NOT_FOUND | Market data not found | ไม่พบข้อมูลตลาด |

### Signal (SIG)

| Code | Name | EN | TH |
|---|---|---|---|
| SIG001 | ENGINE_FAILED | Signal engine failed | ระบบสัญญาณล้มเหลว |
| SIG002 | NOT_FOUND | Signal not found | ไม่พบสัญญาณ |
| SIG003 | INVALID_TYPE | Invalid signal type | ประเภทสัญญาณไม่ถูกต้อง |

### AI (AI)

| Code | Name | EN | TH |
|---|---|---|---|
| AI001 | CLAUDE_TIMEOUT | Claude API timeout | Claude API หมดเวลา |
| AI002 | CLAUDE_FAILED | Claude API failed | Claude API ล้มเหลว |
| AI003 | NOT_FOUND | AI result not found | ไม่พบผลการวิเคราะห์ AI |

### Rule (RUL)

| Code | Name | EN | TH |
|---|---|---|---|
| RUL001 | CONDITION_INVALID | Invalid rule condition | เงื่อนไขไม่ถูกต้อง |
| RUL002 | NOT_FOUND | Rule not found | ไม่พบกฎ |
| RUL003 | DUPLICATE | Rule already exists | กฎนี้มีอยู่แล้ว |

## Response Format

### Success

```json
{
  "status": 200,
  "message": "success",
  "data": {}
}
```

### Error

```json
{
  "status": 400,
  "message": "Invalid symbol",
  "data": {
    "timestamp": "2026-06-12T10:00:00+07:00",
    "status": 400,
    "errorCode": "MKT001",
    "error": "SYMBOL_INVALID",
    "messageEN": "Invalid symbol",
    "messageTH": "สัญลักษณ์ไม่ถูกต้อง",
    "path": "/api/v1/markets/latest"
  }
}
```

## Database Schema

```
Schema: trading_office

Tables:
- market_prices     ← ราคาตลาดจาก Binance
- signals           ← สัญญาณ LONG/SHORT/WAIT
- ai_results        ← ผลการวิเคราะห์จาก Claude
- rules             ← Rule Engine configuration
```

## เพิ่ม Module ใหม่

```
1. สร้าง locales/{module}_error.json   ← error codes
2. สร้าง model/{module}.go
3. สร้าง migrations/000N_{module}.up.sql
4. สร้าง repository/{module}_repository.go
5. สร้าง service/{module}_service.go + test
6. สร้าง handler/{module}_handler.go
7. เพิ่ม route ใน route/route_apiserver.go
```
