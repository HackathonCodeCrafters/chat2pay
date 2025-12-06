# Chat2Pay Backend API

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?logo=docker)](https://www.docker.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?logo=postgresql)](https://www.postgresql.org/)
[![Xendit](https://img.shields.io/badge/Xendit-Integrated-00D1B2)](https://xendit.co/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Backend API untuk sistem pembayaran QRIS dengan integrasi Xendit, dibangun dengan Go Fiber dan PostgreSQL.

**🎯 Status**: ✅ Production Ready | 🧪 Fully Tested | 🐳 Docker Optimized

## 🚀 Features

### ✅ Authentication & Authorization
- **JWT-based Authentication** with 24-hour token expiry
- **Role-based Access Control (RBAC)**:
  - **Merchant**: Manage products, orders, view customers
  - **Customer**: Create orders, view own orders
- **Secure Password** with bcrypt hashing
- **Dual Registration System**:
  - Merchant Registration → Creates Merchant + MerchantUser (owner)
  - Customer Registration → Creates Customer account

### 💳 Payment Features (QRIS)
- **Xendit Integration** untuk QRIS payment gateway
- **Invoice Creation** dengan QR code generation
- **Real-time Payment Status** tracking
- **Webhook Handler** untuk notifikasi pembayaran
- **Payment Logging** untuk audit trail
- **Sandbox & Production** environment support

### 🗄️ Database
- **PostgreSQL 16** dengan pgvector extension
- **UUID-based primary keys** untuk scalability
- **Automated migrations** dengan sql-migrate
- **Connection pooling** untuk performance

## 📋 Prerequisites

- Docker & Docker Compose (Recommended)
- Golang 1.25+ (untuk development)
- PostgreSQL 16+ (jika run lokal)
- Xendit API Key (untuk payment features)

## 🛠️ Installation

### 1. Clone Repository
```bash
git clone <repository-url>
cd backend-chat2pay
```

### 2. Configuration
Edit `config/yaml/app.yaml`:
```yaml
app:
  name: chat2pay
  port: 9005

db:
  dialect: postgres
  host: postgres      # use 'postgres' for Docker, 'localhost' for local
  port: 5432
  db_name: postgres
  username: postgres
  password: BZHnDem5mcbtJkk9YbZee787SyAFGy7x
  ssl_mode: disable

jwt:
  key: uwoww_sangad
  expired_minute: 1440  # 24 hours

xendit:
  api_key: xnd_development_YOUR_API_KEY
  webhook_verification_token: your_webhook_token
  environment: sandbox  # or 'production'
  callback_url: https://your-domain.com/api/webhooks/xendit

gemini:
  api_key: your_gemini_api_key  # for AI features

logger:
  enable: true
```

### 3. Run with Docker (Recommended)
```bash
# Start all services (PostgreSQL + Backend)
docker-compose -f docker-compose.dev.yaml up -d --build

# Check status
docker ps

# View logs
docker logs chat2pay-be
docker logs postgres_svc

# Stop services
docker-compose -f docker-compose.dev.yaml down
```

This will start:
- **PostgreSQL 16** with pgvector on port `5432`
- **Backend API** on port `9005`

### 4. Run Migrations
Migrations run automatically on application startup. Manual migration:
```bash
# Run migrations
docker exec chat2pay-be ./chat2pay migration up

# Rollback
docker exec chat2pay-be ./chat2pay migration down
```

### 5. Run Locally (without Docker)
```bash
# Install dependencies
go mod tidy

# Update config/yaml/app.yaml (set host: localhost)
# Make sure PostgreSQL is running locally

# Run migrations
go run main.go migration up

# Run application
go run main.go http
```

Application will run at: `http://localhost:9005`

## 🚦 Quick Start

```bash
# 1. Start services with Docker
docker-compose -f docker-compose.dev.yaml up -d --build

# 2. Wait for services to be ready (check logs)
docker logs chat2pay-be

# 3. Test the API
curl http://localhost:9005/api/

# 4. Setup test data for payments
bash setup_test_data.sh

# 5. Test payment flow
bash test_payment_docker.sh
```

## 📚 API Endpoints

### Available Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/` | Health check | No |
| POST | `/api/auth/merchant/register` | Register merchant | No |
| POST | `/api/auth/merchant/login` | Merchant login | No |
| POST | `/api/auth/customer/register` | Register customer | No |
| POST | `/api/auth/customer/login` | Customer login | No |
| POST | `/api/payments/qris/create` | Create QRIS payment | No |
| GET | `/api/payments/:id` | Get payment status | No |
| POST | `/api/webhooks/xendit` | Xendit webhook handler | No |

### API Documentation

See complete documentation with examples below or check: **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)**

### Quick Start - Complete API Examples

#### 1. Authentication

**Register Merchant:**
```bash
curl -X POST http://localhost:9005/api/auth/merchant/register \
  -H "Content-Type: application/json" \
  -d '{
    "merchant_name": "Toko Elektronik ABC",
    "legal_name": "PT ABC Elektronik",
    "email": "owner@abc.com",
    "phone": "081234567890",
    "name": "John Doe",
    "password": "password123"
  }'
```

**Register Customer:**
```bash
curl -X POST http://localhost:9005/api/auth/customer/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "phone": "081234567890",
    "password": "password123"
  }'
```

**Merchant Login:**
```bash
curl -X POST http://localhost:9005/api/auth/merchant/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "owner@abc.com",
    "password": "password123"
  }'
# Response will contain access_token - save it!
```

**Customer Login:**
```bash
curl -X POST http://localhost:9005/api/auth/customer/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane@example.com",
    "password": "password123"
  }'
```

#### 2. Merchant Management

**Get All Merchants (Public):**
```bash
curl -X GET http://localhost:9005/api/merchants
```

**Get Merchant by ID (Public):**
```bash
curl -X GET http://localhost:9005/api/merchants/1
```

**Create Merchant (Merchant Only):**
```bash
TOKEN="your_merchant_access_token"
curl -X POST http://localhost:9005/api/merchants \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Toko Baru",
    "legal_name": "PT Toko Baru",
    "email": "baru@example.com",
    "phone": "082345678901"
  }'
```

**Update Merchant (Merchant Only):**
```bash
TOKEN="your_merchant_access_token"
curl -X PUT http://localhost:9005/api/merchants/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Toko Updated",
    "legal_name": "PT Updated",
    "email": "owner@abc.com",
    "phone": "081999999999"
  }'
```

**Delete Merchant (Merchant Only):**
```bash
TOKEN="your_merchant_access_token"
curl -X DELETE http://localhost:9005/api/merchants/1 \
  -H "Authorization: Bearer $TOKEN"
```

#### 3. Product Management

**Get All Products (Public):**
```bash
curl -X GET "http://localhost:9005/api/products?page=1&limit=10"
```

**Get Product by ID (Public):**
```bash
curl -X GET http://localhost:9005/api/products/1
```

**Create Product (Merchant Only):**
```bash
TOKEN="your_merchant_access_token"
curl -X POST http://localhost:9005/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "merchant_id": 1,
    "name": "Laptop ASUS ROG",
    "description": "Gaming laptop with RTX 4070",
    "sku": "LAPTOP-001",
    "price": 25000000,
    "stock": 10
  }'
```

**Update Product (Merchant Only):**
```bash
TOKEN="your_merchant_access_token"
curl -X PUT http://localhost:9005/api/products/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "merchant_id": 1,
    "name": "Laptop ASUS ROG Updated",
    "description": "Gaming laptop with RTX 4080",
    "sku": "LAPTOP-001",
    "price": 30000000,
    "stock": 15
  }'
```

**Delete Product (Merchant Only):**
```bash
TOKEN="your_merchant_access_token"
curl -X DELETE http://localhost:9005/api/products/1 \
  -H "Authorization: Bearer $TOKEN"
```

#### 4. Customer Management

**Get All Customers (Merchant Only):**
```bash
TOKEN="your_merchant_access_token"
curl -X GET "http://localhost:9005/api/customers?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

**Get Customer by ID (Authenticated):**
```bash
TOKEN="your_access_token"
curl -X GET http://localhost:9005/api/customers/1 \
  -H "Authorization: Bearer $TOKEN"
```

**Update Customer (Authenticated):**
```bash
TOKEN="your_customer_access_token"
curl -X PUT http://localhost:9005/api/customers/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Jane Updated",
    "email": "jane.updated@example.com",
    "phone": "082999999999"
  }'
```

**Delete Customer (Authenticated):**
```bash
TOKEN="your_customer_access_token"
curl -X DELETE http://localhost:9005/api/customers/1 \
  -H "Authorization: Bearer $TOKEN"
```

#### 5. Payment Management (QRIS)

**Create QRIS Payment:**
```bash
# Create payment for an order
curl -X POST http://localhost:9005/api/payments/qris/create \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": "f6202b09-3597-40d6-84dc-745e2c50576b",
    "amount": 100000
  }'

