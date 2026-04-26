# Donezo — API Contract Detail

> Dokumentasi lengkap untuk seluruh endpoint API Donezo.
> Base URL: `{BASE_URL}/api/{group-service}/v1/{resource}`
> Semua endpoint (kecuali yang diberi label Public) memerlukan `Authorization: Bearer <access_token>`

---

## Daftar Isi

1. [Auth Endpoints](#1-auth-endpoints)
2. [User Endpoints](#2-user-endpoints)
3. [Todo Endpoints](#3-todo-endpoints)
4. [Sub-task Endpoints](#4-sub-task-endpoints)
5. [Tag Endpoints](#5-tag-endpoints)
6. [Shared Patterns](#6-shared-patterns)

---

## 1. Auth Endpoints

Group: `/api/auth/v1`

---

### 1.1 Register

Membuat akun pengguna baru. Setelah berhasil, server mengembalikan access token dan refresh token.

**Logic & Detail:**
- Server menerima email, password, dan full_name dari client.
- **Validasi Input:**
  - `email`: harus format email valid (regex: `^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$`).
  - `password`: minimal 8 karakter, harus mengandung minimal 1 huruf besar, 1 huruf kecil, 1 angka.
  - `full_name`: minimal 2 karakter, maksimal 100 karakter.
- **Cek Duplikat:** Server melakukan query `SELECT id FROM users WHERE email = $1`. Jika row ditemukan, kembalikan error 409.
- **Hash Password:** Password di-hash menggunakan `bcrypt` dengan cost factor 10. Hasil hash disimpan di kolom `password_hash`.
- **Generate UUID:** User ID di-generate otomatis oleh database (`gen_random_uuid()`).
- **Insert ke DB:** Query `INSERT INTO users (email, password_hash, full_name) VALUES ($1, $2, $3)`.
- **Generate JWT:** Setelah user berhasil dibuat, server generate:
  - **Access Token** (JWT): disimpan di memory client (localStorage), TTL 15 menit (production).
  - **Refresh Token** (JWT): disimpan di HTTP-only cookie, TTL 7 hari.
- **Token Claims:**
  - `sub`: user_id (UUID)
  - `email`: email user
  - `iat`: issued at timestamp
  - `exp`: expiration timestamp
  - `type`: `"access"` atau `"refresh"`
- **Simpen di mana:**
  - Access token → client memory (localStorage / Zustand store).
  - Refresh token → HTTP-only cookie (secure, sameSite=strict).
  - User data → database table `users`.

```http
POST /api/auth/v1/register
Content-Type: application/json

Request Body:
{
  "email": "john.doe@example.com",
  "password": "SecurePass123",
  "full_name": "John Doe"
}
```

**Validation Rules:**

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `email` | string | Yes | Valid email format, max 255 chars |
| `password` | string | Yes | Min 8 chars, 1 uppercase, 1 lowercase, 1 digit |
| `full_name` | string | Yes | Min 2 chars, max 100 chars |

**Response 201 Created:**
```json
{
  "message": "User registered successfully",
  "code": 201,
  "result": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john.doe@example.com",
      "full_name": "John Doe",
      "avatar_url": null,
      "created_at": "2026-04-26T20:00:00Z"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

**Response 400 Bad Request (Validation Error):**
```json
{
  "message": "Validation failed",
  "code": 400,
  "error": "password must contain at least one uppercase letter, one lowercase letter, and one digit"
}
```

**Response 409 Conflict (Email Exists):**
```json
{
  "message": "Email already registered",
  "code": 409,
  "error": "email john.doe@example.com is already in use"
}
```

**Response 422 (Business Rule):**
```json
{
  "message": "Registration failed",
  "code": 422,
  "error": "unable to create user account"
}
```

---

### 1.2 Login

Autentikasi pengguna dengan email dan password. Mengembalikan access token dan refresh token baru.

**Logic & Detail:**
- Server menerima `email` dan `password`.
- **Cari User:** Query `SELECT * FROM users WHERE email = $1 LIMIT 1`. Jika tidak ditemukan, kembalikan 401 (jangan kasih tahu apakah email salah atau password salah — security reason).
- **Verifikasi Password:** Bandingkan password plain dengan `password_hash` menggunakan `bcrypt.CompareHashAndPassword()`. Jika tidak cocok, kembalikan 401.
- **Generate Token:** Sama dengan register — generate access token + refresh token.
- **Update Last Login (opsional):** Bisa update kolom `updated_at` atau buat kolom `last_login_at` (jika ada di iterasi berikutnya).
- **Simpen di mana:**
  - Access token → client memory.
  - Refresh token → HTTP-only cookie.

```http
POST /api/auth/v1/login
Content-Type: application/json

Request Body:
{
  "email": "john.doe@example.com",
  "password": "SecurePass123"
}
```

**Validation Rules:**

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `email` | string | Yes | Valid email format |
| `password` | string | Yes | Min 8 chars |

**Response 200 OK:**
```json
{
  "message": "Login successful",
  "code": 200,
  "result": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john.doe@example.com",
      "full_name": "John Doe",
      "avatar_url": null
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 900
  }
}
```

**Response 401 Unauthorized:**
```json
{
  "message": "Invalid credentials",
  "code": 401,
  "error": "email or password is incorrect"
}
```

**Response 429 Too Many Requests:**
```json
{
  "message": "Too many login attempts",
  "code": 429,
  "error": "please try again after 5 minutes"
}
```

---

### 1.3 Refresh Token

Mendapatkan access token baru menggunakan refresh token. Dipanggil oleh client ketika access token expired (menerima 401 dari API).

**Logic & Detail:**
- Client mengirim refresh token via HTTP-only cookie (atau header `X-Refresh-Token` tergantung implementasi).
- **Validasi Refresh Token:** Verifikasi JWT signature, cek expiration (`exp` claim), cek `type` claim harus `"refresh"`.
- **Extract Claims:** Ambil `sub` (user_id) dari token.
- **Cek User:** Query `SELECT id FROM users WHERE id = $1`. Jika user tidak ditemukan (misal: dihapus), kembalikan 401.
- **Generate New Access Token:** Buat access token baru dengan TTL fresh (15 menit).
- **Rotate Refresh Token (opsional, recommended):** Generate refresh token baru juga, invalidate yang lama. Ini mencegah refresh token reuse attack.
- **Simpen di mana:**
  - Access token baru → replace di client memory.
  - Refresh token baru (jika di-rotate) → replace di HTTP-only cookie.

```http
POST /api/auth/v1/refresh
Cookie: refresh_token=eyJhbGciOiJIUzI1NiIs...
```

**Response 200 OK:**
```json
{
  "message": "Token refreshed successfully",
  "code": 200,
  "result": {
    "access_token": "eyJhbGciOiJIUzI1NiIsNEW...",
    "expires_in": 900
  }
}
```

**Response 401 Unauthorized (Invalid Refresh Token):**
```json
{
  "message": "Invalid refresh token",
  "code": 401,
  "error": "refresh token is invalid or expired"
}
```

---

### 1.4 Logout

Logout pengguna dan invalidate token.

**Logic & Detail:**
- Client mengirim access token di header `Authorization: Bearer <token>`.
- Server validasi access token.
- **Invalidate Token:** Caranya tergantung strategi:
  - **Stateless (MVP):** Tidak ada invalidasi server-side. Client cukup menghapus token dari memory. Refresh token dihapus dari cookie.
  - **Stateful (Future):** Simpan token di blacklist (Redis/database) dengan TTL sama dengan token expiration. Setiap request cek blacklist dulu.
- **Clear Cookie:** Server kirim response dengan cookie refresh_token yang expired (max-age=0).
- **Simpen di mana:**
  - Stateless: tidak ada data yang disimpan di server. Token tetap valid sampai expired tapi client sudah tidak punya.
  - Stateful (future): token masuk ke Redis blacklist dengan key `blacklist:<token_jti>`.

```http
POST /api/auth/v1/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response 200 OK:**
```json
{
  "message": "Logout successful",
  "code": 200,
  "result": null
}
```

**Response 401 Unauthorized:**
```json
{
  "message": "Unauthorized",
  "code": 401,
  "error": "invalid or missing token"
}
```

---

## 2. User Endpoints

Group: `/api/users/v1`

---

### 2.1 Get Profile

Mengambil data profil pengguna yang sedang login.

**Logic & Detail:**
- Middleware auth sudah mengekstrak user_id dari JWT dan menyimpannya di request context (`ctx.Value("user_id")`).
- Handler ambil user_id dari context.
- **Query DB:** `SELECT id, email, full_name, avatar_url, created_at, updated_at FROM users WHERE id = $1`.
- **Tidak perlu password_hash** di response (security).
- **Simpen di mana:**
  - Data user → table `users` di PostgreSQL.
  - User ID dari JWT → disimpan di request context oleh auth middleware.

```http
GET /api/users/v1/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response 200 OK:**
```json
{
  "message": "Profile retrieved successfully",
  "code": 200,
  "result": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john.doe@example.com",
    "full_name": "John Doe",
    "avatar_url": "https://cdn.donezo.app/avatars/john.jpg",
    "created_at": "2026-04-26T20:00:00Z",
    "updated_at": "2026-04-26T20:00:00Z"
  }
}
```

**Response 401 Unauthorized:**
```json
{
  "message": "Unauthorized",
  "code": 401,
  "error": "invalid or expired token"
}
```

---

### 2.2 Update Profile

Mengupdate data profil pengguna.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi Input:**
  - `full_name`: min 2 chars, max 100 chars.
  - `avatar_url`: valid URL format (opsional).
- **Query DB:** `UPDATE users SET full_name = $1, avatar_url = $2, updated_at = now() WHERE id = $3`.
- **Return updated data:** Query ulang atau return dari parameter.
- **Simpen di mana:**
  - Updated data → table `users` di PostgreSQL.

```http
PUT /api/users/v1/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

Request Body:
{
  "full_name": "John Doe Updated",
  "avatar_url": "https://cdn.donezo.app/avatars/john-new.jpg"
}
```

**Validation Rules:**

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `full_name` | string | No | Min 2 chars, max 100 chars |
| `avatar_url` | string | No | Valid URL format, max 500 chars |

**Response 200 OK:**
```json
{
  "message": "Profile updated successfully",
  "code": 200,
  "result": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john.doe@example.com",
    "full_name": "John Doe Updated",
    "avatar_url": "https://cdn.donezo.app/avatars/john-new.jpg",
    "created_at": "2026-04-26T20:00:00Z",
    "updated_at": "2026-04-26T21:00:00Z"
  }
}
```

**Response 400 Bad Request:**
```json
{
  "message": "Validation failed",
  "code": 400,
  "error": "full_name must be between 2 and 100 characters"
}
```

---

## 3. Todo Endpoints

Group: `/api/todos/v1`

---

### 3.1 Create Todo

Membuat todo baru untuk user yang sedang login.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi Input:**
  - `title`: wajib, min 3 chars, max 255 chars.
  - `description`: opsional, max 5000 chars.
  - `priority`: opsional, default `"medium"`. Harus salah satu dari: `low`, `medium`, `high`, `urgent`.
  - `due_date`: opsional, format ISO 8601 (`2026-05-01T17:00:00Z`). Harus di masa depan (opsional: bisa allow hari ini).
  - `tag_ids`: opsional, array of UUID. Setiap tag harus milik user yang sama (cek ownership).
- **Generate UUID:** Todo ID di-generate oleh DB (`gen_random_uuid()`).
- **Default Status:** `"pending"`.
- **Insert ke DB:**
  ```sql
  INSERT INTO todos (user_id, title, description, status, priority, due_date)
  VALUES ($1, $2, $3, 'pending', $4, $5)
  ```
- **Assign Tags (jika ada):** Untuk setiap `tag_id` di array, insert ke `todo_tags`:
  ```sql
  INSERT INTO todo_tags (todo_id, tag_id) VALUES ($1, $2)
  ```
  - **Ownership Check:** Sebelum insert ke `todo_tags`, cek `SELECT user_id FROM tags WHERE id = $1`. Jika `user_id` != user_id dari JWT, skip tag tersebut (atau return 403 jika strict).
- **Return:** Todo yang baru dibuat, termasuk tags yang berhasil di-assign.
- **Simpen di mana:**
  - Todo data → table `todos`.
  - Tag associations → table `todo_tags`.

```http
POST /api/todos/v1/
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

Request Body:
{
  "title": "Complete API documentation",
  "description": "Write detailed API contract for all endpoints",
  "priority": "high",
  "due_date": "2026-05-01T17:00:00Z",
  "tag_ids": [
    "660e8400-e29b-41d4-a716-446655440001",
    "660e8400-e29b-41d4-a716-446655440002"
  ]
}
```

**Validation Rules:**

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `title` | string | Yes | Min 3 chars, max 255 chars |
| `description` | string | No | Max 5000 chars |
| `priority` | string | No | Enum: `low`, `medium`, `high`, `urgent`. Default: `medium` |
| `due_date` | string | No | ISO 8601 format. Must be future date |
| `tag_ids` | []uuid | No | Max 10 tags. Each tag must belong to current user |

**Response 201 Created:**
```json
{
  "message": "Todo created successfully",
  "code": 201,
  "result": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Complete API documentation",
    "description": "Write detailed API contract for all endpoints",
    "status": "pending",
    "priority": "high",
    "due_date": "2026-05-01T17:00:00Z",
    "tags": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "Work",
        "color": "#3B82F6"
      },
      {
        "id": "660e8400-e29b-41d4-a716-446655440002",
        "name": "Urgent",
        "color": "#EF4444"
      }
    ],
    "sub_tasks": [],
    "created_at": "2026-04-26T21:00:00Z",
    "updated_at": "2026-04-26T21:00:00Z"
  }
}
```

**Response 400 Bad Request:**
```json
{
  "message": "Validation failed",
  "code": 400,
  "error": "title must be at least 3 characters"
}
```

**Response 403 Forbidden (Tag Ownership):**
```json
{
  "message": "Forbidden",
  "code": 403,
  "error": "tag 660e8400-e29b-41d4-a716-446655440001 does not belong to current user"
}
```

---

### 3.2 List Todos

Mengambil daftar todo milik user yang sedang login dengan filter, search, sort, dan pagination.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Query Building:** Server membangun SQL query secara dinamis berdasarkan query parameters:
  - Base query: `SELECT * FROM todos WHERE user_id = $1`
  - Tambahkan `AND` clause untuk setiap filter yang ada.
  - Untuk search: `AND title ILIKE '%' || $search || '%'` (case-insensitive).
  - Untuk due_date range: `AND due_date >= $from AND due_date <= $to`.
  - Untuk tag filter: `JOIN todo_tags tt ON t.id = tt.todo_id WHERE tt.tag_id = $tag_id`.
- **Sort:** `ORDER BY` field yang diminta + direction (`ASC`/`DESC`).
  - Default sort: `created_at DESC`.
  - Valid sort fields: `created_at`, `updated_at`, `due_date`, `priority`, `title`.
- **Pagination:**
  - `LIMIT $page_size OFFSET ($page - 1) * $page_size`.
  - Max page_size: 100. Jika client kirim > 100, clamp ke 100.
- **Count Total:** Query terpisah untuk menghitung total records:
  ```sql
  SELECT COUNT(*) FROM todos WHERE user_id = $1 [AND filters...]
  ```
- **Metadata Calculation:**
  - `total_pages = ceil(total_records / page_size)`
  - `has_next = current_page < total_pages`
  - `has_prev = current_page > 1`
- **Tags & Sub-tasks (opsional):** Untuk setiap todo di list, bisa:
  - **N+1 (tidak recommended):** Query tags dan sub-tasks per todo.
  - **JOIN (recommended):** Gunakan JSON aggregation:
    ```sql
    SELECT t.*,
      COALESCE(json_agg(DISTINCT jsonb_build_object('id', tg.id, 'name', tg.name, 'color', tg.color)) FILTER (WHERE tg.id IS NOT NULL), '[]') as tags,
      COALESCE(json_agg(DISTINCT jsonb_build_object('id', st.id, 'title', st.title, 'is_completed', st.is_completed)) FILTER (WHERE st.id IS NOT NULL), '[]') as sub_tasks
    FROM todos t
    LEFT JOIN todo_tags tt ON t.id = tt.todo_id
    LEFT JOIN tags tg ON tt.tag_id = tg.id
    LEFT JOIN sub_tasks st ON t.id = st.todo_id
    WHERE t.user_id = $1
    GROUP BY t.id
    ```
- **Simpen di mana:**
  - Data todo → table `todos`.
  - Filter parameters → query string (tidak disimpan di server).

```http
GET /api/todos/v1/?page=1&page_size=10&status=pending&priority=high&search=api&sort_by=due_date&sort_order=asc
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Query Parameters:**

