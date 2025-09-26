# Report API Documentation

## Overview
The Report API provides read-only endpoints to retrieve transaction reports by month/year, for today, and for an exact date. These endpoints aggregate transactions and return totals and line items.

Base prefix: `/api/reports`

Authentication: Not required (read-only). If you need to protect them, put behind JWT later.

---

## 1) Report by Month and Year

- Method: GET
- Path: `/api/reports/:bulan/:tahun`
- Description: Returns all transactions in the specified month (1-12) and year (>= 1970), with total count and sum of total_price.

Request
- Path Params:
  - `bulan` — integer month (1-12)
  - `tahun` — integer year (e.g., 2025)

Example
```
GET /api/reports/9/2025
```

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "month": 9,
    "year": 2025,
    "total_transactions": 3,
    "sum_total_price": 450000,
    "transactions": [
      {
        "id_transaction": "tx-uuid-1",
        "total_price": 150000,
        "buyer_contact": "0812-xxxx",
        "timestamp": "2025-09-26T10:30:00Z"
      }
    ]
  }
}
```
- 400 Bad Request
```json
{ "STATUS": "BAD_REQUEST", "MESSAGE": "invalid month/year" }
```
- 500 Internal Server Error
```json
{ "STATUS": "INTERNAL_SERVER_ERROR", "ERROR": "failed to query transactions" }
```

---

## 2) Report for Today

- Method: GET
- Path: `/api/reports/today`
- Description: Returns transactions that occurred today (server local time), with total count and sum.

Example
```
GET /api/reports/today
```

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "date": "2025-09-26",
    "total_transactions": 2,
    "sum_total_price": 300000,
    "transactions": [
      {
        "id_transaction": "tx-uuid-2",
        "total_price": 200000,
        "buyer_contact": "0813-xxxx",
        "timestamp": "2025-09-26T12:00:00Z"
      }
    ]
  }
}
```
- 500 Internal Server Error
```json
{ "STATUS": "INTERNAL_SERVER_ERROR", "ERROR": "failed to query transactions" }
```

---

## 3) Report by Exact Date (dd/mm/yyyy)

- Method: GET
- Path: `/api/reports/date/:dd/:mm/:yyyy`
- Description: Returns transactions for the exact calendar date.

Request
- Path Params:
  - `dd` — day (1-31)
  - `mm` — month (1-12)
  - `yyyy` — year (>= 1970)

Example
```
GET /api/reports/date/26/09/2025
```

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "date": "2025-09-26",
    "total_transactions": 1,
    "sum_total_price": 150000,
    "transactions": [
      {
        "id_transaction": "tx-uuid-3",
        "total_price": 150000,
        "buyer_contact": "0812-xxxx",
        "timestamp": "2025-09-26T08:00:00Z"
      }
    ]
  }
}
```
- 400 Bad Request
```json
{ "STATUS": "BAD_REQUEST", "MESSAGE": "invalid dd/mm/yyyy" }
```
- 500 Internal Server Error
```json
{ "STATUS": "INTERNAL_SERVER_ERROR", "ERROR": "failed to query transactions" }
```

---

## 4) Today Summary (Revenue, Counts)

- Method: GET
- Path: `/api/reports/today/summary`
- Description: Returns today's revenue (sum_total_price), total transactions, and total products sold.

Example
```
GET /api/reports/today/summary
```

Responses
- 200 OK
```json
{
  "MESSAGE": "SUCCESS",
  "STATUS": "OK",
  "DATA": {
    "date": "2025-09-26",
    "total_transactions": 5,
    "total_products_sold": 18,
    "sum_total_price": 750000
  }
}
```
- 500 Internal Server Error
```json
{ "STATUS": "INTERNAL_SERVER_ERROR", "ERROR": "failed to query transactions" }
```

## Data Shapes

Monthly Response
```json
{
  "month": 9,
  "year": 2025,
  "total_transactions": 0,
  "sum_total_price": 0,
  "transactions": [
    { "id_transaction": "", "total_price": 0, "buyer_contact": "", "timestamp": "" }
  ]
}
```

Today/Exact Date Response
```json
{
  "date": "YYYY-MM-DD",
  "total_transactions": 0,
  "sum_total_price": 0,
  "transactions": [
    { "id_transaction": "", "total_price": 0, "buyer_contact": "", "timestamp": "" }
  ]
}
```

## Notes
- These endpoints currently perform in-memory filtering on paginated results; for large datasets, consider DB-side date range queries.
- Timestamps are ISO 8601 strings.