# Chat2Pay

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go 1.25">
  <img src="https://img.shields.io/badge/Fiber-v2-00ACD7?style=for-the-badge&logo=go&logoColor=white" alt="Fiber v2">
  <img src="https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL 16">
  <img src="https://img.shields.io/badge/pgvector-Vector_DB-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="pgvector">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Next.js-16-000000?style=for-the-badge&logo=next.js&logoColor=white" alt="Next.js 16">
  <img src="https://img.shields.io/badge/React-19-61DAFB?style=for-the-badge&logo=react&logoColor=black" alt="React 19">
  <img src="https://img.shields.io/badge/TypeScript-5-3178C6?style=for-the-badge&logo=typescript&logoColor=white" alt="TypeScript">
  <img src="https://img.shields.io/badge/Tailwind_CSS-4-06B6D4?style=for-the-badge&logo=tailwindcss&logoColor=white" alt="Tailwind CSS">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Google-Gemini-8E75B2?style=for-the-badge&logo=googlegemini&logoColor=white" alt="Gemini">
  <img src="https://img.shields.io/badge/OpenAI-GPT-412991?style=for-the-badge&logo=openai&logoColor=white" alt="OpenAI">
  <img src="https://img.shields.io/badge/Mistral-AI-FF7000?style=for-the-badge&logo=mistral&logoColor=white" alt="Mistral">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker">
  <img src="https://img.shields.io/badge/Swagger-OpenAPI_3.0-85EA2D?style=for-the-badge&logo=swagger&logoColor=black" alt="Swagger">
  <img src="https://img.shields.io/badge/JWT-Auth-000000?style=for-the-badge&logo=jsonwebtokens&logoColor=white" alt="JWT">
  <img src="https://img.shields.io/badge/License-MIT-yellow?style=for-the-badge" alt="MIT License">
</p>

---

**Chat2Pay** adalah platform e-commerce dengan fitur **chat-to-pay** - memungkinkan customer mencari dan membeli produk menggunakan percakapan natural language dengan bantuan AI/LLM.

## 🎯 Fitur Utama

- **AI-Powered Product Search** - Cari produk dengan bahasa natural menggunakan LLM (Gemini/OpenAI/Mistral)
- **Vector Search** - Pencarian produk berbasis semantic similarity dengan pgvector
- **Multi-Role Authentication** - JWT-based auth dengan role Merchant dan Customer
- **Complete E-Commerce Flow** - Manajemen produk, order, customer, dan merchant
- **Swagger Documentation** - API docs interaktif untuk frontend integration

## 🚀 Quick Start

### Prerequisites

- [Go 1.21+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Docker & Docker Compose](https://www.docker.com/)

### 1. Clone Repository

```bash
git clone https://github.com/HackathonCodeCrafters/chat2pay.git
cd chat2pay
```

### 2. Setup Backend

```bash
cd backend-chat2pay

# Start PostgreSQL
docker-compose up -d postgres

# Create database
docker exec -it postgres_svc psql -U postgres -c "CREATE DATABASE chat2pay;"

# Copy & edit config
cp config/yaml/app.example.yaml config/yaml/app.yaml
# Edit app.yaml: tambahkan API keys (Gemini/OpenAI/Mistral)

# Run backend
go mod tidy
go run main.go migrate up
go run main.go http
```

Backend: `http://localhost:9005`

### 3. Setup Frontend

```bash
cd frontend-chat2pay

npm install
cp .env.example .env.local
# Edit .env.local: NEXT_PUBLIC_API_BASE_URL=http://localhost:9005

npm run dev
```

Frontend: `http://localhost:3000`

## 📚 API Documentation

Swagger UI: **http://localhost:9005/swagger/index.html**

## 🤖 AI Chat-to-Pay

```bash
curl -X POST http://localhost:9005/api/products/ask \
  -H "Content-Type: application/json" \
  -d '{"prompt": "Saya butuh laptop gaming budget 25 juta"}'
```

## 🐳 Docker

```bash
cd backend-chat2pay
docker-compose up -d
```

## 👥 Team

**Code Crafters** - Hackathon Team

## 📄 License

MIT License