| Parameter | Type | Default | Rules |
|-----------|------|---------|-------|
| `page` | int | `1` | Min 1 |
| `page_size` | int | `10` | Min 1, max 100 |
| `status` | string | - | Enum: `pending`, `in_progress`, `completed`, `cancelled` |
| `priority` | string | - | Enum: `low`, `medium`, `high`, `urgent` |
| `tag_id` | uuid | - | Valid UUID of user's tag |
| `search` | string | - | Max 100 chars. Search in title (case-insensitive) |
| `due_date_from` | date | - | ISO 8601 date format |
| `due_date_to` | date | - | ISO 8601 date format |
| `sort_by` | string | `created_at` | Enum: `created_at`, `updated_at`, `due_date`, `priority`, `title` |
| `sort_order` | string | `desc` | Enum: `asc`, `desc` |

**Response 200 OK:**
```json
{
  "message": "Todos retrieved successfully",
  "code": 200,
  "result": {
    "data": [
      {
        "id": "770e8400-e29b-41d4-a716-446655440000",
        "title": "Complete API documentation",
        "description": "Write detailed API contract for all endpoints",
        "status": "pending",
        "priority": "high",
        "due_date": "2026-05-01T17:00:00Z",
        "tags": [
          { "id": "660e8400-e29b-41d4-a716-446655440001", "name": "Work", "color": "#3B82F6" }
        ],
        "sub_task_count": 3,
        "completed_sub_task_count": 1,
        "created_at": "2026-04-26T21:00:00Z",
        "updated_at": "2026-04-26T21:00:00Z"
      }
    ],
    "metadata": {
      "current_page": 1,
      "page_size": 10,
      "total_pages": 5,
      "total_records": 50,
      "has_next": true,
      "has_prev": false
    }
  }
}
```

