# Warung POS - Backend (Go)

Minimal backend implemented in Go with layered architecture (controllers, services, repositories, models).

Requirements

- Go 1.20+
- MySQL (XAMPP) running on localhost:3306 (default)

Quick start

1. Start MySQL (XAMPP). Create a database (default name used is `warung_pos`). Example SQL:

```sql
CREATE DATABASE IF NOT EXISTS warung_pos CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
```

2. (Optional) Create a `.env` file in `backend/` with values. Defaults are shown below.

```
DB_USER=root
DB_PASS=
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=warung_pos
JWT_SECRET=your-secret
PORT=8080
```

3. Install dependencies and build:

```bash
cd backend
go mod tidy
go build ./...
```

4. Run the server:

```bash
go run ./main.go
```

The server exposes simple endpoints under `/api`:

- POST /api/auth/register
- POST /api/auth/login
- GET /api/categories
- POST /api/categories
- GET /api/menus
- POST /api/menus
- GET /api/transactions
- POST /api/transactions
- POST /api/uploads (multipart, field `file`) -> upload image, returns {data:{url: "/uploads/<name>"}}
- GET /uploads/\*filepath (static file serving)

Notes & next steps

- Add authentication middleware to protect routes.
- Add validation and error handling improvements.
- Add tests and CI.

Schema / SQL

- See `backend/schema.sql` for full CREATE DATABASE, table definitions and example queries.