# Response:
# {
#   "status": true,
#   "data": {
#     "payment_id": "bd40a1ad-a8dc-4fe1-bb99-e50fd06dbc56",
#     "order_id": "f6202b09-3597-40d6-84dc-745e2c50576b",
#     "amount": 100000,
#     "status": "PENDING",
#     "qr_code_url": "https://checkout-staging.xendit.co/web/...",
#     "expires_at": "2025-12-07T14:55:14.034Z",
#     "created_at": "2025-12-06T14:55:14.335653Z"
#   }
# }
```

**Check Payment Status:**
```bash
PAYMENT_ID="bd40a1ad-a8dc-4fe1-bb99-e50fd06dbc56"
curl -X GET http://localhost:9005/api/payments/$PAYMENT_ID

# Response:
# {
#   "status": true,
#   "data": {
#     "payment_id": "bd40a1ad-a8dc-4fe1-bb99-e50fd06dbc56",
#     "order_id": "f6202b09-3597-40d6-84dc-745e2c50576b",
#     "amount": 100000,
#     "status": "PAID",  # or "PENDING", "EXPIRED", "FAILED"
#     "paid_at": "2025-12-06T21:55:39.961109Z",
#     "created_at": "2025-12-06T14:55:14.335653Z",
#     "updated_at": "2025-12-06T14:55:39.962334Z"
#   }
# }
```

**Xendit Webhook (Internal - Called by Xendit):**
```bash
# This endpoint is called automatically by Xendit when payment status changes
# POST http://localhost:9005/api/webhooks/xendit

