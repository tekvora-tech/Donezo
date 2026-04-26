# Donezo — Architecture

## 1. High-Level System Architecture

```mermaid
flowchart TB
    subgraph Client["Client Layer"]
        direction TB
        Browser["Browser"]
        FE["Vite + React + Tailwind
SPA (Single Page Application)"]
    end

    subgraph Network["Network Layer"]
        HTTPS["HTTPS / TLS 1.3"]
        CORS["CORS Middleware"]
    end

    subgraph API["API Gateway Layer"]
        Router["HTTP Router
Gin / Echo"]
        Middleware["Middleware Stack
• Logger
• Recovery
• Auth (JWT)
• Rate Limiter
• CORS"]
    end

    subgraph Backend["Backend Layer (Go)"]
        direction TB
        H["HTTP Handlers
(DTO Validation, Response Formatting)"]
        S["Service Layer
(Business Logic, Rules, Orchestration)"]
        R["Repository Layer
(SQL Queries, pgx)"]
    end

    subgraph Data["Data Layer"]
        DB[("Supabase PostgreSQL
pgxpool Connection Pool")]
    end

    subgraph External["External Services (Future)"]
        Redis[("Redis
Cache / Session (Opsional)")]
        Storage[("Supabase Storage
File Upload (Opsional)")]
    end

    Browser --> FE
    FE -->|API Request| HTTPS
    HTTPS --> CORS
    CORS --> Router
    Router --> Middleware
    Middleware --> H
    H --> S
    S --> R
    R -->|SQL Query| DB
    S -.->|Cache Read/Write| Redis
    S -.->|File Upload| Storage
```

## 2. Backend Layer Architecture (Go)

Backend mengikuti prinsip **Clean Architecture** dengan interface-driven design. Setiap layer hanya bergantung pada interface dari layer di bawahnya, bukan implementasi konkret.

```mermaid
flowchart TB
    subgraph Transport["Transport Layer"]
        direction TB
        Handler["HTTP Handler
• Request parsing
• DTO validation
• Response formatting
• Error handling"]
        DTO["DTO (Data Transfer Objects)
• Request structs
• Response structs
• Validation tags"]
    end

    subgraph Service["Service Layer"]
        direction TB
        ServiceInterface["Service Interface
(Contract)"]
        ServiceImpl["Service Implementation
• Business logic
• Rule validation
• Orchestration
• Transaction management"]
    end

    subgraph Repository["Repository Layer"]
        direction TB
        RepoInterface["Repository Interface
(Contract)"]
        RepoImpl["Repository Implementation
• SQL queries
• pgx execution
• Row mapping
• Transaction handling"]
    end

    subgraph Domain["Domain Layer"]
        Model["Domain Models
• Structs
• Constants
• Enums
• Custom types"]
    end

    Handler -->|uses| ServiceInterface
    ServiceImpl -->|implements| ServiceInterface
    ServiceImpl -->|uses| RepoInterface
    RepoImpl -->|implements| RepoInterface
    RepoImpl -->|uses| Model
    Handler -->|uses| DTO
    ServiceImpl -->|uses| Model
    ServiceImpl -->|returns| DTO
```

### 2.1 Layer Responsibilities

| Layer | Tanggung Jawab | Tidak Boleh |
|-------|-----------------|-------------|
| **Transport (Handler)** | Menerima HTTP request, parse input, validasi DTO, format response, handle HTTP-specific concerns | Tidak boleh mengandung business logic |
| **Service** | Menjalankan business logic, validasi aturan bisnis, orkestrasi antar repository, transaction management | Tidak boleh langsung query database |
| **Repository** | Eksekusi query SQL, mapping row ke struct, transaction handling di level database | Tidak boleh mengandung business logic |
| **Domain** | Definisi struct, constant, enum, custom types | Tidak boleh bergantung pada framework/library external |

### 2.2 Dependency Rule

