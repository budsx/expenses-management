# Expenses Management System

## Demo Testing
[![Watch demo video](https://img.youtube.com/vi/MGztMOVYpkY/maxresdefault.jpg)](https://www.youtube.com/watch?v=MGztMOVYpkY)


## Prerequisites

- Docker dan Docker Compose

## Cara Menjalankan

Menggunakan Docker Compose

```bash
git clone https://github.com/budsx/expenses-management.git
git clone https://github.com/budsx/expenses-management-web.git
```

Atau menggunakan Makefile:

```
make compose-up
```

## Akses Aplikasi

Setelah berhasil dijalankan, Anda dapat mengakses:

- **Web**: http://localhost:3000
- **API Backend**: http://localhost:8000
- **Database Admin (Adminer)**: http://localhost:8081
  - Server: postgres
  - Username: postgres
  - Password: postgres
  - Database: expenses_management
- **RabbitMQ Management**: http://localhost:15672
  - Username: development
  - Password: development

## API Endpoints

### Health Check

- **GET** `/api/health` - Check service health
```bash
curl --location 'http://localhost:8080/api/health' \
--header 'Accept: application/json'
```

### Authentication

- **POST** `/api/auth/login` - User authentication
```bash
curl --location 'http://localhost:8080/api/auth/login' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--data '{
    "email": "manager@company.com",
    "password": "password"
}'
```

**Response Example:**
```json
{
    "message": "success",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "expires_at": 1640995200
    }
}
```

### Expense Management

- **POST** `/api/expenses` - Create new expense
```bash
curl --location 'http://localhost:8080/api/expenses' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN' \
--data '{
    "description": "Business lunch",
    "amount_idr": 100000,
    "receipt_url": "https://example.com/receipt.jpg"
}'
```

**Response Example:**
```json
{
    "message": "success",
    "data": {
        "id": 1,
        "user_id": 1,
        "description": "Business lunch",
        "amount_idr": 100000,
        "receipt_url": "https://example.com/receipt.jpg",
        "status": "pending",
        "auto_approved": false
    }
}
```

- **GET** `/api/expenses` - Get expenses with pagination and filters
```bash
# Get all expenses (page 1, 10 items per page)
curl --location 'http://localhost:8080/api/expenses?page=1&page_size=10' \
--header 'Accept: application/json' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN'

# Get expenses with filters
curl --location 'http://localhost:8080/api/expenses?page=1&page_size=10&status=1&user_id=1' \
--header 'Accept: application/json' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN'
```

**Response Example:**
```json
{
    "message": "success",
    "data": {
        "expenses": [
            {
                "id": 1,
                "user_id": 1,
                "description": "Business lunch",
                "amount_idr": 100000,
                "receipt_url": "https://example.com/receipt.jpg",
                "status": "pending",
                "auto_approved": false
            }
        ],
        "page": 1,
        "page_size": 10,
        "total": 100,
        "total_pages": 10
    }
}
```

- **GET** `/api/expenses/{id}` - Get specific expense by ID
```bash
curl --location 'http://localhost:8080/api/expenses/1' \
--header 'Accept: application/json' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN'
```

**Response Example:**
```json
{
    "message": "success",
    "data": {
        "id": 1,
        "user_id": 1,
        "description": "Business lunch",
        "amount_idr": 100000,
        "receipt_url": "https://example.com/receipt.jpg",
        "status": "pending",
        "auto_approved": false
    }
}
```

- **PUT** `/api/expenses/{id}/approve` - Approve expense
```bash
curl --location --request PUT 'http://localhost:8080/api/expenses/1/approve' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN' \
--data '{
    "expense_id": 1,
    "approver_id": 2,
    "status": 1,
    "notes": "Approved by manager"
}'
```

**Response Example:**
```json
{
    "message": "success",
    "data": {
        "message": "Expense approved successfully"
    }
}
```

- **PUT** `/api/expenses/{id}/reject` - Reject expense
```bash
curl --location --request PUT 'http://localhost:8080/api/expenses/1/reject' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Authorization: Bearer YOUR_JWT_TOKEN' \
--data '{
    "expense_id": 1,
    "approver_id": 2,
    "status": 0,
    "notes": "Receipt is not clear enough"
}'
```

**Response Example:**
```json
{
    "message": "success",
    "data": {
        "message": "Expense rejected successfully"
    }
}
```

### Error Response Format

All endpoints may return errors in the following format:
```json
{
    "error": "Validation failed",
    "message": "Invalid request"
}
```

### Authentication

Most endpoints require JWT token authentication. Include the token in the Authorization header:
```
Authorization: Bearer YOUR_JWT_TOKEN
```

### Status Codes

- **200**: Success
- **400**: Bad Request - Invalid input
- **401**: Unauthorized - Missing or invalid token
- **403**: Forbidden - Insufficient permissions
- **404**: Not Found - Resource not found
- **500**: Internal Server Error



## Menghentikan Aplikasi

```bash
# Hentikan semua services
docker compose down

# Hentikan dan hapus volumes
docker compose down -v
```

Atau menggunakan Makefile:
```bash
make compose-down
```