**Response 401 Unauthorized:**
```json
{
  "message": "Unauthorized",
  "code": 401,
  "error": "invalid or expired token"
}
```

---

### 3.3 Get Todo Detail

Mengambil detail lengkap satu todo termasuk sub-tasks dan tags.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi ID:** Pastikan `id` di path adalah UUID valid.
- **Query DB:**
  ```sql
  SELECT t.*,
    COALESCE(json_agg(DISTINCT jsonb_build_object('id', tg.id, 'name', tg.name, 'color', tg.color)) FILTER (WHERE tg.id IS NOT NULL), '[]') as tags,
    COALESCE(json_agg(DISTINCT jsonb_build_object('id', st.id, 'title', st.title, 'is_completed', st.is_completed, 'created_at', st.created_at)) FILTER (WHERE st.id IS NOT NULL), '[]') as sub_tasks
  FROM todos t
  LEFT JOIN todo_tags tt ON t.id = tt.todo_id
  LEFT JOIN tags tg ON tt.tag_id = tg.id
  LEFT JOIN sub_tasks st ON t.id = st.todo_id
  WHERE t.id = $1 AND t.user_id = $2
  GROUP BY t.id
  ```
- **Ownership Check:** `AND t.user_id = $2` memastikan user hanya bisa lihat todo miliknya sendiri. Jika query tidak return row, kembalikan 404 (jangan kasih tahu apakah todo tidak ada atau bukan miliknya — security).
- **Simpen di mana:**
  - Data todo → table `todos`.
  - Tags → table `tags` via `todo_tags`.
  - Sub-tasks → table `sub_tasks`.