```mermaid
flowchart LR
    subgraph Inner["Inner Layer (Independent)"]
        D["Domain"]
    end

    subgraph Middle["Middle Layer"]
        R["Repository Interface"]
        S["Service Interface"]
    end

    subgraph Outer["Outer Layer (Dependent)"]
        H["Handler"]
        RI["Repository Implementation"]
        SI["Service Implementation"]
    end

    H --> S
    SI --> R
    SI --> S
    RI --> R
    SI --> D
    RI --> D
    H --> D
```

**Aturan:** Dependency mengalir ke dalam (inward). Domain layer tidak bergantung pada layer lain. Interface didefinisikan di layer yang menggunakannya (Dependency Inversion Principle).

---

## 3. Request Lifecycle

```mermaid
sequenceDiagram
    participant C as Client (React)
    participant R as Router (Gin/Echo)
    participant M as Middleware
    participant H as Handler
    participant D as DTO/Validator
    participant S as Service
    participant RI as Repo Interface
    participant RP as Repo Implementation
    participant DB as PostgreSQL

    C->>R: HTTP Request
    R->>M: Pass through middleware
    M->>M: Logger, Recovery, CORS, Rate Limit
    M->>H: Authenticated Request
    H->>H: Parse request body/query/params
    H->>D: Validate DTO
    alt Validation Failed
        D-->>H: Validation Error
        H-->>C: 400 Bad Request
    else Validation Passed
        D-->>H: Validated DTO
        H->>S: Call service method
        S->>S: Apply business rules
        S->>RI: Call repository method
        RI->>RP: Interface dispatch
        RP->>DB: Execute SQL query
        DB-->>RP: Query result
        RP-->>RI: Mapped struct
        RI-->>S: Domain model
        S->>S: Transform to response DTO
        S-->>H: Response DTO
        H->>H: Wrap in envelope
        H-->>C: 200/201 Success Response
    end
```

---

## 4. Dependency Injection Flow

```mermaid
flowchart TB
    subgraph Main["main.go / wire.go"]
        direction TB
        DBPool[("pgxpool.Pool")]
        Config[("App Config")]
    end

    subgraph Repos["Repository Instantiation"]
        UserRepo[("userRepo
implements UserRepository")]
        TodoRepo[("todoRepo
implements TodoRepository")]
        TagRepo[("tagRepo
implements TagRepository")]
    end

    subgraph Services["Service Instantiation"]
        UserSvc[("userService
implements UserService")]
        TodoSvc[("todoService
implements TodoService")]
        AuthSvc[("authService
implements AuthService")]
    end

    subgraph Handlers["Handler Instantiation"]
        UserH[("userHandler")]
        TodoH[("todoHandler")]
        AuthH[("authHandler")]
    end

    subgraph Router["Router Setup"]
        R[("HTTP Router
/api/{group}/v1")]
    end

    DBPool --> UserRepo
    DBPool --> TodoRepo
    DBPool --> TagRepo

    UserRepo --> UserSvc
    TodoRepo --> TodoSvc
    TagRepo --> TodoSvc
    UserRepo --> AuthSvc

    UserSvc --> UserH
    TodoSvc --> TodoH
    AuthSvc --> AuthH

    UserH --> R
    TodoH --> R
    AuthH --> R
```

---

## 5. Authentication Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant H as Auth Handler
    participant S as Auth Service
    participant R as User Repository
    participant DB as PostgreSQL
    participant JWT as JWT Library

    %% Register
    C->>H: POST /api/auth/v1/register
    H->>S: Register(email, password, full_name)
    S->>S: Validate input
    S->>R: FindByEmail(email)
    R->>DB: SELECT * FROM users WHERE email = ?
    DB-->>R: No rows
    R-->>S: nil (email available)
    S->>S: Hash password (bcrypt)
    S->>R: Create(user)
    R->>DB: INSERT INTO users ...
    DB-->>R: User created
    R-->>S: User model
    S->>JWT: Generate tokens
    JWT-->>S: Access + Refresh tokens
    S-->>H: AuthResponse
    H-->>C: 201 Created + tokens

    %% Login
    C->>H: POST /api/auth/v1/login
    H->>S: Login(email, password)
    S->>R: FindByEmail(email)
    R->>DB: SELECT * FROM users WHERE email = ?
    DB-->>R: User row
    R-->>S: User model
    S->>S: Compare password (bcrypt)
    alt Password match
        S->>JWT: Generate tokens
        JWT-->>S: Access + Refresh tokens
        S-->>H: AuthResponse
        H-->>C: 200 OK + tokens
    else Password mismatch
        S-->>H: Unauthorized error
        H-->>C: 401 Unauthorized
    end
