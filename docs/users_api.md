# Users/Profile API

Base path: `/api/profile`

Authentication: All endpoints require a valid Bearer token (JWT) unless noted. The token's `sub` claim is used as the user ID.

## GET /api/profile/me
Returns current user's basic info and their profiles.

Headers:
- Authorization: Bearer <token>

Response 200 (application/json):
```
{
  "user_id": "string",
  "email": "user@example.com",
  "role": "ADMIN|CASHIER|USER",
  "profiles": [
    {
      "id_profile": "string",
      "name": "string",
      "contact": "string",
      "address": "string",
      "image_url": "string"
    }
  ]
}
```

Errors:
- 401 UNAUTHORIZED
- 404 NOT_FOUND (user not found)

## POST /api/profile
Create a profile for the current user.

Headers:
- Authorization: Bearer <token>
- Content-Type: application/json

Body:
```
{
  "name": "string",
  "contact": "string",
  "address": "string",
  "image_url": "string"  // Optional: stored filename or full URL
}
```

Response 201:
```
{
  "MESSAGE": "SUCCESS",
  "STATUS": "created",
  "DATA": {
  "id_profile": "string",
  "id_user": "string",
  "name": "string",
  "contact": "string",
  "address": "string",
  "photo": "string",          // stored filename or URL
  "is_deleted": false,
  "timestamp": "RFC3339"
  }
}
```

Errors:
- 400 BAD_REQUEST
- 401 UNAUTHORIZED
- 500 INTERNAL_SERVER_ERROR (failed to create profile)

## PUT /api/profile/:id
Update a profile belonging to the current user.

Headers:
- Authorization: Bearer <token>
- Content-Type: application/json

Body (any field optional):
```
{
  "name": "string",
  "contact": "string",
  "address": "string",
  "image_url": "string"
}
```

Response 200:
```
{
  "MESSAGE": "SUCCESS",
  "STATUS": "updated",
  "DATA": {
  "id_profile": "string",
  "id_user": "string",
  "name": "string",
  "contact": "string",
  "address": "string",
  "photo": "string",
  "is_deleted": false,
  "timestamp": "RFC3339"
  }
}
```

Errors:
- 400 BAD_REQUEST
- 401 UNAUTHORIZED
- 404 NOT_FOUND (profile not found or not owned by user)
- 500 INTERNAL_SERVER_ERROR (failed to update profile)

## DELETE /api/profile/:id
Delete a profile belonging to the current user.

Headers:
- Authorization: Bearer <token>

Response 200:
```
{
  "MESSAGE": "SUCCESS",
  "STATUS": "deleted",
  "DATA": { "id": "string" }
}
```

Errors:
- 401 UNAUTHORIZED
- 404 NOT_FOUND (profile not found or not owned by user)
- 500 INTERNAL_SERVER_ERROR (failed to delete profile)

## POST /api/users (Admin only)
Create a new user along with an initial profile. Only available to `admin` role.

Headers:
- Authorization: Bearer <token>
- Content-Type: application/json

Body:
```
{
  "email": "user@example.com",
  "password": "min 6 chars",
  "role": "cashier|admin (optional, default cashier)",
  "profile": {
    "name": "string",
    "contact": "string",
    "address": "string",
    "image_url": "string"
  }
}
```

Response 201:
```
{
  "MESSAGE": "SUCCESS",
  "STATUS": "created",
  "DATA": {
    "user": {
      "id_user": "string",
      "email": "user@example.com",
      "role": "cashier",
      "is_deleted": false,
      "timestamp": "RFC3339"
    },
    "profile": {
      "id_profile": "string",
      "id_user": "string",
      "name": "string",
      "contact": "string",
      "address": "string",
      "photo": "string",
      "is_deleted": false,
      "timestamp": "RFC3339"
    }
  }
}
```

Errors:
- 400 BAD_REQUEST
- 401 UNAUTHORIZED (no token or non-admin)
- 500 INTERNAL_SERVER_ERROR (failed to create user/profile)

## GET /api/users (Admin only)
List users with pagination. Only available to `admin` role.

Headers:
- Authorization: Bearer <token>

Query Parameters:
- `count` (optional): items per page, default 10, max 100
- `page` (optional): page number, default 1

Example:
```
GET /api/users?count=20&page=2
```

Response 200:
```
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": [
    {
      "id_user": "string",
      "email": "user@example.com",
      "role": "admin|cashier|user",
      "is_deleted": false,
      "timestamp": "RFC3339"
    }
  ]
}
```

Errors:
- 401 UNAUTHORIZED (no token or non-admin)
- 500 INTERNAL_SERVER_ERROR (failed to list users)

## PUT /api/profile/email
Update the current user's email.

Headers:
- Authorization: Bearer <token>
- Content-Type: application/json

Body:
```
{
  "email": "user@example.com"
}
```

Response 200:
```
{
  "MESSAGE": "SUCCESS",
  "STATUS": "email updated",
  "DATA": { "email": "user@example.com" }
}
```

Errors:
- 400 BAD_REQUEST
- 401 UNAUTHORIZED
- 404 NOT_FOUND (user not found)
- 500 INTERNAL_SERVER_ERROR (failed to update email)

Notes:
- Request body uses `image_url`, but the create/update responses return the entity as-is which uses the `photo` JSON field. The `GET /me` endpoint maps profiles to `image_url`.
- `image_url`/`photo` should be the stored filename returned by the Images API or a static/public URL if you serve images statically.
- The handler uses the JWT `sub` claim as `id_user`.
