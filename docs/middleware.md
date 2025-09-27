# Middleware Documentation

This project uses several middleware components to handle authentication, CORS, and rate limiting. Below are the available middleware, their purpose, and usage.

---

## 1. JWT Middleware

**File:** `http/middleware/jwt.go`

**Purpose:**
- Protects endpoints by requiring a valid Bearer JWT token in the `Authorization` header.
- Verifies the token signature using the configured `JWT_SECRET`.
- On success, stores JWT claims in the Gin context as `claims` (type: `jwt.MapClaims`).
- On failure, returns 401 Unauthorized and aborts the request.

**Usage:**
- Applied to routes that require authentication, e.g.:
  ```go
  router.POST("/api/items", middleware.JWTMiddleware(cfg), handler.Create)
  ```
- The handler can access claims via:
  ```go
  claims, _ := c.Get("claims")
  ```

**Configuration:**
- Set `JWT_SECRET` in your environment or `.env` file.

---

## 2. CORS Middleware

**File:** `http/middleware/cors.go`

**Purpose:**
- Enables Cross-Origin Resource Sharing (CORS) for frontend-backend communication.
- Sets appropriate headers for `Access-Control-Allow-Origin`, `Allow-Methods`, `Allow-Headers`, and handles preflight OPTIONS requests.
- Echoes the request's `Origin` header or uses `*` if not present.
- Aborts OPTIONS requests with status 204 (no content).

**Usage:**
- Should be applied globally to all routes:
  ```go
  router.Use(middleware.CORSMiddleware())
  ```

**Notes:**
- If you use cookies, uncomment the `Access-Control-Allow-Credentials` header in the code.

---

## 3. Rate Limiting Middleware (Login Only)

**File:** `http/middleware/ratelimit.go`

**Purpose:**
- Limits the number of login attempts per client IP to prevent brute-force attacks.
- Default: 5 requests per minute per IP.
- Returns 429 Too Many Requests if the limit is exceeded.

**Usage:**
- Apply to login/auth endpoints only:
  ```go
  router.POST("/api/auth/login", middleware.LoginRateLimiter(), handler.Login)
  ```

**Notes:**
- Uses an in-memory map to track client IPs and their rate limiters.
- Old entries are cleaned up every 5 minutes.

---

## Adding Middleware
- Middleware can be applied globally (`router.Use(...)`) or per-route/group as needed.
- Order matters: CORS should be applied before any route handlers.

## See Also
- See each middleware file in `http/middleware/` for implementation details and customization.