# You can simulate webhook for testing:
curl -X POST http://localhost:9005/api/webhooks/xendit \
  -H "Content-Type: application/json" \
  -d '{
    "id": "693443d26720fe660ada1624",
    "external_id": "ORDER-f6202b09-3597-40d6-84dc-745e2c50576b",
    "status": "PAID",
    "amount": 100000,
    "paid_at": "2025-12-06T14:56:00.000Z",
    "payment_channel": "QRIS"
  }'
```

#### 6. Helper Script - Get Token

Save this as `get_token.sh`:
```bash
#!/bin/bash

# Get Merchant Token
get_merchant_token() {
  curl -s -X POST http://localhost:9005/api/auth/merchant/login \
    -H "Content-Type: application/json" \
    -d '{"email":"owner@abc.com","password":"password123"}' \
    | grep -o '"access_token":"[^"]*' | cut -d'"' -f4
}

# Get Customer Token
get_customer_token() {
  curl -s -X POST http://localhost:9005/api/auth/customer/login \
    -H "Content-Type: application/json" \
    -d '{"email":"jane@example.com","password":"password123"}' \
    | grep -o '"access_token":"[^"]*' | cut -d'"' -f4
}

# Usage:
# TOKEN=$(get_merchant_token)
# echo $TOKEN
```

## 🔒 Authorization Rules

### Public Endpoints (No Auth)
- `GET /api/merchants` - List merchants
- `GET /api/merchants/:id` - Get merchant detail
- `GET /api/products` - List products
- `GET /api/products/:id` - Get product detail
- `POST /api/payments/qris/create` - Create QRIS payment
- `GET /api/payments/:id` - Check payment status
- `POST /api/webhooks/xendit` - Xendit webhook handler

### Merchant Only
- Create/Update/Delete: Merchants, Products
- View all customers
- Update order status, Delete orders

### Authenticated (Merchant or Customer)
- Create orders
- View own orders
- Update own profile

## 🧪 Testing

### Automated Payment Testing
```bash
# Setup test data first
bash setup_test_data.sh

