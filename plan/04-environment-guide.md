# Donezo — Environment Guide

## 1. Environment Overview

Donezo mendukung tiga environment utama: **Local**, **Staging**, dan **Production**. Setiap environment memiliki konfigurasi yang berbeda untuk database, JWT, CORS, dan security settings.

```mermaid
flowchart LR
    subgraph Dev["Development"]
        Local["Local
Developer Machine"]
    end

    subgraph Test["Testing"]
        Staging["Staging
Pre-production"]
    end

    subgraph Prod["Production"]
        Production["Production
Live Users"]
    end

    Local -->|Promote| Staging
    Staging -->|Promote| Production
```

| Environment    | Tujuan             | Database                           | JWT TTL  | CORS             | SSL      |
| -------------- | ------------------ | ---------------------------------- | -------- | ---------------- | -------- |
| **Local**      | Pengembangan lokal | Supabase local / Docker PostgreSQL | 24 jam   | `localhost:5173` | Optional |
| **Staging**    | Testing & QA       | Supabase staging project           | 1 jam    | Staging domain   | Required |
| **Production** | Pengguna aktif     | Supabase production project        | 15 menit | `donezo.app`     | Required |

---

## 2. Backend Environment Variables

### 2.1 File: `.env.local` (Local Development)

```bash
# ==========================================
# Application
# ==========================================
APP_ENV=development
APP_PORT=8080
APP_BASE_URL=http://localhost:8080
APP_NAME=Donezo

# ==========================================
# Database (Supabase PostgreSQL)
# ==========================================
DB_HOST=db.xxx.supabase.co
DB_PORT=5432
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=your_supabase_password
DB_SSL_MODE=require
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=30m

# ==========================================
# JWT Authentication
# ==========================================
JWT_SECRET=your-local-256-bit-secret-key-here-min-32-chars
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h

# ==========================================
# CORS
# ==========================================
ALLOWED_ORIGINS=http://localhost:5173

# ==========================================
# Rate Limiting
# ==========================================
RATE_LIMIT_REQUESTS_PER_MINUTE=60
RATE_LIMIT_AUTH_REQUESTS_PER_MINUTE=10
```

### 2.2 File: `.env.staging`

```bash
# ==========================================
# Application
# ==========================================
APP_ENV=staging
APP_PORT=8080
APP_BASE_URL=https://api-staging.donezo.app
APP_NAME=Donezo

# ==========================================
# Database (Supabase Staging)
# ==========================================
DB_HOST=db.xxx.staging.supabase.co
DB_PORT=5432
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=your_staging_password
DB_SSL_MODE=require
DB_MAX_OPEN_CONNS=50
DB_MAX_IDLE_CONNS=20
DB_CONN_MAX_LIFETIME=30m

# ==========================================
# JWT Authentication
# ==========================================
JWT_SECRET=your-staging-256-bit-secret-key-here
JWT_EXPIRATION=1h
JWT_REFRESH_EXPIRATION=24h

# ==========================================
# CORS
# ==========================================
ALLOWED_ORIGINS=https://staging.donezo.app

# ==========================================
# Rate Limiting
# ==========================================
RATE_LIMIT_REQUESTS_PER_MINUTE=100
RATE_LIMIT_AUTH_REQUESTS_PER_MINUTE=20
```

### 2.3 File: `.env.production`

```bash
# ==========================================
# Application
# ==========================================
APP_ENV=production
APP_PORT=8080
APP_BASE_URL=https://api.donezo.app
APP_NAME=Donezo

# ==========================================
# Database (Supabase Production)
# ==========================================
DB_HOST=db.xxx.production.supabase.co
DB_PORT=5432
DB_NAME=postgres
DB_USER=postgres
DB_PASSWORD=your_production_password
DB_SSL_MODE=require
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=50
DB_CONN_MAX_LIFETIME=30m

# ==========================================
# JWT Authentication
# ==========================================
JWT_SECRET=your-production-256-bit-secret-key-here
JWT_EXPIRATION=15m
JWT_REFRESH_EXPIRATION=168h

# ==========================================
# CORS
# ==========================================
ALLOWED_ORIGINS=https://donezo.app

# ==========================================
# Rate Limiting
# ==========================================
RATE_LIMIT_REQUESTS_PER_MINUTE=120
RATE_LIMIT_AUTH_REQUESTS_PER_MINUTE=20
```

