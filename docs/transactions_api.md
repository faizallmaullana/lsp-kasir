# Transactions API Documentation

## Overview
The Transactions API provides endpoints to manage cashier transactions, including listing, retrieving, creating, updating, and deleting transactions, as well as their purchased items (pivot rows).

- Totals are calculated on the server based on item price × quantity.
- Create/Update/Delete require JWT authentication.

## Base URL

```
/api/transactions
```

## Endpoints

### 1) List Transactions (Paginated)

- Method: GET
- Path: `/api/transactions`
- Description: Retrieves a paginated list of transactions.

Request
- Headers:
  - `Content-Type: application/json`
- Query Parameters:
  - `count` (optional) — items per page, default 10, max 100
  - `page` (optional) — page number, default 1

Example
```
GET /api/transactions?count=20&page=2
```

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": [
    {
      "id_transaction": "uuid-string",
      "id_user": "uuid-user-or-empty",
      "buyer_contact": "0812-xxxx",
      "total_price": 199000,
      "is_deleted": false,
      "timestamp": "2025-09-26T10:30:00Z"
    }
  ]
}
```
- 500 Internal Server Error
```json
{
  "STATUS": "INTERNAL_SERVER_ERROR",
  "ERROR": "failed to list transactions"
}
```

---

### 2) Get Transaction by ID (with Item Details)

- Method: GET
- Path: `/api/transactions/:id`
- Description: Retrieves a single transaction by its ID, including purchased items (pivot rows).

Request
- Headers:
  - `Content-Type: application/json`
- Path Params:
  - `id` — transaction UUID

Example
```
GET /api/transactions/123e4567-e89b-12d3-a456-426614174000
```

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "transaction": {
      "id_transaction": "123e4567-e89b-12d3-a456-426614174000",
      "id_user": "uuid-user-or-empty",
      "buyer_contact": "0812-xxxx",
      "total_price": 199000,
      "is_deleted": false,
      "timestamp": "2025-09-26T10:30:00Z"
    },
    "items": [
      {
        "id_item": "item-uuid-1",
        "item_name": "Product A",
        "image_url": "https://example.com/a.jpg",
        "quantity": 2,
        "price": 99000
      },
      {
        "id_item": "item-uuid-2",
        "item_name": "Product B",
        "image_url": "https://example.com/b.jpg",
        "quantity": 1,
        "price": 1000
      }
    ]
  }
}
```
- 404 Not Found
```json
{
  "STATUS": "NOT_FOUND",
  "MESSAGE": "transaction not found"
}
```

---

### 3) Create Transaction

- Method: POST
- Path: `/api/transactions`
- Auth: Bearer JWT required
- Description: Creates a new cashier transaction. The server validates item IDs and computes `total_price` based on current item prices.

Request
- Headers:
  - `Content-Type: application/json`
  - `Authorization: Bearer <your_jwt_token>`
- Body:
```json
{
  "buyer_contact": "string (optional)",
  "items": [
    { "id_item": "string (required)", "quantity": 1 },
    { "id_item": "string (required)", "quantity": 3 }
  ]
}
```

