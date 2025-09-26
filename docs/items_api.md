# Items API Documentation

## Overview
The Items API provides CRUD operations for managing inventory items in the system.

## Base URL
```
/api/items
```

## Endpoints

### 1. List Items (Paginated, Optional Type Filter)

**Endpoint:** `GET /api/items`

**Description:** Retrieves a paginated list of items.

#### Request

**Headers:**
```
Content-Type: application/json
```

**Query Parameters:**
- `count` (optional): Number of items per page (default: 10, max: 100)
- `page` (optional): Page number (default: 1)
- `type` (optional): Filter by `item_type` (exact match)

**Examples:**
```
GET /api/items?count=20&page=2
GET /api/items?type=beverage&count=10&page=1
```

#### Responses

**Success (200 OK):**
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": [
    {
      "id_item": "uuid-string",
      "item_name": "Product Name",
  "item_type": "beverage",
      "is_available": true,
      "price": 29.99,
      "description": "Product description",
      "image_url": "https://example.com/image.jpg",
      "timestamp": "2025-09-26T10:30:00Z",
      "is_deleted": false
    }
  ]
}
```

**Internal Server Error (500):**
```json
{
  "STATUS": "INTERNAL_SERVER_ERROR",
  "ERROR": "failed to list items"
}
```

---

### 2. Get Item by ID

**Endpoint:** `GET /api/items/:id`

**Description:** Retrieves a single item by its ID.

#### Request

**Headers:**
```
Content-Type: application/json
```

**Path Parameters:**
- `id` (required): Item UUID

**Example:**
```
GET /api/items/123e4567-e89b-12d3-a456-426614174000
```

#### Responses

**Success (200 OK):**
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "id_item": "123e4567-e89b-12d3-a456-426614174000",
    "item_name": "Product Name",
  "item_type": "beverage",
    "is_available": true,
    "price": 29.99,
    "description": "Product description",
    "image_url": "https://example.com/image.jpg",
    "timestamp": "2025-09-26T10:30:00Z",
    "is_deleted": false
  }
}
```

**Not Found (404):**
```json
{
  "STATUS": "NOT_FOUND",
  "MESSAGE": "item not found"
}
```

---

### 3. Create Item

**Endpoint:** `POST /api/items`

**Description:** Creates a new item.

**Authentication:** Requires valid JWT token.

#### Request

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <your_jwt_token>
```

**Body:**
```json
{
  "item_name": "string (required)",
  "item_type": "string (optional)",
  "is_available": "boolean (optional, default: true)",
  "price": "number (required)",
  "description": "string (optional)",
  "image_url": "string (optional)"
}
```

**Example:**
```json
{
  "item_name": "New Product",
  "is_available": true,
  "price": 49.99,
  "description": "A great new product",
  "image_url": "https://example.com/new-product.jpg"
}
```

#### Responses

**Success (201 Created):**
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "created",
  "DATA": {
    "id_item": "generated-uuid",
    "item_name": "New Product",
  "item_type": "beverage",
    "is_available": true,
    "price": 49.99,
    "description": "A great new product",
    "image_url": "https://example.com/new-product.jpg",
    "timestamp": "2025-09-26T10:30:00Z",
    "is_deleted": false
  }
}
```

**Bad Request (400):**
```json
{
  "STATUS": "BAD_REQUEST",
  "ERROR": "validation error message"
}
```

**Unauthorized (401):**
```json
{
  "STATUS": "UNAUTHORIZED"
}
```

**Internal Server Error (500):**
```json
{
  "STATUS": "INTERNAL_SERVER_ERROR",
  "ERROR": "failed to create item"
}
```

---

### 4. Update Item

**Endpoint:** `PUT /api/items/:id`

**Description:** Updates an existing item. Only provided fields will be updated.

**Authentication:** Requires valid JWT token.