```http
GET /api/todos/v1/770e8400-e29b-41d4-a716-446655440000
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response 200 OK:**
```json
{
  "message": "Todo retrieved successfully",
  "code": 200,
  "result": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Complete API documentation",
    "description": "Write detailed API contract for all endpoints",
    "status": "pending",
    "priority": "high",
    "due_date": "2026-05-01T17:00:00Z",
    "tags": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "Work",
        "color": "#3B82F6"
      }
    ],
    "sub_tasks": [
      {
        "id": "880e8400-e29b-41d4-a716-446655440000",
        "title": "Write auth endpoints",
        "is_completed": true,
        "created_at": "2026-04-26T21:30:00Z"
      },
      {
        "id": "880e8400-e29b-41d4-a716-446655440001",
        "title": "Write todo endpoints",
        "is_completed": false,
        "created_at": "2026-04-26T21:35:00Z"
      }
    ],
    "created_at": "2026-04-26T21:00:00Z",
    "updated_at": "2026-04-26T21:00:00Z"
  }
}
```

**Response 404 Not Found:**
```json
{
  "message": "Todo not found",
  "code": 404,
  "error": "todo with id 770e8400-e29b-41d4-a716-446655440000 does not exist"
}
```

---

### 3.4 Update Todo

Mengupdate data todo yang sudah ada.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi ID:** Pastikan `id` di path adalah UUID valid.
- **Ownership Check:** Query `SELECT user_id FROM todos WHERE id = $1`. Jika `user_id` != user_id dari JWT, kembalikan 403.
- **Validasi Input:** Sama dengan Create Todo, tapi semua field opsional (partial update). Hanya field yang dikirim yang akan di-update.
  - `title`: jika dikirim, min 3 chars, max 255 chars.
  - `description`: jika dikirim, max 5000 chars.
  - `status`: jika dikirim, harus valid enum.
  - `priority`: jika dikirim, harus valid enum.
  - `due_date`: jika dikirim, harus future date.
  - `tag_ids`: jika dikirim, replace semua tag yang ada (bukan append). Untuk menghapus semua tag, kirim array kosong `[]`.
- **Update DB:**
  ```sql
  UPDATE todos SET title = $1, description = $2, status = $3, priority = $4, due_date = $5, updated_at = now()
  WHERE id = $6 AND user_id = $7
  ```
- **Update Tags (jika tag_ids dikirim):**
  1. `DELETE FROM todo_tags WHERE todo_id = $1`
  2. Untuk setiap tag_id di array: `INSERT INTO todo_tags (todo_id, tag_id) VALUES ($1, $2)`
  3. Ownership check untuk setiap tag_id (sama seperti create).
- **Return:** Todo yang sudah di-update.
- **Simpen di mana:**
  - Updated data → table `todos`.
  - Tag associations → table `todo_tags` (di-replace, bukan append).

```http
PUT /api/todos/v1/770e8400-e29b-41d4-a716-446655440000
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