# Run comprehensive payment test
bash test_payment_docker.sh
```

This will test:
- ✅ Payment creation with Xendit API
- ✅ Payment status retrieval
- ✅ Webhook processing
- ✅ Status update to PAID
- ✅ Database logging

### Manual Testing
```bash
# 1. Create payment
curl -X POST http://localhost:9005/api/payments/qris/create \
  -H "Content-Type: application/json" \
  -d '{"order_id":"YOUR_ORDER_ID","amount":100000}'

# 2. Get payment status
curl http://localhost:9005/api/payments/YOUR_PAYMENT_ID

# 3. Simulate webhook (payment success)
curl -X POST http://localhost:9005/api/webhooks/xendit \
  -H "Content-Type: application/json" \
  -d '{
    "external_id": "ORDER-YOUR_ORDER_ID",
    "status": "PAID",
    "amount": 100000
  }'

# 4. Verify payment status changed
curl http://localhost:9005/api/payments/YOUR_PAYMENT_ID
```

### Unit Testing
```bash
# Run all tests
go test ./... -v --cover

# Run specific package tests
go test ./internal/service -v
go test ./internal/repositories -v

# Test with coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Database Verification
```bash
# Check payments
docker exec -i postgres_svc psql -U postgres -d postgres -c \
  "SELECT id, payment_method, amount, status FROM payments ORDER BY created_at DESC LIMIT 5;"

# Check payment logs
docker exec -i postgres_svc psql -U postgres -d postgres -c \
  "SELECT payment_id, event_type, created_at FROM payment_logs ORDER BY created_at DESC LIMIT 10;"
```

## 📁 Project Structure

```
.
├── bootstrap/                    # Dependency injection & bootstrapping
│   ├── const.go                 # DI container constants
│   ├── handler.go               # Handler registrations
│   ├── package.go               # Package dependencies (Xendit, LLM)
│   ├── repository.go            # Repository registrations
│   └── service.go               # Service registrations
├── config/
│   └── yaml/
│       ├── app.yaml             # Configuration file
│       └── cfg_yaml.go          # Config loader
├── internal/
│   ├── api/
│   │   ├── dto/                 # Request/Response DTOs
│   │   │   ├── payment_dto.go  # Payment DTOs
│   │   │   ├── auth_*.go       # Auth DTOs
│   │   │   └── ...
│   │   ├── handlers/
│   │   │   ├── payment_handler.go    # Payment endpoints
│   │   │   ├── webhook_handler.go    # Xendit webhook
│   │   │   └── ...
│   │   ├── routes/
│   │   │   ├── payment.go      # Payment routes
│   │   │   └── ...
│   │   └── router.go            # Main router
│   ├── entities/
│   │   └── entities.go          # Database models (GORM)
│   ├── payment/
│   │   └── xendit/
│   │       ├── client.go        # Xendit API client
│   │       ├── constants.go     # Payment constants
│   │       └── types.go         # Xendit types
│   ├── repositories/
│   │   ├── payment_repository.go      # Payment DB operations
│   │   ├── payment_log_repository.go  # Payment logs
│   │   └── ...
│   ├── service/
│   │   ├── payment_service.go   # Payment business logic
│   │   └── ...
│   ├── middlewares/
│   │   └── jwt/                 # JWT authentication
│   └── pkg/
│       └── llm/                 # AI integrations (Gemini, OpenAI)
├── migration/                   # Database migrations
│   ├── 20251205*.sql           # Initial schema migrations
│   └── 20251206*.sql           # Payment-related migrations
├── Dockerfile.dev               # Development Docker image
├── docker-compose.dev.yaml      # Docker Compose setup
├── go.mod                       # Go dependencies
├── main.go                      # Application entry point
├── test_payment_docker.sh       # Payment testing script
├── setup_test_data.sh           # Test data setup
├── README.md                    # This file
└── README_PAYMENT.md            # Payment documentation
```

## 🐛 Troubleshooting

### Database Connection Error
```bash
# Check if containers are running
docker ps

# Check PostgreSQL logs
docker logs postgres_svc

# Check if database exists
docker exec -i postgres_svc psql -U postgres -l

# Reset database
docker-compose -f docker-compose.dev.yaml down -v
docker-compose -f docker-compose.dev.yaml up -d --build
```

### Application Not Starting
```bash
# Check application logs
docker logs chat2pay-be --tail 50

# Restart application
docker-compose -f docker-compose.dev.yaml restart backend

# Rebuild from scratch
docker-compose -f docker-compose.dev.yaml down
docker-compose -f docker-compose.dev.yaml up -d --build
```