Responses
- 201 Created
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "created",
  "DATA": {
    "transaction": {
      "id_transaction": "generated-uuid",
      "id_user": "uuid-user-or-empty",
      "buyer_contact": "0812-xxxx",
      "total_price": 299000,
      "is_deleted": false,
      "timestamp": "2025-09-26T10:30:00Z"
    },
    "items": [
      {
        "id_transaction": "generated-uuid",
        "id_item": "item-uuid-1",
        "is_deleted": false,
        "quantity": 2,
        "price": 99000
      }
    ]
  }
}
```
- 400 Bad Request (validation or invalid item)
```json
{
  "STATUS": "BAD_REQUEST",
  "ERROR": "items required"
}
```
- 401 Unauthorized
```json
{
  "STATUS": "UNAUTHORIZED"
}
```

Security Notes
- The create endpoint derives `id_user` from the JWT claim `sub`. Ensure the token is provided via `Authorization: Bearer <token>` and the server has a valid `JWT_SECRET`.
- 500 Internal Server Error
```json
{
  "STATUS": "INTERNAL_SERVER_ERROR",
  "ERROR": "failed to create transaction"
}
```

---

### 4) Update Transaction

- Method: PUT
- Path: `/api/transactions/:id`
- Auth: Bearer JWT required
- Description: Updates mutable fields of a transaction (currently only `buyer_contact`).

Request
- Headers:
  - `Content-Type: application/json`
  - `Authorization: Bearer <your_jwt_token>`
- Path Params:
  - `id` — transaction UUID
- Body:
```json
{
  "buyer_contact": "string (optional)"
}
```

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "updated",
  "DATA": {
    "id_transaction": "123e4567-e89b-12d3-a456-426614174000",
    "id_user": "uuid-user-or-empty",
    "buyer_contact": "updated contact",
    "total_price": 199000,
    "is_deleted": false,
    "timestamp": "2025-09-26T10:30:00Z"
  }
}
```
- 400 Bad Request
```json
{
  "STATUS": "BAD_REQUEST",
  "ERROR": "validation error message"
}
```
- 401 Unauthorized
```json
{
  "STATUS": "UNAUTHORIZED"
}
```
- 404 Not Found
```json
{
  "STATUS": "NOT_FOUND",
  "MESSAGE": "transaction not found"
}
```
- 500 Internal Server Error
```json
{
  "STATUS": "INTERNAL_SERVER_ERROR",
  "ERROR": "failed to update transaction"
}
```

---

### 5) Delete Transaction (Soft Delete)

- Method: DELETE
- Path: `/api/transactions/:id`
- Auth: Bearer JWT required
- Description: Soft deletes a transaction (sets `is_deleted = true`) and also soft-deletes its pivot items.

Request
- Headers:
  - `Authorization: Bearer <your_jwt_token>`
- Path Params:
  - `id` — transaction UUID

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "deleted",
  "DATA": { "id": "123e4567-e89b-12d3-a456-426614174000" }
}
```
- 401 Unauthorized
```json
{
  "STATUS": "UNAUTHORIZED"
}
```
- 500 Internal Server Error
```json
{
  "STATUS": "INTERNAL_SERVER_ERROR",
  "ERROR": "failed to delete transaction"
}
```

## Data Models

Transaction
```json
{
  "id_transaction": "string (UUID)",
  "id_user": "string (UUID)",
  "buyer_contact": "string",
  "total_price": "number (decimal)",
  "is_deleted": "boolean",
  "timestamp": "string (ISO 8601)"
}
```

Transaction Item (Pivot)
```json
{
  "id_transaction": "string (UUID)",
  "id_item": "string (UUID)",
  "is_deleted": "boolean",
  "quantity": "integer",
  "price": "number (unit price at purchase time)"
}
```

## Notes & Constraints
- `total_price` is computed by the server as Σ(quantity × current item price).
- `quantity` must be >= 1; if omitted or <= 0, it defaults to 1.
- Soft delete is used; records are not physically removed.
- Pagination limit is capped at 100.
- Create/Update/Delete require a valid JWT in the `Authorization` header.

## Examples

List Transactions
```bash
curl -X GET "http://localhost:8000/api/transactions?count=10&page=1"
```

Get Transaction with Items
```bash
curl -X GET http://localhost:8000/api/transactions/123e4567-e89b-12d3-a456-426614174000
```

Create Transaction
```bash
curl -X POST http://localhost:8000/api/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "buyer_contact": "0812-xxxx",
    "items": [
      { "id_item": "item-uuid-1", "quantity": 2 },
      { "id_item": "item-uuid-2", "quantity": 1 }
    ]
  }'
```

Update Transaction
```bash
curl -X PUT http://localhost:8000/api/transactions/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_jwt_token>" \
  -d '{
    "buyer_contact": "updated contact"
  }'
```

Delete Transaction
```bash
curl -X DELETE http://localhost:8000/api/transactions/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer <your_jwt_token>"
```