Request Body:
{
  "title": "Complete API documentation v2",
  "status": "in_progress",
  "priority": "urgent",
  "tag_ids": [
    "660e8400-e29b-41d4-a716-446655440001"
  ]
}
```

**Validation Rules:**

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `title` | string | No | Min 3 chars, max 255 chars |
| `description` | string | No | Max 5000 chars |
| `status` | string | No | Enum: `pending`, `in_progress`, `completed`, `cancelled` |
| `priority` | string | No | Enum: `low`, `medium`, `high`, `urgent` |
| `due_date` | string | No | ISO 8601 format. Future date |
| `tag_ids` | []uuid | No | Max 10 tags. Replace existing tags. Empty array = remove all tags |

**Response 200 OK:**
```json
{
  "message": "Todo updated successfully",
  "code": 200,
  "result": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Complete API documentation v2",
    "description": "Write detailed API contract for all endpoints",
    "status": "in_progress",
    "priority": "urgent",
    "due_date": "2026-05-01T17:00:00Z",
    "tags": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "Work",
        "color": "#3B82F6"
      }
    ],
    "sub_tasks": [
      {
        "id": "880e8400-e29b-41d4-a716-446655440000",
        "title": "Write auth endpoints",
        "is_completed": true,
        "created_at": "2026-04-26T21:30:00Z"
      }
    ],
    "created_at": "2026-04-26T21:00:00Z",
    "updated_at": "2026-04-26T22:00:00Z"
  }
}
```

**Response 403 Forbidden:**
```json
{
  "message": "Forbidden",
  "code": 403,
  "error": "you do not have permission to update this todo"
}
```

**Response 404 Not Found:**
```json
{
  "message": "Todo not found",
  "code": 404,
  "error": "todo with id 770e8400-e29b-41d4-a716-446655440000 does not exist"
}
```

---

### 3.5 Delete Todo

Menghapus todo dan semua data terkait (sub-tasks, tag associations).

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi ID:** Pastikan `id` di path adalah UUID valid.
- **Ownership Check:** Query `SELECT user_id FROM todos WHERE id = $1`. Jika `user_id` != user_id dari JWT, kembalikan 403.
- **Delete DB:**
  ```sql
  DELETE FROM todos WHERE id = $1 AND user_id = $2
  ```
- **Cascading Delete:** Karena foreign keys di `sub_tasks` dan `todo_tags` sudah di-set `ON DELETE CASCADE`, PostgreSQL akan otomatis menghapus:
  - Semua sub_tasks yang terkait dengan todo ini.
  - Semua todo_tags junction records.
  - Tags itu sendiri **TIDAK** dihapus (tags tetap ada di table `tags`).
- **Return:** 204 No Content (tidak ada body response).
- **Simpen di mana:**
  - Data dihapus dari table `todos`, `sub_tasks`, `todo_tags`.
  - Tags tetap ada di table `tags`.

```http
DELETE /api/todos/v1/770e8400-e29b-41d4-a716-446655440000
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response 204 No Content:**
```
(no body)
```

**Response 403 Forbidden:**
```json
{
  "message": "Forbidden",
  "code": 403,
  "error": "you do not have permission to delete this todo"
}
```

**Response 404 Not Found:**
```json
{
  "message": "Todo not found",
  "code": 404,
  "error": "todo with id 770e8400-e29b-41d4-a716-446655440000 does not exist"
}
```

---

## 4. Sub-task Endpoints

Group: `/api/todos/v1/:todo_id/sub-tasks`

---

### 4.1 Create Sub-task

Menambahkan sub-task ke dalam todo tertentu.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi todo_id:** Pastikan UUID valid.
- **Ownership Check:** Query `SELECT user_id FROM todos WHERE id = $1`. Jika `user_id` != user_id dari JWT, kembalikan 403.
- **Validasi Input:**
  - `title`: wajib, min 1 char, max 255 chars.
- **Generate UUID:** Sub-task ID di-generate oleh DB.
- **Default is_completed:** `false`.
- **Insert ke DB:**
  ```sql
  INSERT INTO sub_tasks (todo_id, title, is_completed) VALUES ($1, $2, false)
  ```
- **Return:** Sub-task yang baru dibuat.
- **Simpen di mana:**
  - Sub-task data → table `sub_tasks`.

```http
POST /api/todos/v1/770e8400-e29b-41d4-a716-446655440000/sub-tasks
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

Request Body:
{
  "title": "Write database schema section"
}
```

**Validation Rules:**

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `title` | string | Yes | Min 1 char, max 255 chars |

**Response 201 Created:**
```json
{
  "message": "Sub-task created successfully",
  "code": 201,
  "result": {
    "id": "880e8400-e29b-41d4-a716-446655440002",
    "todo_id": "770e8400-e29b-41d4-a716-446655440000",
    "title": "Write database schema section",
    "is_completed": false,
    "created_at": "2026-04-26T22:30:00Z"
  }
}
```

**Response 403 Forbidden:**
```json
{
  "message": "Forbidden",
  "code": 403,
  "error": "you do not have permission to add sub-tasks to this todo"
}
```

**Response 404 Not Found (Todo tidak ada):**
```json
{
  "message": "Todo not found",
  "code": 404,
  "error": "todo with id 770e8400-e29b-41d4-a716-446655440000 does not exist"
}
```

---

### 4.2 List Sub-tasks

Mengambil daftar sub-task untuk todo tertentu.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi todo_id:** UUID valid.
- **Ownership Check:** Query `SELECT user_id FROM todos WHERE id = $1`.
- **Query DB:**
  ```sql
  SELECT * FROM sub_tasks WHERE todo_id = $1 ORDER BY created_at ASC
  ```