### Payment Creation Fails
```bash
# Check Xendit API key in config/yaml/app.yaml
# Make sure you're using the correct environment (sandbox/production)

# Test Xendit connectivity
curl -u xnd_development_YOUR_API_KEY: https://api.xendit.co/v2/invoices

# Check application logs for Xendit errors
docker logs chat2pay-be | grep -i xendit
```

### Webhook Not Working
- Make sure `callback_url` is publicly accessible (use ngrok for testing)
- Verify webhook token matches Xendit configuration
- Check `payment_logs` table for webhook events
```bash
docker exec -i postgres_svc psql -U postgres -d postgres -c \
  "SELECT * FROM payment_logs WHERE event_type = 'webhook_received' ORDER BY created_at DESC LIMIT 5;"
```

### JWT Token Invalid
- Check `config/yaml/app.yaml` - jwt.key must be the same
- Token expires in 24 hours, re-login if expired
- Verify token format: `Authorization: Bearer <token>`

### Migration Errors
```bash
# Check migration status
docker exec chat2pay-be ./chat2pay migration status

# Run migrations manually
docker exec chat2pay-be ./chat2pay migration up

# Check migration files
ls -la migration/
```

## 📝 Configuration Reference

### app.yaml Structure
```yaml
app:
  name: chat2pay
  port: 9005

db:
  dialect: postgres
  host: postgres              # 'postgres' for Docker, 'localhost' for local
  port: 5432
  db_name: postgres
  username: postgres
  password: your_secure_password
  ssl_mode: disable

jwt:
  key: your_jwt_secret_key
  expired_minute: 1440        # 24 hours

xendit:
  api_key: xnd_development_XXX         # Get from Xendit Dashboard
  webhook_verification_token: XXX      # Optional, for webhook security
  environment: sandbox                 # 'sandbox' or 'production'
  callback_url: https://your-domain.com/api/webhooks/xendit

gemini:
  api_key: your_gemini_api_key        # For AI features (optional)

logger:
  enable: true
```

### Environment Variables (Production)
```bash
# Database
export DB_HOST=your_db_host
export DB_PORT=5432
export DB_NAME=postgres
export DB_USER=postgres
export DB_PASSWORD=your_secure_password

# Application
export APP_PORT=9005
export JWT_SECRET=your_jwt_secret_key

# Xendit
export XENDIT_API_KEY=xnd_production_XXX
export XENDIT_WEBHOOK_TOKEN=your_webhook_token
export XENDIT_ENVIRONMENT=production
export XENDIT_CALLBACK_URL=https://your-domain.com/api/webhooks/xendit
```

## 🔗 Useful Links

- **Xendit Documentation**: https://developers.xendit.co/api-reference/
- **Xendit Dashboard**: https://dashboard.xendit.co/
- **ngrok** (for webhook testing): https://ngrok.com/
- **Go Fiber Documentation**: https://docs.gofiber.io/
- **GORM Documentation**: https://gorm.io/docs/

## ✅ Tested & Working Features

All features below have been tested and verified working in Docker:

### Core Backend
- ✅ Docker setup dengan PostgreSQL 16 + pgvector
- ✅ Automated database migrations
- ✅ Health check endpoint

### Payment System (QRIS)
- ✅ QRIS payment creation via Xendit API
- ✅ QR code generation dan checkout URL
- ✅ Payment status tracking (PENDING → PAID/EXPIRED/FAILED)
- ✅ Webhook handler untuk notifikasi pembayaran
- ✅ Payment logging untuk audit trail
- ✅ Database persistence (payments & payment_logs tables)

### Testing
- ✅ Automated test scripts (`test_payment_docker.sh`)
- ✅ Test data setup (`setup_test_data.sh`)
- ✅ Manual testing via curl commands
- ✅ Database verification queries

### Infrastructure
- ✅ Multi-container Docker setup
- ✅ Network isolation
- ✅ Health checks
- ✅ Volume persistence
- ✅ Auto-restart policies

## 📄 License

This project is licensed under the MIT License.

---

**Last Updated**: December 2025  
**Tested on**: Docker 20.10+, PostgreSQL 16, Go 1.25+