```

---

## 6. API URL Structure

```
/api/{group-service}/v1/{resource}
```

| Group | Base Path | Resources |
|-------|-----------|-----------|
| **Auth** | `/api/auth/v1` | `/register`, `/login`, `/refresh`, `/logout` |
| **Users** | `/api/users/v1` | `/profile`, `/profile/update` |
| **Todos** | `/api/todos/v1` | `/`, `/:id`, `/:id/sub-tasks` |
| **Tags** | `/api/tags/v1` | `/`, `/:id` |

```mermaid
flowchart LR
    subgraph API["API Root"]
        direction TB
        Root["/api"]
    end

    subgraph Groups["Service Groups"]
        direction TB
        Auth["/auth/v1"]
        Users["/users/v1"]
        Todos["/todos/v1"]
        Tags["/tags/v1"]
    end

    subgraph Resources["Resources"]
        direction TB
        A1["/register"]
        A2["/login"]
        A3["/refresh"]
        U1["/profile"]
        T1["/"]
        T2["/:id"]
        T3["/:id/sub-tasks"]
        TG1["/"]
        TG2["/:id"]
    end

    Root --> Auth
    Root --> Users
    Root --> Todos
    Root --> Tags

    Auth --> A1
    Auth --> A2
    Auth --> A3
    Users --> U1
    Todos --> T1
    Todos --> T2
    Todos --> T3
    Tags --> TG1
    Tags --> TG2