- **No Pagination:** Sub-task list tidak dipaginate (biasanya jumlahnya kecil).
- **Simpen di mana:**
  - Data sub-task → table `sub_tasks`.

```http
GET /api/todos/v1/770e8400-e29b-41d4-a716-446655440000/sub-tasks
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response 200 OK:**
```json
{
  "message": "Sub-tasks retrieved successfully",
  "code": 200,
  "result": {
    "data": [
      {
        "id": "880e8400-e29b-41d4-a716-446655440000",
        "todo_id": "770e8400-e29b-41d4-a716-446655440000",
        "title": "Write auth endpoints",
        "is_completed": true,
        "created_at": "2026-04-26T21:30:00Z"
      },
      {
        "id": "880e8400-e29b-41d4-a716-446655440001",
        "todo_id": "770e8400-e29b-41d4-a716-446655440000",
        "title": "Write todo endpoints",
        "is_completed": false,
        "created_at": "2026-04-26T21:35:00Z"
      }
    ]
  }
}
```

---

### 4.3 Update Sub-task

Mengupdate sub-task (biasanya untuk toggle `is_completed` atau ubah title).

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi IDs:** Pastikan `todo_id` dan `sub_task_id` adalah UUID valid.
- **Ownership Check (double):**
  1. Cek `SELECT user_id FROM todos WHERE id = $1` (todo_id dari path).
  2. Cek `SELECT todo_id FROM sub_tasks WHERE id = $1` (sub_task_id dari path).
  3. Pastikan todo_id dari sub_task sama dengan todo_id dari path.
  4. Pastikan user_id dari todo == user_id dari JWT.
- **Partial Update:** Hanya field yang dikirim yang di-update.
  - `title`: jika dikirim, min 1 char, max 255 chars.
  - `is_completed`: jika dikirim, harus boolean.
- **Update DB:**
  ```sql
  UPDATE sub_tasks SET title = $1, is_completed = $2 WHERE id = $3 AND todo_id = $4
  ```
- **Return:** Sub-task yang sudah di-update.
- **Simpen di mana:**
  - Updated data → table `sub_tasks`.

```http
PUT /api/todos/v1/770e8400-e29b-41d4-a716-446655440000/sub-tasks/880e8400-e29b-41d4-a716-446655440001
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

Request Body:
{
  "is_completed": true
}
```

**Validation Rules:**

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `title` | string | No | Min 1 char, max 255 chars |
| `is_completed` | boolean | No | Must be true or false |

**Response 200 OK:**
```json
{
  "message": "Sub-task updated successfully",
  "code": 200,
  "result": {
    "id": "880e8400-e29b-41d4-a716-446655440001",
    "todo_id": "770e8400-e29b-41d4-a716-446655440000",
    "title": "Write todo endpoints",
    "is_completed": true,
    "created_at": "2026-04-26T21:35:00Z"
  }
}
```

**Response 403 Forbidden:**
```json
{
  "message": "Forbidden",
  "code": 403,
  "error": "you do not have permission to update this sub-task"
}
```

**Response 404 Not Found:**
```json
{
  "message": "Sub-task not found",
  "code": 404,
  "error": "sub-task with id 880e8400-e29b-41d4-a716-446655440001 does not exist"
}
```

---

### 4.4 Delete Sub-task

Menghapus sub-task dari todo.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi IDs:** UUID valid untuk todo_id dan sub_task_id.
- **Ownership Check:** Sama dengan Update Sub-task (double check).
- **Delete DB:**
  ```sql
  DELETE FROM sub_tasks WHERE id = $1 AND todo_id = $2
  ```
- **Return:** 204 No Content.
- **Simpen di mana:**
  - Data dihapus dari table `sub_tasks`.

```http
DELETE /api/todos/v1/770e8400-e29b-41d4-a716-446655440000/sub-tasks/880e8400-e29b-41d4-a716-446655440001
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response 204 No Content:**
```
(no body)
```

---

## 5. Tag Endpoints

Group: `/api/tags/v1`

---

### 5.1 Create Tag

Membuat tag baru untuk user yang sedang login.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi Input:**
  - `name`: wajib, min 1 char, max 50 chars.
  - `color`: opsional, default `"#3B82F6"`. Harus format hex 6 digit (`#RRGGBB`).
- **Cek Duplikat:** Query `SELECT id FROM tags WHERE user_id = $1 AND name = $2`. Jika ditemukan, kembalikan 409.
- **Generate UUID:** Tag ID di-generate oleh DB.
- **Insert ke DB:**
  ```sql
  INSERT INTO tags (user_id, name, color) VALUES ($1, $2, $3)
  ```
- **Return:** Tag yang baru dibuat.
- **Simpen di mana:**
  - Tag data → table `tags`.

```http
POST /api/tags/v1/
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

Request Body:
{
  "name": "Personal",
  "color": "#10B981"
}
```

**Validation Rules:**

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `name` | string | Yes | Min 1 char, max 50 chars |
| `color` | string | No | Hex format `#RRGGBB`. Default: `#3B82F6` |

**Response 201 Created:**
```json
{
  "message": "Tag created successfully",
  "code": 201,
  "result": {
    "id": "660e8400-e29b-41d4-a716-446655440003",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Personal",
    "color": "#10B981",
    "created_at": "2026-04-26T23:00:00Z"
  }
}
```

**Response 409 Conflict:**
```json
{
  "message": "Tag already exists",
  "code": 409,
  "error": "tag with name 'Personal' already exists"
}
```

---

### 5.2 List Tags

Mengambil daftar tag milik user yang sedang login.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Query DB:**
  ```sql
  SELECT * FROM tags WHERE user_id = $1 ORDER BY created_at ASC
  ```
- **No Pagination:** Tag list biasanya kecil (< 100), tidak perlu pagination.
- **Return:** Array of tags.
- **Simpen di mana:**
  - Data tag → table `tags`.