#### Request

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <your_jwt_token>
```

**Path Parameters:**
- `id` (required): Item UUID

**Body:**
```json
{
  "item_name": "string (optional)",
  "item_type": "string (optional)",
  "is_available": "boolean (optional)",
  "price": "number (optional)",
  "description": "string (optional)",
  "image_url": "string (optional)"
}
```

**Example:**
```json
{
  "item_name": "Updated Product Name",
  "price": 59.99
}
```

#### Responses

**Success (200 OK):**
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "updated",
  "DATA": {
    "id_item": "123e4567-e89b-12d3-a456-426614174000",
    "item_name": "Updated Product Name",
  "item_type": "beverage",
    "is_available": true,
    "price": 59.99,
    "description": "Original description",
    "image_url": "https://example.com/image.jpg",
    "timestamp": "2025-09-26T10:30:00Z",
    "is_deleted": false
  }
}
```

**Bad Request (400):**
```json
{
  "STATUS": "BAD_REQUEST",
  "ERROR": "validation error message"
}
```

**Unauthorized (401):**
```json
{
  "STATUS": "UNAUTHORIZED"
}
```

**Not Found (404):**
```json
{
  "STATUS": "NOT_FOUND",
  "MESSAGE": "item not found"
}
```

**Internal Server Error (500):**
```json
{
  "STATUS": "INTERNAL_SERVER_ERROR",
  "ERROR": "failed to update item"
}
```

---

### 5. Delete Item (Soft Delete)

**Endpoint:** `DELETE /api/items/:id`

**Description:** Soft deletes an item by setting is_deleted flag to true.

**Authentication:** Requires valid JWT token.

#### Request

**Headers:**
```
Authorization: Bearer <your_jwt_token>
```

**Path Parameters:**
- `id` (required): Item UUID

**Example:**
```
DELETE /api/items/123e4567-e89b-12d3-a456-426614174000
```

#### Responses

**Success (200 OK):**
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "deleted",
  "DATA": {
    "id": "123e4567-e89b-12d3-a456-426614174000"
  }
}
```

**Unauthorized (401):**
```json
{
  "STATUS": "UNAUTHORIZED"
}
```

**Internal Server Error (500):**
```json
{
  "STATUS": "INTERNAL_SERVER_ERROR",
  "ERROR": "failed to delete item"
}
```

## Data Models

### Item Object
```json
{
  "id_item": "string (UUID)",
  "item_name": "string",
  "item_type": "string",
  "is_available": "boolean",
  "price": "number (decimal)",
  "description": "string",
  "image_url": "string",
  "timestamp": "string (ISO 8601)",
  "is_deleted": "boolean"
}
```

## Error Handling

All endpoints follow a consistent error response format:
- `STATUS`: Error type identifier
- `ERROR`: (optional) Error message details
- `MESSAGE`: (optional) Human-readable message

## Security Notes

1. Create, Update, and Delete operations require JWT authentication
2. List and Get operations are publicly accessible
3. Soft delete is implemented - items are not physically removed
4. UUIDs are automatically generated for new items
5. Pagination limits are enforced (max 100 items per page)

## Example Usage

### List Items
```bash
curl -X GET "http://localhost:8000/api/items?count=10&page=1"
```

### Get Item
```bash
curl -X GET http://localhost:8000/api/items/123e4567-e89b-12d3-a456-426614174000
```

### Create Item
```bash
curl -X POST http://localhost:8000/api/items \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_jwt_token" \
  -d '{
    "item_name": "New Product",
    "price": 49.99,
    "description": "Product description"
  }'
```

### Update Item
```bash
curl -X PUT http://localhost:8000/api/items/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer your_jwt_token" \
  -d '{
    "item_name": "Updated Name",
    "price": 59.99
  }'
```

### Delete Item
```bash
curl -X DELETE http://localhost:8000/api/items/123e4567-e89b-12d3-a456-426614174000 \
  -H "Authorization: Bearer your_jwt_token"
```