---

## 3. Frontend Environment Variables

### 3.1 File: `.env.local` (Local Development)

```bash
# ==========================================
# API Configuration
# ==========================================
VITE_API_BASE_URL=http://localhost:8080/api

# ==========================================
# App Configuration
# ==========================================
VITE_APP_NAME=Donezo
VITE_APP_ENV=development
```

### 3.2 File: `.env.staging`

```bash
# ==========================================
# API Configuration
# ==========================================
VITE_API_BASE_URL=https://api-staging.donezo.app/api

# ==========================================
# App Configuration
# ==========================================
VITE_APP_NAME=Donezo
VITE_APP_ENV=staging
```

### 3.3 File: `.env.production`

```bash
# ==========================================
# API Configuration
# ==========================================
VITE_API_BASE_URL=https://api.donezo.app/api

# ==========================================
# App Configuration
# ==========================================
VITE_APP_NAME=Donezo
VITE_APP_ENV=production
```

---

## 4. Environment Variable Reference

### 4.1 Backend Variables

| Variable                              | Type     | Required | Default    | Deskripsi                                                  |
| ------------------------------------- | -------- | -------- | ---------- | ---------------------------------------------------------- |
| `APP_ENV`                             | string   | Yes      | -          | Environment: `development`, `staging`, `production`        |
| `APP_PORT`                            | int      | Yes      | `8080`     | Port server HTTP                                           |
| `APP_BASE_URL`                        | string   | Yes      | -          | Base URL aplikasi                                          |
| `APP_NAME`                            | string   | No       | `Donezo`   | Nama aplikasi                                              |
| `DB_HOST`                             | string   | Yes      | -          | Host database PostgreSQL                                   |
| `DB_PORT`                             | int      | Yes      | `5432`     | Port database                                              |
| `DB_NAME`                             | string   | Yes      | `postgres` | Nama database                                              |
| `DB_USER`                             | string   | Yes      | -          | Username database                                          |
| `DB_PASSWORD`                         | string   | Yes      | -          | Password database                                          |
| `DB_SSL_MODE`                         | string   | Yes      | `require`  | Mode SSL: `disable`, `require`, `verify-ca`, `verify-full` |
| `DB_MAX_OPEN_CONNS`                   | int      | No       | `25`       | Maksimum open connections                                  |
| `DB_MAX_IDLE_CONNS`                   | int      | No       | `10`       | Maksimum idle connections                                  |
| `DB_CONN_MAX_LIFETIME`                | duration | No       | `30m`      | Maksimum lifetime per connection                           |
| `JWT_SECRET`                          | string   | Yes      | -          | Secret key untuk signing JWT (minimal 32 karakter)         |
| `JWT_EXPIRATION`                      | duration | No       | `24h`      | TTL access token                                           |
| `JWT_REFRESH_EXPIRATION`              | duration | No       | `168h`     | TTL refresh token (7 hari)                                 |
| `ALLOWED_ORIGINS`                     | string   | Yes      | -          | Daftar origin yang diizinkan (comma-separated)             |
| `RATE_LIMIT_REQUESTS_PER_MINUTE`      | int      | No       | `60`       | Rate limit umum per menit                                  |
| `RATE_LIMIT_AUTH_REQUESTS_PER_MINUTE` | int      | No       | `10`       | Rate limit auth endpoints per menit                        |

### 4.2 Frontend Variables