```

---

## 7. Middleware Stack

```mermaid
flowchart LR
    Request["HTTP Request"] --> Logger["Logger Middleware
• Request ID
• Timestamp
• Method & Path
• Duration"]
    Logger --> Recovery["Recovery Middleware
• Catch panics
• Return 500
• Log stack trace"]
    Recovery --> CORS["CORS Middleware
• Origin whitelist
• Methods & Headers"]
    CORS --> RateLimit["Rate Limiter
• IP-based limit
• Burst allowance"]
    RateLimit --> Auth["Auth Middleware
• Extract Bearer token
• Validate JWT
• Set user context"]
    Auth --> Handler["Handler"]
    Handler --> Response["HTTP Response"]
```

| Middleware | Urutan | Fungsi |
|------------|--------|--------|
| **Logger** | 1 | Log setiap request dengan detail (method, path, status, duration) |
| **Recovery** | 2 | Tangkap panic, prevent server crash, return 500 |
| **CORS** | 3 | Handle cross-origin requests sesuai whitelist |
| **Rate Limiter** | 4 | Batasi jumlah request per IP (khususnya auth endpoints) |
| **Auth (JWT)** | 5 | Validasi token, extract claims, set user context |

---

## 8. Error Handling Strategy

```mermaid
flowchart TB
    subgraph Errors["Error Types"]
        direction TB
        VE["ValidationError
(400 Bad Request)"]
        AE["AuthError
(401 Unauthorized)"]
        FE["ForbiddenError
(403 Forbidden)"]
        NE["NotFoundError
(404 Not Found)"]
        CE["ConflictError
(409 Conflict)"]
        BE["BusinessError
(422 Unprocessable)"]
        IE["InternalError
(500 Server Error)"]
    end

    subgraph Handler["Handler Layer"]
        SW["switch err.(type)"]
        Format["Format Response Envelope"]
    end

    VE --> SW
    AE --> SW
    FE --> SW
    NE --> SW
    CE --> SW
    BE --> SW
    IE --> SW
    SW --> Format
```

### 8.1 Error Mapping

| Error Type | HTTP Status | Contoh |
|------------|-------------|--------|
| `ValidationError` | `400` | Field required, invalid format |
| `AuthError` | `401` | Token invalid, token expired |
| `ForbiddenError` | `403` | Mengakses todo orang lain |
| `NotFoundError` | `404` | Todo tidak ditemukan |
| `ConflictError` | `409` | Email sudah terdaftar |
| `BusinessError` | `422` | Business rule violation |
| `InternalError` | `500` | Database error, unexpected error |

---

## 9. Key Architectural Decisions

| Keputusan | Implementasi | Alasan |
|-----------|-------------|--------|
| **Interface-driven** | Service & Repository menggunakan interface | Mudah mock untuk unit test; swap implementasi tanpa ubah business logic |
| **No ORM** | Raw SQL dengan pgx + scany (opsional) | Kontrol penuh query, performa optimal, kurva belajar rendah |
| **pgxpool** | Connection pooling built-in pgx | Performa tinggi, async notifications, native PostgreSQL |
| **JWT Stateless** | Access + Refresh token | API stateless, scalable horizontal, tidak perlu session store di MVP |
| **Supabase PostgreSQL** | Managed database | RLS siap aktif, backups otomatis, real-time potential |
| **Clean Architecture** | Layered dengan dependency inversion | Testable, maintainable, tidak terikat framework |
| **DTO Pattern** | Separate request/response structs | Decouple domain model dari API contract |

---

## 10. Scalability Considerations

### 10.1 Horizontal Scaling

```mermaid
flowchart TB
    subgraph LB["Load Balancer"]
        Nginx["Nginx / Cloudflare"]
    end

    subgraph AppServers["Application Servers"]
        App1["Go Instance 1"]
        App2["Go Instance 2"]
        App3["Go Instance 3"]
    end

    subgraph Data["Data Layer"]
        DB[("Supabase PG
Primary")]
        Cache[("Redis
(Opsional)")]
    end

    Nginx --> App1
    Nginx --> App2
    Nginx --> App3
    App1 --> DB
    App2 --> DB
    App3 --> DB
    App1 -.-> Cache
    App2 -.-> Cache
    App3 -.-> Cache
```

### 10.2 Caching Strategy (Future)

| Layer | Cache Key | TTL | Invalidation |
|-------|-----------|-----|--------------|
| **User Profile** | `user:{user_id}` | 5 menit | On update |
| **Todo List** | `todos:{user_id}:{filter_hash}` | 1 menit | On create/update/delete |
| **Tag List** | `tags:{user_id}` | 10 menit | On create/update/delete |

---

## 11. Testing Strategy

```mermaid
flowchart TB
    subgraph Unit["Unit Tests"]
        direction TB
        SvcTest["Service Tests
• Mock repository
• Test business logic
• Test edge cases"]
        HandlerTest["Handler Tests
• Mock service
• Test HTTP layer
• Test validation"]
    end

    subgraph Integration["Integration Tests"]
        direction TB
        RepoTest["Repository Tests
• Real test DB
• Test SQL queries
• Test transactions"]
        APITest["API Tests
• Spin up server
• Real HTTP calls
• End-to-end flow"]
    end

    subgraph E2E["E2E Tests"]
        direction TB
        FETest["Frontend Tests
• Cypress / Playwright
• User journey flows"]
    end

    SvcTest --> RepoTest
    HandlerTest --> APITest
    RepoTest --> APITest
    APITest --> FETest
```

| Tipe Test | Scope | Tools |
|-----------|-------|-------|
| **Unit Test** | Service layer (mock repo), Handler layer (mock service) | testify, mockery |
| **Integration Test** | Repository layer (test DB), API layer (test server) | testify, testcontainers-go |
| **E2E Test** | Full user journey (opsional) | Playwright / Cypress |

---

## 12. Monitoring & Observability (Future)

| Aspek | Tool | Deskripsi |
|-------|------|-----------|
| **Logging** | Structured JSON logs | Zap / Logrus dengan request ID |
| **Metrics** | Prometheus + Grafana | Request rate, latency, error rate |
| **Tracing** | OpenTelemetry / Jaeger | Distributed tracing untuk debug |
| **Health Check** | `/health` endpoint | Liveness & readiness probes |
| **Alerting** | PagerDuty / Slack | Alert untuk error rate > threshold |
