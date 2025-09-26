# Authentication API Documentation

## Overview
The Authentication API provides endpoints for user login and token refresh functionality.

## Base URL
```
/api/auth
```

## Endpoints

### 1. Login

**Endpoint:** `POST /api/auth/login`

**Description:** Authenticates a user and returns an access token along with user information.

**Rate Limiting:** This endpoint is protected by login rate limiting middleware.

#### Request

**Headers:**
```
Content-Type: application/json
```

**Body:**
```json
{
  "email": "string (required)",
  "password": "string (required)"
}
```

**Example:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Responses

**Success (200 OK):**
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "user": {
      "id_user": "string",
      "email": "string",
      "role": "string",
      "timestamp": "2025-09-26T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
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
  "ERROR": "failed to create session" // or "failed to generate token"
}
```

---

### 2. Refresh Token

**Endpoint:** `POST /api/auth/refresh`

**Description:** Refreshes an existing access token using the current valid token.

**Authentication:** Requires valid JWT token in Authorization header.

#### Request

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <your_jwt_token>
```

**Body:** No body required

#### Responses

**Success (200 OK):**
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
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
  "ERROR": "failed to refresh token"
}
```

## JWT Token Structure

The JWT token contains the following claims:
- `sub`: User ID
- `role`: User role
- `session_id`: Session identifier
- `iat`: Issued at timestamp
- `exp`: Expiration timestamp

## Error Handling

All endpoints follow a consistent error response format:
- `STATUS`: Error type identifier
- `ERROR`: (optional) Error message details

## Security Notes

1. Login endpoint is protected by rate limiting to prevent brute force attacks
2. Refresh token endpoint requires valid JWT authentication
3. All passwords are excluded from response payloads
4. JWT tokens have configurable expiration times
5. Sessions are created and tracked for each login

## Example Usage

### Login Example
```bash
curl -X POST http://localhost:8000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }'
```

### Refresh Token Example
```bash
curl -X POST http://localhost:8000/api/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```
