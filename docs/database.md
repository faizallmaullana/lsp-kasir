# Database Schema

## Database Connection

The backend uses PostgreSQL and connects using environment variables. You can set these in a `.env` file or as system environment variables.

**Environment Variables:**

- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `5432`)
- `DB_USER` (default: `postgres`)
- `DB_PASS` (default: `password`)
- `DB_NAME` (default: `company_profile_db`)

**DSN Format:**

```
host=<DB_HOST> user=<DB_USER> password=<DB_PASS> dbname=<DB_NAME> port=<DB_PORT> sslmode=disable
```

**Where to configure:**
- See `conf/setup_env.go` and `conf/setup_database.go` for details.
- The `.env` file (if present) is loaded automatically; otherwise, system environment variables are used.
- On startup, the app will auto-migrate all entities.

**Example .env file:**
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=yourpassword
DB_NAME=company_profile_db
```

This document describes the entities, their fields, and relationships as defined in `models/entity`.

All IDs are UUID (stored as varchar(36)). Timestamps use `autoCreateTime`. Soft delete is implemented with the `is_deleted` boolean across tables.

## users

Fields:
- id_user (varchar(36), PK, unique, not null)
- email (varchar(255), unique, not null)
- password (varchar(255), not null) — hashed
- role (varchar(50), not null)
- is_deleted (boolean, default false)
- timestamp (timestamp, autoCreateTime)

Relationships:
- has many transactions (fk: transactions.id_user → users.id_user, CASCADE on update/delete)
- has many sessions (fk: sessions.id_user → users.id_user, CASCADE on update/delete)
- has many profiles (fk: profiles.id_user → users.id_user, CASCADE on update/delete)

## profiles

Fields:
- id_profile (varchar(36), PK, unique, not null)
- id_user (varchar(36), not null)
- name (varchar(100), not null)
- contact (varchar(120))
- address (varchar(255))
- photo (varchar(255))
- is_deleted (boolean, default false)
- timestamp (timestamp, autoCreateTime)

Relationships:
- belongs to users (fk: profiles.id_user → users.id_user)

## items

Fields:
- id_item (varchar(36), PK, unique, not null)
- item_name (varchar(255), not null)
- item_type (varchar(50), index)
- is_available (boolean, default true)
- price (decimal(10,2), not null)
- description (text)
- image_url (varchar(255))
- timestamp (timestamp, autoCreateTime)
- is_deleted (boolean, default false)

Relationships:
- has many pivot_items_to_transactions (fk: pivot_items_to_transactions.id_item → items.id_item, CASCADE on update/delete)

## transactions

Fields:
- id_transaction (varchar(36), PK, unique, not null)
- id_user (varchar(36), not null, index)
- buyer_contact (varchar(120))
- total_price (decimal)
- is_deleted (boolean, default false)
- timestamp (timestamp, autoCreateTime)

Relationships:
- belongs to users (fk: transactions.id_user → users.id_user)
- has many pivot_items_to_transactions (fk: pivot_items_to_transactions.id_transaction → transactions.id_transaction, CASCADE on update/delete)

## pivot_items_to_transactions

Fields:
- id_transaction (varchar(36), not null, index)
- id_item (varchar(36), not null, index)
- is_deleted (boolean, default false)
- quantity (int)
- price (decimal)

Relationships:
- belongs to transactions (fk: id_transaction → transactions.id_transaction)
- belongs to items (fk: id_item → items.id_item)

## sessions

Fields:
- id_session (varchar(36), PK, unique, not null)
- id_user (varchar(36), not null)
- is_loged_in (boolean)
- is_deleted (boolean, default false)
- timestamp (timestamp, autoCreateTime)

Relationships:
- belongs to users (fk: sessions.id_user → users.id_user)

## images

Fields:
- id_image (varchar(36), PK, unique, not null)
- file_name (varchar(255))
- content_type (varchar(120))
- size (bigint)
- data (bytea) — not used in API responses; files are stored on disk under `storages/images`
- is_deleted (boolean, default false)
- timestamp (timestamp, autoCreateTime)

Notes:
- For images, only metadata is stored in DB; actual image blobs are saved to disk and served/downloaded via the Images API.
- All soft deletes set `is_deleted = true` and most list queries filter `is_deleted = false`.