```http
GET /api/tags/v1/
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response 200 OK:**
```json
{
  "message": "Tags retrieved successfully",
  "code": 200,
  "result": {
    "data": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "Work",
        "color": "#3B82F6",
        "created_at": "2026-04-26T20:00:00Z"
      },
      {
        "id": "660e8400-e29b-41d4-a716-446655440002",
        "name": "Urgent",
        "color": "#EF4444",
        "created_at": "2026-04-26T20:05:00Z"
      },
      {
        "id": "660e8400-e29b-41d4-a716-446655440003",
        "name": "Personal",
        "color": "#10B981",
        "created_at": "2026-04-26T23:00:00Z"
      }
    ]
  }
}
```

---

### 5.3 Update Tag

Mengupdate tag yang sudah ada.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi ID:** UUID valid.
- **Ownership Check:** Query `SELECT user_id FROM tags WHERE id = $1`. Jika `user_id` != user_id dari JWT, kembalikan 403.
- **Validasi Input:**
  - `name`: jika dikirim, min 1 char, max 50 chars.
  - `color`: jika dikirim, hex format.
- **Cek Duplikat (jika name diubah):** Query `SELECT id FROM tags WHERE user_id = $1 AND name = $2 AND id != $3`. Jika ditemukan, kembalikan 409.
- **Update DB:**
  ```sql
  UPDATE tags SET name = $1, color = $2 WHERE id = $3 AND user_id = $4
  ```
- **Return:** Tag yang sudah di-update.
- **Simpen di mana:**
  - Updated data → table `tags`.

```http
PUT /api/tags/v1/660e8400-e29b-41d4-a716-446655440003
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

Request Body:
{
  "color": "#8B5CF6"
}
```

**Validation Rules:**

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `name` | string | No | Min 1 char, max 50 chars |
| `color` | string | No | Hex format `#RRGGBB` |

**Response 200 OK:**
```json
{
  "message": "Tag updated successfully",
  "code": 200,
  "result": {
    "id": "660e8400-e29b-41d4-a716-446655440003",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Personal",
    "color": "#8B5CF6",
    "created_at": "2026-04-26T23:00:00Z"
  }
}
```

**Response 403 Forbidden:**
```json
{
  "message": "Forbidden",
  "code": 403,
  "error": "you do not have permission to update this tag"
}
```

**Response 409 Conflict:**
```json
{
  "message": "Tag already exists",
  "code": 409,
  "error": "tag with name 'Personal' already exists"
}
```

---

### 5.4 Delete Tag

Menghapus tag. Semua asosiasi dengan todo akan otomatis terhapus karena `ON DELETE CASCADE`.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi ID:** UUID valid.
- **Ownership Check:** Query `SELECT user_id FROM tags WHERE id = $1`.
- **Delete DB:**
  ```sql
  DELETE FROM tags WHERE id = $1 AND user_id = $2
  ```
- **Cascading Effect:** Karena `todo_tags.tag_id` punya `ON DELETE CASCADE`, semua record di `todo_tags` yang pakai tag ini akan otomatis terhapus. Todo itu sendiri **TIDAK** terhapus.
- **Return:** 204 No Content.
- **Simpen di mana:**
  - Tag dihapus dari table `tags`.
  - Junction records dihapus dari table `todo_tags`.

```http
DELETE /api/tags/v1/660e8400-e29b-41d4-a716-446655440003
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response 204 No Content:**
```
(no body)
```

**Response 403 Forbidden:**
```json
{
  "message": "Forbidden",
  "code": 403,
  "error": "you do not have permission to delete this tag"
}
```

---

### 5.5 Assign Tag to Todo

Menambahkan tag ke todo (bisa untuk tag yang belum di-assign).

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi IDs:** UUID valid untuk todo_id dan tag_id.
- **Ownership Check (double):**
  1. Cek `SELECT user_id FROM todos WHERE id = $1`.
  2. Cek `SELECT user_id FROM tags WHERE id = $1`.
  3. Pastikan kedua user_id sama dengan user_id dari JWT.
- **Cek Duplikat:** Query `SELECT 1 FROM todo_tags WHERE todo_id = $1 AND tag_id = $2`. Jika ditemukan, kembalikan 409 (atau 200 OK dengan message "already assigned").
- **Insert ke DB:**
  ```sql
  INSERT INTO todo_tags (todo_id, tag_id) VALUES ($1, $2)
  ```
- **Return:** Todo yang sudah di-update dengan tag baru.
- **Simpen di mana:**
  - Junction record → table `todo_tags`.

```http
POST /api/todos/v1/770e8400-e29b-41d4-a716-446655440000/tags
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

Request Body:
{
  "tag_id": "660e8400-e29b-41d4-a716-446655440003"
}
```

**Validation Rules:**

| Field | Type | Required | Rules |
|-------|------|----------|-------|
| `tag_id` | uuid | Yes | Valid UUID, must belong to current user |

**Response 200 OK:**
```json
{
  "message": "Tag assigned to todo successfully",
  "code": 200,
  "result": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "title": "Complete API documentation v2",
    "tags": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "Work",
        "color": "#3B82F6"
      },
      {
        "id": "660e8400-e29b-41d4-a716-446655440003",
        "name": "Personal",
        "color": "#8B5CF6"
      }
    ]
  }
}
```

**Response 409 Conflict:**
```json
{
  "message": "Tag already assigned",
  "code": 409,
  "error": "tag is already assigned to this todo"
}
```

---

### 5.6 Remove Tag from Todo

Menghapus tag dari todo.

**Logic & Detail:**
- Middleware auth ekstrak user_id dari JWT.
- **Validasi IDs:** UUID valid untuk todo_id dan tag_id.
- **Ownership Check:** Sama dengan Assign Tag (double check).
- **Delete DB:**
  ```sql
  DELETE FROM todo_tags WHERE todo_id = $1 AND tag_id = $2
  ```
- **Return:** Todo yang sudah di-update (tanpa tag yang dihapus).
- **Simpen di mana:**
  - Junction record dihapus dari table `todo_tags`.

```http
DELETE /api/todos/v1/770e8400-e29b-41d4-a716-446655440000/tags/660e8400-e29b-41d4-a716-446655440003
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

