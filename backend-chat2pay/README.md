# Chat2Pay Backend API

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

## 📋 Prerequisites

- Docker & Docker Compose
- Golang 1.19 or latest
- PostgreSQL 17

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
  host: localhost  # or host.docker.internal for Docker
  port: 5432
  db_name: chat2pay
  username: postgres
  password: your_password

jwt:
  key: your_secret_key_here
  expired_minute: 1440  # 24 hours

logger:
  enable: true
```

### 3. Run with Docker Compose
```bash
docker-compose up -d
```

This will start:
- PostgreSQL database on port 5432
- Backend API on port 9005

### 4. Run Locally (without Docker)
```bash
# Install dependencies
go mod tidy

# Run migrations (will auto-run on startup)
# Migrations are in folder: migrations/001_create_merchants_table.sql

# Run application
go run main.go
```

Application will run at: `http://localhost:9005`

## 📚 API Documentation

See complete documentation at: **[API_DOCUMENTATION.md](./API_DOCUMENTATION.md)**

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

#### 5. Order Management

**Create Order (Authenticated):**
```bash
TOKEN="your_customer_access_token"
curl -X POST http://localhost:9005/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "customer_id": 1,
    "merchant_id": 1,
    "items": [
      {
        "product_id": 1,
        "quantity": 2,
        "price": 25000000
      }
    ],
    "shipping_cost": 50000,
    "discount": 100000
  }'
```

**Get All Orders (Authenticated):**
```bash
TOKEN="your_access_token"
curl -X GET "http://localhost:9005/api/orders?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN"
```

**Get Order by ID (Authenticated):**
```bash
TOKEN="your_access_token"
curl -X GET http://localhost:9005/api/orders/1 \
  -H "Authorization: Bearer $TOKEN"
```

**Update Order Status (Merchant Only):**
```bash
TOKEN="your_merchant_access_token"
curl -X PATCH http://localhost:9005/api/orders/1/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "status": "paid"
  }'
# Status options: pending, paid, shipped, completed, cancelled
```

**Delete Order (Merchant Only):**
```bash
TOKEN="your_merchant_access_token"
curl -X DELETE http://localhost:9005/api/orders/1 \
  -H "Authorization: Bearer $TOKEN"
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

### Merchant Only
- Create/Update/Delete: Merchants, Products
- View all customers
- Update order status, Delete orders

### Authenticated (Merchant or Customer)
- Create orders
- View own orders
- Update own profile

## 🧪 Testing

```bash
# Run tests
cd internal/service
go test -v --cover

# Run specific test
go test -v -run TestServiceName
```

## 📁 Project Structure

```
.
├── bootstrap/              # Database connection & dependencies
├── config/                 # Configuration files (app.yaml)
│   └── yaml/
├── internal/               # Main application code
│   ├── api/
│   │   ├── dto/           # Data Transfer Objects (request/response)
│   │   ├── handlers/      # HTTP request handlers
│   │   ├── presenter/     # Response formatters
│   │   ├── routes/        # Route definitions
│   │   └── router.go      # Main router setup
│   ├── consts/            # Constants
│   ├── entities/          # Database models (GORM)
│   ├── middlewares/       # JWT auth & other middlewares
│   │   └── jwt/
│   ├── pkg/               # Shared packages (logger, etc)
│   ├── repositories/      # Database operations (Repository pattern)
│   └── service/           # Business logic (Service layer)
├── migrations/            # Database migrations (SQL)
├── Dockerfile
├── docker-compose.yaml
├── go.mod
├── main.go
├── README.md
└── API_DOCUMENTATION.md
```

## 🐛 Troubleshooting

### Database Connection Error
```bash
# Check if PostgreSQL is running
docker ps

# Check connection
psql -h localhost -U postgres -d chat2pay

# Reset database
docker-compose down -v
docker-compose up -d
```

### JWT Token Invalid
- Check `config/yaml/app.yaml` - jwt.key must be the same
- Token expires in 24 hours, re-login if expired

### Migration Not Running
- Migrations auto-run on startup
- Check `migrations/` folder for SQL files
- Make sure database connection success in logs

## 📝 Environment Variables

Application uses `config/yaml/app.yaml` for configuration.

For production, use environment variables:
```bash
export DB_HOST=your_db_host
export DB_PORT=5432
export DB_NAME=chat2pay
export DB_USER=postgres
export DB_PASSWORD=your_password
export JWT_SECRET=your_secret_key
export APP_PORT=9005
```

## 📄 License

This project is licensed under the MIT License.