| Variable            | Type   | Required | Default       | Deskripsi                |
| ------------------- | ------ | -------- | ------------- | ------------------------ |
| `VITE_API_BASE_URL` | string | Yes      | -             | Base URL API backend     |
| `VITE_APP_NAME`     | string | No       | `Donezo`      | Nama aplikasi (untuk UI) |
| `VITE_APP_ENV`      | string | No       | `development` | Environment frontend     |

**Catatan:** Semua variabel frontend harus diawali dengan `VITE_` agar dapat diakses oleh Vite.

---

## 5. Supabase Setup Guide

### 5.1 Membuat Project Supabase

1. Buka [supabase.com](https://supabase.com) dan login
2. Klik "New Project"
3. Pilih organization dan beri nama project: `donezo-dev` (local), `donezo-staging`, `donezo-prod`
4. Pilih region terdekat dengan target audience (misal: `Southeast Asia` untuk Indonesia)
5. Set database password dan simpan dengan aman
6. Tunggu provisioning selesai

### 5.2 Mengambil Connection String

1. Masuk ke project dashboard
2. Pilih menu **Settings** → **Database**
3. Di bagian **Connection string**, pilih tab **URI**
4. Copy connection string dan ganti `[YOUR-PASSWORD]` dengan password database

```
postgresql://postgres:[YOUR-PASSWORD]@db.xxx.supabase.co:5432/postgres
```

### 5.3 Connection Pooling (PgBouncer)

Supabase menyediakan connection pooling via PgBouncer untuk mengelola banyak koneksi:

1. Di **Settings** → **Database**, cari **Connection pooling**
2. Gunakan connection string pooling untuk production:

```
postgresql://postgres:[YOUR-PASSWORD]@db.xxx.supabase.co:6543/postgres
```

**Catatan:** Gunakan port `6543` (PgBouncer) untuk production, bukan port `5432` langsung.

---

## 6. Local Development Setup

### 6.1 Prerequisites

| Tool              | Versi Min | Deskripsi                              |
| ----------------- | --------- | -------------------------------------- |
| Go                | 1.22      | Bahasa backend                         |
| Node.js           | 18+       | Runtime frontend                       |
| pnpm / npm / yarn | Latest    | Package manager                        |
| PostgreSQL        | 15+       | Database (atau gunakan Supabase local) |
| golang-migrate    | Latest    | Database migrations                    |
| Git               | Latest    | Version control                        |

### 6.2 Backend Setup (Local)

```bash
# 1. Clone repository
git clone https://github.com/tekvora-tech/Donezo.git
cd Donezo/backend

# 2. Copy environment file
cp .env.example .env.local

# 3. Edit .env.local dengan konfigurasi lokal
# (Isi DB_HOST, DB_PASSWORD, JWT_SECRET, dll.)

# 4. Install dependencies
go mod init
go mod tidy

# 5. Run database migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/donezo?sslmode=disable" up

# 6. Run server
go run cmd/api/main.go
# Server berjalan di http://localhost:8080
```

### 6.3 Frontend Setup (Local)

```bash
# 1. Navigate to frontend directory
cd donezo/frontend

# 2. Copy environment file
cp .env.example .env.local

# 3. Edit .env.local
# VITE_API_BASE_URL=http://localhost:8080/api

# 4. Install dependencies
pnpm install

# 5. Run development server
pnpm dev
# Server berjalan di http://localhost:5173
```

### 6.4 Docker Setup (Opsional)

```bash
# Docker Compose untuk local development
docker-compose -f docker-compose.dev.yml up -d

# Services yang dijalankan:
# - PostgreSQL (port 5432)
# - Redis (port 6379) - opsional
# - Backend (port 8080)
# - Frontend (port 5173)
```

---

## 7. Environment Security

### 7.1 Secret Management

| Secret                     | Storage                               | Catatan                     |
| -------------------------- | ------------------------------------- | --------------------------- |
| `DB_PASSWORD`              | Environment variable / Secret manager | Jangan commit ke repo       |
| `JWT_SECRET`               | Environment variable / Secret manager | Minimal 32 karakter, random |
| `DB_PASSWORD` (Production) | AWS Secrets Manager / HashiCorp Vault | Rotate secara berkala       |

### 7.2 `.gitignore`

```gitignore
# Environment files
.env
.env.local
.env.staging
.env.production
.env.*.local

# Binary
*.exe
*.dll
*.so
*.dylib
bin/
dist/

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db
```

### 7.3 JWT Secret Generation

```bash
# Generate secure JWT secret
openssl rand -base64 32

# Atau menggunakan pwgen
pwgen -s 32 1
```

---

## 8. Environment Comparison Matrix

| Aspek               | Local                   | Staging                  | Production                  |
| ------------------- | ----------------------- | ------------------------ | --------------------------- |
| **Database**        | Supabase local / Docker | Supabase project staging | Supabase project production |
| **SSL**             | Optional                | Required                 | Required                    |
| **Logging Level**   | Debug                   | Info                     | Info / Warn                 |
| **CORS Origins**    | `localhost:5173`        | Staging domain           | Production domain           |
| **JWT Access TTL**  | 24 jam                  | 1 jam                    | 15 menit                    |
| **JWT Refresh TTL** | 7 hari                  | 1 hari                   | 7 hari                      |
| **Rate Limit**      | 60 req/min              | 100 req/min              | 120 req/min                 |
| **Auth Rate Limit** | 10 req/min              | 20 req/min               | 20 req/min                  |
| **DB Connections**  | 25 open / 10 idle       | 50 open / 20 idle        | 100 open / 50 idle          |
| **Error Detail**    | Full stack trace        | Limited                  | Minimal (user-friendly)     |
| **Auto-reload**     | Yes (air/nodemon)       | No                       | No                          |
| **Monitoring**      | Console logs            | Basic metrics            | Full observability          |

---

## 9. Deployment Pipeline (Overview)

```mermaid
flowchart LR
    subgraph Dev["Development"]
        Code["Code Changes"]
        Lint["Lint & Format"]
        Test["Unit Tests"]
    end

    subgraph CI["CI Pipeline"]
        Build["Build Binary"]
        Docker["Build Docker Image"]
        Push["Push to Registry"]
    end

    subgraph CD["CD Pipeline"]
        StagingDeploy["Deploy to Staging"]
        SmokeTest["Smoke Tests"]
        ProdDeploy["Deploy to Production"]
    end

    Code --> Lint
    Lint --> Test
    Test --> Build
    Build --> Docker
    Docker --> Push
    Push --> StagingDeploy
    StagingDeploy --> SmokeTest
    SmokeTest --> ProdDeploy
```

---

## 10. Troubleshooting

### 10.1 Common Issues

| Issue                  | Penyebab                     | Solusi                                                        |
| ---------------------- | ---------------------------- | ------------------------------------------------------------- |
| `connection refused`   | Database tidak berjalan      | Cek status PostgreSQL / Supabase                              |
| `SSL is not enabled`   | SSL mode mismatch            | Set `DB_SSL_MODE=require` untuk Supabase                      |
| `too many connections` | Connection pool habis        | Naikkan `DB_MAX_OPEN_CONNS` atau gunakan PgBouncer            |
| `CORS error`           | Origin tidak di-whitelist    | Cek `ALLOWED_ORIGINS` di backend                              |
| `JWT invalid`          | Secret mismatch atau expired | Cek `JWT_SECRET` dan expiration time                          |
| `migration failed`     | Schema conflict              | Jalankan `migrate down` lalu `migrate up`, atau periksa versi |

### 10.2 Health Check Endpoint

```
GET /health

Response:
{
  "status": "healthy",
  "timestamp": "2026-04-26T20:50:00Z",
  "version": "1.0.0",
  "database": "connected"
}
```