**Response 200 OK:**
```json
{
  "message": "Tag removed from todo successfully",
  "code": 200,
  "result": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "title": "Complete API documentation v2",
    "tags": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "name": "Work",
        "color": "#3B82F6"
      }
    ]
  }
}
```

**Response 404 Not Found:**
```json
{
  "message": "Tag not assigned",
  "code": 404,
  "error": "tag is not assigned to this todo"
}
```

---

## 6. Shared Patterns

### 6.1 Authentication Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant Auth as Auth Middleware
    participant H as Handler
    participant S as Service
    participant R as Repository
    participant DB as PostgreSQL

    C->>Auth: Request + Authorization: Bearer <token>
    Auth->>Auth: Extract token from header
    Auth->>Auth: Verify JWT signature
    Auth->>Auth: Check expiration (exp claim)
    alt Token Invalid/Expired
        Auth-->>C: 401 Unauthorized
    else Token Valid
        Auth->>Auth: Extract user_id from claims
        Auth->>H: Pass request with user_id in context
        H->>S: Call service with user_id
        S->>R: Query with user_id filter
        R->>DB: Execute SQL
        DB-->>R: Result
        R-->>S: Data
        S-->>H: Response
        H-->>C: 200 OK + data
    end
```

### 6.2 Ownership Check Pattern

Setiap endpoint yang mengakses resource spesifik (by ID) harus melakukan ownership check:

```mermaid
flowchart TD
    A["Extract user_id from JWT"] --> B["Parse resource ID from URL"]
    B --> C{"Resource exists?"}
    C -->|No| D["Return 404"]
    C -->|Yes| E{"resource.user_id == jwt.user_id?"}
    E -->|No| F["Return 403"]
    E -->|Yes| G["Proceed with operation"]
```

**Catatan Keamanan:**
- Selalu return 404 (bukan 403) jika resource tidak ditemukan — jangan kasih tahu attacker bahwa resource itu ada tapi bukan miliknya.
- Ownership check dilakukan di service layer (bukan hanya di query SQL) untuk memisahkan concerns.

### 6.3 Pagination Query Pattern

```mermaid
flowchart TD
    A["Receive page & page_size params"] --> B{"page < 1?"}
    B -->|Yes| C["Set page = 1"]
    B -->|No| D{"page_size > 100?"}
    C --> D
    D -->|Yes| E["Set page_size = 100"]
    D -->|No| F["Calculate offset = (page - 1) * page_size"]
    E --> F
    F --> G["Execute query with LIMIT & OFFSET"]
    G --> H["Execute COUNT query"]
    H --> I["Calculate metadata"]
    I --> J["Return data + metadata"]
```

### 6.4 Error Response Mapping

| Error Type | HTTP Status | Kapan Terjadi |
|------------|-------------|---------------|
| `ValidationError` | `400` | Input tidak valid (format, length, enum) |
| `AuthError` | `401` | Token tidak ada, invalid, atau expired |
| `ForbiddenError` | `403` | Token valid tapi bukan pemilik resource |
| `NotFoundError` | `404` | Resource tidak ditemukan |
| `ConflictError` | `409` | Resource sudah ada (duplicate) |
| `BusinessError` | `422` | Business rule violation |
| `InternalError` | `500` | Database error, unexpected error |

### 6.5 Request ID Tracking

Setiap request HTTP diberi unique request ID oleh Logger middleware:

```
X-Request-ID: req_abc123xyz
```

Request ID ini:
- Dilog di setiap layer (handler, service, repository)
- Dikirim balik ke client di response header
- Digunakan untuk tracing dan debugging

### 6.6 Data Storage Summary

| Data | Tabel | Kolom Utama | Relasi |
|------|-------|-------------|--------|
| User | `users` | `id`, `email`, `password_hash`, `full_name` | - |
| Todo | `todos` | `id`, `user_id`, `title`, `status`, `priority`, `due_date` | FK → `users.id` |
| Sub-task | `sub_tasks` | `id`, `todo_id`, `title`, `is_completed` | FK → `todos.id` |
| Tag | `tags` | `id`, `user_id`, `name`, `color` | FK → `users.id` |
| Todo-Tag | `todo_tags` | `todo_id`, `tag_id` | FK → `todos.id`, FK → `tags.id` |
| JWT Token | Client memory / Cookie | - | - |
| Blacklist Token | Redis (future) | `blacklist:<jti>` | - |

### 6.7 cURL Examples

**Register:**
```bash
curl -X POST http://localhost:8080/api/auth/v1/register   -H "Content-Type: application/json"   -d '{
    "email": "john@example.com",
    "password": "SecurePass123",
    "full_name": "John Doe"
  }'
```

**Create Todo:**
```bash
curl -X POST http://localhost:8080/api/todos/v1/   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."   -H "Content-Type: application/json"   -d '{
    "title": "Complete API docs",
    "priority": "high",
    "due_date": "2026-05-01T17:00:00Z"
  }'
```

**List Todos with Filter:**
```bash
curl "http://localhost:8080/api/todos/v1/?page=1&page_size=10&status=pending&priority=high&search=api"   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Toggle Sub-task:**
```bash
curl -X PUT http://localhost:8080/api/todos/v1/770e8400-e29b-41d4-a716-446655440000/sub-tasks/880e8400-e29b-41d4-a716-446655440001   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."   -H "Content-Type: application/json"   -d '{"is_completed": true}'
```
